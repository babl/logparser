package main

import "github.com/urfave/cli"

func configureCli() (app *cli.App) {
	app = cli.NewApp()
	app.Usage = "Babl QA Logs Parser"
	app.Version = Version
	app.Action = func(c *cli.Context) {
		kafkaBrokers := c.String("kafka-brokers")
		if len(kafkaBrokers) == 0 {
			kafkaBrokers = "sandbox.babl.sh:9092"
		}
		run(kafkaBrokers, c.GlobalBool("debug"), c.GlobalBool("output"), c.GlobalBool("readonly"))
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "kafka-brokers, kb",
			Usage: "Comma separated list of kafka brokers",
			Value: "sandbox.babl.sh:9092",
		},
		cli.BoolFlag{
			Name:   "readonly, ro",
			Usage:  "Will only read from kafka brokers, no write operations",
			EnvVar: "BABL_LOGPARSER_RO",
		},
		cli.BoolFlag{
			Name:   "output, o",
			Usage:  "Outputs JSON data",
			EnvVar: "BABL_LOGPARSER_OUTPUT",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debug mode & verbose logging",
			EnvVar: "BABL_DEBUG",
		},
	}
	return
}
