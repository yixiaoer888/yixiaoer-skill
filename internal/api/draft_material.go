package api

import (
	"net/url"
	"strconv"
)

func (c *Client) SaveDraft(body map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.Put("/taskSets/drafts", body, &result)
	return result, err
}

func (c *Client) Material(body map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.Post("/material", body, &result)
	return result, err
}

func (c *Client) MoveMaterial(materialID string, body map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	if materialID != "" {
		body["materialIds"] = []string{materialID}
	}
	err := c.Post("/material/batch/set-group", body, &result)
	return result, err
}

type MaterialGroupOptions struct {
	Page int
	Size int
}

// MaterialListOptions describes the filters accepted by the material library.
// FileName is passed through to the service and is also matched exactly by the
// higher-level move-by-name workflow before it performs a write.
type MaterialListOptions struct {
	Page     int
	Size     int
	FileName string
	Type     string
	GroupID  string
}

func (c *Client) Materials(opts MaterialListOptions) (interface{}, error) {
	values := url.Values{}
	if opts.Page > 0 {
		values.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size > 0 {
		values.Set("size", strconv.Itoa(opts.Size))
	}
	if opts.FileName != "" {
		values.Set("fileName", opts.FileName)
	}
	if opts.Type != "" {
		values.Set("type", opts.Type)
	}
	if opts.GroupID != "" {
		values.Set("groupId", opts.GroupID)
	}
	return c.queryData(QueryValues("/material", values))
}

func (c *Client) MaterialGroups(opts MaterialGroupOptions) (interface{}, error) {
	values := url.Values{}
	if opts.Page > 0 {
		values.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size > 0 {
		values.Set("size", strconv.Itoa(opts.Size))
	}
	endpoint := "/material/groups"
	if encoded := values.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}
	var result interface{}
	if err := c.Get(endpoint, &result); err != nil {
		return nil, err
	}
	if typed, ok := result.(map[string]interface{}); ok {
		return DataOrSelf(typed), nil
	}
	return result, nil
}
