package main


import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	// "time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var rootCert = "certs/AmazonRootCA1.pem"
var thingCertificate = "certs/certificate.pem.crt"
var privateKey = "certs/private.pem.key"
var brokerUri = "tls://a1r930ukf7ddm6-ats.iot.eu-west-1.amazonaws.com:8883"

// Adapted from https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/ssl/main.go
// Also see https://www.eclipse.org/paho/clients/golang/
func NewTLSConfig() (config *tls.Config, err error) {
	// Import trusted certificates from CAfile.pem.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(rootCert)
	if err != nil {
		return
	}
	certpool.AppendCertsFromPEM(pemCerts)

	// Import client certificate/key pair.
	cert, err := tls.LoadX509KeyPair(thingCertificate, privateKey)
	if err != nil {
		return
	}

	// Create tls.Config with desired tls properties
	config = &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
	return
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func createMqttClient() mqtt.Client {
	tlsconfig, err := NewTLSConfig()
	if err != nil {
		log.Fatalf("failed to create TLS configuration: %v", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerUri)
	opts.SetClientID("clientID").SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	// Start the connection.
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}

	return c
}

/*
func upd() {
		// Send shadow update.
		update := `{
			"state": {
					"desired" : {
							"color" : { "r" : 10 },
							"engine" : "ON"
					}
			}
	}`
		fmt.Println("Sending update.")
		if token := c.Publish("$aws/things/iotexample/shadow/update", 0, false, update); token.Wait() && token.Error() != nil {
			log.Fatalf("failed to send update: %v", token.Error())
		}
	
		fmt.Println("Listening for new events.")
		if token := c.Subscribe("$aws/things/iotexample/shadow/update/accepted", 0, nil); token.Wait() && token.Error() != nil {
			log.Fatalf("failed to create subscription: %v", token.Error())
		}
	
		fmt.Println("Sleeping.")
		time.Sleep(time.Second * 60)
		fmt.Println("Disconnecting.")
		c.Disconnect(250)
}
*/