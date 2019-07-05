package uuid

import uuid "github.com/satori/go.uuid"

// Tool uuid contract interface
type Tool interface {
	CreateV4() string
}

// handler uuid object
type handler struct{}

// NewHandler : create uuid object handler
func NewHandler() Tool {
	return &handler{}
}

// Create : creating uuid
func (h *handler) CreateV4() string {
	return uuid.NewV4().String()
}
