package event

import (
	"fmt"

	kargoapi "github.com/akuity/kargo/api/v1alpha1"
)

// AutoPromotionDenied is event data recorded when an admission webhook denies
// the creation of an auto-promotion. Auto-promotion is evaluated per requested
// Freight origin, so the event identifies the candidate Freight that was denied
// and the Stage it was destined for.
//
// The denial is recorded regardless of which admission policy rejected the
// create. No Promotion exists to reference, so this borrows the shape of the
// other Freight-scoped events.
type AutoPromotionDenied struct {
	Common
	Freight
}

func (a *AutoPromotionDenied) Type() kargoapi.EventType {
	return kargoapi.EventTypeAutoPromotionDenied
}

// NewAutoPromotionDenied creates a new `AutoPromotionDenied` event for the
// candidate Freight that could not be auto-promoted into the named Stage.
func NewAutoPromotionDenied(message, actor, stageName string, freight *kargoapi.Freight,
) *AutoPromotionDenied {
	common, freightEvent := NewFreightCommon(message, actor, stageName, freight)
	return &AutoPromotionDenied{
		Common:  common,
		Freight: freightEvent,
	}
}

func (a *AutoPromotionDenied) MarshalAnnotations() map[string]string {
	annotations := map[string]string{}
	a.Common.MarshalAnnotationsTo(annotations)
	a.Freight.MarshalAnnotationsTo(annotations)
	return annotations
}

// UnmarshalAutoPromotionDeniedAnnotations converts the given annotations into an
// AutoPromotionDenied event. This is used by the main event handler to convert
// the data into a normal structured event, but is exposed for convenience.
func UnmarshalAutoPromotionDeniedAnnotations(
	eventID string,
	annotations map[string]string,
) (*AutoPromotionDenied, error) {
	freight, err := UnmarshalFreightAnnotations(annotations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal freight annotations: %w", err)
	}
	common, err := UnmarshalCommonAnnotations(eventID, annotations)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal common annotations: %w", err)
	}
	evt := AutoPromotionDenied{
		Common:  common,
		Freight: freight,
	}
	return &evt, nil
}
