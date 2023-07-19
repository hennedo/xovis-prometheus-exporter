package main

import "time"

type LineCountRequest struct {
	SensorTime struct {
		Timezone string    `json:"timezone"`
		Time     time.Time `json:"time"`
	} `json:"sensor-time"`
	Status struct {
		Code string `json:"code"`
	} `json:"status"`
	Content struct {
		Element []struct {
			ElementID   int       `json:"element-id"`
			ElementName string    `json:"element-name"`
			SensorType  string    `json:"sensor-type"`
			DataType    string    `json:"data-type"`
			From        time.Time `json:"from"`
			To          time.Time `json:"to"`
			Resolution  string    `json:"resolution"`
			Measurement []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value []struct {
					Value int    `json:"value"`
					Label string `json:"label"`
				} `json:"value"`
			} `json:"measurement"`
		} `json:"element"`
	} `json:"content"`
	SensorInfo struct {
		SerialNumber string `json:"serial-number"`
		IPAddress    string `json:"ip-address"`
		Name         string `json:"name"`
		Group        string `json:"group"`
		DeviceType   string `json:"device-type"`
	} `json:"sensor-info"`
}
