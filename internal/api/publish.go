package api

func (c *Client) Publish(body map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.Post("/taskSets/v2", body, &result)
	return result, err
}
