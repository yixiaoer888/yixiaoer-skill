package domain

const SkillVersion = "3.0.0"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Action  string      `json:"action"`
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success    bool        `json:"success"`
	ErrorCode  string      `json:"errorCode"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	Suggestion string      `json:"suggestion,omitempty"`
}
