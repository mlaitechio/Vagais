package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// ChatHandler handles real-time chat functionality
type ChatHandler struct {
	*BaseHandler
	agentService *services.AgentService
	upgrader     websocket.Upgrader
	clients      map[string]*websocket.Conn
}

// NewChatHandler creates a new chat handler
func NewChatHandler(db *gorm.DB, cfg *config.Config) *ChatHandler {
	return &ChatHandler{
		BaseHandler:  NewBaseHandler(db, cfg),
		agentService: services.AgentServiceInstance,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients: make(map[string]*websocket.Conn),
	}
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Type      string                 `json:"type"`
	AgentID   string                 `json:"agent_id"`
	UserID    string                 `json:"user_id"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ChatResponse represents a chat response from agent
type ChatResponse struct {
	Type      string                 `json:"type"`
	AgentID   string                 `json:"agent_id"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// HandleChatWebSocket handles WebSocket connections for real-time chat
func (h *ChatHandler) HandleChatWebSocket(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		h.sendError(c, http.StatusBadRequest, "Agent ID is required")
		return
	}

	// Get user from context (assuming auth middleware has set it)
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Store client connection
	clientKey := userID + ":" + agentID
	h.clients[clientKey] = conn

	// Send welcome message
	welcomeMsg := ChatResponse{
		Type:      "welcome",
		AgentID:   agentID,
		Message:   "Connected to agent chat. You can start chatting now!",
		Timestamp: time.Now(),
	}
	conn.WriteJSON(welcomeMsg)

	// Handle incoming messages
	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(h.clients, clientKey)
			break
		}

		// Process the message
		h.processChatMessage(conn, agentID, userID, &msg)
	}
}

// processChatMessage processes incoming chat messages
func (h *ChatHandler) processChatMessage(conn *websocket.Conn, agentID, userID string, msg *ChatMessage) {
	// Set message metadata
	msg.AgentID = agentID
	msg.UserID = userID
	msg.Timestamp = time.Now()

	// Send acknowledgment
	ack := ChatResponse{
		Type:      "ack",
		AgentID:   agentID,
		Message:   "Message received",
		Timestamp: time.Now(),
	}
	conn.WriteJSON(ack)

	// Process message with agent
	go h.processWithAgent(conn, agentID, userID, msg)
}

// processWithAgent processes the message with the agent
func (h *ChatHandler) processWithAgent(conn *websocket.Conn, agentID, userID string, msg *ChatMessage) {
	// Simulate agent processing time
	time.Sleep(500 * time.Millisecond)

	// Get agent details
	agent, err := h.agentService.GetAgent(agentID)
	if err != nil {
		errorResp := ChatResponse{
			Type:      "error",
			AgentID:   agentID,
			Message:   "Agent not found or unavailable",
			Timestamp: time.Now(),
		}
		conn.WriteJSON(errorResp)
		return
	}

	// Simulate agent response based on message content
	response := h.generateAgentResponse(agent.Name, msg.Message)

	// Send agent response
	agentResp := ChatResponse{
		Type:      "response",
		AgentID:   agentID,
		Message:   response,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"agent_name":      agent.Name,
			"processing_time": "500ms",
		},
	}

	conn.WriteJSON(agentResp)
}

// generateAgentResponse generates a simple response based on the agent and message
func (h *ChatHandler) generateAgentResponse(agentName, message string) string {
	// This is a simple response generator
	// In a real implementation, this would call the actual agent service
	responses := map[string]string{
		"hello":  "Hello! I'm " + agentName + ". How can I help you today?",
		"help":   "I'm here to assist you. What would you like to know?",
		"thanks": "You're welcome! Is there anything else I can help you with?",
		"bye":    "Goodbye! Have a great day!",
	}

	// Check for exact matches first
	if response, exists := responses[message]; exists {
		return response
	}

	// Default response
	return "I understand you said: \"" + message + "\". I'm " + agentName + " and I'm here to help. Could you please provide more details about what you need?"
}

// BroadcastMessage broadcasts a message to all connected clients for an agent
func (h *ChatHandler) BroadcastMessage(agentID string, message ChatResponse) {
	for clientKey, conn := range h.clients {
		// Check if this client is connected to the specific agent
		if len(clientKey) > len(agentID) && clientKey[len(clientKey)-len(agentID)-1:] == ":"+agentID {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error broadcasting to client %s: %v", clientKey, err)
				delete(h.clients, clientKey)
			}
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (h *ChatHandler) GetConnectedClients() int {
	return len(h.clients)
}

// CloseAllConnections closes all WebSocket connections
func (h *ChatHandler) CloseAllConnections() {
	for clientKey, conn := range h.clients {
		conn.Close()
		delete(h.clients, clientKey)
	}
}
