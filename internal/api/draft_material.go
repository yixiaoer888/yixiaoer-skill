package api

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
