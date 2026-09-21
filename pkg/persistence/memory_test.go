package persistence_test

import (
	"context"
	"sync"
	"testing"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/persistence"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
)

func TestInMemoryRoundTrip(t *testing.T) {
	store := persistence.NewInMemoryStore()
	writer := persistence.NewInMemoryWriter(store)
	reader := persistence.NewInMemoryReader(store)
	ec := handlers.NewEventContext("a", false, 0, nil, 0, 0, 0, nil, nil, handlers.PayloadKindAlert, writer, reader)

	evt := &eventsv1.ImpactEvent{Name: "IMPACT"}
	if err := writer.Save(context.Background(), ec, evt); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded := &eventsv1.ImpactEvent{}
	ok, err := reader.GetEventByID(context.Background(), "a", loaded)
	if err != nil || !ok {
		t.Fatalf("get: ok=%v err=%v", ok, err)
	}

	wrong := &eventsv1.DrowsinessEvent{}
	ok, err = reader.GetEventByID(context.Background(), "a", wrong)
	if err != nil || ok {
		t.Fatalf("expected type mismatch skip, ok=%v err=%v", ok, err)
	}

	missing, err := reader.GetEventByID(context.Background(), "missing", &eventsv1.ImpactEvent{})
	if err != nil || missing {
		t.Fatalf("expected missing, ok=%v err=%v", missing, err)
	}

	list, err := reader.GetEvents(context.Background(), &eventsv1.ImpactEvent{}, -1)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v len=%d", err, len(list))
	}
}

func TestGetRootEvent(t *testing.T) {
	store := persistence.NewInMemoryStore()
	writer := persistence.NewInMemoryWriter(store)
	reader := persistence.NewInMemoryReader(store)

	root := &eventsv1.UploadEvent{Name: "UPLOAD"}
	rootEC := handlers.NewEventContext("root-1", false, 0, nil, 0, 0, 0, nil, nil, handlers.PayloadKindSystem, writer, reader)
	if err := writer.Save(context.Background(), rootEC, root); err != nil {
		t.Fatalf("save root: %v", err)
	}

	child := &eventsv1.UploadEvent{Name: "CHILD", SourceId: "root-1"}
	dest := &eventsv1.UploadEvent{}
	ok, err := reader.GetRootEvent(context.Background(), child, dest)
	if err != nil || !ok || dest.GetName() != "UPLOAD" {
		t.Fatalf("root: ok=%v err=%v name=%s", ok, err, dest.GetName())
	}

	noField := &eventsv1.ImpactEvent{Name: "X"}
	ok, err = reader.GetRootEvent(context.Background(), noField, &eventsv1.ImpactEvent{})
	if err != nil || ok {
		t.Fatalf("expected no source_id path, ok=%v err=%v", ok, err)
	}

	emptySource := &eventsv1.UploadEvent{Name: "CHILD"}
	ok, err = reader.GetRootEvent(context.Background(), emptySource, &eventsv1.UploadEvent{})
	if err != nil || ok {
		t.Fatalf("expected empty source_id path, ok=%v err=%v", ok, err)
	}

	missing := &eventsv1.UploadEvent{Name: "CHILD", SourceId: "nope"}
	_, err = reader.GetRootEvent(context.Background(), missing, &eventsv1.UploadEvent{})
	if err == nil {
		t.Fatal("expected missing root error")
	}
}

func TestInMemoryConcurrentSave(t *testing.T) {
	store := persistence.NewInMemoryStore()
	writer := persistence.NewInMemoryWriter(store)
	reader := persistence.NewInMemoryReader(store)

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := string(rune('a'+(i%26))) + string(rune('0'+i%10))
			ec := handlers.NewEventContext(id, false, 0, nil, 0, 0, 0, nil, nil, handlers.PayloadKindAlert, writer, reader)
			_ = writer.Save(context.Background(), ec, &eventsv1.ImpactEvent{Name: "IMPACT"})
			_, _ = reader.GetEventByID(context.Background(), id, &eventsv1.ImpactEvent{})
			_, _ = reader.GetEvents(context.Background(), &eventsv1.ImpactEvent{}, 5)
		}(i)
	}
	wg.Wait()
}
