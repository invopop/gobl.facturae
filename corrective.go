package facturae

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
)

// Corrective is used to represent the details of a previous invoice that
// this one serves to modify or replace.
// Details here are a bit sketchy as in practical usage, most companies
// don't demand this level of detail.
type Corrective struct {
	InvoiceNumber               string `xml:",omitempty"`
	InvoiceSeriesCode           string `xml:",omitempty"`
	ReasonCode                  string
	ReasonDescription           string
	TaxPeriod                   *PeriodDates
	CorrectionMethod            string
	CorrectionMethodDescription string
	AdditionalReasonDescription string `xml:",omitempty"`
	InvoiceIssueDate            string `xml:",omitempty"`
}

var methodDefs = &tax.KeyDefinition{
	// Take from "CorrectionMethodType"
	Key: "correction-method",
	Name: i18n.String{
		i18n.EN: "Correction Method",
		i18n.ES: "Método de rectificación",
	},
	Codes: []*tax.CodeDefinition{
		{
			// Corrective
			Code: cbc.Code("01"),
			Name: i18n.String{
				i18n.EN: "Full items",
				i18n.ES: "Rectificación íntegra",
			},
		},
		{
			// Credit or Debit notes
			Code: cbc.Code("02"),
			Name: i18n.String{
				i18n.EN: "Corrected items only",
				i18n.ES: "Rectificación por diferencias",
			},
		},
		{
			// Unused
			Code: cbc.Code("03"),
			Name: i18n.String{
				i18n.EN: "Bulk deal",
				i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
			},
		},
		{
			// Unused
			Code: cbc.Code("04"),
			Name: i18n.String{
				i18n.EN: "Authorized by the Tax Agency",
				i18n.ES: "Autorizadas por la Agencia Tributaria",
			},
		},
	},
}

func newCorrective(inv *bill.Invoice) *Corrective {
	p := inv.Preceding[0]
	c := &Corrective{
		InvoiceNumber:               p.Code,
		InvoiceSeriesCode:           p.Series,
		InvoiceIssueDate:            p.IssueDate.String(),
		AdditionalReasonDescription: p.Reason,
	}

	// Add period information
	period := p.Period
	if period == nil {
		// if no period is given, use the issue date as base
		period = &cal.Period{
			Start: *p.IssueDate,
			End:   *p.IssueDate,
		}
	}
	c.TaxPeriod = newPeriodDates(period)

	// determine the reason from the extension
	r := es.New()
	kd := r.ExtensionDef(es.ExtKeyFacturaECorrection)
	cc := p.Ext[es.ExtKeyFacturaECorrection]
	row := kd.CodeDef(cc.Code())
	if row != nil {
		c.ReasonCode = row.Code.String()
		c.ReasonDescription = row.Name[i18n.ES]
	}

	// determine the method from the type of invoice
	cm := cbc.Code("04") // default assume "authorized"
	switch inv.Type {
	case bill.InvoiceTypeCorrective:
		cm = cbc.Code("01")
	case bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote:
		cm = cbc.Code("02")
	}
	md := methodDefs.CodeDef(cm)
	c.CorrectionMethod = md.Code.String()
	c.CorrectionMethodDescription = md.Name[i18n.ES]

	return c
}
