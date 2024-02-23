package facturae

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/xmldsig"
)

// Namespaces used for FacturaE. DSig stuff is handled in the signatures.
const (
	NamespaceFacturaE = "http://www.facturae.gob.es/formato/Versiones/Facturaev3_2_2.xml"
)

// XAdES Signer Roles used for FacturaE
const (
	XAdESSupplier   xmldsig.XAdESSignerRole = "supplier"
	XAdESCustomer   xmldsig.XAdESSignerRole = "customer"
	XAdESThirdParty xmldsig.XAdESSignerRole = "third party"
)

var (
	xadesConfig = &xmldsig.XAdESConfig{
		Role:        XAdESThirdParty,
		Description: "Factura Electrónica",
		Policy: &xmldsig.XAdESPolicyConfig{
			URL:         "http://www.facturae.es/politica_de_firma_formato_facturae/politica_de_firma_formato_facturae_v3_1.pdf",
			Description: "Política de Firma FacturaE v3.1",
			Algorithm:   "http://www.w3.org/2000/09/xmldsig#sha1",
			Hash:        "Ohixl6upD6av8N7pEvDABhEL6hM=",
		},
	}
)

// Document is a pseudo-model for containing the XML document being created.
type Document struct {
	env     *gobl.Envelope `xml:"-"` // Envelope to convert.
	invoice *bill.Invoice  `xml:"-"` // Invoice contained in envelope.

	XMLName     xml.Name `xml:"fe:Facturae"`
	FeNamespace string   `xml:"xmlns:fe,attr"`

	FileHeader *FileHeader
	Parties    *Parties
	Invoices   *Invoices
	Signature  *xmldsig.Signature `xml:"ds:Signature,omitempty"`
}

type options struct {
	certificate  *xmldsig.Certificate
	addTimestamp bool
	thirdParty   *ThirdParty
}

// Option defines a callback configuration option used to customize the
// conversion to XML process
type Option func(*options)

// WithCertificate will use the provided certificate to sign the XML document.
func WithCertificate(cert *xmldsig.Certificate) Option {
	return func(opts *options) {
		opts.certificate = cert
	}
}

// WithTimestamp will ensure the XML document is timestamped
func WithTimestamp(val bool) Option {
	return func(opts *options) {
		opts.addTimestamp = val
	}
}

// WithThirdParty adds optional information about who is manipulating and signing
// the document.
func WithThirdParty(tp *ThirdParty) Option {
	return func(opts *options) {
		opts.thirdParty = tp
	}
}

// LoadGOBL will build a FacturaE Document from the source buffer
func LoadGOBL(src io.Reader, opts ...Option) (*Document, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return nil, err
	}

	return NewInvoice(env, opts...)
}

// NewInvoice expects the base envelope and provides a new Document
// containing the XML version.
func NewInvoice(env *gobl.Envelope, opts ...Option) (*Document, error) {
	// prepare our options
	xmlOpts := new(options)
	for _, opt := range opts {
		opt(xmlOpts)
	}
	invoice, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, errors.New("expected an invoice")
	}

	// Make sure we're dealing with raw data
	var err error
	invoice, err = invoice.RemoveIncludedTaxes()
	if err != nil {
		return nil, fmt.Errorf("removing taxes: %w", err)
	}

	// Invert if we're dealing with a credit note
	if invoice.Type == bill.InvoiceTypeCreditNote {
		invoice.Invert()
		if err := invoice.Calculate(); err != nil {
			return nil, fmt.Errorf("inverting invoice: %w", err)
		}
	}

	// Basic document headers
	d := &Document{
		env:         env,
		invoice:     invoice,
		FeNamespace: NamespaceFacturaE,
		Parties: &Parties{
			Seller: &Party{},
			Buyer:  &Party{},
		},
	}

	d.Parties.Seller.TaxID = NewTaxID(invoice.Supplier.TaxID.Code, invoice.Supplier.TaxID.Country)
	if d.Parties.Seller.TaxID.PersonTypeCode == "F" {
		d.Parties.Seller.Individual = NewIndividual(invoice.Supplier)
	} else {
		d.Parties.Seller.LegalEntity = NewLegalEntity(invoice.Supplier)
	}

	d.Parties.Buyer.TaxID = NewTaxID(invoice.Customer.TaxID.Code, invoice.Customer.TaxID.Country)
	if d.Parties.Buyer.TaxID.PersonTypeCode == "F" {
		d.Parties.Buyer.Individual = NewIndividual(invoice.Customer)
	} else {
		d.Parties.Buyer.LegalEntity = NewLegalEntity(invoice.Customer)
	}

	d.Invoices = &Invoices{
		List: []*Invoice{d.newInvoice(invoice)},
	}

	d.FileHeader = newFileHeader(invoice, xmlOpts.thirdParty)

	if xmlOpts.certificate != nil {
		data, err := d.Canonical()
		if err != nil {
			return nil, fmt.Errorf("converting to canonincal format: %w", err)
		}

		sigopts := []xmldsig.Option{
			xmldsig.WithDocID(env.Head.UUID.String()),
			xmldsig.WithXAdES(xadesConfig),
			xmldsig.WithCertificate(xmlOpts.certificate),
		}
		if xmlOpts.addTimestamp {
			sigopts = append(sigopts, xmldsig.WithTimestamp(xmldsig.TimestampFreeTSA))
		}
		sig, err := xmldsig.Sign(data, sigopts...)
		if err != nil {
			return nil, err
		}
		d.Signature = sig
	}

	return d, nil
}

// Buffer returns a byte buffer representation of the complete XML document.
func (d *Document) Buffer() (*bytes.Buffer, error) {
	return d.buffer(xml.Header)
}

// String converts a struct representation to its string representation
func (d *Document) String() (string, error) {
	buf, err := d.Buffer()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Bytes returns the XML document bytes
func (d *Document) Bytes() ([]byte, error) {
	buf, err := d.Buffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Document) buffer(base string) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)
	// data, err := xml.MarshalIndent(d, "", "  ") // not compatible with certificates
	data, err := xml.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshal document: %w", err)
	}
	if _, err := buf.Write(data); err != nil {
		return nil, fmt.Errorf("writing to buffer: %w", err)
	}
	return buf, nil
}

// Canonical converts a struct representation of facturae to its
// canonical representation as defined in https://www.w3.org/TR/2001/REC-xml-c14n-20010315
// (for a simpler explanation look at https://www.di-mgt.com.au/xmldsig-c14n.html)
// This is used when we need to create a hash for signing, timestamping, ...
func (d *Document) Canonical() ([]byte, error) {
	buf, err := d.buffer("")
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
