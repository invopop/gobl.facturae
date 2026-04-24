package facturae_test

import (
	"testing"

	facturae "github.com/invopop/gobl.facturae"
	"github.com/invopop/gobl.facturae/test"
	"github.com/invopop/gobl/bill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceInvoiceTotals(t *testing.T) {
	t.Run("should contain the total amounts block", func(t *testing.T) {
		doc, err := test.LoadGOBL("invoice-vat.json")
		require.NoError(t, err)

		xmlInvoice := doc.Invoices.List[0]

		assert.Nil(t, err)
		assert.Equal(t, "4202.00", xmlInvoice.InvoiceTotals.TotalGrossAmountBeforeTaxes)
		assert.Equal(t, "", xmlInvoice.InvoiceTotals.TotalGeneralDiscounts)
		assert.Equal(t, "4202.00", xmlInvoice.InvoiceTotals.TotalGrossAmount)
		assert.Equal(t, "882.42", xmlInvoice.InvoiceTotals.TotalTaxOutputs)
		assert.Equal(t, "0.00", xmlInvoice.InvoiceTotals.TotalTaxesWithheld)
		assert.Equal(t, "5084.42", xmlInvoice.InvoiceTotals.InvoiceTotal)
		assert.Equal(t, "5084.42", xmlInvoice.InvoiceTotals.TotalOutstandingAmount)
		assert.Equal(t, "5084.42", xmlInvoice.InvoiceTotals.TotalExecutableAmount)
	})

	t.Run("falls back to invoice IssueDate when an advance has no date", func(t *testing.T) {
		env, err := test.LoadTestEnvelope("invoice-with-advance.json")
		require.NoError(t, err)

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		require.NotEmpty(t, inv.Payment.Advances)
		inv.Payment.Advances[0].Date = nil

		doc, err := facturae.NewInvoice(env)
		require.NoError(t, err)

		poa := doc.Invoices.List[0].InvoiceTotals.PaymentsOnAccount
		require.NotNil(t, poa)
		require.Len(t, poa.PaymentOnAccount, 1)
		assert.Equal(t, inv.IssueDate.String(), poa.PaymentOnAccount[0].PaymentOnAccountDate.String())
	})
}
