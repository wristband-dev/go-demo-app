package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	goauth "github.com/wristband-dev/go-auth"
)

type (
	UserRole struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	// Metadata is the structure used to store user metadata in the session
	Metadata struct {
		Email              string      `json:"email"`
		FullName           string      `json:"fullName"`
		TenantName         string      `json:"tenantName"`
		TenantCustomDomain string      `json:"tenantCustomDomain"`
		Now                string      `json:"now"`
		Roles              []*UserRole `json:"roles"`
	}
)

//go:embed dist
var frontendFS embed.FS

type Middlewares []func(next http.Handler) http.Handler

func (m Middlewares) Apply(next http.Handler) http.Handler {
	for _, middleware := range m {
		next = middleware(next)
	}
	return next
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	httpClient := &http.Client{}

	// Initialize the Wristband Auth configuration
	cfg := goauth.NewAuthConfig(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("APPLICATION_VANITY_DOMAIN"),
	)
	auth, err := cfg.WristbandAuth(goauth.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal(err)
	}

	app := auth.NewApp(NewGorillaSessionManager())
	log.Println("Wristband configuration initialized successfully")
	log.Println("Starting server")

	apiMux := http.NewServeMux()

	middlewares := Middlewares{
		app.RefreshTokenIfExpired,
		app.RequireAuthentication,
		goauth.CacheControlMiddleware,
	}
	sessionHandler := middlewares.Apply(app.SessionHandler(goauth.WithSessionMetadataExtractor(func(sess goauth.Session) any {
		return Metadata{
			Email:              sess.UserInfo.Email,
			FullName:           sess.UserInfo.Email,
			TenantName:         sess.TenantName,
			TenantCustomDomain: sess.CustomTenantDomain,
			Now:                time.Now().Format(time.RFC850),
			Roles: []*UserRole{
				{
					Name:        "app:invotasticb2b:owner",
					ID:          "someId",
					DisplayName: "Owner",
				},
			},
		}
	})))
	protectedHandler := middlewares.Apply(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "This is a protected route", "value": 1}`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

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

	const listenPort = ":6001"

	log.Printf("Listening on %s\n", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiMux.ServeHTTP(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})))
}

const (
	// SessionName is the name used for the session cookie
	SessionName = "session"

	// SessionKey is the key used to store auth data in the session
	SessionKey = "auth_session"
)

// GorillaSessionManager implements the wristbandauth.SessionManager interface
// using gorilla/sessions for session management
type GorillaSessionManager struct {
	store sessions.Store
}

// NewGorillaSessionManager creates a new session manager using gorilla/sessions
func NewGorillaSessionManager() *GorillaSessionManager {
	store := sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	store.Options.Secure = false // Make sure this is true in production
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
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
			Secure:   false, // Make sure this is true in production
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
