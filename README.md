# GOBL to FacturaE Toolkit

Convert GOBL documents into the Spain's FacturaE format.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [GNU Affero General Public License v3.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library, at this time we will require acceptance of transferring copyright to Invopop Ltd.

[![Lint](https://github.com/invopop/gobl.facturae/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl.facturae/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl.facturae/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl.facturae/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl.facturae)](https://goreportcard.com/report/github.com/invopop/gobl.facturae)
[![GoDoc](https://godoc.org/github.com/invopop/gobl.facturae?status.svg)](https://godoc.org/github.com/invopop/gobl.facturae)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl.facturae)

## Usage

### Go

There are a couple of entry points to build a new FacturaE document. If you already have a GOBL Envelope available in Go, you could convert and output to a data file like this:

```golang
doc, err := facturae.NewInvoice(env)
if err != nil {
    panic(err)
}

data, err := doc.Bytes()
if err != nil {
    panic(err)
}

if err = os.WriteFile("./test.xml", data, 0644); err != nil {
    panic(err)
}
```

If you're loading from a file, you can use the `LoadGOBL` convenience method:

```golang
doc, err := facturae.LoadGOBL(file)
if err != nil {
    panic(err)
}
// do something with doc
```

Outputting to a FacturaE XML is most useful when the document is signed. Use a certificate to sign the document as follows:

```golang
// import from github.com/invopop/xmldsig
cert, err := xmldsig.LoadCertificate(filename, password)
if err != nil {
    panic(err)
}

doc, err := facturae.NewInvoice(env, facturae.WithCertificate(cert))
if err != nil {
    panic(err)
}
```

### CLI

The command line interface can be useful for situations when you're using a language other than Golang in your application.

```bash
# install example
```

Simply provide the input GOBL JSON file and output to a file or another application:

```bash
./gobl.facturae convert input.json output.xml
```

If you have a digital certificate, run with:

```bash
./gobl.facturae convert -c cert.p12 -p password input.json > output.xml
```

The command also supports pipes:

```bash
cat input.json > ./gobl.facturae > output.xml
```

## Notes

- To validate the XML output and digital certificates, use https://face.gob.es/es/facturas/validar-visualizar-facturas
- In most cases Go structures have been written using the same naming from the XML style document. This means names are not repeated in tags and generally makes it a bit easier map the XML output to the internal structures.

## Current Conversion Limitations

The FacturaE format is quite complex due to the number of local requirements in Spain. We've put a lot of effort

- _Payment Information_: FacturaE requires each payment instruction to have a Due Date. The GOBL invoice allows these details to be independent. If you require payment instructions to appear on a FacturaE document, there must be a due date.

## Integration Tests

There are some integration and XML generation tests available in the `/test` path. To execute them, there are two [Magefile](https://magefile.org/) commands.

The first will convert YAML source data into GOBL JSON documents:

```
mage -v convertYAML
```

The second will generate the FacturaE XML documents from the GOBL sources, using the digital certificates that are available in the `/test/certificates` path:

```
mage -v convertXML
```

Sample data sources are contained in the `/test/data` directory. YAML documents are stored in the Git repository, but JSON and XML must be generated using the above commands.
