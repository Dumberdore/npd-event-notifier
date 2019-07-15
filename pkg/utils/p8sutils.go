package utils

import (
  "fmt"
  "net/http"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "k8s.io/api/core/v1"
)

//Define the metrics we wish to expose// Prometheus counter to record the number or errors received from Torbit and Infoblox
var npdEventCounter = promauto.NewCounterVec(
  prometheus.CounterOpts{
    Name: "npd_event",
    Help: "Event raised by Node problem detector",
  },
  []string{
    "npdReason",
  },
)

func StartServer() {
  //This section will start the HTTP server and expose
  //any metrics on the /metrics endpoint.
  http.Handle("/metrics", promhttp.Handler())
  fmt.Println("Beginning to serve on port :8080")
  go http.ListenAndServe(":8080", nil)
}

func IncrementCounter(defaultEvent *v1.Event) {
  switch defaultEvent.Reason {
  case "DockerMonitorKilledDocker":
    fmt.Printf("[Interested Event] [%s] : %s\n", defaultEvent.Reason, defaultEvent.Message)
    npdEventCounter.WithLabelValues("DockerMonitorKilledDocker").Inc()

  case "KubeletMonitorKilledKubelet":
    fmt.Printf("[Interested Event] [%s] : %s\n", defaultEvent.Reason, defaultEvent.Message)
    npdEventCounter.WithLabelValues("KubeletMonitorKilledKubelet").Inc()

  default:
    fmt.Printf("[Mehh] [%s] :  (%s)\n", defaultEvent.Reason, defaultEvent.Message)
  }
}

func init() {
  //Set fooMetric to 1
  npdEventCounter.WithLabelValues("DockerMonitorKilledDocker")
  npdEventCounter.WithLabelValues("KubeletMonitorKilledKubelet")

  StartServer()
}
