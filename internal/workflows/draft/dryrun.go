package draft

func PreviewSave(payload map[string]interface{}) map[string]interface{} {
	body := cloneMap(payload)
	delete(body, "action")
	body["isDraft"] = true
	return body
}
