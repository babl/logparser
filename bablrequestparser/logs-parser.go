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
	consumer, err0 := sarama.NewConsumerFromClient(*client)
	Check(err0)
	defer consumer.Close()

	partition := int32(0)
	offsetNewest, err1 := (*client).GetOffset(topic, partition, sarama.OffsetNewest)
	Check(err1)
	offsetOldest, err2 := (*client).GetOffset(topic, partition, sarama.OffsetOldest)
	Check(err2)

	offset := offsetNewest

	log.Warn("Consuming from topic: ", topic, " partition: ", partition, " offset: ", offset, " OffsetNewest:", offsetNewest, " OffsetOldest:", offsetOldest)
	pc, err := consumer.ConsumePartition(topic, partition, offset)
	Check(err)
	defer pc.Close()

	for msg := range pc.Messages() {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("RAW message received")

		rawdata := RAWJsonData{}
		err1 := rawdata.UnmarshalJSON(msg.Value)
		Check(err1)
		//rawdata.DebugJson()
		//rawdata.Debug()
		if checkRegex("\"rid\":.*", rawdata.Message) {
			go func() { chRAWData <- &rawdata }()
		}
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
