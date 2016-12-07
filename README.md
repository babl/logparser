# babl-logparser

### Params
* -kb: kafka-brokers, it is mandatory to use the broker URL sandbox.babl.sh:9092, it will be used for the `cluster` property for the parsed JSON result
* -ro: ReadOnly mode, will not push parsed results into Kafka (debug locally)
* -o: Debug Output, sends the JSON parsed result into STDOUT

```
go build -v && ./babl-logparser -kb sandbox.babl.sh:9092 -ro -o | jq .
```
