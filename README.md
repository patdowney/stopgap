# Stopgap Metrics Forwarder

## What
Can't be arsed implementing a graphite interface for your favourite 
self hosted service? But you have some metrics exposed via JSON?

Use stopgap to parse the JSON, flatten it and send any integer or float values to graphite.


Converts
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
into metrics something like the following
```
metrics.some.thing = 1234
metrics.something_else = 6666
```

Non number types are ignored for now.

## Examples
Flatten and forward elasticsearch jvm metrics to graphite host
```
curl http://elasticsearch:9200/_nodes/stats/jvm | ./stopgap -graphite-host 192.168.205.227
```

Specify a custom prefix for the generated keys
```
curl http://elasticsearch:9200/_nodes/stats/process | ./stopgap -graphite-host 192.168.205.227 -prefix es
```
