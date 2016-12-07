package bablrequestparser

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func ListenToLogsRAW(client *sarama.Client, topic string, chRAWData chan *RAWJsonData) {
	log.Debug("Consuming from topic: ", topic)
	ch := make(chan *kafka.ConsumerData)
	go kafka.Consume(client, topic, ch, kafka.ConsumerOptions{Offset: sarama.OffsetNewest})
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("RAW message received")

		rawdata := RAWJsonData{}
		err1 := rawdata.UnmarshalJSON(msg.Value)
		Check(err1)
		//rawdata.DebugJson()
		//rawdata.Debug()
		if checkRegex("\"rid\":.*", rawdata.Message) {
			go func() { chRAWData <- &rawdata }()
		}
		msg.Processed <- "success"
	}
	panic("listenToRAWMessages: Lost connection to Kafka")
}

func WriteToLogsQA(producer *sarama.SyncProducer, topic string, chRAWData chan *RAWJsonData,
	cluster string, dbgoutput bool, readonly bool) {

	for rawdata := range chRAWData {
		rawdata.Cluster = cluster
		qadata := ParseRawData(rawdata)
		if dbgoutput {
			qadata.DebugJson()
		}
		rhJson, _ := json.Marshal(qadata)
		//fmt.Printf("%s\n", rhJson)
		if !readonly {
			kafka.SendMessage(producer, qadata.RequestId, topic, &rhJson)
		}
	}
}
