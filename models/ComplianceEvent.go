package models

import (
	"encoding/json"
	"errors"
	"time"
)

// ComplianceEventReceiver type
type ComplianceEventReceiver struct {
	EventID     string                      `json:"ExecutionId"`
	ResultLevel ComplianceEventStatusNumber `json:"ResultLevel"`
	Messages    []ComplianceEventMessage    `json:"StatusMessages"`
}

// ComplianceEventMessage type
type ComplianceEventMessage struct {
	EventTimestamp             time.Time                   `json:"ComplianceCheckTimestamp"`
	EventStatus                ComplianceEventStatusNumber `json:"ComplianceCheckLevel"`
	ComplianceCheckFunction    string                      `json:"ComplianceCheckFunction"`
	ComplianceCheckMessageText string                      `json:"ComplianceCheckMessageText"`
	ComplianceCategory         string                      `json:"ComplianceCategory"`
}

// ComplianceEventStatusNumber type
type ComplianceEventStatusNumber struct {
	Status int `json:"Status"`
}

// ComplianceEventStatus type
type ComplianceEventStatus struct {
	Status string
}

// ComplianceEvent type
type ComplianceEvent struct {
	EventID                      string
	IntResultLevel               int
	ResultLevel                  string
	ComplianceCategory           string
	ComplianceCheckFunction      string
	ComplianceCheckResultMessage string
	EventTimestamp               time.Time
	EventStatus                  string
}

// Status map[int]string
var status = map[int]string{
	0: "OK",
	1: "WARNING",
	2: "CRITICAL",
	3: "UNKNOWN",
}

// NewComplianceEvent constructor
func NewComplianceEvent(id string, intResult int, result, complianceCategory, complianceCheckFunction, complianceCheckMessage, eventStatus string, eventTimestamp time.Time) ComplianceEvent {
	return ComplianceEvent{
		EventID:                      id,
		IntResultLevel:               intResult,
		ResultLevel:                  result,
		ComplianceCategory:           complianceCategory,
		ComplianceCheckFunction:      complianceCheckFunction,
		ComplianceCheckResultMessage: complianceCheckMessage,
		EventStatus:                  eventStatus,
		EventTimestamp:               eventTimestamp,
	}
}

// ToComplianceEvents type
func (cer *ComplianceEventReceiver) ToComplianceEvents() ([]ComplianceEvent, error) {
	ce := []ComplianceEvent{}
	overallStatusText := status[cer.ResultLevel.Status]
	id := cer.EventID

	for _, message := range cer.Messages {
		event := NewComplianceEvent(
			id,
			cer.ResultLevel.Status,
			overallStatusText,
			message.ComplianceCategory,
			message.ComplianceCheckFunction,
			message.ComplianceCheckMessageText,
			status[message.EventStatus.Status],
			message.EventTimestamp,
		)
		ce = append(ce, event)
	}

	return ce, nil
}

// ToJson json stuff
func (ce ComplianceEvent) ToJson() (string, error) {
	event, err := json.Marshal(ce)

	if err != nil {
		return "", errors.New("could not marshal ComplianceEvent object to json")
	}
	return string(event), nil
}
