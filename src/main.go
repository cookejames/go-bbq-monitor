package main

import (
	"context"
	"strconv"
	"time"
	"fmt"
	"github.com/go-ble/ble"
	"github.com/mgutz/logxi/v1"
	"github.com/sworisbreathing/go-ibbq/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var logger = log.New("main")
var mqttClient mqtt.Client
var currentTemperatures []int16
var thingName = "inkbird-4xs"

func parseTemperature(temperature float64) int16 {
  if (int16(temperature) == 6552) {
    return -1
  } else {
    return int16(temperature)
  }
}
func temperatureReceived(temperatures []float64) {
	logger.Info("Temperature reading received")
	if (len(currentTemperatures) != len(temperatures)) {
		currentTemperatures = make([]int16, len(temperatures))
	}

	changed := false
	for i, temperature := range temperatures {
    parsedTemp := parseTemperature(temperature)
		if (currentTemperatures[i] != parsedTemp) {
			changed = true
			currentTemperatures[i] = parsedTemp
		}
	}
	
	if (changed) {
		update := fmt.Sprintf("{\"state\": {\"reported\" : {\"temperature1\" : %d, \"temperature2\" : %d, \"temperature3\" : %d, \"temperature4\" : %d } } }", currentTemperatures[0], currentTemperatures[1], currentTemperatures[2], currentTemperatures[3])
		logger.Info("Temperature changed sending update: ", update)
		if token := mqttClient.Publish(fmt.Sprintf("$aws/things/%s/shadow/name/temperature/update", thingName), 0, false, update); token.Wait() && token.Error() != nil {
			logger.Fatal("Failed to publish to MQTT: %v", token.Error())
		}
	} else {
		logger.Info("No temperature change") 
	}
}
func batteryLevelReceived(batteryLevel int) {
	logger.Info("Received battery data", "batteryPct", strconv.Itoa(batteryLevel))
}
func statusUpdated(status ibbq.Status) {
	logger.Info("iBBQ status updated", "status", status)
	update := fmt.Sprintf("{\"state\": {\"reported\" : {\"connection\" : \"%s\" } } }", status)
	if token := mqttClient.Publish(fmt.Sprintf("$aws/things/%s/shadow/name/connection/update", thingName), 0, false, update); token.Wait() && token.Error() != nil {
		logger.Fatal("Failed to publish to MQTT: %v", token.Error())
	}
}

func disconnectedHandler(cancel func(), done chan struct{}) func() {
	return func() {
		logger.Info("iBBQ disconnected, waiting before exiting")
		time.Sleep(2 * time.Second)
		logger.Info("Exiting...")
		mqttClient.Disconnect(250)
		cancel()
		close(done)
	}
}

func main() {
	var err error
	logger.Debug("initializing context")
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()
	registerInterruptHandler(cancel)
	ctx := ble.WithSigHandler(ctx1, cancel)
	logger.Debug("context initialized")

  mqttClient = createMqttClient()
	logger.Info("MQTT Client connected")

	var bbq ibbq.Ibbq
	logger.Debug("instantiating ibbq struct")
	done := make(chan struct{})
	var config ibbq.Configuration
	if config, err = ibbq.NewConfiguration(60*time.Second, 5*time.Minute); err != nil {
		logger.Fatal("Error creating configuration", "err", err)
	}
	if bbq, err = ibbq.NewIbbq(ctx, config, disconnectedHandler(cancel, done), temperatureReceived, batteryLevelReceived, statusUpdated); err != nil {
		logger.Fatal("Error creating iBBQ", "err", err)
	}
	logger.Debug("instantiated ibbq struct")
	logger.Info("Connecting to device")
	if err = bbq.Connect(); err != nil {
		logger.Fatal("Error connecting to device", "err", err)
	}
	logger.Info("Connected to device")
	<-ctx.Done()
	<-done
}
