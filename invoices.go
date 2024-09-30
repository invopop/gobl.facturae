package facturae

import (
	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Invoices contains info about a batch of invoices.
// In our case there will only be one invoice per batch
type Invoices struct {
	List []*Invoice `xml:"Invoice"`
}

// Invoice contains info about a single invoice
type Invoice struct {
	InvoiceHeader    *InvoiceHeader
	InvoiceIssueData *InvoiceIssueData
	TaxesOutputs     *TaxSummary `xml:",omitempty"`
	TaxesWithheld    *TaxSummary `xml:",omitempty"`
	InvoiceTotals    *InvoiceTotals
	Items            *Items
	PaymentDetails   *PaymentDetails `xml:",omitempty"`
	LegalLiterals    *LegalLiterals  `xml:",omitempty"`
	AdditionalData   *AdditionalData `xml:",omitempty"`
}

// InvoiceHeader contains general information of a single invoice
type InvoiceHeader struct {
	InvoiceNumber       string
	InvoiceSeriesCode   string `xml:",omitempty"`
	InvoiceDocumentType string
	InvoiceClass        string
	Corrective          *Corrective `xml:",omitempty"`
}

// PeriodDates is used in corrective tax periods to define a date
// range.
type PeriodDates struct {
	StartDate string
	EndDate   string
}

// InvoiceIssueData contains information about dates, lang and currencies
type InvoiceIssueData struct {
	IssueDate                    string
	OperationDate                string       `xml:",omitempty"`
	PlaceOfIssue                 string       `xml:",omitempty"`
	InvoiceingPeriod             *PeriodDates `xml:",omitempty"`
	InvoiceCurrencyCode          string
	ExchangeRateDetails          *ExchangeRateDetails `xml:",omitempty"` // TODO!
	TaxCurrencyCode              string
	LanguageName                 string // FIXME: [JUANJO] are we going to support multiple languages?
	InvoiceDescription           string `xml:",omitempty"` // TODO
	ReceiverTransactionReference string `xml:",omitempty"` // TODO
	FileReference                string `xml:",omitempty"` // TODO
	ReceiverContractReference    string `xml:",omitempty"` // TODO
}

// LegalLiterals contains an array of legal texts to add to the Invoice in certain situations.
type LegalLiterals struct {
	LegalReference []string
}

// ExchangeRateDetails describes how to exchange from the invoices currency
// to euros.
type ExchangeRateDetails struct {
	ExchangeRate     string
	ExchangeRateDate string
}

// NewInvoice creates a new invoice with facturae format
func (d *Document) newInvoice(invoice *bill.Invoice) *Invoice {
	valueDate := &invoice.IssueDate
	if invoice.ValueDate != nil {
		valueDate = invoice.ValueDate
	}
	xmlInvoice := &Invoice{
		InvoiceHeader: newInvoiceHeader(invoice),
		InvoiceIssueData: &InvoiceIssueData{
			IssueDate:           invoice.IssueDate.String(),
			OperationDate:       valueDate.String(),
			InvoiceCurrencyCode: string(invoice.Currency),
			TaxCurrencyCode:     string(invoice.Currency),
			LanguageName:        "es",
		},
		InvoiceTotals:  newInvoiceTotals(invoice),
		AdditionalData: newAdditionalData(invoice),
		PaymentDetails: newPaymentDetails(invoice.Payment),
		LegalLiterals:  newLegalLiterals(invoice),
		Items:          newItems(invoice.Lines, invoice.Totals.Taxes),
	}
	xmlInvoice.setTaxes(invoice.Totals.Taxes)

	return xmlInvoice
}

// setTaxes performs a set of steps to convert the GOBL tax list into something
// that FacturaE expects.
func (inv *Invoice) setTaxes(taxes *tax.Total) {
	if taxes == nil {
		return
	}
	regular := make([]*Tax, 0)
	retained := make([]*Tax, 0)
	// First loop for bases
	for _, ct := range taxes.Categories {
		for _, rt := range ct.Rates {
			percent := rt.Percent
			if percent == nil {
				percent = num.NewPercentage(0, 0)
			}
			tax := &Tax{
				Code:   categoryTaxCodeMap[ct.Code],
				Rate:   percent.StringWithoutSymbol(),
				Base:   makeAmount(rt.Base),
				Amount: makeAmount(rt.Amount),
			}
			if ct.Retained {
				retained = append(retained, tax)
			} else {
				if rt.Surcharge != nil {
					st := rt.Surcharge
					p := st.Percent
					p = p.Rescale(4) // we need 2 decimal places
					tax.Surcharge = p.StringWithoutSymbol()
					v := makeAmount(st.Amount)
					tax.SurchargeAmount = &v
				}
				regular = append(regular, tax)
			}
		}
	}
	if len(regular) > 0 {
		inv.TaxesOutputs = &TaxSummary{
			List: regular,
		}
	}
	if len(retained) > 0 {
		inv.TaxesWithheld = &TaxSummary{
			List: retained,
		}
	}
}

func newInvoiceHeader(inv *bill.Invoice) *InvoiceHeader {
	h := &InvoiceHeader{
		InvoiceNumber:     inv.Code.String(),
		InvoiceSeriesCode: inv.Series.String(),
	}

	if inv.Tax != nil {
		h.InvoiceDocumentType = inv.Tax.Ext[facturae.ExtKeyDocType].String()
		h.InvoiceClass = inv.Tax.Ext[facturae.ExtKeyInvoiceClass].String()
	}

	// Only one preceding document currently supported
	if len(inv.Preceding) > 0 {
		h.Corrective = newCorrective(inv)
	}

	return h
}

func newLegalLiterals(inv *bill.Invoice) *LegalLiterals {
	lits := make([]string, 0)
	if len(lits) == 0 {
		return nil
	}
	for _, n := range inv.Notes {
		if n.Key == cbc.NoteKeyLegal {
			lits = append(lits, n.Text)
		}
	}
	return &LegalLiterals{
		LegalReference: lits,
	}
}

func newPeriodDates(p *cal.Period) *PeriodDates {
	if p == nil {
		return nil
	}
	return &PeriodDates{
		StartDate: p.Start.String(),
		EndDate:   p.End.String(),
	}
}
