package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/pay"
)

// PaymentDetails contains info about how the invoice should
// be paid
type PaymentDetails struct {
	Installments []*Installment `xml:"Installment"`
}

// Installment contains info about each of the payment terms
type Installment struct {
	InstallmentDueDate              string
	InstallmentAmount               string
	PaymentMeans                    string       `xml:",omitempty"`
	AccountToBeCredited             *BankAccount `xml:",omitempty"`
	DebitReconciliationReference    string       `xml:",omitempty"`
	AccountToBeDebited              *BankAccount `xml:",omitempty"`
	CollectionAdditionalInformation string       `xml:",omitempty"`
}

// BankAccount contains info needed to pay by transfer or direct debit
type BankAccount struct {
	IBAN                  string   `xml:",omitempty"`
	AccountNumber         string   `xml:",omitempty"`
	BranchInSpainAddress  *Address `xml:",omitempty"`
	OverseasBranchAddress *Address `xml:",omitempty"`
	BIC                   string   `xml:",omitempty"`
}

var facturaePaymentMethodCodes = map[pay.MethodKey]string{
	pay.MethodKeyCash:           "01",
	pay.MethodKeyDirectDebit:    "02",
	pay.MethodKeyCreditTransfer: "04",
	pay.MethodKeyCard:           "19",
	pay.MethodKeyOnline:         "13",
}

func newPaymentDetails(paymentInfo *bill.Payment) *PaymentDetails {
	if paymentInfo == nil {
		return nil
	}
	terms := paymentInfo.Terms
	if terms == nil {
		return nil
	}
	if len(terms.DueDates) == 0 {
		return nil
	}

	instructions := paymentInfo.Instructions

	xmlInstallments := make([]*Installment, len(terms.DueDates))
	for i, installment := range terms.DueDates {
		xmlInstallment := &Installment{
			InstallmentDueDate:              installment.Date.String(),
			InstallmentAmount:               installment.Amount.String(),
			PaymentMeans:                    facturaePaymentMethodCodes[instructions.Key],
			CollectionAdditionalInformation: mergeNotes(paymentInfo.Terms.Notes, installment.Notes),
		}

		if instructions.Key == pay.MethodKeyCreditTransfer {
			if len(instructions.CreditTransfer) > 0 {
				xmlInstallment.AccountToBeCredited = newCreditBankAccount(instructions.CreditTransfer[0])
			}
		}

		if instructions.Key == pay.MethodKeyDirectDebit {
			xmlInstallment.AccountToBeDebited = newDebitBankAccount(instructions.DirectDebit)
			xmlInstallment.DebitReconciliationReference = instructions.DirectDebit.Ref
		}

		if instructions.Key == pay.MethodKeyOnline {
			if len(instructions.Online) > 0 {
				if len(xmlInstallment.CollectionAdditionalInformation) > 0 {
					xmlInstallment.CollectionAdditionalInformation += "\n"
				}
				xmlInstallment.CollectionAdditionalInformation += instructions.Online[0].Address
			}
		}

		xmlInstallments[i] = xmlInstallment
	}

	return &PaymentDetails{
		Installments: xmlInstallments,
	}
}

func newCreditBankAccount(info *pay.CreditTransfer) *BankAccount {
	if info == nil {
		return nil
	}
	return &BankAccount{
		IBAN:          info.IBAN,
		BIC:           info.BIC,
		AccountNumber: info.Number,
	}
}

func newDebitBankAccount(info *pay.DirectDebit) *BankAccount {
	if info == nil {
		return nil
	}
	return &BankAccount{
		AccountNumber: info.Account,
	}
}

func mergeNotes(termNotes string, installmentNotes string) string {
	notes := ""
	if termNotes != "" {
		notes = termNotes
		if installmentNotes != "" {
			notes += "\n"
		}
	}

	return notes + installmentNotes
}
