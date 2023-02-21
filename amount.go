package facturae

import "github.com/invopop/gobl/num"

// Amount provides a wrapper around our regular amounts so that we can
// include the `TotalAmount` tag. Eventually this could also handle currency
// conversion.
type Amount struct {
	TotalAmount       string
	EquivalentInEuros string `xml:",omitempty"` // not used yet!
}

// makeAmount always provides an AmountType
func makeAmount(a num.Amount) Amount {
	return Amount{
		TotalAmount: a.MinimalString(),
	}
}
