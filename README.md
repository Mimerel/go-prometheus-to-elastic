# go-prometheus-to-elastic
Go application that sends collected data from prometheus to elasticsearch

This go software is designed read prometheus metrics and send them to elasticsearch.

* Prerequisits

Having a prometheus running
have set up prometheus to accept metric requests
have Go setup on the server
My prometheus requests metrics for exporters set in the target raspberry set in the config file.

* Prometheus settings
Add to you prometheus scrape_configs setting the following configuration that you can adapt to you usage
In my case, all jobs in prometheus are named "prometheus" 
```
scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ['localhost:2112', '192.168.0.100:2112', '192.168.0.101:2112']
  - job_name: 'federate'
    scrape_interval: 15s
    honor_labels: true
    metrics_path: '/federate'
    params:
      'match[]':
        - '{job="prometheus"}'
        - '{__name__=~"job:.*"}'
    static_configs:
      - targets:
        - 'localhost:9090'
```

once this is done, you should be able to request the prometheus url and get see your metrics appear
```
 http://<prometheus_ip>:9090/federate?match[]={job=~"prom.*"}
```


* Using the app.

The app will be running on the port set in the configuration.yaml file.
you need to set an environment variable to specify the full path to you configuration file.

the environment variable is : LOGGER_CONFIGURATION_FILE

If no env variable is set, the application will search for it in the current file from 
which is started the application or the following path : 
```
/home/pi/go/src/go-prometheus-to-elastic/configuration.yaml
```

* RUN : to run the application either : 
```
go run main.go
```
or
```
go build  // to build the application
```
then 
```
./go-prometheus-to-elastic // to run the build
```

You will probably be missing dependencies
```
	github.com/Mimerel/go-logger-client
	gopkg.in/yaml.v2
	...
```
to add a dependency run 

```
go get <name of dependency>
```

for example : 
```
go get github.com/Mimerel/go-logger-client
```

* using the app

Once the app is running, it will automatically, every "scrape_interval" seconds
request the prometheus for the latest metrics and send them to prometheus

* Configuration file

```
elasticSearch:
  url: http://<elasticsearch_ip>:9200
prometheus:
  url: http://<prometheus_ip>:9090/federate?match[]={job=~%22prom.*%22}
host: go-log-to-elastic // name of logs in prometheus
scrape_interval: 20 
```