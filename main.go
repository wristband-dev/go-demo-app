package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	goauth "github.com/wristband-dev/go-auth"
)

type (
	UserRole struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	// Metadata is the structure used to store user session metadata in the Frontend SDK
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

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the Wristband Auth configuration
	authConfig := goauth.NewAuthConfig(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("APPLICATION_VANITY_DOMAIN"),
	)
	wristbandAuth, err := authConfig.WristbandAuth()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the session manager with a secure random secret key and secure cookies enabled
	sessionManager := NewSessionStore(securecookie.GenerateRandomKey(32), true)

	// Create the Wristband Auth app with the session manager
	app := wristbandAuth.NewApp(sessionManager)

	// Define middlewares to apply to protected routes
	authMiddleware := goauth.Middlewares{
		app.RequireAuthentication,
	}

	log.Println("Wristband configuration initialized successfully")
	log.Println("Starting server")
	apiMux := http.NewServeMux()

	sessionHandler := authMiddleware.Apply(app.SessionHandler(goauth.WithSessionMetadataExtractor(func(sess goauth.Session) any {
		return Metadata{
			Email:              sess.UserInfo.Email,
			FullName:           sess.UserInfo.Email,
			TenantName:         sess.TenantName,
			TenantCustomDomain: sess.TenantCustomDomain,
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

	protectedHandler := authMiddleware.Apply(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "This is a protected route", "value": 1}`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))

	// Define unauthenticated auth routes
	apiMux.Handle("/api/auth/login", app.LoginHandler())
	apiMux.Handle("/api/auth/callback", app.CallbackHandler())
	apiMux.Handle("/api/auth/logout", app.LogoutHandler())

	// Define protected auth routes
	apiMux.Handle("/api/auth/session", sessionHandler)
	apiMux.Handle("/api/auth/token", authMiddleware.Apply(app.TokenHandler()))

	// Define a protected API route
	apiMux.Handle("/api/protected", protectedHandler)

	// Apply CORS middleware
	mux := CORS(apiMux)

	// Get a sub-filesystem that starts at the dist directory
	distDir, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(distDir))

	// Start the server with CORS middleware
	const listenPort = ":6001"
	log.Printf("Listening on %s\n", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			mux.ServeHTTP(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})))
}
