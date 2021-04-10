package config

import (
	"flag"
	"strings"
)

type Config struct {
	Addr       string
	Hosts      []string
	CertFile   string
	KeyFile    string
	ContentDir string
}

func (cfg *Config) Flag(f *flag.FlagSet) {
	f.StringVar(&cfg.Addr, "addr", cfg.Addr, `adress to serve`)
	var domainsUsage = "expect requests for that hosts"
	if len(cfg.Hosts) > 0 {
		domainsUsage += ". Default values: " + strings.Join(cfg.Hosts, ", ")
	}
	f.Func("host", domainsUsage, func(domain string) error {
		cfg.Hosts = append(cfg.Hosts, domain)
		return nil
	})
	f.StringVar(&cfg.CertFile, "certfile", cfg.CertFile, `PEM encoded certificate file`)
	f.StringVar(&cfg.KeyFile, "keyfile", cfg.KeyFile, `PEM encoded private key file`)
	f.StringVar(&cfg.ContentDir, "content-dir", cfg.ContentDir, `content dir to serve`)
}
