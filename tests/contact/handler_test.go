package contact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"portfolio-website-backend/internal/dto"
	"portfolio-website-backend/internal/handlers"
)

// setupContactRouter creates a test router wiring the ContactHandler.
// The contact route is public, so no auth middleware parameter is needed.
func setupContactRouter(svc *MockContactService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handlers.NewContactHandler(svc)
	h.RegisterRoutes(r.Group("/"))
	return r
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — happy path
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	userID := uuid.New()
	req := dto.ContactRequest{
		Name:    "Alice Smith",
		Email:   "alice@example.com",
		Subject: "Hello",
		Message: "Great portfolio!",
	}

	mockSvc.On("SendContactEmail", userID, req).Return(nil)

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/"+userID.String(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Emails sent successfully", resp["message"])

	mockSvc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — invalid UUID in path
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_InvalidUserID(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	body, _ := json.Marshal(dto.ContactRequest{
		Name:    "Bob",
		Email:   "bob@example.com",
		Subject: "Hi",
		Message: "Nice work.",
	})
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/not-a-uuid", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Invalid user ID", resp["detail"])

	// No service call should happen
	mockSvc.AssertNotCalled(t, "SendContactEmail", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — missing required fields (binding validation)
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_MissingFields(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	userID := uuid.New()

	// Only name is provided — email, subject, message are missing
	body, _ := json.Marshal(map[string]string{"name": "Charlie"})
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/"+userID.String(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "SendContactEmail", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — invalid email format
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_InvalidEmail(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	userID := uuid.New()
	body, _ := json.Marshal(map[string]string{
		"name":    "Dave",
		"email":   "not-an-email",
		"subject": "Test",
		"message": "Hello",
	})
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/"+userID.String(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "SendContactEmail", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — user not found (404 from service)
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_UserNotFound(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	userID := uuid.New()
	req := dto.ContactRequest{
		Name:    "Eve",
		Email:   "eve@example.com",
		Subject: "Query",
		Message: "Looking good.",
	}

	mockSvc.On("SendContactEmail", userID, req).Return(
		// The service returns this sentinel message for missing users
		fmt.Errorf("user not found"),
	)

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/"+userID.String(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// POST /contact/:user_id — Mailgun error → 500
// ---------------------------------------------------------------------------

func TestContactHandler_SendEmail_MailgunError(t *testing.T) {
	mockSvc := new(MockContactService)
	r := setupContactRouter(mockSvc)

	userID := uuid.New()
	req := dto.ContactRequest{
		Name:    "Frank",
		Email:   "frank@example.com",
		Subject: "Issue",
		Message: "There might be an issue.",
	}

	mockSvc.On("SendContactEmail", userID, req).Return(
		fmt.Errorf("failed to send confirmation email: mailgun: unexpected status 429"),
	)

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, "/contact/"+userID.String(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
