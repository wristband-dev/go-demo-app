module github.com/wristband-dev/go-demo-app

go 1.24.4

require github.com/wristband-dev/go-auth v0.0.0-unpublished

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace github.com/wristband-dev/go-auth v0.0.0-unpublished => ../go-auth
