package facturae

import (
	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
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

var methodDefs = &cbc.KeyDefinition{
	// Take from "CorrectionMethodType"
	Key: "correction-method",
	Name: i18n.String{
		i18n.EN: "Correction Method",
		i18n.ES: "Método de rectificación",
	},
	Values: []*cbc.ValueDefinition{
		{
			// Corrective
			Value: "01",
			Name: i18n.String{
				i18n.EN: "Full items",
				i18n.ES: "Rectificación íntegra",
			},
		},
		{
			// Credit or Debit notes
			Value: "02",
			Name: i18n.String{
				i18n.EN: "Corrected items only",
				i18n.ES: "Rectificación por diferencias",
			},
		},
		{
			// Unused
			Value: "03",
			Name: i18n.String{
				i18n.EN: "Bulk deal",
				i18n.ES: "Rectificación por descuento por volumen de operaciones durante un periodo",
			},
		},
		{
			// Unused
			Value: "04",
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
		InvoiceNumber:               p.Code.String(),
		InvoiceSeriesCode:           p.Series.String(),
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
	kd := tax.ExtensionForKey(facturae.ExtKeyCorrection)
	cc := p.Ext[facturae.ExtKeyCorrection]
	row := kd.ValueDef(cc.String())
	if row != nil {
		c.ReasonCode = row.Value
		c.ReasonDescription = row.Name[i18n.ES]
	}

	// determine the method from the type of invoice
	cm := "04" // default assume "authorized"
	switch inv.Type {
	case bill.InvoiceTypeCorrective:
		cm = "01"
	case bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote:
		cm = "02"
	}
	md := methodDefs.ValueDef(cm)
	c.CorrectionMethod = md.Value
	c.CorrectionMethodDescription = md.Name[i18n.ES]

	return c
}
