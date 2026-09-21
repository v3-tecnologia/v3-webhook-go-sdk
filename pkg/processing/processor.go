package processing

import (
	"context"
	"fmt"
	"strings"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/security"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	notificationsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/notifications"
	ordersv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/orders"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Processor parses webhook JSON into protocol-cloud messages and dispatches by protobuf type.
type Processor struct {
	handlers           map[string]HandlerFunc
	orderHandlers      map[ordersv1.OrderStatus]HandlerFunc
	defaultHandler     DomainEventHandler
	signatureValidator security.SignatureValidator
	writer             handlers.EventWriter
	reader             handlers.EventReader
}

// ProcessWebhook validates an optional signature, parses the webhook envelope as protocol
// protobuf JSON and invokes handlers registered for the concrete payload message types.
func (p *Processor) ProcessWebhook(ctx context.Context, jsonPayload []byte, signature string) handlers.EventHandlingResult {
	if len(strings.TrimSpace(string(jsonPayload))) == 0 {
		return handlers.Failure("webhook payload is empty")
	}

	if p.signatureValidator != nil {
		if strings.TrimSpace(signature) == "" {
			return handlers.Failure("missing webhook signature")
		}
		if err := p.signatureValidator.Validate(jsonPayload, signature); err != nil {
			return handlers.Failure("invalid webhook signature", err)
		}
	}

	webhook := &notificationsv1.Webhook{}
	unmarshaler := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	if err := unmarshaler.Unmarshal(jsonPayload, webhook); err != nil {
		return handlers.Failure("failed to parse webhook JSON as protocol-cloud Webhook", err)
	}

	return p.ProcessEvents(ctx, webhook.GetAttributes())
}

// ProcessEvents routes already-parsed protocol domain events (no JSON).
// Used by producers such as inbound-facade-worker after mapping firmware → domain.
func (p *Processor) ProcessEvents(ctx context.Context, events []*eventsv1.Event) handlers.EventHandlingResult {
	if len(events) == 0 {
		return handlers.Success()
	}

	for _, evt := range events {
		result := p.processAttribute(ctx, evt)
		if !result.IsSuccess() {
			return result
		}
	}

	return handlers.Success()
}

func (p *Processor) processAttribute(ctx context.Context, evt *eventsv1.Event) handlers.EventHandlingResult {
	if evt == nil {
		return handlers.Failure("protocol violation: nil event attribute")
	}

	attrs := evt.GetAttributes()
	if attrs == nil {
		return p.invokeDefault(ctx, evt, nil)
	}

	if data := attrs.GetData(); data != nil {
		resolved, err := resolveFromEventData(data)
		if err != nil {
			return handlers.Failure(err.Error(), err)
		}
		if resolved != nil {
			key := messageKey(resolved.payload)
			if handler, ok := p.handlers[key]; ok {
				result := p.invokeHandler(ctx, handler, evt, resolved.payload, resolved.envelope)
				if !result.IsSuccess() {
					return result
				}
			}
			return p.invokeDefault(ctx, evt, resolved.envelope)
		}
		return p.invokeDefault(ctx, evt, nil)
	}

	if order := attrs.GetOrder(); order != nil {
		status := order.GetStatus()
		if status == ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
			return handlers.Failure("protocol violation: order status is UNSPECIFIED")
		}

		if handler, ok := p.orderHandlers[status]; ok {
			result := p.invokeHandler(ctx, handler, evt, order, order)
			if !result.IsSuccess() {
				return result
			}
		} else if handler, ok := p.handlers[messageKey(order)]; ok {
			result := p.invokeHandler(ctx, handler, evt, order, order)
			if !result.IsSuccess() {
				return result
			}
		}

		return p.invokeDefault(ctx, evt, order)
	}

	return p.invokeDefault(ctx, evt, nil)
}

func (p *Processor) invokeDefault(ctx context.Context, evt *eventsv1.Event, envelope proto.Message) handlers.EventHandlingResult {
	if p.defaultHandler == nil {
		return handlers.Success()
	}

	attrs := evt.GetAttributes()
	var device = (*ordersv1.Device)(nil)
	if attrs != nil {
		device = attrs.GetDevice()
	}

	eventContext := handlers.NewEventContext(
		evt.GetId(),
		evt.GetHasMedia(),
		evt.GetStatus(),
		evt.GetCreatedAt(),
		evt.GetType(),
		evt.GetCategory(),
		evt.GetSub(),
		device,
		extractLocation(envelope),
		handlers.PayloadKindFromEnvelope(envelope),
		p.writer,
		p.reader,
	)

	result := p.defaultHandler(ctx, eventContext, evt)
	if !result.IsSuccess() && result.ErrorMessage == "" {
		return handlers.Failure(fmt.Sprintf("default handler failed for event id=%s", evt.GetId()), result.Err)
	}
	return result
}

func (p *Processor) invokeHandler(
	ctx context.Context,
	handler HandlerFunc,
	evt *eventsv1.Event,
	payload proto.Message,
	envelope proto.Message,
) handlers.EventHandlingResult {
	attrs := evt.GetAttributes()
	device := attrs.GetDevice()

	eventContext := handlers.NewEventContext(
		evt.GetId(),
		evt.GetHasMedia(),
		evt.GetStatus(),
		evt.GetCreatedAt(),
		evt.GetType(),
		evt.GetCategory(),
		evt.GetSub(),
		device,
		extractLocation(payload),
		handlers.PayloadKindFromEnvelope(envelope),
		p.writer,
		p.reader,
	)

	result := handler(ctx, eventContext, payload)
	if !result.IsSuccess() && result.ErrorMessage == "" {
		return handlers.Failure(fmt.Sprintf("handler execution failed for event id=%s", evt.GetId()), result.Err)
	}
	return result
}
