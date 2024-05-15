package powerdns

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
)

// Config describes de configuration interface of this provider
type Config struct {
	ServerURL       string
	APIKey          string
	InsecureHTTPS   bool
	CACertificate   string
	CacheEnable     bool
	CacheMemorySize string
	CacheTTL        int
	ZC              string
}

// Client returns a new client for accessing PowerDNS
func (c *Config) Client() (*Client, error) {

	tlsConfig := &tls.Config{}

	if c.CACertificate != "" {

		caCert, _, err := pathorcontents.Read(c.CACertificate)
		if err != nil {
			return nil, fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsConfig.RootCAs = caCertPool
	}

	tlsConfig.InsecureSkipVerify = c.InsecureHTTPS

	client, err := NewClient(c.ServerURL, c.APIKey, tlsConfig, c.CacheEnable, c.CacheMemorySize, c.CacheTTL, c.ZC)

	if err != nil {
		return nil, fmt.Errorf("Error setting up PowerDNS client: %s", err)
	}

	log.Printf("[INFO] PowerDNS Client configured for server %s", c.ServerURL)

	return client, nil
}
