Cap monetary amounts at two decimals per Orden HAP/1650/2015

FACe was rejecting invoices (error 2600) where line-level `TaxAmount` carried four decimals — e.g. `686.8743`. Orden HAP/1650/2015 Annex II, regla 6, governs decimal precision for every monetary amount in a Facturae document. Each change below is paired with the exact text of the regulation that motivates it.

## Helpers (`amount.go`)

Two helpers, one per Facturae XSD pattern, each encoding a piece of regla 6 at a single boundary:

- `amountString` — `a.RescaleDown(2).MinimalString()` — for fields typed `DoubleUpToEightDecimalType` (pattern `-?[0-9]+(\.[0-9]{1,8})?`). Only trims when precision exceeds two; integer-valued amounts stay integer (e.g. `0` stays `0`, not `0.00`) because the XSD pattern permits any count from 0 to 8 decimals and the regulation only requires the value be *rounded* to two.
- `amountTwoDecimalString` — `a.Rescale(2).String()` — for fields typed `DoubleTwoDecimalType` (pattern `-?[0-9]+\.[0-9]{2}`). Rounds down from higher precision and pads up from lower so the rendering always matches the XSD requirement of exactly two decimals.

`makeAmount` now delegates to `amountString`, so the `AmountType` wrapper (`TaxableBase`, `TaxAmount`, `EquivalenceSurchargeAmount`, `TotalInvoicesAmount`, `TotalOutstandingAmount`, `TotalExecutableAmount`) shares the same cap.

## Line-level amounts (`items.go`, `discounts.go`, `charges.go`)

`TotalCost`, `GrossAmount`, and every line-level `DiscountAmount` / `ChargeAmount` now go through `amountString`:

> En las facturas emitidas en euros, se validará que los importes totales de las líneas relativos al coste total sean numéricos y estén redondeados, de acuerdo con el método común de redondeo, a dos decimales, como resultado del producto del número de unidades por el precio unitario, y que los importes brutos de las líneas sean el resultado de restar del coste total los descuentos, y de sumar los cargos, todos ellos numéricos y con dos decimales.
>
> — Anexo II, regla 6.a

`TaxableBase` and `TaxAmount` at line level are already capped via `makeAmount`, covered by the same rule's catch-all:

> Asimismo se validará que el resto de importes a nivel de línea, con excepción del importe unitario, vengan expresados en euros con dos decimales.
>
> — Anexo II, regla 6.a

`UnitPriceWithoutTax` is deliberately left at its original precision because regla 6.a carves it out:

> No se consideran importes los tipos impositivos o los porcentajes a aplicar que, al igual que el importe unitario, podrán tener los decimales que permita el formato Facturae.
>
> — Anexo II, regla 6.a

## Invoice-level amounts (`invoice_totals.go`)

`TotalGrossAmount`, `TotalGrossAmountBeforeTaxes`, `InvoiceTotal`, `TotalOutstandingAmount`, `TotalExecutableAmount`, `TotalGeneralDiscounts`, `TotalGeneralSurcharges`, `TotalPaymentsOnAccount`, `TotalTaxOutputs`, `TotalTaxesWithheld`, `SubsidyAmount`, and `PaymentOnAccountAmount` all now go through `amountString`:

> En las facturas emitidas en euros, se validará que el total importe bruto de la factura sea numérico y a dos decimales, por suma de los importes brutos de las líneas. Asimismo se validará que el resto de importes vengan expresados en euros con dos decimales.
>
> — Anexo II, regla 6.b

Tax rates and percentage fields (`TaxRate`, `EquivalenceSurcharge`, `DiscountRate`, `ChargeRate`, `SubsidyRate`) are untouched because regla 6.b reiterates the same exemption:

> No se consideran importes los tipos impositivos o los porcentajes a aplicar que podrán tener los decimales que permita el formato Facturae.
>
> — Anexo II, regla 6.b

## Installments (`payments.go`)

`InstallmentAmount` now goes through `amountTwoDecimalString`. The regulation treats it as any other monetary amount and therefore requires two decimals:

> Asimismo se validará que el resto de importes vengan expresados en euros con dos decimales.
>
> — Anexo II, regla 6.b

The XSD type for this element is `DoubleTwoDecimalType`, whose pattern mandates the value *display* with exactly two decimal characters — stripping trailing zeros is invalid — so this field uses the padding helper rather than the trimming one. For EUR invoices the wire format is unchanged; the helper hardens the path against any future higher-precision or zero-decimal-currency input.

## Tests and fixtures

Unit-test assertions in `items_test.go` and `invoice_totals_test.go` were updated for values that now round (`3140.4967` → `3140.5`) or shed trailing zeros (`4202.00` → `4202`, `200.00` → `200`). The seven example XMLs under `test/data/out/` were regenerated via

```
go test -run TestXMLGeneration -update -tags xsdvalidate .
```

so XSD validation against `facturaev3_2_2.xsd` gates the updated fixtures. Full suite passes both with and without the tag.
