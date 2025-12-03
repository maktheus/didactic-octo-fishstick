package queue

import (
	"context"
	"sync"
	"time"

	"github.com/example/back-end-tcc/pkg/logger"
	"github.com/example/back-end-tcc/pkg/observability/metrics"
)

// Message represents queue payload metadata and data.
type Message struct {
	Type string
	Data interface{}
}

// Handler is invoked when a message is consumed.
type Handler func(context.Context, Message) error

// Publisher publishes messages to the queue.
type Publisher interface {
	Publish(ctx context.Context, msg Message) error
}

// Subscriber subscribes handlers to message types.
type Subscriber interface {
	Subscribe(msgType string, handler Handler)
}

// Bus is a simple in-memory implementation of Publisher and Subscriber used for tests and local development.
type Bus struct {
	subscribers map[string][]Handler
	mu          sync.RWMutex
	log         logger.Logger
	metrics     metrics.Recorder
}

// Option configures Bus instrumentation.
type Option func(*Bus)

// WithLogger enables logging for the bus.
func WithLogger(l logger.Logger) Option {
	return func(b *Bus) {
		b.log = l
	}
}

// WithMetrics attaches a metrics recorder to the bus.
func WithMetrics(rec metrics.Recorder) Option {
	return func(b *Bus) {
		b.metrics = rec
	}
}

// NewBus creates a new in-memory bus.
func NewBus(opts ...Option) *Bus {
	bus := &Bus{subscribers: make(map[string][]Handler)}
	for _, opt := range opts {
		opt(bus)
	}
	if bus.log == nil {
		bus.log = logger.New()
	}
	return bus
}

// Publish notifies all subscribers of the message type.
func (b *Bus) Publish(ctx context.Context, msg Message) error {
	start := time.Now()
	b.mu.RLock()
	handlers := append([]Handler(nil), b.subscribers[msg.Type]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, msg); err != nil {
			b.recordPublish(msg.Type, len(handlers), time.Since(start), err)
			return err
		}
	}
	b.recordPublish(msg.Type, len(handlers), time.Since(start), nil)
	return nil
}

// Subscribe registers a handler for a message type.
func (b *Bus) Subscribe(msgType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[msgType] = append(b.subscribers[msgType], handler)
	if b.metrics != nil {
		b.metrics.AddCounter("queue_subscribers_total", map[string]string{"message": msgType}, 1)
	}
	if b.log != nil {
		b.log.Printf("queue: subscribed handler for %s", msgType)
	}
}

func (b *Bus) recordPublish(msgType string, handlers int, duration time.Duration, err error) {
	labels := map[string]string{"message": msgType}
	if err != nil {
		labels["result"] = "error"
	} else {
		labels["result"] = "ok"
	}
	if b.metrics != nil {
		b.metrics.AddCounter("queue_messages_total", labels, 1)
		b.metrics.ObserveHistogram("queue_publish_duration_ms", map[string]string{"message": msgType}, float64(duration.Milliseconds()))
	}
	if b.log != nil {
		if err != nil {
			b.log.Printf("queue: publish %s failed: %v", msgType, err)
		} else {
			b.log.Printf("queue: publish %s handled by %d subscribers in %s", msgType, handlers, duration)
		}
	}
}
