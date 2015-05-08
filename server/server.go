package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var (
	certFile = flag.String("cert", "myserver.crt", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "myserver.key", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "MyCA.crt", "A PEM eoncoded CA's certificate file.")
)

func hello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello")
}

var config *tls.Config

func main() {
	flag.Parse()

	// load the CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// create a config object that requires verified client certs
	config = &tls.Config{
		ClientCAs:      caCertPool,
		ClientAuth:     tls.RequireAndVerifyClientCert,
		GetCertificate: clientCert,
	}
	// register the handler
	http.HandleFunc("/", hello)
	// create the server object
	srvr := http.Server{Addr: ":8080", TLSConfig: config}
	// set a callback for connection state changes
	srvr.ConnState = connState
	// TODO: figure out how to access the client cert
	err = srvr.ListenAndServeTLS(*certFile, *keyFile)
	panic(err)
}

func connState(conn net.Conn, state http.ConnState) {
	log.Printf("in conn state %s\n", state)
	tlsConn := tls.Server(conn, config)
	log.Printf("%+v", tlsConn)
}

// this seems like a way to retrieve the right cert based on SNI
func clientCert(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Print("in negotiation/SNI callback")
	return nil, nil
}
