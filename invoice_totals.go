package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
)

// InvoiceTotals contains the summary of the amounts in an invoice
type InvoiceTotals struct {
	TotalGrossAmount            string     // Without included taxes
	GeneralDiscounts            *Discounts `xml:",omitempty"`
	GeneralSurcharges           *Charges   `xml:",omitempty"`
	TotalGeneralDiscounts       string     `xml:",omitempty"`
	TotalGeneralSurcharges      string     `xml:",omitempty"`
	TotalGrossAmountBeforeTaxes string
	TotalTaxOutputs             string
	TotalTaxesWithheld          string
	InvoiceTotal                string
	Subsidies                   *Subsidies            `xml:",omitempty"`
	PaymentsOnAccount           *PaymentsOnAccount    `xml:",omitempty"`
	ReimbursableExpenses        *ReimbursableExpenses `xml:",omitempty"`
	TotalOutstandingAmount      string                // InvoiceTotal - (Total subvenciones + TotalPaymentsOnAccount)
	TotalPaymentsOnAccount      string                `xml:",omitempty"`
	TotalExecutableAmount       string
	TotalReimbursableExpenses   string `xml:",omitmpty"`
}

// Subsidies is currently a placeholder, as we don't use it yet.
type Subsidies struct {
	Subsidy []*Subsidy
}

// Subsidy is a single subsidy entry
type Subsidy struct {
	SubsidyDescription string
	SubsidyRate        string `xml:",omitempty"`
	SubsidyAmount      string
}

// PaymentsOnAccount stores payments made in advance.
type PaymentsOnAccount struct {
	PaymentOnAccount []*PaymentOnAccount
}

// PaymentOnAccount is a single advance payment.
type PaymentOnAccount struct {
	PaymentOnAccountDate   cal.Date
	PaymentOnAccountAmount string
}

// ReimbursableExpenses is a group of reimbursable expenses.
type ReimbursableExpenses struct {
	ReimbursableExpenses []*ReimbursableExpense
}

// ReimbursableExpense stores info about a single reimbursable expense.
type ReimbursableExpense struct {
	ReimbursableExpensesSellerParty *TaxID    `xml:",omitempty"`
	ReimbursableExpensesBuyerParty  *TaxID    `xml:",omitempty"` // never used!
	IssueDate                       *cal.Date `xml:",omitempty"`
	InvoiceNumber                   string    `xml:",omitempty"`
	InvoiceSeriesCode               string    `xml:",omitempty"`
	ReimbursableExpensesAmount      string
}

func newInvoiceTotals(invoice *bill.Invoice) *InvoiceTotals {
	totals := invoice.Totals

	sum := totals.Sum
	if totals.TaxIncluded != nil {
		sum = sum.Subtract(*totals.TaxIncluded)
	}

	due := totals.Payable
	outstanding := totals.Payable
	if totals.Due != nil {
		due = *totals.Due
		outstanding = *totals.Due
	}
	if totals.Outlays != nil {
		outstanding = outstanding.Subtract(*totals.Outlays)
	}

	xmlTotals := &InvoiceTotals{
		TotalGrossAmount:            sum.String(),
		GeneralDiscounts:            newDiscounts(invoice.Discounts),
		GeneralSurcharges:           newCharges(invoice.Charges),
		TotalGrossAmountBeforeTaxes: totals.Total.String(),
		InvoiceTotal:                totals.TotalWithTax.String(),
		TotalOutstandingAmount:      outstanding.String(),
		TotalExecutableAmount:       due.String(),
	}

	if totals.Discount != nil {
		xmlTotals.TotalGeneralDiscounts = totals.Discount.String()
	}
	if totals.Charge != nil {
		xmlTotals.TotalGeneralSurcharges = totals.Charge.String()
	}
	if totals.Advances != nil {
		xmlTotals.TotalPaymentsOnAccount = totals.Advances.String()
	}
	if totals.Outlays != nil {
		xmlTotals.TotalReimbursableExpenses = totals.Outlays.String()
	}

	if invoice.Payment != nil {
		xmlTotals.setAdvances(invoice.Payment.Advances)
	}
	xmlTotals.setOutlays(invoice.Outlays)
	xmlTotals.setTaxTotals(totals.Taxes)

	return xmlTotals
}

func (it *InvoiceTotals) setTaxTotals(taxes *tax.Total) {
	regular := num.MakeAmount(0, 2)
	retained := num.MakeAmount(0, 2)
	for _, ct := range taxes.Categories {
		if ct.Retained {
			retained = retained.Add(ct.Amount)
		} else {
			regular = regular.Add(ct.Amount)
			if ct.Surcharge != nil {
				regular = regular.Add(*ct.Surcharge)
			}
		}
	}
	it.TotalTaxOutputs = regular.String()
	it.TotalTaxesWithheld = retained.String()
}

func (it *InvoiceTotals) setAdvances(advances []*pay.Advance) {
	regular := make([]*PaymentOnAccount, 0)
	grants := make([]*Subsidy, 0)
	for _, a := range advances {
		if a.Grant {
			g := &Subsidy{
				SubsidyDescription: a.Description,
				SubsidyAmount:      a.Amount.String(),
			}
			if a.Percent != nil {
				g.SubsidyRate = a.Percent.String()
			}
			grants = append(grants, g)
		} else {
			na := &PaymentOnAccount{
				PaymentOnAccountAmount: a.Amount.String(),
			}
			if a.Date != nil {
				na.PaymentOnAccountDate = *a.Date
			}
		}
	}
	if len(regular) > 0 {
		it.PaymentsOnAccount = &PaymentsOnAccount{
			PaymentOnAccount: regular,
		}
	}
	if len(grants) > 0 {
		it.Subsidies = &Subsidies{
			Subsidy: grants,
		}
	}
}

func (it *InvoiceTotals) setOutlays(outlays []*bill.Outlay) {
	os := make([]*ReimbursableExpense, 0)
	for _, o := range outlays {
		no := &ReimbursableExpense{
			IssueDate:                  o.Date,
			InvoiceNumber:              o.Code,
			InvoiceSeriesCode:          o.Series,
			ReimbursableExpensesAmount: o.Amount.String(),
		}
		if o.Supplier != nil {
			tid := o.Supplier.TaxID
			if tid != nil {
				no.ReimbursableExpensesSellerParty = NewTaxID(tid.Code, tid.Country)
			}
		}
		os = append(os, no)
	}
	if len(os) > 0 {
		it.ReimbursableExpenses = &ReimbursableExpenses{
			ReimbursableExpenses: os,
		}
	}
}
