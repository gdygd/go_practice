package sse

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

// EventData represents an SSE event structure
type EventData struct {
	Msgtype string `json:"event"`
	Data    string `json:"data"`
	Id      string `json:"id"`
}

// Client represents a connected SSE client
type Client struct {
	ID      string
	Channel chan EventData
}

// Hub manages all SSE connections and broadcasts
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan EventData
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan EventData),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Channel)
			}
			h.mu.Unlock()

		case event := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Channel <- event:
				default:
					// Skip if client channel is full
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a new client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends an event to all connected clients
func (h *Hub) Broadcast(event EventData) {
	h.broadcast <- event
}

// SendToClient sends an event to a specific client by ID
func (h *Hub) SendToClient(clientID string, event EventData) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, ok := h.clients[clientID]
	if !ok {
		return false
	}

	select {
	case client.Channel <- event:
		return true
	default:
		return false
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HasClient checks if a client exists
func (h *Hub) HasClient(clientID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[clientID]
	return ok
}

// formatSSE formats an EventData into SSE wire format
func formatSSE(event EventData) string {
	var result string

	if event.Id != "" {
		result += fmt.Sprintf("id: %s\n", event.Id)
	}
	if event.Msgtype != "" {
		result += fmt.Sprintf("event: %s\n", event.Msgtype)
	}
	result += fmt.Sprintf("data: %s\n\n", event.Data)

	return result
}

// ServeSSE handles SSE connection for Gin framework
func ServeSSE(hub *Hub, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Accel-Buffering", "no")

		client := &Client{
			ID:      clientID,
			Channel: make(chan EventData, 256),
		}

		hub.Register(client)
		defer hub.Unregister(client)

		// Send initial comment to establish connection
		_, _ = c.Writer.WriteString(": connected\n\n")
		c.Writer.Flush()

		ctx := c.Request.Context()

		for {
			select {
			case event, ok := <-client.Channel:
				if !ok {
					return
				}
				c.SSEvent(event.Msgtype, event.Data)
				c.Writer.Flush()
			case <-ctx.Done():
				return
			}
		}
	}
}

// ServeSSERaw handles SSE with raw format (id, event, data fields)
func ServeSSERaw(hub *Hub, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Accel-Buffering", "no")

		client := &Client{
			ID:      clientID,
			Channel: make(chan EventData, 256),
		}

		hub.Register(client)
		defer hub.Unregister(client)

		// Send initial comment to establish connection
		_, _ = c.Writer.WriteString(": connected\n\n")
		c.Writer.Flush()

		ctx := c.Request.Context()

		for {
			select {
			case event, ok := <-client.Channel:
				if !ok {
					return
				}
				_, _ = c.Writer.WriteString(formatSSE(event))
				c.Writer.Flush()
			case <-ctx.Done():
				return
			}
		}
	}
}

// NewEventData creates an EventData with JSON marshaled data
func NewEventData(msgType string, data any, id string) (EventData, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return EventData{}, err
	}
	return EventData{
		Msgtype: msgType,
		Data:    string(jsonData),
		Id:      id,
	}, nil
}

// NewEventDataString creates an EventData with string data
func NewEventDataString(msgType string, data string, id string) EventData {
	return EventData{
		Msgtype: msgType,
		Data:    data,
		Id:      id,
	}
}
