package main

import (
	"context"
	"strconv"
	"time"
	"github.com/go-ble/ble"
	"github.com/mgutz/logxi/v1"
	"github.com/sworisbreathing/go-ibbq/v2"

	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var logger = log.New("main")
var mqttClient mqtt.Client

func temperatureReceived(temperatures []float64) {
	// logger.Info("Received temperature data", "temperatures", temperatures)
	// Send shadow update.
	update := fmt.Sprintf("{\"state\": {\"reported\" : {\"temperature1\" : \"%d\", \"temperature2\" : \"%d\", \"temperature3\" : \"%d\", \"temperature4\" : \"%d\" } } }", int16(temperatures[0]), int16(temperatures[1]), int16(temperatures[2]), int16(temperatures[3]) )
	logger.Info("Sending update: ", update)
	if token := mqttClient.Publish("$aws/things/inkbird-4xs/shadow/update", 0, false, update); token.Wait() && token.Error() != nil {
		logger.Fatal("failed to send update: %v", token.Error())
	}
}
func batteryLevelReceived(batteryLevel int) {
	logger.Info("Received battery data", "batteryPct", strconv.Itoa(batteryLevel))
}
func statusUpdated(status ibbq.Status) {
	logger.Info("Status updated", "status", status)
}

func disconnectedHandler(cancel func(), done chan struct{}) func() {
	return func() {
		logger.Info("Disconnected")
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
	mqttClient = createMqttClient()
	logger.Info("MQTT Client connected")
	<-ctx.Done()
	<-done
	logger.Info("Exiting")
}
