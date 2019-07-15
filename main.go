package main

import (
  "fmt"
  "log"
  "os"
  "path/filepath"
  "time"

  "github.com/npd-event-notifier/pkg/utils"
  "k8s.io/api/core/v1"
  "k8s.io/apimachinery/pkg/fields"
  "k8s.io/client-go/kubernetes"
  _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/cache"
  "k8s.io/client-go/tools/clientcmd"
)

func main() {

  // Create Incluster/OutOfCluster config

  config, err := rest.InClusterConfig()
  if err != nil {
    log.Println("Is this running OutSide the cluster ?")
    // Bootstrap k8s configuration from local   Kubernetes config file
    kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
    log.Println("Using kubeconfig file: ", kubeconfig)
    config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
      log.Fatal(err)
    }
  }

  // creates the clientset
  kubeClient, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }

  watchlist := cache.NewListWatchFromClient(
    kubeClient.CoreV1().RESTClient(),
    "Events",
    "default",
    fields.Everything(),
  )
  _, controller := cache.NewInformer( // also take a look at NewSharedIndexInformer
    watchlist,
    &v1.Event{},
    0, //Duration is int64
    cache.ResourceEventHandlerFuncs{
      AddFunc: func(obj interface{}) {
        utils.IncrementCounter(obj.(*v1.Event))
      },
      DeleteFunc: func(obj interface{}) {

      },
      UpdateFunc: func(oldObj, newObj interface{}) {
        utils.IncrementCounter(newObj.(*v1.Event))
        fmt.Printf("Event Updated: %T %s\n", newObj, newObj)
      },
    },
  )

  stop := make(chan struct{})
  defer close(stop)
  go controller.Run(stop)
  for {
    time.Sleep(time.Second)
  }
}
