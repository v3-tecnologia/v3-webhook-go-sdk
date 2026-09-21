# v3-webhook-go-sdk

Go SDK for processing **V3 Tecnologia IoT Webhooks**, mirrored from [`v3-webhook-dotnet-sdk`](https://github.com/v3-tecnologia/v3-webhook-dotnet-sdk).

Transport-agnostic (no HTTP server). Routing is driven by **protocol-cloud protobuf types** — not magic strings.

## Features

- Strongly-typed payloads from `protocol-cloud` (protobuf + protojson)
- Type-based routing: `OnEvent[*eventsv1.DrowsinessEvent](builder, handler)`
- Order routing by protocol enum: `OnOrderStatus(builder, orders.OrderStatus_ORDER_STATUS_ACK, handler)`
- Fails on protocol violations (empty `event_group` / missing oneof payload)
- Optional HMAC SHA-256 signature validation
- Optional in-memory persistence helpers

## Install

```bash
export GOPRIVATE=github.com/v3-tecnologia/*
go get github.com/v3-tecnologia/v3-webhook-go-sdk@latest
```

## Quick start

```go
package main

import (
	"context"
	"io"
	"net/http"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/processing"
	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	ordersv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/orders"
)

func main() {
	builder := processing.NewBuilder().
		WithHMACSHA256("your-secret-key")

	processing.OnEvent(builder,
		func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.DrowsinessEvent) handlers.EventHandlingResult {
			_ = evt.GetConfidence()
			return handlers.Success()
		},
	)

	processing.OnEvent(builder,
		func(ctx context.Context, ec handlers.EventContext, evt *eventsv1.UploadEvent) handlers.EventHandlingResult {
			return handlers.Success()
		},
	)

	processing.OnOrderStatus(builder, ordersv1.OrderStatus_ORDER_STATUS_ACK,
		func(ctx context.Context, ec handlers.EventContext, order *eventsv1.OrderStatus) handlers.EventHandlingResult {
			return handlers.Success()
		},
	)

	processor := builder.Build()

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		result := processor.ProcessWebhook(r.Context(), body, r.Header.Get("X-V3-Signature"))
		if !result.IsSuccess() {
			http.Error(w, result.ErrorMessage, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})

	_ = http.ListenAndServe(":8080", nil)
}
```

## Routing model

The protocol is the source of truth:

1. Parse JSON → `domain.notifications.v1.Webhook` (protojson)
2. For each `attributes[]` event:
   - If `data` is set → require `trip_event` or `standalone_event`
   - Resolve `event_group` oneof (envelope: `Dms`, `Alert`, `System`, ...)
   - Resolve envelope event oneof (payload: `DrowsinessEvent`, `ImpactEvent`, ...)
   - Dispatch by **protobuf message full name** of the payload
   - If `order` is set → dispatch by `orders.v1.OrderStatus` enum (or `*events.OrderStatus` type)
3. Protocol violations (missing oneofs, empty standalone, UNSPECIFIED order status) return failure

Unregistered payload types are skipped (success). Invalid protocol shape fails.

## EventContext

| Field | Description |
|---|---|
| `ID` | Event id |
| `HasMedia` | Media flag |
| `PayloadKind` | Derived from protocol envelope type |
| `Status` / `Type` / `Category` / `Sub` | Event metadata |
| `Device` | Device from attributes |
| `Location` | Nested location when present |
| `Save` / `GetEventByID` / ... | Persistence helpers when configured |

## Signature validation

```go
builder := processing.NewBuilder().WithHMACSHA256("secret")
result := builder.Build().ProcessWebhook(ctx, body, signatureHeader)
```

## Persistence (optional)

```go
store := persistence.NewInMemoryStore()
builder := processing.NewBuilder().WithPersistence(
	persistence.NewInMemoryReader(store),
	persistence.NewInMemoryWriter(store),
)
```

## Tests

```bash
go test ./pkg/processing ./pkg/security ./pkg/handlers ./pkg/persistence \
  -race -coverprofile=coverage.out -covermode=atomic

go tool cover -func=coverage.out | tail -1
```

Core SDK packages target **>= 80%** statement coverage with the race detector enabled.

Fixtures live under `test/events/`.

## Legacy helpers

`pkg/types/*` contains older hand-rolled wrappers kept for compatibility. Prefer `pkg/processing` + `protocol-cloud` types for new code.
