package event

import (
	kargoapi "github.com/akuity/kargo/api/v1alpha1"
)

// AutoPromotionDenied is event data recorded when an admission webhook denies
// the creation of an auto-promotion for a Stage. Because no Promotion is
// created, the event references the Stage itself rather than a Promotion.
//
// The denial is recorded regardless of which admission policy rejected the
// create. It is a lightweight event and, unlike the Promotion and Freight
// lifecycle events, is not registered in KnownEventTypes; consumers receive it
// as a Custom event.
type AutoPromotionDenied struct {
	ID          string  `json:"id,omitempty"`
	Project     string  `json:"project"`
	Actor       *string `json:"actor,omitempty"`
	StageName   string  `json:"stageName"`
	FreightName string  `json:"freightName"`
	Message     string  `json:"message,omitempty"`
}

// NewAutoPromotionDenied creates a new AutoPromotionDenied event for the given
// Stage and Freight. The actor is recorded when non-empty.
func NewAutoPromotionDenied(
	message, actor, stageName, project, freightName string,
) *AutoPromotionDenied {
	evt := &AutoPromotionDenied{
		Project:     project,
		StageName:   stageName,
		FreightName: freightName,
		Message:     message,
	}
	if actor != "" {
		evt.Actor = &actor
	}
	return evt
}

func (e *AutoPromotionDenied) Type() kargoapi.EventType {
	return kargoapi.EventTypeAutoPromotionDenied
}

func (e *AutoPromotionDenied) Kind() string {
	return "Stage"
}

func (e *AutoPromotionDenied) GetName() string {
	return e.StageName
}

func (e *AutoPromotionDenied) GetProject() string {
	return e.Project
}

func (e *AutoPromotionDenied) GetID() string {
	return e.ID
}

func (e *AutoPromotionDenied) GetMessage() string {
	return e.Message
}

func (e *AutoPromotionDenied) SetMessage(msg string) {
	e.Message = msg
}
