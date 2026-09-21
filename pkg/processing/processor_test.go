package processing

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/persistence"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	locationv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/location"
	ordersv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/orders"
	"google.golang.org/protobuf/proto"
)

func TestProcessWebhook_DMSEvents(t *testing.T) {
	cases := []struct {
		file string
		reg  func(*Builder, func())
	}{
		{"vision-yawning.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.YawningEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-drowsiness.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.DrowsinessEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-drinking.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.DrinkingEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-eating.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.EatingEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-eye-closure.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.EyeClosureEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-gaze-distraction.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.GazeDistractionEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-gaze-fixation.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.GazeFixationEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-pose-distraction-pitch.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.PoseDistractionPitchEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-pose-distraction-yaw.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.PoseDistractionYawEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-smoking.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.SmokingEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
		{"vision-on-phone.json", func(b *Builder, on func()) {
			OnEvent(b, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.OnPhoneEvent) handlers.EventHandlingResult {
				on()
				return handlers.Success()
			})
		}},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			payload, err := wrapWebhookFixture(filepath.Join("dms-events", tc.file))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}

			handled := false
			builder := NewBuilder()
			tc.reg(builder, func() { handled = true })

			result := builder.Build().ProcessWebhook(context.Background(), payload, "")
			if !result.IsSuccess() {
				t.Fatalf("process failed: %s (%v)", result.ErrorMessage, result.Err)
			}
			if !handled {
				t.Fatal("handler was not called")
			}
		})
	}
}

func TestProcessEvents_DefaultHandler(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	builder := NewBuilder()
	called := false
	builder.WithDefaultHandler(func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.Event) handlers.EventHandlingResult {
		called = true
		if evt.GetId() == "" {
			t.Fatal("expected event id")
		}
		if ec.PayloadKind != handlers.PayloadKindDms {
			t.Fatalf("expected dms kind, got %s", ec.PayloadKind)
		}
		return handlers.Success()
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("process failed: %s", result.ErrorMessage)
	}
	if !called {
		t.Fatal("default handler was not called")
	}
}

func TestProcessWebhook_SkipsUnknownRoute(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	result := NewBuilder().Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("expected success for unmatched route, got: %s", result.ErrorMessage)
	}
}

func TestProcessWebhook_ProtocolViolation_EmptyStandalone(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("ack-events", "order-status-event.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	result := NewBuilder().Build().ProcessWebhook(context.Background(), payload, "")
	if result.IsSuccess() {
		t.Fatal("expected protocol violation for empty standalone_event")
	}
}

func TestProcessWebhook_EmptyPayload(t *testing.T) {
	result := NewBuilder().Build().ProcessWebhook(context.Background(), []byte("   "), "")
	if result.IsSuccess() {
		t.Fatal("expected failure for empty payload")
	}
}

func TestProcessWebhook_InvalidJSON(t *testing.T) {
	result := NewBuilder().Build().ProcessWebhook(context.Background(), []byte("{"), "")
	if result.IsSuccess() {
		t.Fatal("expected failure for invalid json")
	}
}

func TestProcessWebhook_EmptyAttributes(t *testing.T) {
	payload := []byte(`{"id":"w1","attributes":[]}`)
	result := NewBuilder().Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("expected success, got %s", result.ErrorMessage)
	}
}

func TestProcessWebhook_HMAC(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	secret := "test-secret"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	builder := NewBuilder().WithHMACSHA256(secret)
	OnEvent(builder, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.DrowsinessEvent) handlers.EventHandlingResult {
		return handlers.Success()
	})

	ok := builder.Build().ProcessWebhook(context.Background(), payload, signature)
	if !ok.IsSuccess() {
		t.Fatalf("expected valid signature success: %s", ok.ErrorMessage)
	}

	missing := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if missing.IsSuccess() {
		t.Fatal("expected missing signature failure")
	}

	invalid := builder.Build().ProcessWebhook(context.Background(), payload, "deadbeef")
	if invalid.IsSuccess() {
		t.Fatal("expected invalid signature failure")
	}
}

func TestProcessWebhook_OrderAck(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("ack-events", "ack-order-event.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	handled := false
	builder := NewBuilder()
	OnOrderStatus(builder, ordersv1.OrderStatus_ORDER_STATUS_ACK,
		func(ctx context.Context, ec handlers.EventContext, order *eventsv1.OrderStatus) handlers.EventHandlingResult {
			handled = true
			if ec.PayloadKind != handlers.PayloadKindOrder {
				t.Fatalf("unexpected payload kind: %s", ec.PayloadKind)
			}
			if order.GetId() == "" {
				t.Fatal("expected order id")
			}
			return handlers.Success()
		},
	)

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("process failed: %s (%v)", result.ErrorMessage, result.Err)
	}
	if !handled {
		t.Fatal("order handler was not called")
	}
}

func TestProcessWebhook_OrderByMessageType(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("ack-events", "ack-order-event.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	handled := false
	builder := NewBuilder()
	OnEvent(builder, func(ctx context.Context, ec handlers.EventContext, order *eventsv1.OrderStatus) handlers.EventHandlingResult {
		handled = true
		return handlers.Success()
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() || !handled {
		t.Fatalf("expected typed order handler, success=%v handled=%v err=%s", result.IsSuccess(), handled, result.ErrorMessage)
	}
}

func TestProcessWebhook_SystemReboot(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("hardware-events", "hardware-reboot.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	handled := false
	builder := NewBuilder()
	OnEvent(builder, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.RebootEvent) handlers.EventHandlingResult {
		handled = true
		if ec.PayloadKind != handlers.PayloadKindSystem {
			t.Fatalf("unexpected kind %s", ec.PayloadKind)
		}
		return handlers.Success()
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("process failed: %s", result.ErrorMessage)
	}
	if !handled {
		t.Fatal("reboot handler was not called")
	}
}

func TestProcessWebhook_AlertImpactWithPersistence(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("alert-events", "alert-impact.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	store := persistence.NewInMemoryStore()
	writer := persistence.NewInMemoryWriter(store)
	reader := persistence.NewInMemoryReader(store)

	builder := NewBuilder().WithPersistence(reader, writer)
	OnEvent(builder, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.ImpactEvent) handlers.EventHandlingResult {
		if ec.Location == nil {
			t.Fatal("expected location on impact event")
		}
		if err := ec.Save(ctx, evt); err != nil {
			return handlers.Failure("save failed", err)
		}
		loaded := &eventsv1.ImpactEvent{}
		ok, err := ec.GetEventByID(ctx, ec.ID, loaded)
		if err != nil || !ok {
			return handlers.Failure("get by id failed", err)
		}
		list, err := ec.GetEvents(ctx, &eventsv1.ImpactEvent{}, 5)
		if err != nil || len(list) == 0 {
			return handlers.Failure("get events failed", err)
		}
		return handlers.Success()
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("process failed: %s (%v)", result.ErrorMessage, result.Err)
	}
}

func TestProcessWebhook_HandlerFailure(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	builder := NewBuilder()
	OnEvent(builder, func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.DrowsinessEvent) handlers.EventHandlingResult {
		return handlers.Failure("boom")
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if result.IsSuccess() {
		t.Fatal("expected handler failure")
	}
}

func TestProcessWebhook_HandlerFailureWithoutMessage(t *testing.T) {
	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	builder := NewBuilder()
	key := messageKey(&eventsv1.DrowsinessEvent{})
	builder.On(key, func(ctx context.Context, ec handlers.EventContext, payload proto.Message) handlers.EventHandlingResult {
		return handlers.EventHandlingResult{Success: false, Err: context.Canceled}
	})

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if result.IsSuccess() || result.ErrorMessage == "" {
		t.Fatal("expected synthesized error message")
	}
}

func TestProcessWebhook_EventWithoutAttributes(t *testing.T) {
	payload := []byte(`{"id":"w1","attributes":[{"id":"e1","status":"STATUS_RECEIVED"}]}`)
	result := NewBuilder().Build().ProcessWebhook(context.Background(), payload, "")
	if !result.IsSuccess() {
		t.Fatalf("expected success, got %s", result.ErrorMessage)
	}
}

func TestOnEvent_TypeMismatch(t *testing.T) {
	builder := NewBuilder()
	key := messageKey(&eventsv1.DrowsinessEvent{})
	builder.On(key, func(ctx context.Context, ec handlers.EventContext, payload proto.Message) handlers.EventHandlingResult {
		return handlers.Failure("handler payload type mismatch for " + key)
	})

	payload, err := wrapWebhookFixture(filepath.Join("dms-events", "vision-drowsiness.json"))
	if err != nil {
		t.Fatalf("load fixture: %v", err)
	}

	result := builder.Build().ProcessWebhook(context.Background(), payload, "")
	if result.IsSuccess() {
		t.Fatal("expected type mismatch failure")
	}
}

func TestBuilder_WithWriterAndReader(t *testing.T) {
	store := persistence.NewInMemoryStore()
	builder := NewBuilder().
		WithEventWriter(persistence.NewInMemoryWriter(store)).
		WithEventReader(persistence.NewInMemoryReader(store)).
		WithSignatureValidator(nil)
	if builder.Build() == nil {
		t.Fatal("expected processor")
	}
}

func TestBuilder_Panics(t *testing.T) {
	assertPanic(t, func() { NewBuilder().WithHMACSHA256("") })
	assertPanic(t, func() { NewBuilder().On("", func(context.Context, handlers.EventContext, proto.Message) handlers.EventHandlingResult { return handlers.Success() }) })
	assertPanic(t, func() { NewBuilder().On("x", nil) })
	assertPanic(t, func() {
		OnEvent[*eventsv1.DrowsinessEvent](NewBuilder(), nil)
	})
	assertPanic(t, func() {
		OnOrderStatus(NewBuilder(), ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED, func(context.Context, handlers.EventContext, *eventsv1.OrderStatus) handlers.EventHandlingResult {
			return handlers.Success()
		})
	})
	assertPanic(t, func() {
		OnOrderStatus(NewBuilder(), ordersv1.OrderStatus_ORDER_STATUS_ACK, nil)
	})
}

func TestExtractLocation_Direct(t *testing.T) {
	loc := &locationv1.Location{}
	if got := extractLocation(loc); got != loc {
		t.Fatal("expected same location pointer")
	}
	if got := extractLocation(nil); got != nil {
		t.Fatal("expected nil")
	}
}

func TestResolveFromEventData_Errors(t *testing.T) {
	if _, err := resolveFromEventData(nil); err != nil {
		t.Fatal("nil data should not error")
	}
	if _, err := resolveFromEventData(&eventsv1.EventData{}); err == nil {
		t.Fatal("expected missing trip/standalone error")
	}
	if _, err := resolveContainer(nil); err == nil {
		t.Fatal("expected empty container error")
	}
	if _, err := resolveContainer(&eventsv1.TripEvent{}); err == nil {
		t.Fatal("expected unset event_group error")
	}
}

func TestMessageKey(t *testing.T) {
	key := messageKey(&eventsv1.DrowsinessEvent{})
	if key == "" || key != string((&eventsv1.DrowsinessEvent{}).ProtoReflect().Descriptor().FullName()) {
		t.Fatalf("unexpected key %s", key)
	}
	if messageKey(nil) != "" {
		t.Fatal("expected empty key for nil")
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

func wrapWebhookFixture(rel string) ([]byte, error) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "test", "events", rel))
	if err != nil {
		return nil, err
	}

	var event json.RawMessage
	if err := json.Unmarshal(raw, &event); err != nil {
		return nil, err
	}

	wrapper := map[string]any{
		"id":         "webhook-test-id",
		"created_at": "2025-12-15T19:02:27.853477693Z",
		"attributes": []json.RawMessage{event},
	}
	return json.Marshal(wrapper)
}
