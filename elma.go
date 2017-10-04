package elma

import (
	"net/http"
	"log"
	"time"
)

type Result struct {
	Resp *http.Response
	Time time.Time
}

type Client struct {
	client *http.Client
	logger *log.Logger
	healthyStatusCodes []int
	body *http.Request
	url string
	interval time.Duration
}

func (c *Client)Monitoring() chan Result {
	ret := make(chan Result)

	go func() {
		ticker := time.NewTicker(c.interval)

		for {
			select {
			case <-ticker.C:
				resp, err := c.client.Do(c.body)
				if err != nil {
					panic(err)
				}
				ret <- Result{
					resp,
					time.Now(),
				}
			}
		}
	}()

	return ret
}