package facturae

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

// TaxSummary contains a list with the info of each tax type
type TaxSummary struct {
	List []*Tax `xml:"Tax"`
}

// Tax contains the info of a particular tax type
type Tax struct {
	Code            string  `xml:"TaxTypeCode"`
	Rate            string  `xml:"TaxRate"`
	Base            Amount  `xml:"TaxableBase"`
	Amount          Amount  `xml:"TaxAmount"`
	Surcharge       string  `xml:"EquivalenceSurcharge,omitempty"` // must be with two decimal places
	SurchargeAmount *Amount `xml:"EquivalenceSurchargeAmount,omitempty"`
}

var categoryTaxCodeMap = map[cbc.Code]string{
	tax.CategoryVAT:    "01",
	es.TaxCategoryIPSI: "02", // Ceuta, Melilla
	es.TaxCategoryIGIC: "03", // Canary Islands
	es.TaxCategoryIRPF: "04",
}
