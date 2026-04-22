package facturae

import "github.com/invopop/gobl/num"

// Amount wraps the Facturae AmountType. Per Orden HAP/1650/2015 Annex II
// rule 6, monetary amounts must not carry more than two decimals at either
// line or invoice level. Unit prices and tax rates are the only exceptions
// and therefore do not flow through this type.
type Amount struct {
	TotalAmount       string
	EquivalentInEuros string `xml:",omitempty"` // not used yet!
}

// amountString renders a GOBL amount as a Facturae monetary string, capped at
// two decimals per Orden HAP/1650/2015 Annex II rule 6. Do not use for unit
// prices or rate/percentage fields, which are exempt from the cap.
func amountString(a num.Amount) string {
	return a.RescaleDown(2).MinimalString()
}

// amountTwoDecimalString renders an amount with exactly two decimals, as
// required by Facturae XSD fields typed DoubleTwoDecimalType (pattern
// -?[0-9]+\.[0-9]{2}), e.g. InstallmentAmount. Values are rounded down from
// higher precision and padded up from lower precision.
func amountTwoDecimalString(a num.Amount) string {
	return a.Rescale(2).String()
}

// makeAmount always provides an AmountType, capped at two decimals.
func makeAmount(a num.Amount) Amount {
	return Amount{
		TotalAmount: amountString(a),
	}
}
