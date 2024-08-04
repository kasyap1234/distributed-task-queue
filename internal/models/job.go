package models

import "time"

type Job struct {
	ID         string `json:"id,omitempty"`
	Status     string `json:"status,omitempty"`
	Payload    string `json:"payload,omitempty"`
	Retries    int    `json:"retries,omitempty"`
	MaxRetries int    `json:"maxRetries,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	DelayUntil time.Time `json:"delayUntil,omitempty"`
	RateLimit int    `json:"rateLimit,omitempty"`
	Dependenceies []string `json:"dependencies,omitempty"`
	
}
