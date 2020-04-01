package ingress

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client/test"
	"github.com/google/go-cmp/cmp"
)

// TestHandlerChannelInbound uses a channel based mock client for both inbound and decouple.
func TestHandlerChannelInbound(t *testing.T) {
	inbound, ic := test.NewMockReceiverClient(t, 1)
	decouple, dc := test.NewMockSenderClient(t, 1)
	_, cleanup := createAndStartIngress(t, inbound, decouple)
	defer cleanup()

	input := createTestEvent()
	// Send an event to the inbound receiver client.
	ic <- input
	// Retrieve the event from the decouple sink.
	output := <-dc

	if dif := cmp.Diff(input, output); dif != "" {
		t.Errorf("Output event doesn't match input, dif: %v", dif)
	}
}

// TestHandlerDefaultHTTPInbound uses default HTTP inbound client and a channel based mock decouple
// sink.
func TestHandlerDefaultHTTPInbound(t *testing.T) {
	decouple, dc := test.NewMockSenderClient(t, 1)
	_, cleanup := createAndStartIngress(t, nil, decouple)
	defer cleanup()

	input := createTestEvent()
	// Send an event to the inbound receiver client.
	sendEventOverHTTP(t, "http://localhost:8080", input)
	// Retrieve the event from the decouple sink.
	output := <-dc

	if dif := cmp.Diff(input, output); dif != "" {
		t.Errorf("Output event doesn't match input, dif: %v", dif)
	}
}

func sendEventOverHTTP(t *testing.T, url string, events ...cloudevents.Event) {
	p, err := cloudevents.NewHTTP(cloudevents.WithTarget(url))
	if err != nil {
		t.Fatalf("Failed to create HTTP protocol: %+v", err)
	}
	ce, err := cloudevents.NewClient(p)
	if err != nil {
		t.Fatalf("Failed to create HTTP client: %+v", err)
	}
	for _, event := range events {
		ce.Send(context.Background(), event)
	}
}

// createAndStartIngress creates an ingress and calls its Start() method in a goroutine.
func createAndStartIngress(t *testing.T, inbound cloudevents.Client, decouple DecoupleSink) (h *handler, cleanup func()) {
	ctx := context.Background()
	h, err := NewHandler(ctx,
		WithInboundClient(inbound),
		WithDecoupleSink(decouple))
	if err != nil {
		t.Fatalf("Failed to create ingress handler: %+v", err)
	}
	go h.Start(ctx)
	cleanup = func() {
		// Any cleanup steps should go here. For now none.
	}
	return h, cleanup
}

func createTestEvent() cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID("test-id")
	event.SetSource("test-source")
	event.SetType("test-type")
	return event
}