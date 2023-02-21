package facturae

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/xmldsig"
)

// Standard error messages for the repo
var (
	ErrCertificateNotFound = errors.New("certificate-not-found")
)

// CertificateRepo is an implementation of the CertificateRepo
// that loads the certificate from a local file
type CertificateRepo struct {
	certs map[string]*xmldsig.Certificate
}

// NewCertificateRepo instantiates a new repository of certificates stored
// in a local directory. Each certificate found will be parsed
// and stored in a map of file names without extension to certificate instances.
func NewCertificateRepo(path string) (*CertificateRepo, error) {
	certs := make(map[string]*xmldsig.Certificate)

	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if filepath.Ext(file) == ".p12" {
			cert, err := xmldsig.LoadCertificate(file, "invopop")
			if err != nil {
				return fmt.Errorf("loading certificate %s: %w", file, err)
			}
			fn := filepath.Base(file)
			fn = strings.TrimSuffix(fn, ".p12")
			certs[fn] = cert
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CertificateRepo{certs: certs}, nil
}

// Get returns the local pre-loaded certificate
func (repo *CertificateRepo) Get(id string) (*xmldsig.Certificate, error) {
	cert, ok := repo.certs[id]
	if !ok {
		return nil, ErrCertificateNotFound
	}
	return cert, nil
}
