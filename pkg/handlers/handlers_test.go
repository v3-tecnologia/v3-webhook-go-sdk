package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/persistence"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	"google.golang.org/protobuf/proto"
)

func TestResultHelpers(t *testing.T) {
	ok := handlers.Success()
	if !ok.IsSuccess() {
		t.Fatal("expected success")
	}

	fail := handlers.Failure("x", errors.New("boom"))
	if fail.IsSuccess() || fail.ErrorMessage != "x" || fail.Err == nil {
		t.Fatalf("unexpected failure: %+v", fail)
	}
}

func TestPayloadKindFromEnvelope(t *testing.T) {
	cases := []struct {
		env  proto.Message
		want handlers.EventPayloadKind
	}{
		{&eventsv1.Dms{}, handlers.PayloadKindDms},
		{&eventsv1.Alert{}, handlers.PayloadKindAlert},
		{&eventsv1.Vision{}, handlers.PayloadKindVision},
		{&eventsv1.Connection{}, handlers.PayloadKindConnection},
		{&eventsv1.Health{}, handlers.PayloadKindHardware},
		{&eventsv1.Telemetry{}, handlers.PayloadKindTelemetry},
		{&eventsv1.DriverBehavior{}, handlers.PayloadKindDriverBehavior},
		{&eventsv1.Vehicle{}, handlers.PayloadKindVehicle},
		{&eventsv1.DriverIdentification{}, handlers.PayloadKindDriverIdentification},
		{&eventsv1.Geofence{}, handlers.PayloadKindGeofence},
		{&eventsv1.Adas{}, handlers.PayloadKindAdas},
		{&eventsv1.System{}, handlers.PayloadKindSystem},
		{&eventsv1.OrderStatus{}, handlers.PayloadKindOrder},
		{&eventsv1.DrowsinessEvent{}, handlers.PayloadKindUnknown},
	}

	for _, tc := range cases {
		if got := handlers.PayloadKindFromEnvelope(tc.env); got != tc.want {
			t.Fatalf("%T: got %s want %s", tc.env, got, tc.want)
		}
	}
}

func TestEventContextWithoutPersistence(t *testing.T) {
	ctx := handlers.NewEventContext("id", false, 0, nil, 0, 0, 0, nil, nil, handlers.PayloadKindSystem, nil, nil)
	if err := ctx.Save(context.Background(), &eventsv1.ImpactEvent{}); err == nil {
		t.Fatal("expected save error")
	}
	if _, err := ctx.GetEventByID(context.Background(), "id", &eventsv1.ImpactEvent{}); err == nil {
		t.Fatal("expected get by id error")
	}
	if _, err := ctx.GetEvents(context.Background(), &eventsv1.ImpactEvent{}, 1); err == nil {
		t.Fatal("expected get events error")
	}
	if _, err := ctx.GetRootEvent(context.Background(), &eventsv1.ImpactEvent{}, &eventsv1.ImpactEvent{}); err == nil {
		t.Fatal("expected get root error")
	}
}

func TestEventContextWithPersistence(t *testing.T) {
	store := persistence.NewInMemoryStore()
	writer := persistence.NewInMemoryWriter(store)
	reader := persistence.NewInMemoryReader(store)
	ec := handlers.NewEventContext("evt-1", true, 0, nil, 0, 0, 0, nil, nil, handlers.PayloadKindAlert, writer, reader)

	evt := &eventsv1.ImpactEvent{Name: "IMPACT"}
	if err := ec.Save(context.Background(), evt); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded := &eventsv1.ImpactEvent{}
	ok, err := ec.GetEventByID(context.Background(), "evt-1", loaded)
	if err != nil || !ok || loaded.GetName() != "IMPACT" {
		t.Fatalf("get by id failed ok=%v err=%v name=%s", ok, err, loaded.GetName())
	}

	list, err := ec.GetEvents(context.Background(), &eventsv1.ImpactEvent{}, 0)
	if err != nil || len(list) != 1 {
		t.Fatalf("get events failed: %v len=%d", err, len(list))
	}
}
