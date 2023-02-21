package facturae

import (
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
)

// AdditionalData stores notes that the provider of the invoice
// estimates are needed
type AdditionalData struct {
	InvoiceAdditionalInformation string `xml:"InvoiceAdditionalInformation,omitempty"`
}

func newAdditionalData(inv *bill.Invoice) *AdditionalData {
	txt := make([]string, 0)
	for _, n := range inv.Notes {
		// Skip legal codes, we have already dealt with them in Legal Literals
		if n.Key != cbc.NoteKeyLegal {
			txt = append(txt, n.Text)
		}
	}
	if len(txt) == 0 {
		return nil
	}
	return &AdditionalData{
		InvoiceAdditionalInformation: strings.Join(txt, "; "),
	}
}
