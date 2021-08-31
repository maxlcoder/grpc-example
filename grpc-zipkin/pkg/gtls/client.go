package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

type Client struct {
	CaFile string
	CertFile string
	KeyFile string
	ServerName string
}

func (t *Client) GetCredentialsByCA() (credentials.TransportCredentials, error)  {

	cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(t.CaFile)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("certPool.AppendCertsFromPEM fail")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName: t.ServerName,
		RootCAs: certPool,
	})
	return c, err

}