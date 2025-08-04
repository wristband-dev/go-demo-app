package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/wristband-dev/go-auth"
)

type (
	UserRole struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	// Metadata is the structure used to store user metadata in the session
	Metadata struct {
		Email            string      `json:"email"`
		FullName         string      `json:"fullName"`
		TenantDomainName string      `json:"tenantDomainName"`
		Roles            []*UserRole `json:"roles"`
	}
)

//go:embed dist
var frontendFS embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wristbandClientID := os.Getenv("WRISTBAND_CLIENT_ID")
	wristbandClientSecret := os.Getenv("WRISTBAND_CLIENT_SECRET")
	tenantID := os.Getenv("WRISTBAND_TENANT_ID")
	var tenantDomains *goauth.TenantDomains
	if tenantID != "" {
		tenantDomains = &goauth.TenantDomains{
			TenantDomain: tenantID,
		}
	}
	auth, err := goauth.NewWristbandAuth(goauth.WristbandAuthConfig{
		Client: goauth.ConfidentialClient{
			ClientID:     wristbandClientID,
			ClientSecret: wristbandClientSecret,
		},
		Domains: goauth.AppDomains{
			RootDomain:      "localhost:8080",
			WristbandDomain: os.Getenv("WRISTBAND_DOMAIN"),
			DefaultDomains:  tenantDomains,
		},
	}, goauth.WithLogoutRedirectURL("/"))
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	app := goauth.NewApp(auth, goauth.AppInput{
		LoginPath:      "/api/auth/login",
		CallbackURL:    "http://localhost:8080/api/auth/callback",
		SessionManager: NewGorillaSessionManager(store),
		SessionMetadataExtractor: func(sess goauth.Session) any {
			return Metadata{
				Email:            sess.UserInfo.Email,
				FullName:         sess.UserInfo.Email,
				TenantDomainName: "global",
				Roles: []*UserRole{
					{
						Name:        "app:invotasticb2b:owner",
						ID:          "someId",
						DisplayName: "Owner",
					},
				},
			}
		},
	})

	apiMux := http.NewServeMux()

	var sessionHandler http.Handler = app.SessionHandler()
	var protectedHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "This is a protected route", "value": 1}`))
	})
	middlewares := []func(next http.Handler) http.Handler{
		app.RefreshTokenIfExpired,
		app.RequireAuthentication,
	}
	for _, middleware := range middlewares {
		sessionHandler = middleware(sessionHandler)
		protectedHandler = middleware(protectedHandler)
	}

	apiMux.Handle("/api/auth/login", app.LoginHandler())
	apiMux.Handle("/api/auth/callback", app.CallbackHandler())
	apiMux.Handle("/api/auth/logout", app.LogoutHandler())
	apiMux.Handle("/api/session", sessionHandler)
	apiMux.Handle("/api/protected", protectedHandler)

	// Get a sub-filesystem that starts at the dist directory
	distDir, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(distDir))

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nRequest received for: %q\n", r.URL.String())
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiMux.ServeHTTP(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})))
}

const (
	// SessionName is the name used for the session cookie
	SessionName = "wristband_auth_session"

	// SessionKey is the key used to store auth data in the session
	SessionKey = "auth_session"
)

// GorillaSessionManager implements the wristbandauth.SessionManager interface
// using gorilla/sessions for session management
type GorillaSessionManager struct {
	store sessions.Store
}

// NewGorillaSessionManager creates a new session manager using gorilla/sessions
func NewGorillaSessionManager(store sessions.Store) *GorillaSessionManager {
	return &GorillaSessionManager{
		store: store,
	}
}

// StoreSession implements the goauth.SessionManager interface
func (m *GorillaSessionManager) StoreSession(_ context.Context, w http.ResponseWriter, r *http.Request, session *goauth.Session) error {
	// Get existing session or create a new one
	sess, err := m.store.Get(r, SessionName)
	if err != nil {
		// If there's an error getting the session, create a new one
		// This can happen if the session was tampered with or is invalid
		sess = sessions.NewSession(m.store, SessionName)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30, // 30 days
			HttpOnly: true,
			Secure:   r.TLS != nil, // Set to true if connection is HTTPS
			SameSite: http.SameSiteLaxMode,
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
func (m *GorillaSessionManager) GetSession(_ context.Context, r *http.Request) (*goauth.Session, error) {
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
func (m *GorillaSessionManager) ClearSession(_ context.Context, w http.ResponseWriter, r *http.Request) error {
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
