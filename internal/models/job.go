package models 

import (
    "encoding/json"
    "time"
)

type Job struct {
    ID        string          `json:"id"`
    Type      string          `json:"type"`
    Payload   json.RawMessage `json:"payload"`
    Status    string          `json:"status"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

type EmailJob struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

type ImageResizeJob struct {
    InputPath  string `json:"input_path"`
    OutputPath string `json:"output_path"`
    Width      int    `json:"width"`
    Height     int    `json:"height"`
}

type DataProcessingJob struct {
    Data []int `json:"data"`
}



