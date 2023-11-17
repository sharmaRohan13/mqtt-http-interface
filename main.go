package main

import (
	"flag"
	"fmt"

	"mqtt-http-interface/config"
	"mqtt-http-interface/httpClient"
	"mqtt-http-interface/mqttClient"
)

func main() {
	fmt.Print("Started..\n")
	configFile := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	cfg := config.ParseConfig(*configFile)

	m := mqttClient.ConnectMqtt(
		cfg.Mqtt.Broker,
		cfg.Mqtt.ClientId,
		cfg.Mqtt.Username,
		cfg.Mqtt.Password,
	)
	m.Connect()

	httpClient.ListenAndServe(cfg.Http.Server, m)

	select {}
}
