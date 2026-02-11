package facturae

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// Address stores location information for addresses in Spain
type Address struct {
	Address     string `xml:",omitempty"`
	PostCode    string `xml:",omitempty"`
	Town        string `xml:",omitempty"`
	Province    string `xml:",omitempty"`
	CountryCode string `xml:",omitempty"`
}

// OverseasAddress stores location information for addresses outside Spain
type OverseasAddress struct {
	Address         string `xml:",omitempty"`
	PostCodeAndTown string `xml:",omitempty"`
	Province        string `xml:",omitempty"`
	CountryCode     string `xml:",omitempty"`
}

func newAddress(address *org.Address) *Address {
	return &Address{
		Address:     addressLine1(address),
		PostCode:    address.Code.String(),
		Town:        address.Locality,
		Province:    address.Region,
		CountryCode: countryCode(address.Country),
	}
}

func newOverseasAddress(address *org.Address) *OverseasAddress {
	postCodeAndTown := address.Locality
	if address.Code != "" {
		postCodeAndTown = address.Code.String() + " " + address.Locality
	}
	return &OverseasAddress{
		Address:         addressLine1(address),
		PostCodeAndTown: postCodeAndTown,
		Province:        address.Region,
		CountryCode:     countryCode(address.Country),
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
