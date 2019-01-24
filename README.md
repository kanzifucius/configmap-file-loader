# configmap file loader

A kubernetes controller to load config map key values as files into a volume.
Config Maps are sourced by label and/or namespace

## Usage

```
Usage of ./configmapfileloader_unix:
  -development
    	development flag will allow to run the operator outside a kubernetes cluster
  -key-name-filter string
    	regex the key for the config map data key must contain
  -kubeconfig string
    	kubernetes configuration path, only used when development mode enabled (default "/home/kanzi/.kube/config")
  -label-filter string
    	kubernets label filer eg org.kanzi=test. Only configmaps with this label will be detected
  -namespace string
    	kubernetes namespace where this app is running
  -resync-seconds int
    	The number of seconds the controller will resync the resources (default 60)
  -volume-dir string
    	volume directory where the config map data will be written
  -webhook-method string
    	the HTTP method url to use to send the webhook (default "POST")
  -webhook-status-code int
    	the HTTP status code indicating successful triggering of reload (default 200)
  -webhook-url string
  	    the HTTP url to use to send the webhook notification
```

## Installation

```
go get -u github.com/golang/dep/cmd/dep
make
```


