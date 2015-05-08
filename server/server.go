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
	srvr.TLSConfig.BuildNameToCertificate()
	// set a callback for connection state changes
	srvr.ConnState = connState
	// TODO: figure out how to access the client cert
	err = srvr.ListenAndServeTLS(*certFile, *keyFile)
	panic(err)
}

func connState(conn net.Conn, state http.ConnState) {
	tlscon, ok := conn.(*tls.Conn)
	if !ok {
		panic("not a tls connection")
	}
	if state == http.StateActive {
		connState := tlscon.ConnectionState()
		sub := connState.PeerCertificates[0].Subject.CommonName
		log.Printf("Verified client cert common name: %s\n", sub)
	}
}

// this seems like a way to retrieve the right cert based on SNI
func clientCert(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Print("in negotiation/SNI callback")
	return nil, nil
}
