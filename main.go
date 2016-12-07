package main

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	. "github.com/larskluge/babl-logparser/bablrequestparser"
	"github.com/larskluge/babl-server/kafka"
)

type server struct {
	kafkaClient   *sarama.Client
	kafkaProducer *sarama.SyncProducer
}

const Version = "1.0.0"
const clientID = "babl-qa-logsparser"

var debug bool

func main() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)

	log.Warn("App START")

	app := configureCli()
	app.Run(os.Args)
}

func run(kafkaBrokers string, dbg bool, dbgoutput bool, readonly bool) {
	debug = dbg
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	s := server{}

	cluster := strings.Split(kafkaBrokers, ":")[0]
	const kafkaTopicRAW = "logs.raw"
	const kafkaTopicQA = "logs.qa"
	brokers := strings.Split(kafkaBrokers, ",")
	s.kafkaClient = kafka.NewClient(brokers, kafkaTopicRAW, debug)
	defer (*s.kafkaClient).Close()
	s.kafkaProducer = kafka.NewProducer(brokers, clientID+".producer")
	defer (*s.kafkaProducer).Close()

	chRAWData := make(chan *RAWJsonData)

	log.Warn("App Run ListenToLogsRAW")
	go ListenToLogsRAW(s.kafkaClient, kafkaTopicRAW, chRAWData)

	WriteToLogsQA(s.kafkaProducer, kafkaTopicQA, chRAWData, cluster, dbgoutput, readonly)
}
