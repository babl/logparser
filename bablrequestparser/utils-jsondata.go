package bablrequestparser

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

func isValidField(y interface{}, ytype reflect.Kind) bool {
	if y != nil && (reflect.TypeOf(y).Kind() == ytype) {
		return true
	}
	return false
}

func getFieldData(x interface{}, xtype reflect.Kind) interface{} {
	var result interface{} = nil
	if isValidField(x, xtype) {
		return x
	} else {
		switch xtype {
		case reflect.String:
			result = ""
		case reflect.Float64:
			result = float64(0.0)
		}
	}
	return result
}

/*
func getFieldDataStringArray(x map[string]interface{}, fieldname string) string {
	var result string
	list := x[fieldname].([]interface{})
	for _, v := range list {
		result += v.(string) + ","
	}
	fmt.Println("Result: ", result)
	return result
}*/

func getFieldDataString(x interface{}) string {
	return getFieldData(x, reflect.String).(string)
}

func getFieldDataInt(x interface{}) int {
	return int(getFieldData(x, reflect.Float64).(float64))
}

func getFieldDataInt32(x interface{}) int32 {
	return int32(getFieldData(x, reflect.Float64).(float64))
}

func getFieldDataFloat64(x interface{}) float64 {
	return getFieldData(x, reflect.Float64).(float64)
}

func checkRegex(pattern string, message string) bool {
	/*
		str := "{\"code\":\"req-downloading\",\"level\":\"info\",\"msg\":\"Downloading external payload\",\"payload_url\":\"http://sandbox.babl.sh:4442/ba223c69a40e0add\",\"rid\":\"cmmg71on7rugs\",\"time\":\"2016-12-06T11:20:18Z\"}"
		r := regexp.MustCompile(`\"rid\":.*`)
		fmt.Printf("%#v\n", r.FindStringSubmatch(str))
	*/
	result, _ := regexp.MatchString(pattern, message)
	return result
}

func ParseRawData(rawdata *RAWJsonData) QAJsonData {
	qadata := QAJsonData{}
	// Top Level => if not mapped use: getFieldDataString(rawdata.Z["_HOSTNAME"])
	qadata.Cluster = rawdata.Cluster
	qadata.Host = rawdata.HostName
	qadata.ImageName = rawdata.ContainerName
	// Message => Y
	qadata.RequestId = getFieldDataString(rawdata.Y["rid"])
	qadata.Key = getFieldDataString(rawdata.Y["key"])
	qadata.Message = getFieldDataString(rawdata.Y["msg"])
	qadata.Error = getFieldDataString(rawdata.Y["error"])
	qadata.Level = getFieldDataString(rawdata.Y["level"])
	qadata.Code = getFieldDataString(rawdata.Y["code"])
	qadata.Status = getFieldDataString(rawdata.Y["status"])
	qadata.Stderr = getFieldDataString(rawdata.Y["stderr"])
	qadata.Topic = getFieldDataString(rawdata.Y["topic"])
	qadata.Partition = getFieldDataInt32(rawdata.Y["partition"])
	qadata.Offset = getFieldDataInt32(rawdata.Y["offset"])
	qadata.Duration = getFieldDataFloat64(rawdata.Y["duration_ms"])

	// custom fields conversion
	qadata.AtVersion = "1"

	if isValidField(rawdata.Y["time"], reflect.String) {
		t1, _ := time.Parse(time.RFC3339, rawdata.Y["time"].(string))
		qadata.Timestamp = t1
		qadata.AtTimestamp = t1
	}

	if qadata.Error == "" && qadata.Stderr != "" {
		qadata.Error = qadata.Stderr
	}

	// supervisor2 OR babl-server specific data
	if checkRegex("supervisor.*", rawdata.ContainerName) {
		qadata.Service = "supervisor2"
		qadata.Supervisor = "supervisor2"

		qadata.Module = getFieldDataString(rawdata.Y["module"])
		if qadata.Module != "" {
			qadata.ModuleVersion = "v0"
		}
	} else {
		qadata.Service = "babl-server"

		module := strings.Split(rawdata.ContainerName, ".")
		qadata.Module = module[0]
		qadata.Module = strings.Replace(qadata.Module, "--", "/", -1)
		qadata.ModuleVersion = "v0"
	}

	// topics for groupconsumer : "New Group Message Received"
	if isValidField(rawdata.Y["topics"], reflect.Slice) {
		str := ""
		vals := rawdata.Y["topics"].([]interface{})
		for _, val := range vals {
			str += " " + val.(string)
		}
		qadata.Topic = strings.Replace(strings.Trim(str, " "), " ", ",", -1)
	}
	return qadata
}
