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

func (c *Client) Patch(endpoint string, body interface{}, out interface{}) error {
	return c.Do(http.MethodPatch, endpoint, body, out)
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
		return remoteErrorFromBody(resp.StatusCode, raw)
	}
	if err := assertBusinessOK(raw); err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

// remoteErrorFromBody builds a remote error from a non-2xx response, preferring
// the gateway's structured {statusCode, message, code} envelope over the raw
// body so callers (and the retry heuristics) see a clean message.
func remoteErrorFromBody(status int, raw []byte) error {
	message, code := parseResponseEnvelope(raw)
	if message == "" {
		message = strings.TrimSpace(string(raw))
	}
	if message == "" {
		message = fmt.Sprintf("HTTP %d", status)
	}
	details := map[string]interface{}{"httpStatus": status}
	if code != "" {
		details["code"] = code
	}
	if body := strings.TrimSpace(string(raw)); body != "" {
		details["body"] = body
	}
	return yxerrors.Remote(message, details)
}

// assertBusinessOK rejects 2xx responses whose envelope carries a non-zero
// business statusCode. The gateway wraps success as {statusCode:0,data:...} and
// reports failures via non-2xx, but this guards against endpoints that surface
// a business error inside an HTTP 200 response. Responses without a numeric
// statusCode (bare arrays, {data:...} test fixtures) pass through untouched.
func assertBusinessOK(raw []byte) error {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return nil
	}
	var env struct {
		StatusCode *float64    `json:"statusCode"`
		Message    interface{} `json:"message"`
		Code       interface{} `json:"code"`
	}
	if err := json.Unmarshal(trimmed, &env); err != nil {
		return nil
	}
	if env.StatusCode == nil || *env.StatusCode == 0 {
		return nil
	}
	message := stringifyEnvelopeValue(env.Message)
	if message == "" {
		message = fmt.Sprintf("business error statusCode=%d", int(*env.StatusCode))
	}
	details := map[string]interface{}{"statusCode": int(*env.StatusCode)}
	if code := stringifyEnvelopeValue(env.Code); code != "" {
		details["code"] = code
	}
	return yxerrors.Remote(message, details)
}

// parseResponseEnvelope extracts message and code from a JSON object body,
// returning empty strings when the body is not a JSON object.
func parseResponseEnvelope(raw []byte) (message, code string) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return "", ""
	}
	var env struct {
		Message interface{} `json:"message"`
		Code    interface{} `json:"code"`
	}
	if err := json.Unmarshal(trimmed, &env); err != nil {
		return "", ""
	}
	return stringifyEnvelopeValue(env.Message), stringifyEnvelopeValue(env.Code)
}

// stringifyEnvelopeValue renders a message/code field that may be a string,
// number, or array (the gateway flattens arrays, but older paths may not).
func stringifyEnvelopeValue(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case []interface{}:
		if len(typed) > 0 {
			return stringifyEnvelopeValue(typed[0])
		}
		return ""
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "<nil>" {
			return ""
		}
		return text
	}
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

func QueryValues(endpoint string, values url.Values) string {
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
