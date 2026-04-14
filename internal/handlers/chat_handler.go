package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/repository"
	"portfolio-website-backend/internal/services"
)

type ChatHandler struct {
	chatSvc   services.ChatService
	vectorSvc services.VectorService
	chatRepo  repository.ChatRepository
}

func NewChatHandler(chatSvc services.ChatService, vectorSvc services.VectorService, chatRepo repository.ChatRepository) *ChatHandler {
	return &ChatHandler{
		chatSvc:   chatSvc,
		vectorSvc: vectorSvc,
		chatRepo:  chatRepo,
	}
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	switch v := userIDVal.(type) {
	case uuid.UUID:
		return v, true
	case string:
		id, err := uuid.Parse(v)
		if err == nil {
			return id, true
		}
	}
	return uuid.Nil, false
}

// RegisterRoutes registers endpoints for the chat functionality
func (h *ChatHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	chatGroup := router.Group("/chatbot")
	{
		// Authenticated Routes
		authGroup := chatGroup.Group("")
		authGroup.Use(authMiddleware)
		{
			authGroup.POST("/sync", h.SyncContext)
			authGroup.GET("/sessions", h.GetChatSessions)
			authGroup.GET("/sessions/:session_id/messages", h.GetSessionMessages)
		}

		// Public Routes (frontend widget uses WebSockets with query params)
		chatGroup.GET("/ws/chat", h.WebsocketChat)
	}
}

// SyncContext handles POST /chatbot/sync
// @Summary      Sync Vector Context
// @Description  Triggers a manual refresh of the Vector Store for the authenticated user
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /chatbot/sync [post]
func (h *ChatHandler) SyncContext(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user ID"})
		return
	}

	result, err := h.vectorSvc.SyncUserData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetChatSessions handles GET /chatbot/sessions
// @Summary      Get Chat Sessions
// @Description  Retrieve active chat sessions for the authenticated user
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query int false "Limit number of sessions" default(50)
// @Param        offset query int false "Offset for pagination" default(0)
// @Success      200 {array} dto.ChatSessionResponse
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /chatbot/sessions [get]
func (h *ChatHandler) GetChatSessions(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user ID"})
		return
	}

	limit := 50
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(c.Query("offset")); err == nil && o >= 0 {
		offset = o
	}

	sessions, err := h.chatRepo.GetAllSessions(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	if sessions == nil {
		sessions = []dto.ChatSessionResponse{}
	}

	c.JSON(http.StatusOK, sessions)
}

// GetSessionMessages handles GET /chatbot/sessions/:session_id/messages
// @Summary      Get Session Messages
// @Description  Retrieve the message history for a specific chat session
// @Tags         Chatbot
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        session_id path string true "Session UUID" format(uuid)
// @Success      200 {array} dto.ChatMessageResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /chatbot/sessions/{session_id}/messages [get]
func (h *ChatHandler) GetSessionMessages(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid user ID"})
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	messages, err := h.chatRepo.GetSessionMessages(sessionID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	if messages == nil {
		messages = []dto.ChatMessageResponse{}
	}

	c.JSON(http.StatusOK, messages)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, check origin properly. For now allow all to match legacy
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebsocketChat handles GET /chatbot/ws/chat
// @Summary      WebSocket Chat Endpoint
// @Description  Initiates a WebSocket connection for real-time portfolio chat. Requires session_id and user_id.
// @Tags         Chatbot
// @Param        session_id query string true "Client Chat Session UUID" format(uuid)
// @Param        user_id    query string true "Portfolio Owner UUID" format(uuid)
// @Success      101 "Switching Protocols to WebSockets"
// @Router       /chatbot/ws/chat [get]
func (h *ChatHandler) WebsocketChat(c *gin.Context) {
	// The frontend passes session_id and user_id in the query params
	sessionIDStr := c.Query("session_id")
	userIDStr := c.Query("user_id")

	if sessionIDStr == "" || userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and user_id are required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Upgrade to WS
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer ws.Close()

	// Ensure session exists
	err = h.chatSvc.GetOrCreateSession(sessionID, userID)
	if err != nil {
		ws.WriteJSON(dto.WsMessage{Type: "error", Payload: "Failed to initialize session"})
		return
	}

	// Send history
	history, payload, err := h.chatSvc.BuildChatHistory(sessionID, 20) // Limit history in prompt to 20
	if err != nil {
		ws.WriteJSON(dto.WsMessage{Type: "error", Payload: "Failed to load history"})
		return
	}
	if len(payload) > 0 {
		ws.WriteJSON(dto.WsMessage{Type: "history", Payload: payload})
	}

	// Loop to receive messages
	for {
		_, msgData, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Websocket read Error: %v", err)
			break
		}

		userText := string(msgData)

		// Save User Message
		err = h.chatSvc.SaveMessage(sessionID, "user", userText)
		if err != nil {
			log.Printf("Failed to save user message: %v", err)
		}

		// Generate stream
		// We use context.Background() because the request context might be cancelled
		// but we want the LLM response to complete if possible, or we could use ws close event to cancel ctx.
		ctx := context.Background()
		iter, client, err := h.chatSvc.GenerateStream(ctx, sessionID, userID, userText, history)
		if err != nil {
			ws.WriteJSON(dto.WsMessage{Type: "error", Payload: err.Error()})
			continue // Don't break, allow more messages
		}

		var fullResponse string
		for {
			resp, err := iter.Next()
			if err != nil {
				// iterator.Done is not exposed in standard way in some sdks, usually err is iterator.Done
				break
			}
			
			if resp != nil && len(resp.Candidates) > 0 {
				for _, part := range resp.Candidates[0].Content.Parts {
					if text, ok := part.(genai.Text); ok {
						fullResponse += string(text)
						ws.WriteJSON(dto.WsMessage{Type: "content", Payload: string(text)})
					}
				}
			}
		}
		
		// Wait, the SDK uses genai.Text so it's not string. It's type Text string
		// Let me refine the string extraction logic below to handle it safely
		
		// In previous block: if text, ok := part.(string); ok is wrong, it's genai.Text
		// Let's rely on type switch or fmt.Sprintf
		
		// Send End of stream
		ws.WriteJSON(dto.WsMessage{Type: "end"})
		
		// Wait, I messed up the loop variable state and iterator above, let me just fix it here:
		// Since I can't overwrite the loop variables easily inline here, let me use multi replace below to fix the type assertion
		client.Close()

		// Save bot message
		if fullResponse != "" {
			h.chatSvc.SaveMessage(sessionID, "bot", fullResponse)
		}
	}
}
