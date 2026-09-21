<p align="center">
    <img src="./.github/logo.png" width="200px">
</p>

<h1 align="center" style="font-weight: bold;">v3-webhook-go-sdk</h1>

<p align="center">
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#requirements">Requirements / Pré-requisitos</a> •
  <a href="#get-started">Get Started / Como executar</a> •
  <a href="#routing">Routing / Roteamento</a> •
  <a href="#testing">Testing / Testes</a> •
  <a href="#contribute">Contributing / Como contribuir</a>
</p>

<p align="center">
<b>Go SDK for processing V3 Tecnologia IoT Webhooks, mirrored from <a href="https://github.com/v3-tecnologia/v3-webhook-dotnet-sdk">v3-webhook-dotnet-sdk</a>. Transport-agnostic (no HTTP server). Routing is driven by protocol-cloud protobuf types — not magic strings.</b>
</p>

<p align="center">
<b>SDK em Go para processar Webhooks IoT da V3 Tecnologia, espelhada do <a href="https://github.com/v3-tecnologia/v3-webhook-dotnet-sdk">v3-webhook-dotnet-sdk</a>. Agnóstica de transporte (sem servidor HTTP). O roteamento é guiado pelos tipos protobuf do protocol-cloud — sem magic strings.</b>
</p>

<h2 id="tech-stack">💻 Tech Stack</h2>

This project uses the following technologies / Este projeto utiliza as seguintes tecnologias:

- [Go 1.25+](https://golang.org/)
- [protocol-cloud](https://github.com/v3-tecnologia/protocol-cloud) (protobuf domain models)
- [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf) (protojson)

<h2 id="requirements">❗ Requirements / Pré-requisitos</h2>

To use this SDK, you need / Para usar esta SDK, você precisa:

- [Go 1.25+](https://golang.org/doc/install)
- [Git](https://git-scm.com/downloads)
- Access to private V3 modules (`GOPRIVATE`) / Acesso aos módulos privados da V3 (`GOPRIVATE`)

<h2 id="get-started">🚀 Get Started / Como executar</h2>

### 0. Clone the repository / Clonar o repositório

```bash
git clone git@github.com:v3-tecnologia/v3-webhook-go-sdk.git
cd v3-webhook-go-sdk
```

### 1. Private modules / Módulos privados

```bash
export GOPRIVATE=github.com/v3-tecnologia/*
```

### 2. Install / Instalar

```bash
go get github.com/v3-tecnologia/v3-webhook-go-sdk@latest
```

Or, from the repository / Ou, a partir do repositório:

```bash
go mod tidy
```

### 3. Quick start / Uso rápido

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

A runnable example lives at `examples/webhook` / Um exemplo executável está em `examples/webhook`:

```bash
go run ./examples/webhook
```

### 4. Signature validation (optional) / Validação de assinatura (opcional)

```go
builder := processing.NewBuilder().WithHMACSHA256("secret")
result := builder.Build().ProcessWebhook(ctx, body, signatureHeader)
```

### 5. Persistence (optional) / Persistência (opcional)

```go
store := persistence.NewInMemoryStore()
builder := processing.NewBuilder().WithPersistence(
	persistence.NewInMemoryReader(store),
	persistence.NewInMemoryWriter(store),
)
```

<h2 id="routing">🧭 Routing / Roteamento</h2>

The protocol is the source of truth / O protocolo é a fonte da verdade:

1. Parse JSON → `domain.notifications.v1.Webhook` (protojson)
2. For each `attributes[]` event / Para cada evento em `attributes[]`:
   - If `data` is set → require `trip_event` or `standalone_event` / Se `data` estiver definido → exige `trip_event` ou `standalone_event`
   - Resolve `event_group` oneof (envelope: `Dms`, `Alert`, `System`, ...)
   - Resolve envelope event oneof (payload: `DrowsinessEvent`, `ImpactEvent`, ...)
   - Dispatch by **protobuf message full name** of the payload / Despacha pelo **nome completo protobuf** do payload
   - If `order` is set → dispatch by `orders.v1.OrderStatus` enum / Se `order` estiver definido → despacha pelo enum `orders.v1.OrderStatus`
3. Protocol violations return failure / Violações de protocolo retornam falha

Unregistered payload types are skipped (success). Invalid protocol shape fails.  
Tipos de payload não registrados são ignorados (sucesso). Formato inválido do protocolo falha.

Producers that already hold domain events can call `ProcessEvents` / Produtores que já possuem eventos de domínio podem chamar `ProcessEvents`:

```go
result := processor.ProcessEvents(ctx, events)
```

#### EventContext

| Field / Campo | Description / Descrição |
|---|---|
| `ID` | Event id / Id do evento |
| `HasMedia` | Media flag / Flag de mídia |
| `PayloadKind` | Derived from protocol envelope / Derivado do envelope do protocolo |
| `Status` / `Type` / `Category` / `Sub` | Event metadata / Metadados do evento |
| `Device` | Device from attributes / Device dos attributes |
| `Location` | Nested location when present / Localização aninhada quando presente |
| `Save` / `GetEventByID` / ... | Persistence helpers when configured / Helpers de persistência quando configurados |

#### Project structure / Estrutura do projeto

```
.
├── examples/webhook/     # HTTP sample consumer
├── pkg/
│   ├── handlers/         # EventContext, results, persistence ports
│   ├── persistence/      # In-memory reader/writer
│   ├── processing/       # Builder, Processor, OnEvent routing
│   ├── security/         # HMAC-SHA256 signature validation
│   └── types/            # Legacy helpers (prefer protocol-cloud)
├── test/events/          # Webhook JSON fixtures
└── README.md
```

> **Note / Nota:** `pkg/types/*` keeps older hand-rolled wrappers for compatibility. Prefer `pkg/processing` + `protocol-cloud` types for new code.  
> `pkg/types/*` mantém wrappers antigos por compatibilidade. Prefira `pkg/processing` + tipos do `protocol-cloud` em código novo.

<h2 id="testing">🧪 Testing / Testes</h2>

To run tests / Para executar os testes:

```bash
go test ./... -v
```

Core packages with race detector and coverage / Pacotes core com race detector e cobertura:

```bash
go test ./pkg/processing ./pkg/security ./pkg/handlers ./pkg/persistence \
  -race -coverprofile=coverage.out -covermode=atomic

go tool cover -func=coverage.out | tail -1
go tool cover -html=coverage.out
```

Core SDK packages target **>= 80%** statement coverage with the race detector enabled.  
Os pacotes core da SDK miram **>= 80%** de cobertura com o race detector habilitado.

Fixtures live under `test/events/` / Fixtures em `test/events/`.

<h2 id="contribute">📫 Contributing / Como contribuir</h2>

1. Fork the project / Faça um fork do projeto
2. Create a feature branch / Crie uma branch para sua feature
   ```bash
   git checkout -b feature/nome-da-feature
   ```
3. Follow conventional commits / Siga o padrão de commits convencional:
   - `feat:` new features / novas features
   - `fix:` bug fixes / correção de bugs
   - `docs:` documentation / atualização de documentação
   - `test:` tests / adição ou modificação de testes
   - `refactor:` refactoring / refatoração de código

4. Commit your changes / Faça commit das suas alterações:
   ```bash
   git commit -m "feat: add new functionality"
   ```

5. Push your branch / Faça push para sua branch:
   ```bash
   git push origin feature/nome-da-feature
   ```

6. Open a Pull Request describing the change and wait for review / Abra um Pull Request explicando a mudança e aguarde a revisão
