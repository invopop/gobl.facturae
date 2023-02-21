package facturae

import (
	"fmt"

	"github.com/invopop/gobl/bill"
)

const (
	version = "3.2.2"
)

// List of accepted Mobility Types
const (
	ModalityIndividual = "I" // Individual
	ModalityBatch      = "L" // Batch or multiple invoices in on document
)

// List of accepted Issuer Types
const (
	InvoiceIssuerSeller     = "EM"  // Seller
	InvoiceIssuerBuyer      = " RE" // Buyer
	InvoiceIssuerThirdParty = "TE"  // Third party
)

// FileHeader contains the FileHeader element of a facturae invoice
type FileHeader struct {
	SchemaVersion     string
	Modality          string
	InvoiceIssuerType string      // Who is signing the invoice?
	ThirdParty        *ThirdParty `xml:",omitempty"`
	Batch             *Batch
}

// ThirdParty is used to identify an intermediary building and issuing the invoice.
type ThirdParty struct {
	TaxIdentification *TaxID
	LegalEntity       *LegalEntity
}

// Batch would contain info about a group of invoices that are submitted
// at the same time, but we only allow one invoice at a time
type Batch struct {
	BatchIdentifier        string
	InvoicesCount          int
	TotalInvoicesAmount    Amount
	TotalOutstandingAmount Amount
	TotalExecutableAmount  Amount
	InvoiceCurrencyCode    string
}

func newFileHeader(invoice *bill.Invoice, tp *ThirdParty) *FileHeader {
	outstanding := invoice.Totals.Payable
	if invoice.Totals.Outlays != nil {
		outstanding = outstanding.Subtract(*invoice.Totals.Outlays)
	}
	ii := InvoiceIssuerSeller
	if tp != nil {
		ii = InvoiceIssuerThirdParty
	}
	return &FileHeader{
		SchemaVersion:     version,
		Modality:          ModalityIndividual,
		InvoiceIssuerType: ii,
		ThirdParty:        tp,
		Batch: &Batch{
			BatchIdentifier:        fmt.Sprintf("%s-%s", invoice.Supplier.TaxID.Code, invoice.Code),
			InvoicesCount:          1,
			TotalInvoicesAmount:    makeAmount(invoice.Totals.TotalWithTax),
			TotalOutstandingAmount: makeAmount(outstanding),
			TotalExecutableAmount:  makeAmount(invoice.Totals.Payable),
			InvoiceCurrencyCode:    string(invoice.Currency),
		},
	}
}
