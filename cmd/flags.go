package main

import (
"flag"
"os"
"path/filepath"
"time"

"github.com/kanzifucius/configmap_file_loader/pkg/controller"
"k8s.io/client-go/util/homedir"
)

// Flags are the controller flags.
type Flags struct {
	flagSet *flag.FlagSet

	Namespace     string
	ResyncSec     int
	KubeConfig    string
	Development   bool
	Directory     string
	KeyNameFilter string
	LabelFilter   string
	Webhook string
	WebhookMethod string
	WebhookStatusCode int
}

// ControllerConfig converts the command line flag arguments to controller configuration.
func (f *Flags) ControllerConfig() controller.Config {
	return controller.Config{
		Namespace:    f.Namespace,
		ResyncPeriod: time.Duration(f.ResyncSec) * time.Second,
		Directory:f.Directory,
		Name: "ControllerConfig",
		KeyNameFilter: f.KeyNameFilter,
		LabelFilter : f.LabelFilter,
		Webhook: f.Webhook,
		WebhookMethod: f.WebhookMethod,
		WebhookStatusCode: f.WebhookStatusCode,


	}
}

// NewFlags returns a new Flags.
func NewFlags() *Flags {
	f := &Flags{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
	// Get the user kubernetes configuration in it's home directory.
	kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Init flags.
	f.flagSet.StringVar(&f.Namespace, "namespace", "", "kubernetes namespace where this app is running")
	f.flagSet.IntVar(&f.ResyncSec, "resync-seconds", 60, "The number of seconds the controller will resync the resources")
	f.flagSet.StringVar(&f.KubeConfig, "kubeconfig", kubehome, "kubernetes configuration path, only used when development mode enabled")
	f.flagSet.BoolVar(&f.Development, "development", false, "development flag will allow to run the operator outside a kubernetes cluster")
	f.flagSet.StringVar(&f.Directory, "volume-dir", "", "volume directory where the config map data will be written too")
	f.flagSet.StringVar(&f.KeyNameFilter, "key-name-filter", "", "regex the key for the config map must contain")
	f.flagSet.StringVar(&f.LabelFilter, "label-filter", "", "lable for the config map to filter")
	f.flagSet.StringVar(&f.WebhookMethod,"webhook-method", "POST", "the HTTP method url to use to send the webhook")
	f.flagSet.IntVar(&f.WebhookStatusCode,"webhook-status-code", 200, "the HTTP status code indicating successful triggering of reload")
	f.flagSet.StringVar(&f.Webhook,"webhook-url", "", "the HTTP url to use to send the webhook")


	f.flagSet.Parse(os.Args[1:])

	return f
}