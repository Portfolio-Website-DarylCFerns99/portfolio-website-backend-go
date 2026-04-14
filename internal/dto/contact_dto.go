package dto

// ContactRequest represents the JSON body for the contact form submission.
type ContactRequest struct {
	Name    string `json:"name"    binding:"required"`
	Email   string `json:"email"   binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}
