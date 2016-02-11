package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetPrefix("[kube2lb] ")
}

func main() {
	var apiserver, kubecfg, config, template, notify string
	flag.StringVar(&apiserver, "apiserver", "", "Kubernetes API server URL")
	flag.StringVar(&kubecfg, "kubecfg", "", "Path to kubernetes client configuration (Optional)")
	flag.StringVar(&config, "config", "", "Configuration path to generate")
	flag.StringVar(&template, "template", "", "Configuration source template")
	flag.StringVar(&notify, "notify", "", "Kubernetes API server URL")
	flag.Parse()

	if _, err := os.Stat(template); err != nil {
		log.Fatalf("Template not defined or doesn't exist")
	}

	if notify == "" {
		log.Fatalf("Notifier cannot be empty")
	}

	if f, err := os.OpenFile(config, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		log.Fatalf("Cannot open configuration file to write: %v", err)
	} else {
		f.Close()
	}

	notifier, err := NewNotifier(notify)
	if err != nil {
		log.Fatalf("Couldn't initialize notifier: %s", err)
	}

	client, err := NewKubernetesClient(kubecfg, apiserver)
	if err != nil {
		log.Fatalf("Couldn't connect with Kubernetes API server: %s", err)
	}

	client.AddTemplate(NewTemplate(template, config))
	client.AddNotifier(notifier)

	if err := client.Update(); err != nil {
		log.Fatalf("Couldn't update with initial state: %s", err)
	}

	if err := client.Watch(); err != nil {
		log.Fatalf("Couldn't watch Kubernetes API server: %s", err)
	}
}
