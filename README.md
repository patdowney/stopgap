# Stopgap Metrics Forwarder [![Build Status](https://travis-ci.org/patdowney/stopgap.svg?branch=master)](https://travis-ci.org/patdowney/stopgap)

## What
Can't be arsed implementing a graphite interface for your favourite 
self hosted service? But you have some metrics exposed via JSON?

Use `stopgap` to parse the JSON, flatten it and send any integer or float values to graphite.


Turns
```json
{
  "metrics": {
    "some": {
      "thing": 1234
    },
    "something_else" :  6666
  },
  "name": "some text"
}
```
into metrics that look like the following
```
metrics.some.thing = 1234
metrics.something_else = 6666
```

Non-number types and numbers as string values are ignored for now.

## Usage
```
Usage of ./stopgap:
  -dry-run=false: dry run
  -graphite-host="localhost": graphite host
  -graphite-port=2003: graphite port
  -metric-timestamp=2014-06-17T21:52:33+01:00: metric time stamp (defaults to now)
  -prefix="": metric prefix
```

*Note:*
* `metric-timestamp` - allows the metric timestamp to be customised
* `prefix` - sets the base prefix for the metrics.
  * e.g. given `-prefix=base.prefix` and the metric `something.value = 46` the output will be `base.prefix.something.value = 46`

## Examples
Flatten and forward elasticsearch jvm metrics to graphite host
```
curl http://elasticsearch:9200/_nodes/stats/jvm | stopgap -graphite-host 192.168.205.227
```
or
```
stopgap -graphite-host 192.168.205.227 http://elasticsearch:9200/_nodes/stats/jvm 
```
or
```
stopgap -graphite-host 192.168.205.227 http://elasticsearch:9200/_nodes/stats/{jvm,http,os,process}
```

Add world cup results to graphite (why not?!)
```
stopgap -list-item-key country -prefix worldcup -dry-run http://worldcup.sfg.io/teams/results
```

Specify a custom prefix for the generated keys
```
curl http://elasticsearch:9200/_nodes/stats/process | stopgap -graphite-host 192.168.205.227 -prefix es
```
