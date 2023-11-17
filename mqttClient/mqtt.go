package mqttClient

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttJson struct {
	Topic string `json:"topic"`
	Cmd   string `json:"cmd"`
	Data  string `json:"data"`
}

type PayloadCtrl struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type MqttMessage struct {
	Cmd       string      `json:"cmd"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Payload   PayloadCtrl `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	Id        string      `json:"id"`
}

type Client struct {
	client mqtt.Client
}

func ConnectMqtt(broker string, clientId string, username string, password string) *Client {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connection lost: %v", err)
	}
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Printf("Reconnecting:\n")
	}

	return &Client{
		client: mqtt.NewClient(opts),
	}
}

func (m *Client) Connect() {
	go func() {
		if token := m.client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}()
	fmt.Printf("MQTT Connected \n")
}

func (m *Client) Publish(topic string, payload []byte, qos byte) error {
	token := m.client.Publish(topic, qos, false, string(payload))
	token.Wait()
	return token.Error()
}
