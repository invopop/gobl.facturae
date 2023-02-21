package facturae_test

import (
	"testing"

	"github.com/invopop/gobl.facturae/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartiesInfo(t *testing.T) {
	t.Run("should contain the provider party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat-irpf.json")
		require.NoError(t, err)

		assert.Equal(t, "F", doc.Parties.Seller.TaxID.PersonTypeCode)
		assert.Equal(t, "R", doc.Parties.Seller.TaxID.ResidenceTypeCode)
		assert.Equal(t, "37221735F", doc.Parties.Seller.TaxID.TaxIdentificationNumber)
		assert.Equal(t, "Maria Remedios", doc.Parties.Seller.Individual.Name)
		assert.Equal(t, "Sanchez", doc.Parties.Seller.Individual.FirstSurname)
		assert.Equal(t, "Nuñez", doc.Parties.Seller.Individual.SecondSurname)
		assert.Equal(t, "Campo Real, 74", doc.Parties.Seller.Individual.AddressInSpain.Address)
		assert.Equal(t, "28023", doc.Parties.Seller.Individual.AddressInSpain.PostCode)
		assert.Equal(t, "Torrejón De La Calzada", doc.Parties.Seller.Individual.AddressInSpain.Town)
		assert.Equal(t, "Madrid", doc.Parties.Seller.Individual.AddressInSpain.Province)
		assert.Equal(t, "ESP", doc.Parties.Seller.Individual.AddressInSpain.CountryCode)
		assert.Equal(t, "msohrjnb3@caramail.com", doc.Parties.Seller.Individual.ContactDetails.ElectronicMail)
		assert.Equal(t, "+34612123123", doc.Parties.Seller.Individual.ContactDetails.Telephone)
	})
}

func TestPartiesInfoCustomer(t *testing.T) {
	t.Run("should contain the customer party info", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json")
		require.NoError(t, err)

		assert.Equal(t, "J", doc.Parties.Buyer.TaxID.PersonTypeCode)
		assert.Equal(t, "R", doc.Parties.Buyer.TaxID.ResidenceTypeCode)
		assert.Equal(t, "B77436020", doc.Parties.Buyer.TaxID.TaxIdentificationNumber)
		assert.Equal(t, "Moniward Sl", doc.Parties.Buyer.LegalEntity.CorporateName)
		assert.Equal(t, "Plaza Horno, 35", doc.Parties.Buyer.LegalEntity.AddressInSpain.Address)
		assert.Equal(t, "45083", doc.Parties.Buyer.LegalEntity.AddressInSpain.PostCode)
		assert.Equal(t, "Nombela", doc.Parties.Buyer.LegalEntity.AddressInSpain.Town)
		assert.Equal(t, "Toledo", doc.Parties.Buyer.LegalEntity.AddressInSpain.Province)
		assert.Equal(t, "ESP", doc.Parties.Buyer.LegalEntity.AddressInSpain.CountryCode)
		assert.Equal(t, "bfn25xf3p@lycos.co.uk", doc.Parties.Buyer.LegalEntity.ContactDetails.ElectronicMail)
	})

}
