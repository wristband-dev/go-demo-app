package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	goauth "github.com/wristband-dev/go-auth"
)

const (
	// SessionName is the name used for the session cookie
	SessionName = "session"

	// SessionKey is the key used to store auth data in the session
	SessionKey = "auth_session"
)

// GorillaSessionManager implements the goauth.SessionManager interface
// using gorilla/sessions for session management
type GorillaSessionManager struct {
	store sessions.Store
}

func NewSessionStore(secret []byte, secureCookies bool) goauth.SessionManager {
	store := sessions.NewCookieStore(secret, nil)
	store.Options.Secure = secureCookies
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.MaxAge = 3600 // 1 hour
	return &GorillaSessionManager{
		store: store,
	}
}

// StoreSession implements the goauth.SessionManager interface
func (m *GorillaSessionManager) StoreSession(w http.ResponseWriter, r *http.Request, session *goauth.Session) error {
	// Get existing session or create a new one
	sess, err := m.store.Get(r, SessionName)
	if err != nil {
		// If there's an error getting the session, create a new one
		// This can happen if the session was tampered with or is invalid
		sess, err = m.store.New(r, SessionName)
		if err != nil {
			return err
		}
	}

	// Serialize the session to JSON
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Store the serialized session in the session store
	sess.Values[SessionKey] = string(sessionJSON)

	// Save the session
	return sess.Save(r, w)
}

// GetSession implements the goauth.SessionManager interface
func (m *GorillaSessionManager) GetSession(r *http.Request) (*goauth.Session, error) {
	// Get existing session
	sess, err := m.store.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	// Check if the session contains auth data
	sessionJSON, ok := sess.Values[SessionKey]
	if !ok {
		return nil, errors.New("no auth session found")
	}

	// Parse the serialized session
	var authSession goauth.Session
	err = json.Unmarshal([]byte(sessionJSON.(string)), &authSession)
	if err != nil {
		return nil, err
	}

	return &authSession, nil
}

// ClearSession implements the goauth.SessionManager interface
func (m *GorillaSessionManager) ClearSession(w http.ResponseWriter, r *http.Request) error {
	// Get existing session
	sess, err := m.store.Get(r, SessionName)
	if err != nil {
		// If we can't get the session, that's fine - we wanted to clear it anyway
		return nil
	}

	// Remove the auth data from the session
	delete(sess.Values, SessionKey)

	// Set session to expire
	sess.Options.MaxAge = -1

	// Save the session
	return sess.Save(r, w)
}
