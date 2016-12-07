# babl-logparser

### Params
* -ro: ReadOnly mode, will not push parsed results into Kafka (debug locally)
* -o: Debug Output, sends the JSON parsed result into STDOUT

```
go build -v && ./babl-logparser -kb sandbox.babl.sh:9092 -ro -o | jq .
```
