package tlscommon

import (
	"crypto/tls"

	"github.com/joeshaw/multierror"
)

// Config defines the user configurable options in the yaml file.
type Config struct {
	Enabled          *bool                   `config:"enabled"`
	VerificationMode TLSVerificationMode     `config:"verification_mode"` // one of 'none', 'full'
	Versions         []TLSVersion            `config:"supported_protocols"`
	CipherSuites     []tlsCipherSuite        `config:"cipher_suites"`
	CAs              []string                `config:"certificate_authorities"`
	Certificate      CertificateConfig       `config:",inline"`
	CurveTypes       []tlsCurveType          `config:"curve_types"`
	Renegotiation    tlsRenegotiationSupport `config:"renegotiation"`
}

// LoadTLSConfig will load a certificate from config with all TLS based keys
// defined. If Certificate and CertificateKey are configured, client authentication
// will be configured. If no CAs are configured, the host CA will be used by go
// built-in TLS support.
func LoadTLSConfig(config *Config) (*TLSConfig, error) {
	if !config.IsEnabled() {
		return nil, nil
	}

	fail := multierror.Errors{}
	logFail := func(es ...error) {
		for _, e := range es {
			if e != nil {
				fail = append(fail, e)
			}
		}
	}

	var cipherSuites []uint16
	for _, suite := range config.CipherSuites {
		cipherSuites = append(cipherSuites, uint16(suite))
	}

	var curves []tls.CurveID
	for _, id := range config.CurveTypes {
		curves = append(curves, tls.CurveID(id))
	}

	cert, err := LoadCertificate(&config.Certificate)
	logFail(err)

	cas, errs := LoadCertificateAuthorities(config.CAs)
	logFail(errs...)

	// fail, if any error occurred when loading certificate files
	if err = fail.Err(); err != nil {
		return nil, err
	}

	var certs []tls.Certificate
	if cert != nil {
		certs = []tls.Certificate{*cert}
	}

	// return config if no error occurred
	return &TLSConfig{
		Versions:         config.Versions,
		Verification:     config.VerificationMode,
		Certificates:     certs,
		RootCAs:          cas,
		CipherSuites:     cipherSuites,
		CurvePreferences: curves,
		Renegotiation:    tls.RenegotiationSupport(config.Renegotiation),
	}, nil
}

// Validate valies the TLSConfig struct making sure certificate sure we have both a certificate and
// a key.
func (c *Config) Validate() error {
	return c.Certificate.Validate()
}

// IsEnabled returns true if the `enable` field is set to true in the yaml.
func (c *Config) IsEnabled() bool {
	return c != nil && (c.Enabled == nil || *c.Enabled)
}
