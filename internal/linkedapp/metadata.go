package linkedapp

type Metadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Entrypoints []string  `json:"entrypoints"`
	Controls    []Control `json:"controls"`
}

type Control struct {
	Type        string `json:"type"`
	Key         string `json:"key"`
	Label       string `json:"label"`
	Action      string `json:"action"`
	OnAction    string `json:"onAction,omitempty"`
	OffAction   string `json:"offAction,omitempty"`
	Description string `json:"description,omitempty"`
}

func DefaultMetadata() Metadata {
	return Metadata{
		ID:          "yixiaoer",
		Name:        "yixiaoer",
		DisplayName: "蚁小二",
		Description: "蚁小二多平台内容分发链接应用",
		Type:        "linked_app",
		Entrypoints: []string{"yxer skill show", "yxer linked-app status"},
		Controls: []Control{
			{
				Type:        "toggle",
				Key:         "connected",
				Label:       "连接蚁小二",
				Action:      "yxer linked-app toggle",
				OnAction:    "yxer linked-app connect",
				OffAction:   "yxer linked-app disconnect",
				Description: "控制 yixiaoer-skill 在 QClaw 中的链接状态",
			},
		},
	}
}
