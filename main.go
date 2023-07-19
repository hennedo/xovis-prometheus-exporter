package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var Entries = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "xovis",
	Name:      "entries",
	Help:      "Number of people that entered the venue",
}, []string{"sensor_serial_number", "sensor_name", "element_name"})
var Exits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "xovis",
	Name:      "exits",
	Help:      "Number of people that exited the venue",
}, []string{"sensor_serial_number", "sensor_name", "element_name"})
var Sum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "xovis",
	Name:      "sum",
	Help:      "Number of people that are in the venue",
}, []string{"sensor_serial_number", "sensor_name", "element_name"})

var lastReadings map[string]time.Time
var readingsLock sync.RWMutex

func main() {
	lastReadings = make(map[string]time.Time)
	readingsLock = sync.RWMutex{}
	prometheus.MustRegister(Entries, Exits, Sum)
	http.HandleFunc("/xovis", postData)
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func postData(w http.ResponseWriter, r *http.Request) {
	var data LineCountRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Error(err)
		return
	}
	for _, element := range data.Content.Element {
		sum := 0
		for _, measurement := range element.Measurement {
			readingsLock.RLock()
			if lastReading, ok := lastReadings[data.SensorInfo.SerialNumber+element.ElementName]; ok {
				if !measurement.To.After(lastReading) {
					logrus.Infof("already got measurement %s from %s(%s)", measurement.To, data.SensorInfo.Name, element.ElementName)
					continue
				}
			}
			readingsLock.RUnlock()
			for _, value := range measurement.Value {
				if value.Label == "fw" {
					Entries.WithLabelValues(data.SensorInfo.SerialNumber, data.SensorInfo.Name, element.ElementName).Add(float64(value.Value))
					sum += value.Value
				} else {
					Exits.WithLabelValues(data.SensorInfo.SerialNumber, data.SensorInfo.Name, element.ElementName).Add(float64(value.Value))
					sum -= value.Value
				}
			}
			readingsLock.Lock()
			lastReadings[data.SensorInfo.SerialNumber+element.ElementName] = measurement.To
			readingsLock.Unlock()
		}
		Sum.WithLabelValues(data.SensorInfo.SerialNumber, data.SensorInfo.Name, element.ElementName).Add(float64(sum))
	}
}
