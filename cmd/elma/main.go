package main

import (
	"github.com/ieee0824/elma"
	"flag"
	"os"
	"encoding/json"
	"log"
	"io"
)

var (
	configPath *string
)

func readConfig()([]elma.ClientSetting, error) {
	ret := []elma.ClientSetting{}
	f, err := os.Open(*configPath)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(f).Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func post(w io.Writer, r *elma.Result) error {
	if r.Edge {
		_, err := w.Write([]byte(r.String()))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	killer := make(chan bool)
	configPath = flag.String("f", "", "conf path")
	flag.Parse()

	configs, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	for _, config := range configs {
		go func() {
			client := config.Client()
			if client == nil {
				killer <- true
				return
			}
			r := client.Monitoring()
			for {
				post(os.Stdout, <-r)
			}
		}()
	}

	<-killer
}
