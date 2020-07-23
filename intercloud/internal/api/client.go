package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	retry "github.com/sethvargo/go-retry"
)

type Client struct {
	endpoint   string
	token      string
	userAgents UserAgentProducts
	client     *http.Client
}

type ErrQuery struct {
	s       string
	Err     error
	bodyErr []byte
}

func (e *ErrQuery) Error() string {
	return fmt.Sprintf("%s Body: %s", e.s, e.bodyErr)
}

func (e *ErrQuery) Unwrap() error {
	return e.Err
}
func (e *ErrQuery) Body() []byte {
	return e.bodyErr
}

func NewErrQuery(s string, bodyErr []byte, err error) error {
	return &ErrQuery{
		s:       s,
		bodyErr: bodyErr,
		Err:     err,
	}
}

var (
	ErrNotFound        = errors.New("not found")
	ErrForbidden       = errors.New("forbidden")
	ErrOther           = errors.New("other")
	ErrTooManyRequests = errors.New("too many requests")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrPaymentRequired = errors.New("payment required")
)

func getHttpClient() *http.Client {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(5 * time.Second),
		}).Dial,
		TLSHandshakeTimeout: time.Duration(5 * time.Second),
	}

	return &http.Client{
		Timeout:   time.Duration(10 * time.Second),
		Transport: netTransport,
	}
}

func NewClient(
	endpoint,
	token string,
	userAgents UserAgentProducts,
	httpClient *http.Client,
) (client *Client) {

	if httpClient == nil {
		httpClient = getHttpClient()
	}

	client = &Client{
		endpoint:   endpoint,
		token:      token,
		userAgents: userAgents,
		client:     httpClient,
	}
	return
}

func (c *Client) BuildRequest(ctx context.Context, method, url string, body interface{}) (req *http.Request, err error) {
	var data []byte
	if body != nil {
		if data, err = json.Marshal(body); err != nil {
			return
		}
	}
	req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	return
}

func (c *Client) DoRequest(ctx context.Context, method, requestPath string, body interface{}, scope string, queryParams *url.Values) (resp *http.Response, err error) {

	var api, path *url.URL

	if ctx == nil {
		ctx = context.Background()
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second) // max backoff retries = 1 + 2 + 4 + 3 * 0.5 = 8.5s
	defer cancel()

	if path, err = url.Parse(requestPath); err != nil {
		return nil, err
	}
	if api, err = url.Parse(c.endpoint); err != nil {
		return nil, err
	}

	target := api.ResolveReference(path)

	if queryParams != nil {
		for key, values := range *queryParams {
			for _, v := range values {
				target.Query().Add(key, v)
			}
		}
	}

	req, err := c.BuildRequest(ctx, method, target.String(), body)

	if err != nil {
		return
	}

	if scope != "" {
		req.Header.Set("Intercloud-Scope", scope)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Intercloud Terraform Provider "+c.userAgents.String())

	// exponential backoff mechanism with jitter (+/- 500ms)
	b, _ := retry.NewExponential(1 * time.Second)
	b = retry.WithJitter(500*time.Millisecond, b)
	b = retry.WithMaxRetries(3, b)
	err = retry.Do(ctx, b, func(ctx context.Context) error {

		resp, err = c.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return nil
		}

		// Never Get Gb of data
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewReader(data))
		err = json.NewDecoder(bytes.NewReader(data)).Decode(&resp)
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			log.Printf("[DEBUG] rate limit has been hit, waiting before next retry")
			return retry.RetryableError(NewErrQuery("TooManyRequests", data, ErrTooManyRequests))
		case http.StatusUnauthorized:
			return NewErrQuery("Unauthorized", data, ErrUnauthorized)
		case http.StatusPaymentRequired:
			return NewErrQuery("Unauthorized", data, ErrUnauthorized)
		case http.StatusForbidden:
			return NewErrQuery("Not allowed", data, ErrForbidden)
		case http.StatusNotFound:
			return NewErrQuery("Not found", data, ErrNotFound)
		default:
			return NewErrQuery("Other", data, ErrOther)
		}
	})

	if err != nil {
		log.Printf("[DEBUG] request failure (err = %+v)", err)
		return nil, err
	}

	return
}
