package httpClient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mqtt-http-interface/mqttClient"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func ListenAndServe(serverUrl string, m *mqttClient.Client) {

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		publishMqttMessage(w, r, m)
	})

	go func() {
		if err := http.ListenAndServe(serverUrl, nil); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Printf("HTTP Server running: %v\n", serverUrl)
}

func publishMqttMessage(w http.ResponseWriter, r *http.Request, m *mqttClient.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Declare a variable of type Data
	var requestBody mqttClient.MqttJson

	// Unmarshal the JSON data into the Data struct
	if err = json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Request body format incorrect", http.StatusBadRequest)
	}

	mqttPacket := mqttClient.MqttMessage{
		Id:   uuid.New().String(),
		Cmd:  "ctrl",
		From: "ieco-services",
		To:   "ieco-support",
		Payload: mqttClient.PayloadCtrl{
			Cmd:  requestBody.Cmd,
			Data: requestBody.Data,
		},
		Timestamp: time.Now(),
	}
	fmt.Printf("mqttPacket: %v\n", mqttPacket)
	packetConv, err := json.Marshal(mqttPacket)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := m.Publish(requestBody.Topic, packetConv, 2); err != nil {
		http.Error(w, fmt.Sprintf("Mqtt Publish Failed: %s", err.Error()), http.StatusConflict)
	}

	w.WriteHeader(http.StatusOK)
}
