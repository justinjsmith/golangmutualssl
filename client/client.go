package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	certFile = flag.String("cert", "myclient.crt", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "myclient.key", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "MyCA.crt", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()

	// Load client cert
	clientCert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// load the CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client, explicitly setting the servername
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   "myserver",
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	// TODO: figure out of this is required
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	// create the client
	client := &http.Client{Transport: transport}

	// make the request
	resp, err := client.Get("https://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))

}
