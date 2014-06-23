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
  -list-item-key="": use item key instead of item index
  -metric-timestamp=1403551787: metric time stamp (defaults to now)
  -prefix="": metric prefix
```

*Note:*
* `metric-timestamp` - allows the metric timestamp to be customised
* `prefix` - sets the base prefix for the metrics.
  * e.g. given `-prefix=base.prefix` and the metric `something.value = 46` the output will be `base.prefix.something.value = 46`
* `list-item-key` - allows array items to have properties as their key part rather than the index number.


## Examples
`-dry-run` is used to illustrate sent metrics as when sending details to graphite no output is given

Flatten and forward elasticsearch jvm metrics to graphite host
```
$ curl http://elasticsearch:9200/\_nodes/stats/jvm | stopgap -dry-run
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1200  100  1200    0     0   286k      0 --:--:-- --:--:-- --:--:--  390k
2014/06/23 20:46:30 Graphite: nodes.3jmurslbswqiksyf5pwkag.jvm.mem.heap_used_in_bytes 57544312 2014-06-23 20:46:30
2014/06/23 20:46:30 Graphite: nodes.3jmurslbswqiksyf5pwkag.timestamp 1403552790564 2014-06-23 20:46:30
...
$
```
or
```
$ stopgap -graphite-host 192.168.205.227 http://elasticsearch:9200/\_nodes/stats/jvm
2014/06/23 20:45:15 Graphite: nodes.3jmurslbswqiksyf5pwkag.jvm.mem.non_heap_used_in_bytes 27461944 2014-06-23 20:45:15
2014/06/23 20:45:15 Graphite: nodes.3jmurslbswqiksyf5pwkag.jvm.gc.collectors.young.collection_time_in_millis 60 2014-06-23 20:45:15
2014/06/23 20:45:15 Graphite: nodes.3jmurslbswqiksyf5pwkag.jvm.mem.pools.survivor.used_in_bytes 8912888 2014-06-23 20:45:15
...
$
```

or

```
$ stopgap -dry-run http://elasticsearch:9200/_nodes/stats/{jvm,http,os,process}
2014/06/23 20:43:26 Graphite: nodes.3jmurslbswqiksyf5pwkag.jvm.mem.pools.survivor.max_in_bytes 34865152 2014-06-23 20:43:26
...
2014/06/23 20:43:26 Graphite: nodes.3jmurslbswqiksyf5pwkag.http.current_open 2 2014-06-23 20:43:26
...
2014/06/23 20:43:26 Graphite: nodes.3jmurslbswqiksyf5pwkag.os.mem.free_in_bytes 149159936 2014-06-23 20:43:26
...
2014/06/23 20:43:26 Graphite: nodes.3jmurslbswqiksyf5pwkag.process.mem.total_virtual_in_bytes 3859607552 2014-06-23 20:43:26
...
$
```

Specify a custom prefix for the generated keys
```
$ stopgap -graphite-host 192.168.205.227 -prefix es http://elasticsearch:9200/_nodes/stats/jvm
2014/06/23 20:51:25 Graphite: es.nodes.3jmurslbswqiksyf5pwkag.jvm.mem.pools.survivor.peak_max_in_bytes 34865152 2014-06-23 20:51:25
2014/06/23 20:51:25 Graphite: es.nodes.3jmurslbswqiksyf5pwkag.jvm.mem.heap_committed_in_bytes 259522560 2014-06-23 20:51:25
...
$
```

Add world cup results to graphite (why not?!)
```
$ stopgap -dry-run -list-item-key country -prefix worldcup http://worldcup.sfg.io/teams/results
2014/06/23 20:59:26 Graphite: worldcup.korea_republic.draws 1 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.cameroon.draws 0 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.netherlands.losses 0 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.costa_rica.draws 0 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.portugal.goal_differential -4 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.ghana.draws 1 2014-06-23 20:59:24
2014/06/23 20:59:26 Graphite: worldcup.algeria.goals_for 5 2014-06-23 20:59:24
...
$
```

