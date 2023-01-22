package config

import (
	"crypto/tls"
	"crypto/x509"
)

type TLSConfig struct {
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
	Server        bool
}

func SetupTLSConfig(cfg TLSConfig) (*tls.Config, error) {
	// var err error
	tlsConfig := &tls.Config{}
	if cfg.CertFile != "" && cfg.KeyFile != "" && cfg.CAFile != "" {
		tlsConfig.ClientCAs = x509.NewCertPool()
		tlsConfig.ClientCAs.AppendCertsFromPEM([]byte(cfg.CAFile))

		cert, err := tls.X509KeyPair([]byte(cfg.CertFile), []byte(cfg.KeyFile))
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	// if cfg.CAFile != "" {
	if cfg.CertFile != "" && cfg.KeyFile != "" && cfg.CAFile != "" {
		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM([]byte(cfg.CAFile))

		cert, err := tls.X509KeyPair([]byte(cfg.CertFile), []byte(cfg.KeyFile))
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	return tlsConfig, nil
}

// func SetupTLSConfig(cfg TLSConfig) (*tls.Config, error) {
// 	var err error
// 	tlsConfig := &tls.Config{}
// 	if cfg.CertFile != "" && cfg.KeyFile != "" {
// 		tlsConfig.Certificates = make([]tls.Certificate, 1)
// 		tlsConfig.Certificates[0], err = tls.LoadX509KeyPair(
// 			cfg.CertFile,
// 			cfg.KeyFile,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if cfg.CAFile != "" {
// 		b, err := ioutil.ReadFile(cfg.CAFile)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ca := x509.NewCertPool()
// 		ok := ca.AppendCertsFromPEM([]byte(b))
// 		if !ok {
// 			return nil, fmt.Errorf(
// 				"failed to parse root certificate: %q",
// 				cfg.CAFile,
// 			)
// 		}
// 		if cfg.Server {
// 			tlsConfig.ClientCAs = ca
// 			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
// 		} else {
// 			tlsConfig.RootCAs = ca
// 		}
// 		tlsConfig.ServerName = cfg.ServerAddress
// 	}
// 	return tlsConfig, nil
// }
