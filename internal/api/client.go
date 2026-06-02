package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type Client struct {
	cfg        config.Config
	httpClient *http.Client
}

var (
	testBaseURLMu       sync.RWMutex
	testBaseURLOverride string
)

func NewClient(cfg config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) Get(endpoint string, out interface{}) error {
	return c.Do(http.MethodGet, endpoint, nil, out)
}

func (c *Client) Post(endpoint string, body interface{}, out interface{}) error {
	return c.Do(http.MethodPost, endpoint, body, out)
}

func (c *Client) Put(endpoint string, body interface{}, out interface{}) error {
	return c.Do(http.MethodPut, endpoint, body, out)
}

func (c *Client) Do(method, endpoint string, body interface{}, out interface{}) error {
	if err := c.cfg.RequireAPIKey(); err != nil {
		return err
	}
	target := endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		target = baseURL(c.cfg) + "/" + strings.TrimPrefix(endpoint, "/")
	}

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, target, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return yxerrors.Remote("HTTP request failed", err.Error())
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return yxerrors.Remote(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(raw)), string(raw))
	}
	if out == nil {
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func baseURL(cfg config.Config) string {
	testBaseURLMu.RLock()
	override := testBaseURLOverride
	testBaseURLMu.RUnlock()
	if override != "" {
		return override
	}
	return cfg.APIURL
}

func SetBaseURLForTest(rawURL string) func() {
	testBaseURLMu.Lock()
	previous := testBaseURLOverride
	testBaseURLOverride = strings.TrimRight(rawURL, "/")
	testBaseURLMu.Unlock()
	return func() {
		testBaseURLMu.Lock()
		testBaseURLOverride = previous
		testBaseURLMu.Unlock()
	}
}

func Query(endpoint string, params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Set(key, value)
		}
	}
	if values.Encode() == "" {
		return endpoint
	}
	return endpoint + "?" + values.Encode()
}

func DataOrSelf(value map[string]interface{}) interface{} {
	if data, ok := value["data"]; ok {
		return data
	}
	return value
}
