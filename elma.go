package elma

import (
	"net/http"
	"time"
	"errors"
	"bytes"
	"fmt"
)

var (
	NoTargetErr = errors.New("not target")
)

type Result struct {
	Resp *http.Response `json:"resp"`
	Time time.Time `json:"time"`
	Error error `json:"error"`
	healthyStatusCode []int
	Edge bool `json:"edge"`
	target string
}

func (r Result) String() string {
	err := r.IsHealthy()
	if  err == nil {
		return fmt.Sprintf("%v: healthy, status: %v, target: %v", r.Time.Format("2006/01/02 15:04:05 MST"), r.Resp.StatusCode, r.target)
	}
	return fmt.Sprintf("%v: unhealthy: status:%v, error: %v, target: %v", r.Time.Format("2006/01/02 15:04:05 MST"), r.Resp.StatusCode, err.Error(), r.target)
}

func (r *Result) IsHealthy() error {
	if r.Error != nil {
		return r.Error
	}
	if r.Resp == nil {
		return errors.New("response is nil")
	}

	for _, c := range r.healthyStatusCode {
		if c == r.Resp.StatusCode {
			return nil
		}
	}
	return fmt.Errorf("status code is not match, want: %v, but: %v", r.healthyStatusCode, r.Resp.StatusCode)
}

type ClientSetting struct {
	Target string `json:"target"`
	Rate string `json:"rate"`
	UA string `json:"ua"`
	HealthyStatusCodeList []int `json:"healthy_status_code_list"`
	RequestMethod string `json:"request_method"`
	RequestBody []byte `json:"request_body"`
}

func (c *ClientSetting) Client() *Client {
	var req *http.Request
	if c.RequestBody != nil {
		var err error
		req, err = http.NewRequest(c.RequestMethod, c.Target, bytes.NewReader(c.RequestBody))
		if err != nil {
			return nil
		}
	} else {
		var err error
		req, err = http.NewRequest(c.RequestMethod, c.Target, nil)
		if err != nil {
			return nil
		}
	}

	rate, err := time.ParseDuration(c.Rate)
	if err != nil {
		return nil
	}
	healthyStatusCodeList := c.HealthyStatusCodeList
	if len(healthyStatusCodeList) == 0{
		healthyStatusCodeList = []int{200}
	}
	return &Client {
		c.Target,
		rate,
		c.UA,
		new(http.Client),
		healthyStatusCodeList,
		req,
		nil,
	}
}

type Client struct {
	target string
	rate time.Duration
	userAgent string
	client *http.Client
	healthyStatusCode []int
	request *http.Request
	before *bool
}

func New(t string)(*Client, error){
	if t == "" {
		return nil, NoTargetErr
	}
	c := new(Client)

	c.target = t
	c.rate = 10 * time.Second
	c.userAgent = "Elma http/1.1"
	c.client = &http.Client{}
	c.healthyStatusCode = []int{http.StatusOK}
	req, err := http.NewRequest("GET", t, nil)
	if err != nil {
		return nil, err
	}
	c.request = req


	return c, nil
}

func (c *Client) SetRate(d time.Duration) {
	c.rate = d
}

func (c *Client) SetUA(ua string) {
	c.userAgent = ua
}

func (c *Client) SetStatusCode(s []int) {
	c.healthyStatusCode = s
}

func (c *Client) SetRequest(req *http.Request) {
	c.request = req
}


func (c *Client) Monitoring() (chan *Result) {
	ret := make(chan *Result)
	t := time.NewTicker(c.rate)
	if c.userAgent == "" {
		c.request.Header.Set("User-Agent", "Elma http/1.1")
	} else {
		c.request.Header.Set("User-Agent", c.userAgent)
	}

	go func () {
		for {
			select {
			case <-t.C:
				resp, err := c.client.Do(c.request)
				result := Result{
					resp,
					time.Now(),
					err,
					c.healthyStatusCode,
					false,
					c.target,
				}
				if c.before != nil {
					result.Edge = *c.before != (result.IsHealthy() == nil)
				} else {
					result.Edge = true
				}
				b := result.IsHealthy() == nil
				c.before = &b
				ret <- &result
			}
		}
	}()
	return ret
}
