package facturae

import (
	"github.com/invopop/gobl/bill"
)

// Charges contains a list of charge instances. This is used both at a line
// and document level as per the FacturaE specification. Unfortunately, this
// implies that taxes will not be applied in the context of a global charge.
// We recommend avoiding using this in Spanish invoices.
type Charges struct {
	Items []*Charge `xml:"Charge"`
}

// Charge defines the basic data for a single surcharge.
type Charge struct {
	Reason string `xml:"ChargeReason"`
	Rate   string `xml:"ChargeRate,omitempty"`
	Amount string `xml:"ChargeAmount"`
}

func newCharges(charges []*bill.Charge) *Charges {
	if len(charges) == 0 {
		return nil
	}
	m := &Charges{
		Items: make([]*Charge, len(charges)),
	}
	for i, v := range charges {
		m.Items[i] = newCharge(v)
	}
	return m
}

func newCharge(charge *bill.Charge) *Charge {
	nc := &Charge{
		Reason: charge.Reason,
		Amount: charge.Amount.String(),
	}
	if charge.Percent != nil {
		nc.Rate = charge.Percent.StringWithoutSymbol()
	}
	return nc
}

func newLineCharges(lds []*bill.LineCharge) *Charges {
	if len(lds) == 0 {
		return nil
	}
	nlds := &Charges{
		Items: make([]*Charge, len(lds)),
	}
	for i, v := range lds {
		nlds.Items[i] = newLineCharge(v)
	}
	return nlds
}

func newLineCharge(lc *bill.LineCharge) *Charge {
	nc := &Charge{
		Reason: lc.Reason,
		Amount: lc.Amount.String(),
	}
	if lc.Percent != nil {
		nc.Rate = lc.Percent.StringWithoutSymbol()
	}
	return nc
}
