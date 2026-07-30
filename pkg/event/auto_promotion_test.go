package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kargoapi "github.com/akuity/kargo/api/v1alpha1"
)

func TestAutoPromotionDenied(t *testing.T) {
	evt := &AutoPromotionDenied{}
	require.Equal(t, kargoapi.EventTypeAutoPromotionDenied, evt.Type())
}

func TestNewAutoPromotionDenied(t *testing.T) {
	freight := &kargoapi.Freight{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-freight",
			Namespace:         "test-project",
			CreationTimestamp: metav1.Time{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		Alias: "v1.0.0",
		Origin: kargoapi.FreightOrigin{
			Kind: kargoapi.FreightOriginKindWarehouse,
			Name: "test-warehouse",
		},
	}

	evt := NewAutoPromotionDenied("Auto-promotion denied", "test-actor", "test-stage", freight)

	require.Equal(t, kargoapi.EventTypeAutoPromotionDenied, evt.Type())
	// No Promotion was created, so the event is scoped to the candidate Freight
	// and carries the destination Stage alongside it.
	require.Equal(t, "Freight", evt.Kind())
	require.Equal(t, "test-freight", evt.GetName())
	require.Equal(t, "test-stage", evt.StageName)
	require.Equal(t, "test-project", evt.GetProject())
	require.Equal(t, "Auto-promotion denied", evt.GetMessage())
	require.Equal(t, "test-warehouse", evt.WarehouseName)
	require.NotNil(t, evt.Actor)
	require.Equal(t, "test-actor", *evt.Actor)
}

func TestAutoPromotionDeniedAnnotationsRoundTrip(t *testing.T) {
	freight := &kargoapi.Freight{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-freight",
			Namespace:         "test-project",
			CreationTimestamp: metav1.Time{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		Alias: "v1.0.0",
		Origin: kargoapi.FreightOrigin{
			Kind: kargoapi.FreightOriginKindWarehouse,
			Name: "test-warehouse",
		},
	}

	evt := NewAutoPromotionDenied("Auto-promotion denied", "test-actor", "test-stage", freight)

	annotations := evt.MarshalAnnotations()
	require.Equal(t, map[string]string{
		kargoapi.AnnotationKeyEventProject:              "test-project",
		kargoapi.AnnotationKeyEventActor:                "test-actor",
		kargoapi.AnnotationKeyEventFreightName:          "test-freight",
		kargoapi.AnnotationKeyEventFreightCreateTime:    "2024-01-01T00:00:00Z",
		kargoapi.AnnotationKeyEventFreightWarehouseName: "test-warehouse",
		kargoapi.AnnotationKeyEventFreightAlias:         "v1.0.0",
		kargoapi.AnnotationKeyEventStageName:            "test-stage",
	}, annotations)

	decoded, err := UnmarshalAutoPromotionDeniedAnnotations("test-event-id", annotations)
	require.NoError(t, err)
	require.Equal(t, "test-event-id", decoded.GetID())
	require.Equal(t, "test-project", decoded.GetProject())
	require.Equal(t, "test-freight", decoded.GetName())
	require.Equal(t, "test-stage", decoded.StageName)
	require.Equal(t, "test-warehouse", decoded.WarehouseName)
	require.NotNil(t, decoded.Alias)
	require.Equal(t, "v1.0.0", *decoded.Alias)
	require.NotNil(t, decoded.Actor)
	require.Equal(t, "test-actor", *decoded.Actor)
	// The message rides on the Kubernetes Event itself, not in annotations.
	require.Empty(t, decoded.GetMessage())
}
