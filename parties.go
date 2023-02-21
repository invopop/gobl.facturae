package facturae

import (
	"github.com/invopop/gobl/org"
)

// Parties stores emitter and receiver of the invoice
type Parties struct {
	Seller *Party `xml:"SellerParty"`
	Buyer  *Party `xml:"BuyerParty"`
}

// Party stores the info for one of the participants of
// the invoice
type Party struct {
	TaxID       *TaxID       `xml:"TaxIdentification,omitempty"`
	LegalEntity *LegalEntity `xml:"LegalEntity,omitempty"`
	Individual  *Individual  `xml:"Individual,omitempty"`
}

// LegalEntity stores party info for companies, associations, ...
type LegalEntity struct {
	CorporateName   string
	AddressInSpain  *Address        `xml:",omitempty"`
	OverseasAddress *Address        `xml:",omitempty"`
	ContactDetails  *ContactDetails `xml:",omitempty"`
}

// Individual stores party info for people
type Individual struct {
	Name            string
	FirstSurname    string          `xml:",omitempty"`
	SecondSurname   string          `xml:",omitempty"`
	AddressInSpain  *Address        `xml:",omitempty"`
	OverseasAddress *Address        `xml:",omitempty"`
	ContactDetails  *ContactDetails `xml:",omitempty"`
}

// NewIndividual creates info for people
func NewIndividual(party *org.Party) *Individual {
	if len(party.People) == 0 {
		return nil
	}

	personInfo := party.People[0]
	individual := &Individual{
		Name:           personInfo.Name.Given,
		FirstSurname:   personInfo.Name.Surname,
		SecondSurname:  personInfo.Name.Surname2,
		ContactDetails: newContactDetails(party),
	}

	if len(party.Addresses) > 0 {
		address := newAddress(party.Addresses[0])
		if address.CountryCode == ES {
			individual.AddressInSpain = address
		} else {
			individual.OverseasAddress = address
		}
	}

	return individual
}

// NewLegalEntity creates info for companies, associations, ...
func NewLegalEntity(party *org.Party) *LegalEntity {
	entity := &LegalEntity{
		CorporateName:  party.Name,
		ContactDetails: newContactDetails(party),
	}

	if len(party.Addresses) > 0 {
		address := newAddress(party.Addresses[0])
		if address.CountryCode == ES {
			entity.AddressInSpain = address
		} else {
			entity.OverseasAddress = address
		}
	}

	return entity
}
