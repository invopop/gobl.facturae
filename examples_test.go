package facturae_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	facturae "github.com/invopop/gobl.facturae"
	"github.com/invopop/gobl.facturae/test"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateOut = flag.Bool("update", false, "Update the XML files in the test/data/out directory")

func TestXMLGeneration(t *testing.T) {
	schema, err := loadSchema()
	require.NoError(t, err)

	examples, err := lookupExamples()
	require.NoError(t, err)

	opts := prepareOptions()

	for _, example := range examples {
		name := fmt.Sprintf("should convert %s example file successfully", example)

		t.Run(name, func(t *testing.T) {
			data, err := convertExample(example, opts...)
			require.NoError(t, err)

			outPath := filepath.Join(test.GetDataPath(), "out", strings.TrimSuffix(example, ".json")+".xml")

			if *updateOut {
				errs := validateDoc(schema, data)
				for _, e := range errs {
					assert.NoError(t, e)
				}
				if len(errs) > 0 {
					assert.Fail(t, "Invalid XML:\n"+string(data))
					return
				}

				err = os.WriteFile(outPath, data, 0644)
				require.NoError(t, err)

				return
			}

			expected, err := os.ReadFile(outPath)

			require.False(t, os.IsNotExist(err), "output file %s missing, run tests with `--update` flag to create", filepath.Base(outPath))
			require.NoError(t, err)
			require.Equal(t, string(expected), string(data), "output file %s does not match, run tests with `--update` flag to update", filepath.Base(outPath))
		})
	}
}

func loadSchema() (*xsd.Schema, error) {
	schemaPath := filepath.Join(test.GetTestPath(), "schema", "facturaev3_2_2.xsd")
	schema, err := xsd.ParseFromFile(schemaPath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func lookupExamples() ([]string, error) {
	examples, err := filepath.Glob(filepath.Join(test.GetDataPath(), "*.json"))
	if err != nil {
		return nil, err
	}

	for i, example := range examples {
		examples[i] = filepath.Base(example)
	}

	return examples, nil
}

func prepareOptions() []facturae.Option {
	opts := []facturae.Option{
		facturae.WithThirdParty(test.ThirdParty()),
	}
	return opts
}

func convertExample(example string, opts ...facturae.Option) ([]byte, error) {
	doc, err := test.NewDocumentFrom(example, opts...)
	if err != nil {
		return nil, err
	}

	return doc.Bytes()
}

func validateDoc(schema *xsd.Schema, doc []byte) []error {
	xmlDoc, err := libxml2.ParseString(string(doc))
	if err != nil {
		return []error{err}
	}

	err = schema.Validate(xmlDoc)
	if err != nil {
		return err.(xsd.SchemaValidationError).Errors()
	}

	return nil
}
