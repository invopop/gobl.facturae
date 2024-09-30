package facturae

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// Address stores location information
// FIXME: para direcciones extranjeras hay un campo diferentes PostCodeAndTown
type Address struct {
	Address     string `xml:",omitempty"`
	PostCode    string `xml:",omitempty"`
	Town        string `xml:",omitempty"`
	Province    string `xml:",omitempty"`
	CountryCode string `xml:",omitempty"`
}

func newAddress(address *org.Address) *Address {
	return &Address{
		Address:     addressLine1(address),
		PostCode:    address.Code,
		Town:        address.Locality,
		Province:    address.Region,
		CountryCode: countryCode(address.Country),
	}
}

func addressLine1(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}

	return address.Street +
		", " + address.Number +
		addressMaybe(address.Block) +
		addressMaybe(address.Floor) +
		addressMaybe(address.Door)
}

func addressMaybe(element string) string {
	if element != "" {
		return ", " + element
	}
	return ""
}

// ES is the ISO Alpha-3 code used for Spain
const ES string = "ESP"

func countryCode(country l10n.ISOCountryCode) string {
	cd := l10n.Countries().Code(country.Code())
	if cd == nil {
		return ""
	}
	return cd.Alpha3
}
