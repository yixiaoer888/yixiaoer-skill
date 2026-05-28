package domain

const SkillVersion = "3.0.1"

type SuccessResponse struct {
	OK      bool        `json:"ok"`
	Action  string      `json:"action,omitempty"`
	Version string      `json:"version"`
	Data    interface{} `json:"data,omitempty"`
	Next    interface{} `json:"next,omitempty"`
	Notice  interface{} `json:"_notice,omitempty"`
}

type ErrorResponse struct {
	OK      bool          `json:"ok"`
	Version string        `json:"version"`
	Error   ErrorEnvelope `json:"error"`
}

type ErrorEnvelope struct {
	Type        string      `json:"type,omitempty"`
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Hint        string      `json:"hint,omitempty"`
	Retryable   bool        `json:"retryable,omitempty"`
	NextCommand string      `json:"nextCommand,omitempty"`
	Details     interface{} `json:"details,omitempty"`
}
