package main

import (
	"github.com/ieee0824/elma"
	"flag"
	"os"
	"encoding/json"
	"log"
	"io"
	"github.com/ieee0824/sakuya"
	"github.com/ieee0824/getenv"
	"image/color"
	"github.com/joho/godotenv"
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
	godotenv.Load(".env")
	killer := make(chan bool)
	configPath = flag.String("f", "", "conf path")
	flag.Parse()
	slackClient := sakuya.New(getenv.String("SLACK_API_KEY"), getenv.String("CHANNEL"), "", getenv.String("SLACK_NAME"))
	slackClient.SetIconURL(getenv.String("ICON_URL"))


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
				result := <- r
				if result.IsHealthy() == nil {
					slackClient.SetColor(color.RGBA{0x00, 0xff, 0x00, 0xff})
					post(slackClient, result)
				} else {
					slackClient.SetColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
					post(slackClient, result)
				}
			}
		}()
	}

	<-killer
}
