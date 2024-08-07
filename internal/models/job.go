package models

import "time"

type Job struct {
	ID         string `json:"id,omitempty"`
	Status     string `json:"status,omitempty"`
	Type  	  string `json:"type,omitempty"`
	Payload    interface {}
	Retries    int    `json:"retries,omitempty"`
	MaxRetries int    `json:"maxRetries,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	DelayUntil time.Time `json:"delayUntil,omitempty"`
	RateLimit int    `json:"rateLimit,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	
}
