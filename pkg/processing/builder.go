package processing

import (
	"context"
	"reflect"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/security"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	ordersv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/orders"
	"google.golang.org/protobuf/proto"
)

// HandlerFunc is the untyped handler signature used by the processor registry.
type HandlerFunc func(ctx context.Context, eventContext handlers.EventContext, payload proto.Message) handlers.EventHandlingResult

// DomainEventHandler receives the full protocol domain event after routing resolution.
type DomainEventHandler func(ctx context.Context, eventContext handlers.EventContext, evt *eventsv1.Event) handlers.EventHandlingResult

// Builder configures a Processor with typed protocol routes, signature validation and persistence.
type Builder struct {
	handlers           map[string]HandlerFunc
	orderHandlers      map[ordersv1.OrderStatus]HandlerFunc
	defaultHandler     DomainEventHandler
	signatureValidator security.SignatureValidator
	writer             handlers.EventWriter
	reader             handlers.EventReader
}

// NewBuilder creates an empty processor builder.
func NewBuilder() *Builder {
	return &Builder{
		handlers:      make(map[string]HandlerFunc),
		orderHandlers: make(map[ordersv1.OrderStatus]HandlerFunc),
	}
}

// WithHMACSHA256 enables HMAC-SHA256 signature validation using secret.
// Panics if secret is empty.
func (b *Builder) WithHMACSHA256(secret string) *Builder {
	validator, err := security.NewHMACSHA256Validator(secret)
	if err != nil {
		panic(err)
	}
	b.signatureValidator = validator
	return b
}

// WithSignatureValidator sets a custom signature validator.
func (b *Builder) WithSignatureValidator(validator security.SignatureValidator) *Builder {
	b.signatureValidator = validator
	return b
}

// WithPersistence configures both EventReader and EventWriter.
func (b *Builder) WithPersistence(reader handlers.EventReader, writer handlers.EventWriter) *Builder {
	b.reader = reader
	b.writer = writer
	return b
}

// WithEventWriter configures only the EventWriter.
func (b *Builder) WithEventWriter(writer handlers.EventWriter) *Builder {
	b.writer = writer
	return b
}

// WithEventReader configures only the EventReader.
func (b *Builder) WithEventReader(reader handlers.EventReader) *Builder {
	b.reader = reader
	return b
}

// On registers an untyped handler for a concrete protocol message full name.
func (b *Builder) On(messageFullName string, handler HandlerFunc) *Builder {
	if normalizeWhitespace(messageFullName) == "" {
		panic("message full name is required")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}
	b.handlers[messageFullName] = handler
	return b
}

// OnEvent registers a typed handler keyed by the protocol protobuf message type of T.
// T must be a pointer message from protocol-cloud (e.g. *eventsv1.DrowsinessEvent).
func OnEvent[T proto.Message](
	b *Builder,
	handler func(ctx context.Context, eventContext handlers.EventContext, evt T) handlers.EventHandlingResult,
) *Builder {
	if handler == nil {
		panic("handler cannot be nil")
	}
	sample := newProto[T]()
	key := messageKey(sample)
	return b.On(key, func(ctx context.Context, eventContext handlers.EventContext, payload proto.Message) handlers.EventHandlingResult {
		typed, ok := payload.(T)
		if !ok {
			return handlers.Failure("handler payload type mismatch for " + key)
		}
		return handler(ctx, eventContext, typed)
	})
}

// OnOrderStatus registers a handler for a specific protocol orders.OrderStatus enum value.
func OnOrderStatus(
	b *Builder,
	status ordersv1.OrderStatus,
	handler func(ctx context.Context, eventContext handlers.EventContext, order *eventsv1.OrderStatus) handlers.EventHandlingResult,
) *Builder {
	if status == ordersv1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		panic("order status must be a concrete protocol enum value")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}
	b.orderHandlers[status] = func(ctx context.Context, eventContext handlers.EventContext, payload proto.Message) handlers.EventHandlingResult {
		order, ok := payload.(*eventsv1.OrderStatus)
		if !ok {
			return handlers.Failure("handler payload type mismatch for OrderStatus")
		}
		return handler(ctx, eventContext, order)
	}
	return b
}

// WithDefaultHandler registers a handler invoked for every domain event after typed routing.
// Useful for producers (e.g. inbound-facade) that map to protocol events and need a single sink.
func (b *Builder) WithDefaultHandler(handler DomainEventHandler) *Builder {
	if handler == nil {
		panic("default handler cannot be nil")
	}
	b.defaultHandler = handler
	return b
}

// Build creates an immutable Processor from the builder configuration.
func (b *Builder) Build() *Processor {
	handlersCopy := make(map[string]HandlerFunc, len(b.handlers))
	for k, v := range b.handlers {
		handlersCopy[k] = v
	}
	ordersCopy := make(map[ordersv1.OrderStatus]HandlerFunc, len(b.orderHandlers))
	for k, v := range b.orderHandlers {
		ordersCopy[k] = v
	}
	return &Processor{
		handlers:           handlersCopy,
		orderHandlers:      ordersCopy,
		defaultHandler:     b.defaultHandler,
		signatureValidator: b.signatureValidator,
		writer:             b.writer,
		reader:             b.reader,
	}
}

func newProto[T proto.Message]() T {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() != reflect.Ptr {
		panic("OnEvent type parameter must be a pointer protobuf message")
	}
	return reflect.New(typ.Elem()).Interface().(T)
}
