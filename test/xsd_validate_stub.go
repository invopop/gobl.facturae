//go:build !xsdvalidate

package test

import "testing"

// ValidateAgainstSchema skips schema validation unless built with `-tags xsdvalidate`.
func ValidateAgainstSchema(t *testing.T, _ []byte) {
	t.Skip("XSD validation requires libxml2; run with `go test -tags xsdvalidate ./...`")
}
