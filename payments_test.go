package facturae_test

import (
	"testing"

	"github.com/invopop/gobl.facturae/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentsInfo(t *testing.T) {
	t.Run("should contain payment details information", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json")
		require.NoError(t, err)

		xmlInvoice := doc.Invoices.List[0]

		pd := xmlInvoice.PaymentDetails
		require.NotNil(t, pd)
		require.NotNil(t, pd.Installments)
		require.NotNil(t, pd.Installments[0])
		pi := pd.Installments[0]
		assert.Equal(t, "2021-12-30", pi.InstallmentDueDate)
		assert.Equal(t, "5084.42", pi.InstallmentAmount)
		assert.Equal(t, "04", pi.PaymentMeans)
		assert.Equal(t, "ES25 0188 2570 7185 4470 4761", pi.AccountToBeCredited.IBAN)
		assert.Contains(t, pi.CollectionAdditionalInformation, "payment term note")
	})

	t.Run("should render InstallmentAmount net of advances", func(t *testing.T) {
		// invoice-with-advance: 4,840 EUR invoice, 1,000 EUR non-grant
		// advance, one installment. GOBL's due_date records the full
		// 4,840 (gross); Facturae's InstallmentAmount must be the
		// 3,840 outstanding after the advance.
		doc, err := test.LoadGOBL("invoice-with-advance.json")
		require.NoError(t, err)

		pd := doc.Invoices.List[0].PaymentDetails
		require.NotNil(t, pd)
		require.Len(t, pd.Installments, 1)
		assert.Equal(t, "3840.00", pd.Installments[0].InstallmentAmount)
	})
}
