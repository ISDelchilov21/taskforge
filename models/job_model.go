package models

import "encoding/json"

type Job struct {
	ID             int             `json:"id"`
	Status         string          `json:"status"`
	Type           string          `json:"type"`
	Description    string          `json:"description"`
	Payload        json.RawMessage `json:"payload"`
	PriorityWeight int             `json:"priority_weight"`
}

type CreateJobRequest struct {
	Type           string          `json:"type" binding:"required,max=50"`
	Description    string          `json:"description" binding:"max=500"`
	Payload        json.RawMessage `json:"payload"`
	PriorityWeight int             `json:"priority_weight" binding:"min=0,max=100"`
}

type UpdateJobStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending running done failed"`
}
