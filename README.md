# go-demo-app
A Wristband multi-tenant auth demo app with a Go backend and a React frontend.

## Prerequisites

- Go 1.24.4 or later
- Node.js 20.0.0 or later
- Git

## Local Development Setup

This application depends on a local version of the `go-auth` library. Follow these steps to set up the development environment:

### 1. Clone Both Repositories

```bash
# Create a parent directory for both repositories
mkdir wristband-dev
cd wristband-dev

# Clone the go-auth library
git clone https://github.com/wristband-dev/go-auth.git

# Clone this demo application
git clone https://github.com/wristband-dev/go-demo-app.git
```

Your directory structure should look like:
```
wristband-dev/
├── go-auth/          # The go-auth library
└── go-demo-app/      # This demo application
```

### 2. Configure Environment Variables

Create a `.env` file in the `go-demo-app` directory:

```bash
cd go-demo-app
cp .env.example .env  # If available, or create manually
```

Add your Wristband configuration to `.env`:
```
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
APPLICATION_VANITY_DOMAIN=your_domain.wristband.dev
```

### 3. Install Dependencies

**Go Dependencies:**
```bash
# From the go-demo-app directory
go mod download
```

**Frontend Dependencies:**
```bash
# Install React frontend dependencies
cd clientapp
npm install
cd ..
```

### 4. Build the Frontend

```bash
cd clientapp
npm run build
cd ..
```

This will write to the `dist/` directory used by `main.go` to serve the frontend.

### 5. Run the Application

```bash
# From the go-demo-app root directory
go run main.go
```

The application will be available at `http://localhost:8080`.

## Project Structure

- `main.go` - Go backend server with embedded frontend
- `clientapp/` - React + TypeScript + Vite frontend application
- `clientapp/dist/` - Built frontend files (embedded in Go binary)
- `.env` - Environment configuration (not committed to git)

## Authentication Flow

1. Visit `http://localhost:8080`
2. Click "Login" to initiate OAuth flow
3. Authenticate with Wristband
4. Return to application with session established
5. Access protected routes and API endpoints
