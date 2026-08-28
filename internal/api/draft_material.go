package api

import (
	"fmt"
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
	err := c.Patch(fmt.Sprintf("/material/%s", materialID), body, &result)
	return result, err
}

type MaterialGroupOptions struct {
	Page int
	Size int
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
