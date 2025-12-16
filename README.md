# Go Event Library

A Go library for processing and analyzing IoT device events, facilitating work with telemetry data, orders, connections, and system events.

## Features

- ✅ JSON event parsing using protocol-cloud (protojson)
- ✅ Protocol Buffers-based event validation
- ✅ **Agnostic Webhook SDK** - no HTTP server dependencies
- ✅ Integration with any HTTP framework (Gin, Echo, net/http, etc.)
- ✅ Handlers with callbacks to facilitate event consumption
- ✅ Modular architecture with processors per event context
- ✅ Full support for 46+ real event types:
  - Orders and confirmations (ACK)
  - Driver Monitoring System (DMS) - 11 types
  - Driver behavior - 7 types
  - Hardware and system - 14 types
  - Telemetry and location - 6 types
  - Basic vision - 5 types
- ✅ Strongly-typed types for all data structures
- ✅ Based on real production events
- ✅ Compatible with protocol-cloud (v3-tecnologia)

## Installation

```bash
go get github.com/your-org/go-eventlib
```

## Event Structure

The library supports the following event types (based on 46+ real examples):

### Event Categories

- **EVENT_CATEGORY_ORDER**: Events related to orders and configurations
  - `EVENT_SUB_ORDER_STATUS`: Order status (ACK, SENT, FAILED)

- **EVENT_CATEGORY_CONNECTION**: Connectivity events
  - WiFi connected/disconnected
  - SIM card changes
  - Connection errors

- **EVENT_CATEGORY_VISION**: Vision and camera events
  - Face detected/lost/tracked
  - Camera obstructed
  - No face detected

- **EVENT_CATEGORY_DMS**: Driver Monitoring System
  - Drowsiness (DROWSINESS)
  - Visual distraction (GAZE_DISTRACTION, GAZE_FIXATION)
  - Eye closure (EYE_CLOSURE)
  - Phone usage (ON_PHONE)
  - Eating/Drinking (EATING, DRINKING)
  - Distracted posture (POSE_DISTRACTION)

- **EVENT_CATEGORY_DRIVER_BEHAVIOR**: Driver behavior
  - Harsh acceleration (HARSH_ACCELERATION)
  - Harsh braking (HARSH_BRAKING)
  - Speed violation (SPEED_VIOLATION)
  - Persistent maximum speed

- **EVENT_CATEGORY_HEALTH**: Hardware/device health events
  - Device restart (RESTART, REBOOT, R2_RESTART)
  - Device state
  - SD Card mounted/unmounted
  - SIM Card inserted/removed
  - Vehicle battery connected/disconnected

- **EVENT_CATEGORY_SYSTEM**: System events
  - File uploads
  - System alerts

- **EVENT_CATEGORY_TELEMETRY**: Telemetry events
  - `EVENT_SUB_TELEMETRY_IGNITION`: Ignition status (ON/OFF)
  - `EVENT_SUB_TELEMETRY_BATTERY`: Battery status (device/vehicle)
  - `EVENT_SUB_TELEMETRY_LOCATION`: GPS location

- **EVENT_CATEGORY_VEHICLE**: Vehicle events
  - Ignition status
  - Vehicle metrics

- **EVENT_CATEGORY_ALERT**: General alerts
  - Critical, warnings, and informational

## SDK Usage

The SDK processes events automatically through the `EventProcessor`. Parsing is done using `protocol-cloud` with `protojson`, ensuring full compatibility with the V3 protocol. You don't need to do manual parsing - just configure the handlers and process the events.

### Builder Pattern (Recommended)

Use `EventProcessorBuilder` to configure the processor fluently:

```go
import (
    "go-eventlib/pkg/types/base"
    "go-eventlib/pkg/webhook"
)

processor := webhook.NewEventProcessorBuilder().
    WithEventHandler(eventHandler).
    WithConnectionHandler(connectionHandler).
    WithVisionHandler(visionHandler).
    Build()

event, err := processor.ProcessEvent(ctx, jsonBytes)
// The event is already parsed and callbacks are called automatically
// event is of type *base.BaseEvent
```

### Processing Multiple Events

To process multiple events at once (array of events):

```go
events, err := processor.ProcessEvents(ctx, jsonBytes)
// events is of type []*base.BaseEvent
// Callbacks are called for each event automatically
```

## Data Structure

### Package Structure

The library is organized into contextual packages for better modularity:

- **`pkg/types/base`**: Shared base types (`BaseEvent`, `EventCategory`, `EventSub`, etc.)
- **`pkg/types/common`**: Common types (`Location`, `Coordinates`, `GNSS`, `Connectivity`, `Fix`)
- **`pkg/types/order`**: Order events (`order.Event`)
- **`pkg/types/connection`**: Connection events (`connection.Event`)
- **`pkg/types/vision`**: Vision events (`vision.Event`)
- **`pkg/types/hardware`**: Hardware events (`hardware.Event`) - mapped from `EVENT_CATEGORY_HEALTH`
- **`pkg/types/system`**: System events (`system.Event`)
- **`pkg/webhook/processors`**: Modular processors per event context
- **`pkg/types/telemetry`**: Telemetry events (`telemetry.Event`)
- **`pkg/types/alert`**: Alert events (`alert.Event`)
- **`pkg/types/dms`**: DMS events (`dms.Event`)
- **`pkg/types/driverbehavior`**: Behavior events (`driverbehavior.Event`)
- **`pkg/types/vehicle`**: Vehicle events (`vehicle.Event`)

### Base Event
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

// Available helper methods:
// GetID() string
// GetCategory() base.EventCategory
// GetSubType() base.EventSub
// GetDeviceID() string
// GetCreatedAt() time.Time
```

### Specific Event Types

Each event category has its own type with helper methods:

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

## Event Type Verification

```go
import "go-eventlib/pkg/types/base"

// Check event category
if event.Category == base.EventCategoryOrder {
    // Handle order event
}

if event.Category == base.EventCategoryConnection {
    // Handle connection event
}

if event.Category == base.EventCategoryVision {
    // Handle vision event
}

// Access event information using helper methods
deviceID := event.GetDeviceID()
accountID := event.Attributes.Device.AccountID
```

## Agnostic Webhook SDK

The **Agnostic Webhook SDK** allows you to consume IoT events simply. **It does not include an HTTP server** - you create your own HTTP server and use the SDK only to parse events.

### Basic Usage with Callbacks

```go
package main

import (
    "context"
    "encoding/json"
        "io"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "go-eventlib/pkg/types/connection"
    "go-eventlib/pkg/types/order"
    "go-eventlib/pkg/webhook"
)

func main() {
    logger := log.New(os.Stdout, "[WEBHOOK] ", log.LstdFlags)

    eventHandler := webhook.NewEventHandler()
    eventHandler.OnOrderReceived = func(ctx context.Context, event *order.Event) error {
        if ord := event.GetOrder(); ord != nil {
            logger.Printf("Order received: ID=%s, Status=%s, Type=%s", event.ID, ord.Status, ord.Type)
        } else {
            logger.Printf("Order received: ID=%s, Status=%s", event.ID, event.Status)
        }
        return nil
    }
    eventHandler.OnOrderAck = func(ctx context.Context, event *order.Event) error {
        if ord := event.GetOrder(); ord != nil {
            logger.Printf("Order confirmed: ID=%s, Status=%s", event.ID, ord.Status)
        } else {
            logger.Printf("Order confirmed: ID=%s", event.ID)
        }
        return nil
    }

    connectionHandler := webhook.NewConnectionHandler()
    connectionHandler.OnWifiConnected = func(ctx context.Context, event *connection.Event) error {
        if wifi := event.GetWifiConnection(); wifi != nil {
            logger.Printf("WiFi connected: Device=%s, Status=%s", event.Attributes.Device.ID, wifi.Status)
        } else {
            logger.Printf("WiFi connected: Device=%s", event.Attributes.Device.ID)
        }
        return nil
    }
    connectionHandler.OnWifiDisconnected = func(ctx context.Context, event *connection.Event) error {
        logger.Printf("WiFi disconnected: Device=%s", event.Attributes.Device.ID)
        return nil
    }
    connectionHandler.OnSimCardChanged = func(ctx context.Context, event *connection.Event) error {
        if sim := event.GetSimCard(); sim != nil {
            logger.Printf("SIM card changed: Device=%s, Status=%s", event.Attributes.Device.ID, sim.Status)
        }
        return nil
    }
    connectionHandler.OnConnectionError = func(ctx context.Context, event *connection.Event) error {
        logger.Printf("Connection error: Device=%s", event.Attributes.Device.ID)
        return nil
    }

    processor := webhook.NewEventProcessorBuilder().
        WithEventHandler(eventHandler).
        WithConnectionHandler(connectionHandler).
        Build()

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        handleWebhook(w, r, processor, logger)
    })

    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    go func() {
        logger.Printf("Webhook server started on port 8080")
        logger.Printf("Endpoint: http://localhost:8080/")

        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    logger.Printf("Webhook server running. Press Ctrl+C to stop.")

    <-sigChan
    logger.Printf("Interrupt signal received, shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Error stopping server: %v", err)
    }

    logger.Printf("Server shut down successfully")
}

func handleWebhook(w http.ResponseWriter, r *http.Request, processor webhook.EventProcessor, logger *log.Logger) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        logger.Printf("Error reading body: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    ctx := context.Background()
    event, err := processor.ProcessEvent(ctx, body)
    if err != nil {
        logger.Printf("Error processing event: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    eventJSON, err := json.MarshalIndent(event, "", "  ")
    if err != nil {
        logger.Printf("Error serializing event: %v", err)
        logger.Printf("Event processed: ID=%s, Category=%s, Type=%s", event.ID, event.Category, event.Type)
    } else {
        logger.Printf("Event processed:\n%s", string(eventJSON))
    }

    w.WriteHeader(http.StatusAccepted)
}
```

### Available Handlers

The SDK provides predefined handlers with callbacks to facilitate event consumption:

#### EventHandler
```go
import "go-eventlib/pkg/types/order"

eventHandler := webhook.NewEventHandler()
eventHandler.OnOrderReceived = func(ctx context.Context, event *order.Event) error {
    if ord := event.GetOrder(); ord != nil {
        // Process received order with specific data
        log.Printf("Order: %s, Status: %s", ord.ID, ord.Status)
    }
    return nil
}
eventHandler.OnOrderAck = func(ctx context.Context, event *order.Event) error {
    // Process confirmed order
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithEventHandler(eventHandler).
    Build()
```

#### ConnectionHandler
```go
import "go-eventlib/pkg/types/connection"

connectionHandler := webhook.NewConnectionHandler()
connectionHandler.OnWifiConnected = func(ctx context.Context, event *connection.Event) error {
    if wifi := event.GetWifiConnection(); wifi != nil {
        log.Printf("WiFi connected: %s", wifi.Name)
    }
    return nil
}
connectionHandler.OnWifiDisconnected = func(ctx context.Context, event *connection.Event) error {
    // WiFi disconnected
    return nil
}
connectionHandler.OnSimCardChanged = func(ctx context.Context, event *connection.Event) error {
    if sim := event.GetSimCard(); sim != nil {
        log.Printf("SIM card: %s", sim.Status)
    }
    return nil
}
connectionHandler.OnConnectionError = func(ctx context.Context, event *connection.Event) error {
    // Connection error
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithConnectionHandler(connectionHandler).
    Build()
```

#### VisionHandler
```go
import "go-eventlib/pkg/types/vision"

visionHandler := webhook.NewVisionHandler()
visionHandler.OnFaceDetected = func(ctx context.Context, event *vision.Event) error {
    if faceData := event.GetFaceDetectedData(); faceData != nil {
        // Process detected face data
    }
    return nil
}
visionHandler.OnFaceLost = func(ctx context.Context, event *vision.Event) error {
    // Face lost
    return nil
}
visionHandler.OnFaceTracked = func(ctx context.Context, event *vision.Event) error {
    // Face tracked
    return nil
}
visionHandler.OnNoFaceDetected = func(ctx context.Context, event *vision.Event) error {
    // No face detected
    return nil
}
visionHandler.OnCameraObstructed = func(ctx context.Context, event *vision.Event) error {
    // Camera obstructed
    return nil
}
visionHandler.OnVisionAlert = func(ctx context.Context, event *vision.Event) error {
    // Generic vision alert
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithVisionHandler(visionHandler).
    Build()
```

#### HardwareHandler
```go
import "go-eventlib/pkg/types/hardware"

hardwareHandler := webhook.NewHardwareHandler()
hardwareHandler.OnDeviceRestart = func(ctx context.Context, event *hardware.Event) error {
    if sysData := event.GetSystemEventData(); sysData != nil {
        log.Printf("Device restarted: %s", sysData.EventName)
    }
    return nil
}
hardwareHandler.OnVehicleBatteryConnected = func(ctx context.Context, event *hardware.Event) error {
    // Vehicle battery connected
    return nil
}
hardwareHandler.OnVehicleBatteryDisconnected = func(ctx context.Context, event *hardware.Event) error {
    // Vehicle battery disconnected
    return nil
}
hardwareHandler.OnSDCardMounted = func(ctx context.Context, event *hardware.Event) error {
    // SD card mounted
    return nil
}
hardwareHandler.OnSDCardUnmounted = func(ctx context.Context, event *hardware.Event) error {
    // SD card unmounted
    return nil
}
hardwareHandler.OnSimCardInserted = func(ctx context.Context, event *hardware.Event) error {
    // SIM card inserted
    return nil
}
hardwareHandler.OnSimCardRemoved = func(ctx context.Context, event *hardware.Event) error {
    // SIM card removed
    return nil
}
hardwareHandler.OnHardwareAlert = func(ctx context.Context, event *hardware.Event) error {
    // Hardware alert
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithHardwareHandler(hardwareHandler).
    Build()
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
    // System alert
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithSystemHandler(systemHandler).
    Build()
```

#### TelemetryHandler
```go
import "go-eventlib/pkg/types/telemetry"

telemetryHandler := webhook.NewTelemetryHandler()
telemetryHandler.OnBatteryEvent = func(ctx context.Context, event *telemetry.Event) error {
    if metrics := event.GetBatteryMetrics(); metrics != nil {
        // Process battery metrics
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
    // Location event
    return nil
}
telemetryHandler.OnTelemetryEvent = func(ctx context.Context, event *telemetry.Event) error {
    // Any telemetry event
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithTelemetryHandler(telemetryHandler).
    Build()
```

#### AlertHandler
```go
import "go-eventlib/pkg/types/alert"

alertHandler := webhook.NewAlertHandler()
alertHandler.OnCriticalAlert = func(ctx context.Context, event *alert.Event) error {
    log.Printf("Critical alert: %s", event.GetAlertLevel())
    return nil
}
alertHandler.OnWarningAlert = func(ctx context.Context, event *alert.Event) error {
    // Warning alert
    return nil
}
alertHandler.OnInfoAlert = func(ctx context.Context, event *alert.Event) error {
    // Informational alert
    return nil
}
alertHandler.OnAnyAlert = func(ctx context.Context, event *alert.Event) error {
    if alertData := event.GetAlertEventData(); alertData != nil {
        log.Printf("Alert: %s", alertData.EventName)
    }
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithAlertHandler(alertHandler).
    Build()
```

#### DMSHandler
```go
import "go-eventlib/pkg/types/dms"

dmsHandler := webhook.NewDMSHandler()
dmsHandler.OnDrowsiness = func(ctx context.Context, event *dms.Event) error {
    if drowsiness := event.GetDrowsinessData(); drowsiness != nil {
        // Process drowsiness data
    }
    return nil
}
dmsHandler.OnDrinking = func(ctx context.Context, event *dms.Event) error {
    if drinking := event.GetDrinkingData(); drinking != nil {
        // Process drinking data
    }
    return nil
}
dmsHandler.OnEating = func(ctx context.Context, event *dms.Event) error {
    // Eating detected
    return nil
}
dmsHandler.OnEyeClosure = func(ctx context.Context, event *dms.Event) error {
    // Eye closure
    return nil
}
dmsHandler.OnGazeDistraction = func(ctx context.Context, event *dms.Event) error {
    // Visual distraction
    return nil
}
dmsHandler.OnGazeFixation = func(ctx context.Context, event *dms.Event) error {
    // Gaze fixation
    return nil
}
dmsHandler.OnPhone = func(ctx context.Context, event *dms.Event) error {
    // Phone usage
    return nil
}
dmsHandler.OnPoseDistraction = func(ctx context.Context, event *dms.Event) error {
    // Distracted posture
    return nil
}
dmsHandler.OnSmoking = func(ctx context.Context, event *dms.Event) error {
    // Smoking detected
    return nil
}
dmsHandler.OnYawning = func(ctx context.Context, event *dms.Event) error {
    // Yawning detected
    return nil
}
dmsHandler.OnDMSAlert = func(ctx context.Context, event *dms.Event) error {
    // Any DMS alert
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithDMSHandler(dmsHandler).
    Build()
```

#### DriverBehaviorHandler
```go
import "go-eventlib/pkg/types/driverbehavior"

driverBehaviorHandler := webhook.NewDriverBehaviorHandler()
driverBehaviorHandler.OnHarshAcceleration = func(ctx context.Context, event *driverbehavior.Event) error {
    if accel := event.GetHarshAccelerationData(); accel != nil {
        // Process harsh acceleration data
    }
    return nil
}
driverBehaviorHandler.OnHarshBraking = func(ctx context.Context, event *driverbehavior.Event) error {
    if braking := event.GetHarshBrakingData(); braking != nil {
        // Process harsh braking data
    }
    return nil
}
driverBehaviorHandler.OnMaxSpeedFault = func(ctx context.Context, event *driverbehavior.Event) error {
    // Maximum speed violation
    return nil
}
driverBehaviorHandler.OnNormalSpeedReturn = func(ctx context.Context, event *driverbehavior.Event) error {
    // Return to normal speed
    return nil
}
driverBehaviorHandler.OnPersistentMaxSpeed = func(ctx context.Context, event *driverbehavior.Event) error {
    // Persistent maximum speed
    return nil
}
driverBehaviorHandler.OnSharpTurn = func(ctx context.Context, event *driverbehavior.Event) error {
    // Sharp turn
    return nil
}
driverBehaviorHandler.OnStartOvertaking = func(ctx context.Context, event *driverbehavior.Event) error {
    // Start overtaking
    return nil
}
driverBehaviorHandler.OnDriverBehaviorAlert = func(ctx context.Context, event *driverbehavior.Event) error {
    // Any behavior alert
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithDriverBehaviorHandler(driverBehaviorHandler).
    Build()
```

#### VehicleHandler
```go
import "go-eventlib/pkg/types/vehicle"

vehicleHandler := webhook.NewVehicleHandler()
vehicleHandler.OnIgnitionOff = func(ctx context.Context, event *vehicle.Event) error {
    // Ignition off
    return nil
}
vehicleHandler.OnVehicleEvent = func(ctx context.Context, event *vehicle.Event) error {
    // Any vehicle event
    return nil
}

processor := webhook.NewEventProcessorBuilder().
    WithVehicleHandler(vehicleHandler).
    Build()
```

### Supported Event Types

The SDK supports automatic callback processing for the following event types:

#### Order Events
- `ORDER_STATUS_ACK` - Order confirmed
- `ORDER_STATUS_SENT` - Order sent
- `ORDER_STATUS_FAILED` - Order failed

#### Connection Events
- `WIFI_CONNECTED` - WiFi connected
- `WIFI_DISCONNECTED` - WiFi disconnected
- SIM card changes
- Critical connection errors

#### DMS Events (Driver Monitoring System)
- `DROWSINESS` - Drowsiness detected
- `GAZE_DISTRACTION` - Visual distraction
- `GAZE_FIXATION` - Gaze fixation
- `EYE_CLOSURE` - Eye closure
- `ON_PHONE` - Phone usage
- `EATING` - Eating
- `DRINKING` - Drinking
- `POSE_DISTRACTION` - Distracted posture

#### Driver Behavior Events
- `HARSH_ACCELERATION` - Harsh acceleration
- `HARSH_BRAKING` - Harsh braking
- `SPEED_VIOLATION` - Speed violation
- Persistent maximum speed

#### Hardware Events (EVENT_CATEGORY_HEALTH)
- `RESTART` / `REBOOT` / `R2_RESTART` - Device restart
- `SD_CARD_MOUNTED` / `SD_CARD_UNMOUNTED` - SD Card mounted/unmounted
- `SIMCARD_INSERTED` / `SIMCARD_REMOVED` / `SIMCARD_PRESENT` - SIM Card changes
- `VEHICLE_BATTERY_CONNECTED` / `VEHICLE_BATTERY_DISCONNECTED` - Vehicle battery

#### Telemetry Events
- `EVENT_SUB_TELEMETRY_IGNITION` - Ignition status
- `EVENT_SUB_TELEMETRY_BATTERY` - Battery status
- `EVENT_SUB_TELEMETRY_LOCATION` - GPS location

## Examples

Check out the examples in `examples/`:

- `webhook/main.go`: Complete webhook server example using net/http with all handlers

To run the example:

```bash
# Using Air (hot reload)
cd go-eventlib
air

# Or direct execution
cd examples/webhook
go run main.go

# Or build and run
go build -o ./tmp/webhook-example ./examples/webhook/main.go
./tmp/webhook-example
```

The server starts on port 8080. See `examples/webhook/README.md` for more details.

### HTTP Framework Integration

#### With net/http (default)
```go
processor := webhook.NewEventProcessor()
// Configure handlers...

http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    ctx := context.Background()
    event, err := processor.ProcessEvent(ctx, body)
    if err != nil {
        http.Error(w, "Error", 500)
        return
    }
    // Callbacks are called automatically
})
```

#### With Gin
```go
processor := webhook.NewEventProcessor()
// Configure handlers...

r := gin.Default()
r.POST("/webhook", func(c *gin.Context) {
    body, _ := c.GetRawData()
    event, err := processor.ProcessEvent(c.Request.Context(), body)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    // Callbacks are called automatically
})
```

#### With Echo
```go
processor := webhook.NewEventProcessor()
// Configure handlers...

e := echo.New()
e.POST("/webhook", func(c echo.Context) error {
    body, _ := io.ReadAll(c.Request().Body)
    event, err := processor.ProcessEvent(c.Request().Context(), body)
    if err != nil {
        return c.JSON(500, map[string]string{"error": err.Error()})
    }
    // Callbacks are called automatically
    return c.JSON(200, event)
})
```


### Test Payload Examples

#### Order Event (ACK)
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

#### Telemetry Event (Ignition)
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

#### DMS Event (Drowsiness)
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

### Example Test Files

The library includes 46+ real JSON event examples organized by category:

- **ack-events/** (3 files): Order confirmation and upload events
- **dms-events/** (11 files): Driver Monitoring System events
- **driver-behavior-events/** (7 files): Driver behavior events
- **hardware-events/** (14 files): Hardware and system events
- **telemetry-events/** (6 files): Telemetry and location events
- **vision-basic-events/** (5 files): Basic vision events

These files are available in `event-mapping-viewer/src/mindmap-sections/serialized-jsons/` and can be used for testing and development.

## Architecture

The library uses a modular architecture with specialized processors:

- **EventProcessor**: Main interface for event processing
- **EventProcessorBuilder**: Builder pattern for fluent configuration
- **Processors**: Modular processors per context (OrderProcessor, ConnectionProcessor, VisionProcessor, etc.)
- **Handlers**: Handlers with specific callbacks for each event type

Parsing is done using `protocol-cloud` with `protojson`, ensuring:
- Full compatibility with V3 protocol
- Protocol Buffers-based validation
- Support for enums and complex types from protocol-cloud
- Same MarshalOptions as event-handler-worker

## Dependencies

- `github.com/v3-tecnologia/protocol-cloud`: V3 event protocol
- `google.golang.org/protobuf/encoding/protojson`: JSON parsing for Protocol Buffers

## Testing

The library has **55.2%** test coverage:

```bash
# Run all tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic

# View detailed coverage
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

**Coverage by package:**
- `pkg/types/base`: 100.0%
- `pkg/types/order`: 100.0%
- `pkg/types/vehicle`: 100.0%
- `pkg/types/telemetry`: 92.3%
- `pkg/types/alert`: 88.9%
- `pkg/types/connection`: 88.0%
- `pkg/webhook`: 82.6%
- Other packages: 85-87%

## Contributing

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
