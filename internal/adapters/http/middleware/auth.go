package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
	Data      map[string]interface{}
}

// SessionStore manages user sessions
type SessionStore struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

// AuthMiddleware provides session-based authentication functionality
type AuthMiddleware struct {
	store      *SessionStore
	cookieName string
	maxAge     time.Duration
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
		mutex:    sync.RWMutex{},
	}
}

// NewAuthMiddleware creates a new session-based authentication middleware
func NewAuthMiddleware(cookieName string, maxAge time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		store:      NewSessionStore(),
		cookieName: cookieName,
		maxAge:     maxAge,
	}
}

// CreateSession creates a new session for a user
func (a *AuthMiddleware) CreateSession(userID string, data map[string]interface{}) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(a.maxAge),
		Data:      data,
	}

	a.store.mutex.Lock()
	a.store.sessions[sessionID] = session
	a.store.mutex.Unlock()

	return session, nil
}

// GetSession retrieves a session by ID
func (a *AuthMiddleware) GetSession(sessionID string) (*Session, bool) {
	a.store.mutex.RLock()
	defer a.store.mutex.RUnlock()

	session, exists := a.store.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		delete(a.store.sessions, sessionID)
		return nil, false
	}

	return session, true
}

// DeleteSession removes a session
func (a *AuthMiddleware) DeleteSession(sessionID string) {
	a.store.mutex.Lock()
	defer a.store.mutex.Unlock()
	delete(a.store.sessions, sessionID)
}

// RequireAuth validates session cookies for protected routes
func (a *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Cookies(a.cookieName)
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusUnauthorized,
					"message": "Authentication required",
				},
			})
		}

		session, exists := a.GetSession(sessionID)
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusUnauthorized,
					"message": "Invalid or expired session",
				},
			})
		}

		// Store session and user info in context
		c.Locals("session", session)
		c.Locals("user_id", session.UserID)
		c.Locals("authenticated", true)

		return c.Next()
	}
}

// OptionalAuth validates session but doesn't require it
func (a *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Cookies(a.cookieName)
		if sessionID != "" {
			if session, exists := a.GetSession(sessionID); exists {
				c.Locals("session", session)
				c.Locals("user_id", session.UserID)
				c.Locals("authenticated", true)
			}
		}

		return c.Next()
	}
}

// RequireRole checks if user has required role
func (a *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Locals("session")
		if session == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusUnauthorized,
					"message": "Authentication required",
				},
			})
		}

		sessionData := session.(*Session)
		userRole, exists := sessionData.Data["role"]
		if !exists || userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    fiber.StatusForbidden,
					"message": "Insufficient permissions",
				},
			})
		}

		return c.Next()
	}
}

// Login creates a session and sets cookie
func (a *AuthMiddleware) Login(c fiber.Ctx, userID string, userData map[string]interface{}) error {
	session, err := a.CreateSession(userID, userData)
	if err != nil {
		return err
	}

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     a.cookieName,
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	return nil
}

// Logout removes session and clears cookie
func (a *AuthMiddleware) Logout(c fiber.Ctx) error {
	sessionID := c.Cookies(a.cookieName)
	if sessionID != "" {
		a.DeleteSession(sessionID)
	}

	// Clear session cookie
	c.Cookie(&fiber.Cookie{
		Name:     a.cookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	return nil
}

// CleanupExpiredSessions removes expired sessions (should be called periodically)
func (a *AuthMiddleware) CleanupExpiredSessions() {
	a.store.mutex.Lock()
	defer a.store.mutex.Unlock()

	now := time.Now()
	for id, session := range a.store.sessions {
		if now.After(session.ExpiresAt) {
			delete(a.store.sessions, id)
		}
	}
}

// generateSessionID creates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}