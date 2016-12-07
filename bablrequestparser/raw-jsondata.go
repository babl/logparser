package bablrequestparser

import (
	"encoding/json"
	"fmt"
)

type RAWJsonData struct {
	Message       string                 `json:"MESSAGE"`
	HostName      string                 `json:"_HOSTNAME"`
	Cluster       string                 `json:"cluster"`
	ContainerName string                 `json:"CONTAINER_NAME"`
	Y             map[string]interface{} // MESSAGE properties
	Z             map[string]interface{} // Root properties
}

func (rawdata *RAWJsonData) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &rawdata.Z)
	rawdata.Message = getFieldDataString(rawdata.Z["MESSAGE"])
	rawdata.HostName = getFieldDataString(rawdata.Z["_HOSTNAME"])
	rawdata.ContainerName = getFieldDataString(rawdata.Z["CONTAINER_NAME"])

	json.Unmarshal([]byte(rawdata.Message), &rawdata.Y)
	return err
}

func (rawdata *RAWJsonData) Debug() {
	rhJson, _ := json.Marshal(rawdata)
	fmt.Printf("%s\n", rhJson)
}

func (rawdata *RAWJsonData) GetData() []byte {
	rhJson, _ := json.Marshal(rawdata)
	return []byte(rhJson)
}

func (rawdata *RAWJsonData) DebugJson() {
	var result map[string]interface{}
	rawtemp := rawdata
	rawtemp.Z = result
	rhJson, _ := json.Marshal(rawtemp)
	fmt.Printf("%s\n", rhJson)
}
