package models

import "encoding/json"

type Job struct {
	ID             int
	Status         string
	Type           string
	Description    string
	Payload        json.RawMessage
	PriorityWeight int
}
