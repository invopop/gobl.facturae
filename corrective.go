package facturae

import (
	"github.com/invopop/gobl/bill"
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

func newCorrective(inv *bill.Invoice) *Corrective {
	p := inv.Preceding[0]
	c := &Corrective{
		InvoiceNumber:               p.Code,
		InvoiceSeriesCode:           p.Series,
		InvoiceIssueDate:            p.IssueDate.String(),
		TaxPeriod:                   newPeriodDates(p.Period), // Todo - autocalculate
		AdditionalReasonDescription: p.Reason,
	}

	// find the reason
	if len(p.Corrections) > 0 {
		kd := correctionKeyDefinition(p.Corrections[0])
		if kd != nil {
			c.ReasonCode = kd.Meta[es.KeyFacturaE]
			c.ReasonDescription = kd.Desc[i18n.ES]
		}
	}

	// find the method
	kd := correctionMethodKeyDefinition(p.CorrectionMethod)
	if kd != nil {
		c.CorrectionMethod = kd.Meta[es.KeyFacturaE]
		c.CorrectionMethodDescription = kd.Desc[i18n.ES]
	}

	return c
}

func correctionKeyDefinition(key cbc.Key) *tax.KeyDefinition {
	for _, row := range regime.Preceding.Corrections {
		if row.Key == key {
			return row
		}
	}
	return nil
}

func correctionMethodKeyDefinition(key cbc.Key) *tax.KeyDefinition {
	for _, row := range regime.Preceding.CorrectionMethods {
		if row.Key == key {
			return row
		}
	}
	return nil
}
