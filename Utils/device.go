package Utils

import (
	"encoding/json"
	"fmt"
)

//DeviceEvent  - Represent all the event with device information.
type DeviceEvent struct {
	Timestamp int64    `json:"timestamp"`
	EventType string   `json:"eventType"`
	Devices   []Device `json:"devices"`
}

//Device - Represent the device information.
type Device struct {
	Compliant        bool   `json:"compliant"`
	Status           string `json:"status"`
	LastCheckInTime  int64  `json:"lastCheckInTime"`
	RegistrationTime int64  `json:"registrationTime"`
	Identifier       string `json:"identifier"`
	MacAddress       string `json:"macAddress"`
	Manufacturer     string `json:"manufacturer"`
	Model            string `json:"model"`
	Os               string `json:"os"`
	OsVersion        string `json:"osVersion"`
	SerialNumber     string `json:"serialNumber"`
	UserID           string `json:"userId"`
	UserUUID         string `json:"userUuid"`
}

//ReadDevice - reads the device json from file specified
func ReadDevice(device *Device) {
	jsonString := ReadFile("./device.json")
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, &device)
	if err != nil {
		fmt.Println("Error")
	}
}

//SetTimeStamp - set the  current timestamp for deviceEvent information.
func (deviceEvent *DeviceEvent) SetTimeStamp(timestamp int64) {
	deviceEvent.Timestamp = timestamp
}

//AddDevice - add device.
func (deviceEvent *DeviceEvent) AddDevice(device Device) {
	deviceEvent.Devices = append(deviceEvent.Devices, device)
}

//TimeStamp - returns the  current timestamp for deviceEvent information.
func (deviceEvent *DeviceEvent) TimeStamp() int64 {
	return deviceEvent.Timestamp
}
