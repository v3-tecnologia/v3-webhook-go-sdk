# Go Event Library

Uma biblioteca Go para processamento e análise de eventos de dispositivos IoT, facilitando o trabalho com dados de telemetria, ordens, conexões e eventos de sistema.

## Funcionalidades

- ✅ Parsing de eventos JSON estruturados
- ✅ Validação de eventos
- ✅ **SDK Webhook agnóstica** - sem dependências de servidor HTTP
- ✅ Integração com qualquer framework HTTP (Gin, Echo, net/http, etc.)
- ✅ Handlers com callbacks para facilitar consumo de eventos
- ✅ Suporte completo a 46+ tipos de eventos reais:
  - Ordens e confirmações (ACK)
  - Driver Monitoring System (DMS) - 11 tipos
  - Comportamento do motorista - 7 tipos
  - Hardware e sistema - 14 tipos
  - Telemetria e localização - 6 tipos
  - Visão básica - 5 tipos
- ✅ Tipos strongly-typed para todas as estruturas de dados
- ✅ Baseado em eventos reais de produção

## Instalação

```bash
go get github.com/your-org/go-eventlib
```

## Estrutura dos Eventos

A biblioteca suporta os seguintes tipos de eventos (baseado em 46+ exemplos reais):

### Categorias de Eventos

- **EVENT_CATEGORY_ORDER**: Eventos relacionados a ordens e configurações
  - `EVENT_SUB_ORDER_STATUS`: Status de ordens (ACK, SENT, FAILED)

- **EVENT_CATEGORY_CONNECTION**: Eventos de conectividade
  - WiFi conectado/desconectado
  - Mudanças de SIM card
  - Erros de conexão

- **EVENT_CATEGORY_VISION**: Eventos de visão e câmera
  - Face detectada/perdida/rastreada
  - Câmera obstruída
  - Sem face detectada

- **EVENT_CATEGORY_DMS**: Driver Monitoring System
  - Sonolência (DROWSINESS)
  - Distração visual (GAZE_DISTRACTION, GAZE_FIXATION)
  - Fechamento de olhos (EYE_CLOSURE)
  - Uso de celular (ON_PHONE)
  - Comendo/Bebendo (EATING, DRINKING)
  - Postura distraída (POSE_DISTRACTION)

- **EVENT_CATEGORY_DRIVER_BEHAVIOR**: Comportamento do motorista
  - Aceleração brusca (HARSH_ACCELERATION)
  - Frenagem brusca (HARSH_BRAKING)
  - Violação de velocidade (SPEED_VIOLATION)
  - Velocidade máxima persistente

- **EVENT_CATEGORY_HARDWARE**: Eventos de hardware
  - Reinício do dispositivo (RESTART, REBOOT)
  - Estado do dispositivo
  - Relatórios enviados
  - Processamento de ordens

- **EVENT_CATEGORY_SYSTEM**: Eventos de sistema
  - Uploads de arquivos
  - Alertas de sistema

- **EVENT_CATEGORY_TELEMETRY**: Eventos de telemetria
  - `EVENT_SUB_TELEMETRY_IGNITION`: Status de ignição (ON/OFF)
  - `EVENT_SUB_TELEMETRY_BATTERY`: Status de bateria (device/vehicle)
  - `EVENT_SUB_TELEMETRY_LOCATION`: Localização GPS

- **EVENT_CATEGORY_VEHICLE**: Eventos do veículo
  - Status de ignição
  - Métricas do veículo

- **EVENT_CATEGORY_ALERT**: Alertas gerais
  - Críticos, avisos e informativos

## Uso da SDK

A SDK processa eventos automaticamente através do `EventProcessor`. Você não precisa fazer parsing manual - apenas configure os handlers e processe os eventos:

```go
import (
    "go-eventlib/pkg/types/base"
    "go-eventlib/pkg/webhook"
)

processor := webhook.NewEventProcessor()
// Configurar handlers...
event, err := processor.ProcessEvent(ctx, jsonBytes)
// O evento já está parseado e os callbacks são chamados automaticamente
// event é do tipo *base.BaseEvent
```

## Estrutura dos Dados

### Estrutura de Pacotes

A biblioteca está organizada em pacotes contextuais para melhor modularidade:

- **`pkg/types/base`**: Tipos base compartilhados (`BaseEvent`, `EventCategory`, `EventSub`, etc.)
- **`pkg/types/common`**: Tipos comuns (`Location`, `Coordinates`, `GNSS`, `Connectivity`, `Fix`)
- **`pkg/types/order`**: Eventos de ordem (`order.Event`)
- **`pkg/types/connection`**: Eventos de conexão (`connection.Event`)
- **`pkg/types/vision`**: Eventos de visão (`vision.Event`)
- **`pkg/types/hardware`**: Eventos de hardware (`hardware.Event`)
- **`pkg/types/system`**: Eventos de sistema (`system.Event`)
- **`pkg/types/telemetry`**: Eventos de telemetria (`telemetry.Event`)
- **`pkg/types/alert`**: Eventos de alerta (`alert.Event`)
- **`pkg/types/dms`**: Eventos DMS (`dms.Event`)
- **`pkg/types/driverbehavior`**: Eventos de comportamento (`driverbehavior.Event`)
- **`pkg/types/vehicle`**: Eventos de veículo (`vehicle.Event`)

### Evento Base
```go
import "go-eventlib/pkg/types/base"

type BaseEvent struct {
    ID        string            `json:"id"`
    Status    base.EventStatus `json:"status"`
    CreatedAt time.Time         `json:"created_at"`
    Type      base.EventType    `json:"type"`
    Category  base.EventCategory `json:"category"`
    Sub       base.EventSub     `json:"sub"`
    Attributes base.Attributes  `json:"attributes"`
}

// Métodos helper disponíveis:
// GetID() string
// GetCategory() base.EventCategory
// GetSubType() base.EventSub
// GetDeviceID() string
// GetCreatedAt() time.Time
```

### Tipos Específicos de Eventos

Cada categoria de evento tem seu próprio tipo com métodos helper:

```go
import (
    "go-eventlib/pkg/types/order"
    "go-eventlib/pkg/types/connection"
    "go-eventlib/pkg/types/telemetry"
)

// Order Event
orderEvent := order.New(baseEvent)
order := orderEvent.GetOrder()

// Connection Event
connEvent := connection.New(baseEvent)
wifi := connEvent.GetWifiConnection()
simCard := connEvent.GetSimCard()

// Telemetry Event
telemetryEvent := telemetry.New(baseEvent)
telemetryData := telemetryEvent.GetTelemetryData()
batteryMetrics := telemetryEvent.GetBatteryMetrics()
```

## Verificação de Tipos de Evento

```go
import "go-eventlib/pkg/types/base"

// Verificar categoria do evento
if event.Category == base.EventCategoryOrder {
    // Lidar com evento de ordem
}

if event.Category == base.EventCategoryConnection {
    // Lidar com evento de conexão
}

if event.Category == base.EventCategoryVision {
    // Lidar com evento de visão
}

// Acessar informações do evento usando métodos helper
deviceID := event.GetDeviceID()
accountID := event.Attributes.Device.AccountID
```

## SDK Webhook Agnóstica

O **SDK Webhook agnóstico** permite consumir eventos IoT de forma simples. **Não inclui servidor HTTP** - você cria seu próprio servidor HTTP e usa a SDK apenas para parsear eventos.

### Uso Básico com Callbacks

```go
package main

import (
    "context"
    "encoding/json"
    "io"
    "log"
    "net/http"

    "go-eventlib/pkg/types/connection"
    "go-eventlib/pkg/types/order"
    "go-eventlib/pkg/webhook"
)

func main() {
    processor := webhook.NewEventProcessor()

    orderHandler := webhook.NewOrderHandler()
    orderHandler.OnOrderReceived = func(ctx context.Context, event *order.Event) error {
        if ord := event.GetOrder(); ord != nil {
            log.Printf("Ordem recebida: ID=%s, Status=%s", event.ID, ord.Status)
        } else {
            log.Printf("Ordem recebida: ID=%s", event.ID)
        }
        return nil
    }
    orderHandler.OnOrderAck = func(ctx context.Context, event *order.Event) error {
        if ord := event.GetOrder(); ord != nil {
            log.Printf("Ordem confirmada: ID=%s, Status=%s", event.ID, ord.Status)
        } else {
            log.Printf("Ordem confirmada: ID=%s", event.ID)
        }
        return nil
    }
    processor.SetOrderHandler(orderHandler)

    connectionHandler := webhook.NewConnectionHandler()
    connectionHandler.OnWifiConnected = func(ctx context.Context, event *connection.Event) error {
        if wifi := event.GetWifiConnection(); wifi != nil {
            log.Printf("WiFi conectado: Device=%s, Status=%s", event.GetDeviceID(), wifi.Status)
        } else {
            log.Printf("WiFi conectado: Device=%s", event.GetDeviceID())
        }
        return nil
    }
    processor.SetConnectionHandler(connectionHandler)

    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        ctx := context.Background()
        event, err := processor.ProcessEvent(ctx, body)
        if err != nil {
            http.Error(w, "Error processing event", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":   "processed",
            "event_id": event.ID,
            "category": event.Category,
            "type":     event.Type,
        })
    })

    log.Println("Servidor iniciado na porta 8080")
    http.ListenAndServe(":8080", nil)
}
```

### Handlers Disponíveis

A SDK oferece handlers pré-definidos com callbacks para facilitar o consumo de eventos:

#### OrderHandler
```go
import "go-eventlib/pkg/types/order"

orderHandler := webhook.NewOrderHandler()
orderHandler.OnOrderReceived = func(ctx context.Context, event *order.Event) error {
    if ord := event.GetOrder(); ord != nil {
        // Processar ordem recebida com dados específicos
        log.Printf("Ordem: %s, Status: %s", ord.ID, ord.Status)
    }
    return nil
}
orderHandler.OnOrderAck = func(ctx context.Context, event *order.Event) error {
    // Processar ordem confirmada
    return nil
}
processor.SetOrderHandler(orderHandler)
```

#### ConnectionHandler
```go
import "go-eventlib/pkg/types/connection"

connectionHandler := webhook.NewConnectionHandler()
connectionHandler.OnWifiConnected = func(ctx context.Context, event *connection.Event) error {
    if wifi := event.GetWifiConnection(); wifi != nil {
        log.Printf("WiFi conectado: %s", wifi.Name)
    }
    return nil
}
connectionHandler.OnWifiDisconnected = func(ctx context.Context, event *connection.Event) error {
    // WiFi desconectado
    return nil
}
connectionHandler.OnSimCardChanged = func(ctx context.Context, event *connection.Event) error {
    if sim := event.GetSimCard(); sim != nil {
        log.Printf("SIM card: %s", sim.Status)
    }
    return nil
}
connectionHandler.OnConnectionError = func(ctx context.Context, event *connection.Event) error {
    // Erro de conexão
    return nil
}
processor.SetConnectionHandler(connectionHandler)
```

#### VisionHandler
```go
import "go-eventlib/pkg/types/vision"

visionHandler := webhook.NewVisionHandler()
visionHandler.OnFaceDetected = func(ctx context.Context, event *vision.Event) error {
    if faceData := event.GetFaceDetectedData(); faceData != nil {
        // Processar dados de face detectada
    }
    return nil
}
visionHandler.OnFaceLost = func(ctx context.Context, event *vision.Event) error {
    // Face perdida
    return nil
}
visionHandler.OnFaceTracked = func(ctx context.Context, event *vision.Event) error {
    // Face rastreada
    return nil
}
visionHandler.OnNoFaceDetected = func(ctx context.Context, event *vision.Event) error {
    // Sem face detectada
    return nil
}
visionHandler.OnCameraObstructed = func(ctx context.Context, event *vision.Event) error {
    // Câmera obstruída
    return nil
}
visionHandler.OnVisionAlert = func(ctx context.Context, event *vision.Event) error {
    // Alerta de visão genérico
    return nil
}
processor.SetVisionHandler(visionHandler)
```

#### HardwareHandler
```go
import "go-eventlib/pkg/types/hardware"

hardwareHandler := webhook.NewHardwareHandler()
hardwareHandler.OnDeviceRestart = func(ctx context.Context, event *hardware.Event) error {
    if sysData := event.GetSystemEventData(); sysData != nil {
        log.Printf("Dispositivo reiniciado: %s", sysData.EventName)
    }
    return nil
}
hardwareHandler.OnVehicleBatteryConnected = func(ctx context.Context, event *hardware.Event) error {
    // Bateria do veículo conectada
    return nil
}
hardwareHandler.OnVehicleBatteryDisconnected = func(ctx context.Context, event *hardware.Event) error {
    // Bateria do veículo desconectada
    return nil
}
hardwareHandler.OnSDCardMounted = func(ctx context.Context, event *hardware.Event) error {
    // SD card montado
    return nil
}
hardwareHandler.OnSDCardUnmounted = func(ctx context.Context, event *hardware.Event) error {
    // SD card desmontado
    return nil
}
hardwareHandler.OnSimCardInserted = func(ctx context.Context, event *hardware.Event) error {
    // SIM card inserido
    return nil
}
hardwareHandler.OnSimCardRemoved = func(ctx context.Context, event *hardware.Event) error {
    // SIM card removido
    return nil
}
hardwareHandler.OnHardwareAlert = func(ctx context.Context, event *hardware.Event) error {
    // Alerta de hardware
    return nil
}
processor.SetHardwareHandler(hardwareHandler)
```

#### SystemHandler
```go
import "go-eventlib/pkg/types/system"

systemHandler := webhook.NewSystemHandler()
systemHandler.OnUploadEvent = func(ctx context.Context, event *system.Event) error {
    if upload := event.GetUploadData(); upload != nil {
        log.Printf("Upload: %s", upload.Name)
    }
    return nil
}
systemHandler.OnSystemAlert = func(ctx context.Context, event *system.Event) error {
    // Alerta de sistema
    return nil
}
processor.SetSystemHandler(systemHandler)
```

#### TelemetryHandler
```go
import "go-eventlib/pkg/types/telemetry"

telemetryHandler := webhook.NewTelemetryHandler()
telemetryHandler.OnBatteryEvent = func(ctx context.Context, event *telemetry.Event) error {
    if metrics := event.GetBatteryMetrics(); metrics != nil {
        // Processar métricas de bateria
    }
    return nil
}
telemetryHandler.OnIgnitionEvent = func(ctx context.Context, event *telemetry.Event) error {
    if telData := event.GetTelemetryData(); telData != nil {
        log.Printf("Ignition: %s", telData.Status)
    }
    return nil
}
telemetryHandler.OnLocationEvent = func(ctx context.Context, event *telemetry.Event) error {
    // Evento de localização
    return nil
}
telemetryHandler.OnTelemetryEvent = func(ctx context.Context, event *telemetry.Event) error {
    // Qualquer evento de telemetria
    return nil
}
processor.SetTelemetryHandler(telemetryHandler)
```

#### AlertHandler
```go
import "go-eventlib/pkg/types/alert"

alertHandler := webhook.NewAlertHandler()
alertHandler.OnCriticalAlert = func(ctx context.Context, event *alert.Event) error {
    log.Printf("Alerta crítico: %s", event.GetAlertLevel())
    return nil
}
alertHandler.OnWarningAlert = func(ctx context.Context, event *alert.Event) error {
    // Alerta de aviso
    return nil
}
alertHandler.OnInfoAlert = func(ctx context.Context, event *alert.Event) error {
    // Alerta informativo
    return nil
}
alertHandler.OnAnyAlert = func(ctx context.Context, event *alert.Event) error {
    if alertData := event.GetAlertEventData(); alertData != nil {
        log.Printf("Alerta: %s", alertData.EventName)
    }
    return nil
}
processor.SetAlertHandler(alertHandler)
```

#### DMSHandler
```go
import "go-eventlib/pkg/types/dms"

dmsHandler := webhook.NewDMSHandler()
dmsHandler.OnDrowsiness = func(ctx context.Context, event *dms.Event) error {
    if drowsiness := event.GetDrowsinessData(); drowsiness != nil {
        // Processar dados de sonolência
    }
    return nil
}
dmsHandler.OnDrinking = func(ctx context.Context, event *dms.Event) error {
    if drinking := event.GetDrinkingData(); drinking != nil {
        // Processar dados de bebida
    }
    return nil
}
dmsHandler.OnEating = func(ctx context.Context, event *dms.Event) error {
    // Comendo detectado
    return nil
}
dmsHandler.OnEyeClosure = func(ctx context.Context, event *dms.Event) error {
    // Fechamento de olhos
    return nil
}
dmsHandler.OnGazeDistraction = func(ctx context.Context, event *dms.Event) error {
    // Distração visual
    return nil
}
dmsHandler.OnGazeFixation = func(ctx context.Context, event *dms.Event) error {
    // Fixação do olhar
    return nil
}
dmsHandler.OnPhone = func(ctx context.Context, event *dms.Event) error {
    // Uso de celular
    return nil
}
dmsHandler.OnPoseDistraction = func(ctx context.Context, event *dms.Event) error {
    // Postura distraída
    return nil
}
dmsHandler.OnSmoking = func(ctx context.Context, event *dms.Event) error {
    // Fumando detectado
    return nil
}
dmsHandler.OnYawning = func(ctx context.Context, event *dms.Event) error {
    // Bocejando detectado
    return nil
}
dmsHandler.OnDMSAlert = func(ctx context.Context, event *dms.Event) error {
    // Qualquer alerta DMS
    return nil
}
processor.SetDMSHandler(dmsHandler)
```

#### DriverBehaviorHandler
```go
import "go-eventlib/pkg/types/driverbehavior"

driverBehaviorHandler := webhook.NewDriverBehaviorHandler()
driverBehaviorHandler.OnHarshAcceleration = func(ctx context.Context, event *driverbehavior.Event) error {
    if accel := event.GetHarshAccelerationData(); accel != nil {
        // Processar dados de aceleração brusca
    }
    return nil
}
driverBehaviorHandler.OnHarshBraking = func(ctx context.Context, event *driverbehavior.Event) error {
    if braking := event.GetHarshBrakingData(); braking != nil {
        // Processar dados de frenagem brusca
    }
    return nil
}
driverBehaviorHandler.OnMaxSpeedFault = func(ctx context.Context, event *driverbehavior.Event) error {
    // Violação de velocidade máxima
    return nil
}
driverBehaviorHandler.OnNormalSpeedReturn = func(ctx context.Context, event *driverbehavior.Event) error {
    // Retorno à velocidade normal
    return nil
}
driverBehaviorHandler.OnPersistentMaxSpeed = func(ctx context.Context, event *driverbehavior.Event) error {
    // Velocidade máxima persistente
    return nil
}
driverBehaviorHandler.OnSharpTurn = func(ctx context.Context, event *driverbehavior.Event) error {
    // Curva brusca
    return nil
}
driverBehaviorHandler.OnStartOvertaking = func(ctx context.Context, event *driverbehavior.Event) error {
    // Início de ultrapassagem
    return nil
}
driverBehaviorHandler.OnDriverBehaviorAlert = func(ctx context.Context, event *driverbehavior.Event) error {
    // Qualquer alerta de comportamento
    return nil
}
processor.SetDriverBehaviorHandler(driverBehaviorHandler)
```

#### VehicleHandler
```go
import "go-eventlib/pkg/types/vehicle"

vehicleHandler := webhook.NewVehicleHandler()
vehicleHandler.OnIgnitionOff = func(ctx context.Context, event *vehicle.Event) error {
    // Ignição desligada
    return nil
}
vehicleHandler.OnVehicleEvent = func(ctx context.Context, event *vehicle.Event) error {
    // Qualquer evento do veículo
    return nil
}
processor.SetVehicleHandler(vehicleHandler)
```

### Tipos de Eventos Suportados

A SDK suporta processamento automático de callbacks para os seguintes tipos de eventos:

#### Eventos de Ordem
- `ORDER_STATUS_ACK` - Ordem confirmada
- `ORDER_STATUS_SENT` - Ordem enviada
- `ORDER_STATUS_FAILED` - Ordem falhou

#### Eventos de Conexão
- `WIFI_CONNECTED` - WiFi conectado
- `WIFI_DISCONNECTED` - WiFi desconectado
- Mudanças de SIM card
- Erros críticos de conexão

#### Eventos DMS (Driver Monitoring System)
- `DROWSINESS` - Sonolência detectada
- `GAZE_DISTRACTION` - Distração visual
- `GAZE_FIXATION` - Fixação do olhar
- `EYE_CLOSURE` - Fechamento de olhos
- `ON_PHONE` - Uso de celular
- `EATING` - Comendo
- `DRINKING` - Bebendo
- `POSE_DISTRACTION` - Postura distraída

#### Eventos de Comportamento do Motorista
- `HARSH_ACCELERATION` - Aceleração brusca
- `HARSH_BRAKING` - Frenagem brusca
- `SPEED_VIOLATION` - Violação de velocidade
- Velocidade máxima persistente

#### Eventos de Hardware
- `RESTART` / `REBOOT` - Reinício do dispositivo
- Estados do dispositivo
- Relatórios enviados

#### Eventos de Telemetria
- `EVENT_SUB_TELEMETRY_IGNITION` - Status de ignição
- `EVENT_SUB_TELEMETRY_BATTERY` - Status de bateria
- `EVENT_SUB_TELEMETRY_LOCATION` - Localização GPS

## Exemplos

Confira os exemplos em `examples/`:

- `webhook/main.go`: Exemplo de servidor webhook usando net/http

Para executar o exemplo:

```bash
cd examples/webhook
go run main.go
```

### Integração com Frameworks HTTP

#### Com net/http (padrão)
```go
processor := webhook.NewEventProcessor()
// Configurar handlers...

http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    ctx := context.Background()
    event, err := processor.ProcessEvent(ctx, body)
    if err != nil {
        http.Error(w, "Error", 500)
        return
    }
    // Callbacks são chamados automaticamente
})
```

#### Com Gin
```go
processor := webhook.NewEventProcessor()
// Configurar handlers...

r := gin.Default()
r.POST("/webhook", func(c *gin.Context) {
    body, _ := c.GetRawData()
    event, err := processor.ProcessEvent(c.Request.Context(), body)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    // Callbacks são chamados automaticamente
})
```

#### Com Echo
```go
processor := webhook.NewEventProcessor()
// Configurar handlers...

e := echo.New()
e.POST("/webhook", func(c echo.Context) error {
    body, _ := io.ReadAll(c.Request().Body)
    event, err := processor.ProcessEvent(c.Request().Context(), body)
    if err != nil {
        return c.JSON(500, map[string]string{"error": err.Error()})
    }
    // Callbacks são chamados automaticamente
    return c.JSON(200, event)
})
```

### Teste Prático

```bash
# Executar exemplo com servidor HTTP customizado
cd examples/webhook
go run main.go

# Em outro terminal, testar webhook com evento de ordem
curl -X POST http://localhost:8080/webhook/events \
  -H "Content-Type: application/json" \
  -d @../../event-mapping-viewer/src/mindmap-sections/serialized-jsons/ack-events/ack-order-event.json

# Testar com evento de telemetria
curl -X POST http://localhost:8080/webhook/events \
  -H "Content-Type: application/json" \
  -d @../../event-mapping-viewer/src/mindmap-sections/serialized-jsons/telemetry-events/telemetry-ignition.json

# Testar com evento DMS
curl -X POST http://localhost:8080/webhook/events \
  -H "Content-Type: application/json" \
  -d @../../event-mapping-viewer/src/mindmap-sections/serialized-jsons/dms-events/vision-drowsiness.json
```

### Exemplos de Payloads de Teste

#### Evento de Ordem (ACK)
```json
{
  "id": "01KCHNX0HKAAAZD4ZMTY35P8XJ",
  "status": "STATUS_RECEIVED",
  "created_at": "2025-12-15T18:55:59.748719972Z",
  "type": "EVENT_TYPE_ORDER",
  "category": "EVENT_CATEGORY_ORDER",
  "sub": "EVENT_SUB_ORDER_STATUS",
  "attributes": {
    "device": {
      "id": "01KBJ38J7DG4831CW327H1C908",
      "correlation_id": "01KBJ38J7D369CAFVGNWM6H2QX",
      "uid": "862798052131337",
      "account_id": "01GZXXCVVPEKM7E830XAMJKA14"
    },
    "order": {
      "id": "01KCHNWZGS5VXG8EY5SDZ9EVHW",
      "correlation_id": "01KCHNX0HKAAAZD4ZMTY35P8XJ",
      "group": "ORDER_GROUP_CONFIG",
      "status": "ORDER_STATUS_ACK",
      "type": "CONFIG",
      "created_at": "2025-12-15T18:55:53Z",
      "updated_at": "2025-12-15T18:55:54Z"
    }
  }
}
```

#### Evento de Telemetria (Ignition)
```json
{
  "id": "01HXXXXXXXXXXXXXXXXXXXXX",
  "status": "STATUS_RECEIVED",
  "type": "EVENT_TYPE_GENERAL",
  "category": "EVENT_CATEGORY_VEHICLE",
  "sub": "EVENT_SUB_TELEMETRY_IGNITION",
  "created_at": "2024-12-15T10:30:00Z",
  "attributes": {
    "device": {
      "id": "device-123",
      "correlation_id": "device-corr-123",
      "uid": "IMEI123456789012345",
      "account_id": "account-456"
    },
    "data": {
      "telemetry": {
        "id": "telemetry-789",
        "status": "IGNITION_STATUS_ON",
        "hardware": {
          "model": {
            "name": "V3_DEVICE_PRO",
            "vendor": "V3_TECHNOLOGIA"
          }
        },
        "connection": {
          "type": "CELLULAR",
          "signal_strength": "-75dBm"
        },
        "metrics": {
          "device_battery": {
            "component": "BATTERY_COMPONENT_DEVICE",
            "status": "BATTERY_ONLINE",
            "voltage": 4.2
          }
        },
        "timestamp": "2024-12-15T10:30:00Z"
      }
    }
  }
}
```

#### Evento DMS (Sonolência)
```json
{
  "id": "01KCHNFZWN4YPSM0A4YFMH09T2",
  "status": "STATUS_RECEIVED",
  "type": "EVENT_TYPE_GENERAL",
  "category": "EVENT_CATEGORY_DMS",
  "sub": "EVENT_SUB_DMS_BASIC",
  "created_at": "2025-12-15T19:02:27.853477693Z",
  "attributes": {
    "device": {
      "id": "01KBZJ4WWBE5N78S4F3DASW1WJ",
      "correlation_id": "01KBZJ4WWBEA7TK7JGR5CD7S20",
      "uid": "862798051074124",
      "account_id": "01GZXXCVVPEKM7E830XAMJKA14"
    },
    "data": {
      "trip_event": {
        "event_group_name": "DMS",
        "dms": {
          "event_name": "DROWSINESS",
          "drowsiness": {
            "name": "DROWSINESS",
            "confidence": 0.90,
            "bounding_box": {
              "x": 0.3,
              "y": 0.2,
              "width": 0.4,
              "height": 0.5
            },
            "attributes": {
              "perclos": "0.25",
              "blinks_per_min": "8"
            }
          }
        }
      }
    }
  }
}
```

### Arquivos de Exemplo para Testes

A biblioteca inclui 46+ exemplos reais de eventos JSON organizados por categoria:

- **ack-events/** (3 arquivos): Eventos de confirmação de ordens e uploads
- **dms-events/** (11 arquivos): Eventos de Driver Monitoring System
- **driver-behavior-events/** (7 arquivos): Eventos de comportamento do motorista
- **hardware-events/** (14 arquivos): Eventos de hardware e sistema
- **telemetry-events/** (6 arquivos): Eventos de telemetria e localização
- **vision-basic-events/** (5 arquivos): Eventos básicos de visão

Esses arquivos estão disponíveis em `event-mapping-viewer/src/mindmap-sections/serialized-jsons/` e podem ser usados para testes e desenvolvimento.

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.
