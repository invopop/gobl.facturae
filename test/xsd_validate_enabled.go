//go:build xsdvalidate

package test

import (
	"path/filepath"
	"testing"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/require"
)

// ValidateAgainstSchema validates the given data against the FacturaE schema.
func ValidateAgainstSchema(t *testing.T, data []byte) {
	// Load the XSD schema
	schemaPath := filepath.Join(GetTestPath(), "schema", "facturaev3_2_2.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	require.NoError(t, err)

	// Parse the XML document
	doc, err := libxml2.ParseString(string(data))
	require.NoError(t, err)

	// Validate the document against the schema
	err = schema.Validate(doc)
	if err != nil {
		if validationErr, ok := err.(xsd.SchemaValidationError); ok {
			t.Errorf("Schema validation failed with %d errors:", len(validationErr.Errors()))
			for _, e := range validationErr.Errors() {
				t.Errorf("  - %s", e.Error())
			}
			t.FailNow()
		} else {
			require.NoError(t, err)
		}
	}
}
