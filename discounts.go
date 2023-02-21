package facturae

import (
	"github.com/invopop/gobl/bill"
)

// Discounts list discounts applied on a line or document level.
type Discounts struct {
	Items []*Discount `xml:"Discount"`
}

// Discount is used in general and line discounts. Unlike in GOBL,
// FacturaE does not differentiate between the two, which could potentially
// lead to confusion for handling taxes.
type Discount struct {
	Reason string `xml:"DiscountReason"`
	Rate   string `xml:"DiscountRate,omitempty"`
	Amount string `xml:"DiscountAmount"`
}

func newDiscounts(discounts []*bill.Discount) *Discounts {
	if len(discounts) == 0 {
		return nil
	}
	m := &Discounts{
		Items: make([]*Discount, len(discounts)),
	}
	for i, v := range discounts {
		m.Items[i] = newDiscount(v)
	}
	return m
}

func newDiscount(discount *bill.Discount) *Discount {
	nd := &Discount{
		Reason: discount.Reason,
		Amount: discount.Amount.String(),
	}
	if discount.Percent != nil {
		nd.Rate = discount.Percent.StringWithoutSymbol()
	}
	return nd
}

func newLineDiscounts(lds []*bill.LineDiscount) *Discounts {
	if len(lds) == 0 {
		return nil
	}
	nlds := &Discounts{
		Items: make([]*Discount, len(lds)),
	}
	for i, v := range lds {
		nlds.Items[i] = newLineDiscount(v)
	}
	return nlds
}

func newLineDiscount(ld *bill.LineDiscount) *Discount {
	if ld.Percent != nil {
		return &Discount{
			Reason: ld.Reason,
			Rate:   ld.Percent.StringWithoutSymbol(),
			Amount: ld.Amount.String(),
		}
	}
	return &Discount{
		Reason: ld.Reason,
		Amount: ld.Amount.String(),
	}
}
