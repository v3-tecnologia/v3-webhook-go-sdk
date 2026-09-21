<p align="center">
    <img src="./.github/logo.png" width="200px">
</p>

<h1 align="center" style="font-weight: bold;">v3-webhook-go-sdk</h1>

<p align="center">
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#requirements">Pré-requisitos</a> •
  <a href="#get-started">Como Executar</a> •
  <a href="#routing">Roteamento</a> •
  <a href="#testing">Testes</a> •
  <a href="#contribute">Como Contribuir</a>
</p>

<p align="center">
<b>SDK em Go para processar Webhooks IoT da V3 Tecnologia. Agnóstica de transporte (sem servidor HTTP). O roteamento é guiado pelos tipos protobuf do protocol-cloud — sem magic strings.</b>
</p>

<h2 id="tech-stack">💻 Tech Stack</h2>

Este projeto utiliza as seguintes tecnologias:

- [Go 1.25+](https://golang.org/)
- [protocol-cloud](https://github.com/v3-tecnologia/protocol-cloud) (modelos de domínio protobuf)
- [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf) (protojson)

<h2 id="requirements">❗ Pré-requisitos</h2>

Para usar esta SDK, você precisará ter instalado:

- [Go 1.25+](https://golang.org/doc/install)
- [Git](https://git-scm.com/downloads)
- Acesso aos módulos privados da V3 (`GOPRIVATE`)

<h2 id="get-started">🚀 Como executar?</h2>

### 0. Clonar o repositório

```bash
git clone git@github.com:v3-tecnologia/v3-webhook-go-sdk.git
cd v3-webhook-go-sdk
```

### 1. Módulos privados

```bash
export GOPRIVATE=github.com/v3-tecnologia/*
```

### 2. Instalar

```bash
go get github.com/v3-tecnologia/v3-webhook-go-sdk@latest
```

Ou, a partir do repositório:

```bash
go mod tidy
```

### 3. Uso rápido

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

Um exemplo executável está em `examples/webhook`:

```bash
go run ./examples/webhook
```

### 4. Validação de assinatura (opcional)

```go
builder := processing.NewBuilder().WithHMACSHA256("secret")
result := builder.Build().ProcessWebhook(ctx, body, signatureHeader)
```

### 5. Persistência (opcional)

```go
store := persistence.NewInMemoryStore()
builder := processing.NewBuilder().WithPersistence(
	persistence.NewInMemoryReader(store),
	persistence.NewInMemoryWriter(store),
)
```

<h2 id="routing">🧭 Roteamento</h2>

O protocolo é a fonte da verdade:

1. Parse do JSON → `domain.notifications.v1.Webhook` (protojson)
2. Para cada evento em `attributes[]`:
   - Se `data` estiver definido → exige `trip_event` ou `standalone_event`
   - Resolve o oneof `event_group` (envelope: `Dms`, `Alert`, `System`, ...)
   - Resolve o oneof do evento no envelope (payload: `DrowsinessEvent`, `ImpactEvent`, ...)
   - Despacha pelo **nome completo protobuf** do payload
   - Se `order` estiver definido → despacha pelo enum `orders.v1.OrderStatus`
3. Violações de protocolo retornam falha

Tipos de payload não registrados são ignorados (sucesso). Formato inválido do protocolo falha.

Produtores que já possuem eventos de domínio podem chamar `ProcessEvents`:

```go
result := processor.ProcessEvents(ctx, events)
```

#### EventContext

| Campo | Descrição |
|---|---|
| `ID` | Id do evento |
| `HasMedia` | Flag de mídia |
| `PayloadKind` | Derivado do envelope do protocolo |
| `Status` / `Type` / `Category` / `Sub` | Metadados do evento |
| `Device` | Device dos attributes |
| `Location` | Localização aninhada quando presente |
| `Save` / `GetEventByID` / ... | Helpers de persistência quando configurados |

#### Estrutura do projeto

```
.
├── examples/webhook/     # consumidor HTTP de exemplo
├── pkg/
│   ├── handlers/         # EventContext, results, ports de persistência
│   ├── persistence/      # reader/writer em memória
│   ├── processing/       # Builder, Processor, roteamento OnEvent
│   ├── security/         # validação HMAC-SHA256
│   └── types/            # helpers legados (prefira protocol-cloud)
├── test/events/          # fixtures JSON de webhook
└── README.md
```

> **Nota:** `pkg/types/*` mantém wrappers antigos por compatibilidade. Prefira `pkg/processing` + tipos do `protocol-cloud` em código novo.

<h2 id="testing">🧪 Executando Testes</h2>

Para executar os testes:

```bash
go test ./... -v
```

Pacotes core com race detector e cobertura:

```bash
go test ./pkg/processing ./pkg/security ./pkg/handlers ./pkg/persistence \
  -race -coverprofile=coverage.out -covermode=atomic

go tool cover -func=coverage.out | tail -1
go tool cover -html=coverage.out
```

Os pacotes core da SDK miram **>= 80%** de cobertura com o race detector habilitado.

Fixtures em `test/events/`.

<h2 id="contribute">📫 Como contribuir</h2>

1. Faça um fork do projeto
2. Crie uma branch para sua feature
   ```bash
   git checkout -b feature/nome-da-feature
   ```
3. Siga o padrão de commits convencional:
   - `feat:` para novas features
   - `fix:` para correção de bugs
   - `docs:` para atualização de documentação
   - `test:` para adição ou modificação de testes
   - `refactor:` para refatoração de código

4. Faça commit das suas alterações:
   ```bash
   git commit -m "feat: adiciona nova funcionalidade"
   ```

5. Faça push para sua branch:
   ```bash
   git push origin feature/nome-da-feature
   ```

6. Abra um Pull Request explicando a mudança e aguarde a revisão
