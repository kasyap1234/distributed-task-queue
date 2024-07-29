package models 
type Job struct {
	ID string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Payload string `json:"payload,omitempty"`
	Retries int `json:"retries,omitempty"`
	MaxRetries int `json:"maxRetries,omitempty"`

}
