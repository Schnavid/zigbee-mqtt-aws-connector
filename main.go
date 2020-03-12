package main

import (
	"crypto/tls"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"os"
	"strings"
	"time"
)

type AWSMessage struct {
	Identifier string `json:"page"`
	Timestamp  int64  `json:"timestamp"`
	Payload    string `json:"payload"`
}

func main() {

	var awsHost string

	if len(os.Args) > 3 {
		awsHost = os.Args[3]
	}

	mqttLocalClient := connect("192.168.1.32", 1883, "", "")
	mqttAWSClient := connect(awsHost, 8883, "cert.cert", "key.key")

	if token := mqttLocalClient.Subscribe("#", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("[LOCAL MQTT] Connected")

	if token := mqttAWSClient.Subscribe("#", 0, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("[AWS MQTT] Connected")

	select {}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	time := time.Now().UnixNano() / int64(time.Millisecond)

	if (strings.Contains(topic, "homeassistant/")) {
		//will be implemented later on.
	} else if (strings.Contains(topic, "zigbee2mqtt/")) {
		split := strings.Split(topic, "/")

		message := AWSMessage{
			Identifier: split[1],
			Timestamp:  time,
			Payload:    string(payload),
		}

		fmt.Println(message)

	}

}

func connect(host string, port int, certPath, keyPath string) mqtt.Client {

	var opts mqtt.ClientOptions
	var brokerURL string

	opts.ClientID = "hassio-aws-connector"
	opts.CleanSession = true
	opts.AutoReconnect = true
	opts.MaxReconnectInterval = 3 * time.Second

	if certPath == "" && keyPath == "" {
		opts.Username = os.Args[1]
		opts.Password = os.Args[2]
		brokerURL = fmt.Sprintf("tcp://%s:%d", host, port)
	} else {
		cer, _ := tls.LoadX509KeyPair("cert.crt", "key.key")
		tlsConf := tls.Config{Certificates: []tls.Certificate{cer}}
		opts.TLSConfig = &tlsConf
		brokerURL = fmt.Sprintf("tcps://%s:%d", host, port)

	}

	opts.AddBroker(brokerURL)

	mqttClient := mqtt.NewClient(&opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	return mqttClient
}
