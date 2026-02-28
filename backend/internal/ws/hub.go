package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Jayyk09/CUHackIt/internal/agents"
	"github.com/Jayyk09/CUHackIt/internal/pantry"
	"github.com/Jayyk09/CUHackIt/internal/users"
	"github.com/Jayyk09/CUHackIt/pkg/logger"
	"github.com/Jayyk09/CUHackIt/services/gemini"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check against allowed origins
		return true
	},
}

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeConnect        MessageType = "connect"
	MessageTypeGenerate       MessageType = "generate"
	MessageTypeRecipeStart    MessageType = "recipe_start"
	MessageTypeRecipeProgress MessageType = "recipe_progress"
	MessageTypeRecipeComplete MessageType = "recipe_complete"
	MessageTypeError          MessageType = "error"
	MessageTypePing           MessageType = "ping"
	MessageTypePong           MessageType = "pong"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// GeneratePayload is the payload for generate requests
type GeneratePayload struct {
	UserID      string `json:"user_id"`
	Mode        string `json:"mode"`         // "pantry_only", "flexible", "both"
	RecipeCount int    `json:"recipe_count"` // 1-3
}

// RecipeStartPayload is sent when recipe generation starts
type RecipeStartPayload struct {
	TotalRecipes int    `json:"total_recipes"`
	Mode         string `json:"mode"`
}

// RecipeProgressPayload is sent for each generated recipe
type RecipeProgressPayload struct {
	RecipeIndex int            `json:"recipe_index"`
	TotalCount  int            `json:"total_count"`
	Recipe      agents.Recipe  `json:"recipe"`
}

// RecipeCompletePayload is sent when all recipes are generated
type RecipeCompletePayload struct {
	TotalGenerated int             `json:"total_generated"`
	FilteredCount  int             `json:"filtered_count"`
	Recipes        []agents.Recipe `json:"recipes"`
}

// ErrorPayload is sent on errors
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	UserID   string
	Conn     *websocket.Conn
	Send     chan []byte
	hub      *Hub
	mu       sync.Mutex
}

// Hub maintains the set of active clients
type Hub struct {
	clients      map[string]*Client
	register     chan *Client
	unregister   chan *Client
	orchestrator *agents.Orchestrator
	pantryRepo   *pantry.Repository
	userRepo     *users.Repository
	log          *logger.Logger
	mu           sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub(geminiClient *gemini.Client, pantryRepo *pantry.Repository, userRepo *users.Repository, log *logger.Logger) *Hub {
	var orchestrator *agents.Orchestrator
	if geminiClient != nil {
		orchestrator = agents.NewOrchestrator(geminiClient, log)
	}

	return &Hub{
		clients:      make(map[string]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		orchestrator: orchestrator,
		pantryRepo:   pantryRepo,
		userRepo:     userRepo,
		log:          log,
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
			h.log.Info("WebSocket client connected: %s", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			h.log.Info("WebSocket client disconnected: %s", client.ID)
		}
	}
}

// HandleWebSocket handles WebSocket connection upgrades
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade error: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
		hub:  h,
	}

	h.register <- client

	// Send connection confirmation
	client.sendMessage(MessageTypeConnect, map[string]string{
		"client_id": clientID,
		"status":    "connected",
	})

	// Start read/write goroutines
	go client.writePump()
	go client.readPump()
}

// sendMessage sends a typed message to the client
func (c *Client) sendMessage(msgType MessageType, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.Send <- msgBytes:
		return nil
	default:
		return nil // Channel full, message dropped
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(code, message string) {
	c.sendMessage(MessageTypeError, ErrorPayload{
		Code:    code,
		Message: message,
	})
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.log.Error("WebSocket read error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Flush queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages
func (c *Client) handleMessage(rawMessage []byte) {
	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		c.sendError("parse_error", "invalid message format")
		return
	}

	switch msg.Type {
	case MessageTypePing:
		c.sendMessage(MessageTypePong, nil)

	case MessageTypeGenerate:
		var payload GeneratePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			c.sendError("invalid_payload", "invalid generate payload")
			return
		}
		go c.handleGenerate(payload)

	default:
		c.sendError("unknown_type", "unknown message type")
	}
}

// handleGenerate handles recipe generation requests
func (c *Client) handleGenerate(payload GeneratePayload) {
	if c.hub.orchestrator == nil {
		c.sendError("service_unavailable", "recipe generation not available")
		return
	}

	userID := payload.UserID
	if userID == "" {
		c.sendError("invalid_user_id", "missing user ID")
		return
	}

	c.UserID = userID

	// Get user's pantry items
	pantryItems, err := c.hub.pantryRepo.ListByUserID(context.Background(), userID)
	if err != nil {
		c.sendError("pantry_error", "failed to get pantry items")
		return
	}

	if len(pantryItems) == 0 {
		c.sendError("empty_pantry", "pantry is empty - add some items first")
		return
	}

	// Get user preferences
	user, err := c.hub.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		c.sendError("user_error", "failed to get user preferences")
		return
	}

	// Convert pantry items to agent format
	agentPantryItems := make([]agents.PantryItem, len(pantryItems))
	for i, item := range pantryItems {
		category := ""
		if item.Category != nil {
			category = *item.Category
		}
		agentPantryItems[i] = agents.PantryItem{
			ID:       strconv.Itoa(item.ID),
			Name:     item.ProductName,
			Category: category,
			Quantity: float64(item.Quantity),
			Unit:     "item",
		}
	}

	// Determine mode
	var mode agents.OrchestratorMode
	switch payload.Mode {
	case "flexible":
		mode = agents.ModeFlexible
	case "both":
		mode = agents.ModeBoth
	default:
		mode = agents.ModePantryOnly
	}

	recipeCount := payload.RecipeCount
	if recipeCount <= 0 {
		recipeCount = 2
	}
	if recipeCount > 3 {
		recipeCount = 3
	}

	// Send start notification
	c.sendMessage(MessageTypeRecipeStart, RecipeStartPayload{
		TotalRecipes: recipeCount,
		Mode:         string(mode),
	})

	// Generate recipes
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	result, err := c.hub.orchestrator.Generate(ctx, agents.GenerateRequest{
		RecipeRequest: agents.RecipeRequest{
			PantryItems:        agentPantryItems,
			Allergens:          user.Allergens,
			DietaryPreferences: user.DietaryPreferences,
			NutritionalGoals:   user.NutritionalGoals,
			CookingSkill:       user.CookingSkill,
			CuisinePreferences: user.CuisinePreferences,
			RecipeCount:        recipeCount,
		},
		Mode: mode,
	})

	if err != nil {
		c.sendError("generation_error", err.Error())
		return
	}

	// Send each recipe as progress
	for i, recipe := range result.AllRecipes {
		c.sendMessage(MessageTypeRecipeProgress, RecipeProgressPayload{
			RecipeIndex: i + 1,
			TotalCount:  len(result.AllRecipes),
			Recipe:      recipe,
		})
		// Small delay between recipes for UX
		time.Sleep(100 * time.Millisecond)
	}

	// Send completion
	c.sendMessage(MessageTypeRecipeComplete, RecipeCompletePayload{
		TotalGenerated: result.TotalCount,
		FilteredCount:  result.FilteredCount,
		Recipes:        result.AllRecipes,
	})
}
