package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/invopop/gobl"
	"github.com/invopop/gobl.facturae"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/xmldsig"
)

const certificateFile = "facturae.p12"
const certificatePassword = "invopop"

var signingKey = dsig.NewES256Key()

// Random third party example data.
var thirdParty = &facturae.ThirdParty{
	TaxIdentification: &facturae.TaxID{
		PersonTypeCode:          "J",
		ResidenceTypeCode:       "R",
		TaxIdentificationNumber: "B23103039",
	},
	LegalEntity: &facturae.LegalEntity{
		CorporateName: "Hypeprop S.L.",
		AddressInSpain: &facturae.Address{
			Address:     "Calle Campo Real 74",
			PostCode:    "28023",
			Town:        "Torrejón De La Calzada",
			Province:    "Madrid",
			CountryCode: "ESP",
		},
	},
}

// LoadGOBL loads a GoBL test file into structs
func LoadGOBL(name string, opts ...facturae.Option) (*facturae.Document, error) {
	envelopeReader, _ := os.Open(GetDataPath() + name)
	doc, err := facturae.LoadGOBL(envelopeReader, opts...)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ConvertYAML takes the YAML test data and converts into useful json gobl documents.
func ConvertYAML() error {
	var files []string
	err := filepath.Walk(GetDataPath(), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range files {
		fmt.Printf("processing file: %v\n", path)

		// attempt to load and convert
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		// TODO: gobl should have a more direct way to do this soon!
		env := new(gobl.Envelope)
		if err := yaml.Unmarshal(data, env); err != nil {
			return fmt.Errorf("invalid contents: %w", err)
		}

		if err := env.Calculate(); err != nil {
			return fmt.Errorf("failed to complete: %w", err)
		}

		if err := env.Sign(signingKey); err != nil {
			return fmt.Errorf("failed to sign the doc: %w", err)
		}

		// Output to the filesystem
		np := strings.TrimSuffix(path, filepath.Ext(path)) + ".json"
		out, err := json.MarshalIndent(env, "", "	")
		if err != nil {
			return fmt.Errorf("marshalling output: %w", err)
		}
		if err := os.WriteFile(np, out, 0644); err != nil {
			return fmt.Errorf("saving file data: %w", err)
		}

		fmt.Printf("wrote file: %v\n", np)
	}

	return nil
}

// ConvertToXML takes the .json invoices generated previously and converts them
// into XML FacturaE documents.
func ConvertToXML() error {
	var files []string
	err := filepath.Walk(GetDataPath(), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return err
	}

	cert, err := LoadCertificate()
	if err != nil {
		return err
	}

	// Load the config so we have an example third party

	for _, file := range files {
		fmt.Printf("processing file: %v\n", file)

		doc, err := LoadGOBL(file, facturae.WithCertificate(cert), facturae.WithThirdParty(thirdParty))
		if err != nil {
			return err
		}

		data, err := doc.Bytes()
		if err != nil {
			return fmt.Errorf("extracting document bytes: %w", err)
		}

		np := strings.TrimSuffix(file, filepath.Ext(file)) + ".xml"
		err = os.WriteFile(GetDataPath()+"/"+np, data, 0644)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}

	return nil
}

// GetDataPath returns the path where test can find data files
// to be used in tests
func GetDataPath() string {
	return getRootFolder() + "/test/data/"
}

func getRootFolder() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}

// GetCertificatesPath return the path where a test can find crypto
// certificates that can be used in tests
func GetCertificatesPath() string {
	return getRootFolder() + "/test/certificates/"
}

// LoadCertificate will return the standard test certificate info
func LoadCertificate() (*xmldsig.Certificate, error) {
	f := path.Join(GetCertificatesPath(), certificateFile)
	return xmldsig.LoadCertificate(f, certificatePassword)
}
