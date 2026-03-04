<div align="center">
  <a href="https://wristband.dev">
    <picture>
      <img src="https://assets.wristband.dev/images/email_branding_logo_v1.png" alt="Github" width="297" height="64">
    </picture>
  </a>
  <p align="center">
    Enterprise-ready auth that is secure by default, truly multi-tenant, and ungated for small businesses.
  </p>
  <p align="center">
    <b>
      <a href="https://wristband.dev">Website</a> •
      <a href="https://docs.wristband.dev">Documentation</a>
    </b>
  </p>
</div>

<br/>

---

# Wristband Multi-Tenant Demo App for Go

This demo app consists of:

- **Go Backend**: A Go backend with Wristband authentication integration
- **React Frontend**: A React frontend with authentication context

When an unauthenticated user attempts to access the frontend, it will redirect to the Go backend's Login Endpoint, which in turn redirects the user to Wristband to authenticate. Wristband then redirects the user back to your Go backend's Callback Endpoint which sets a session cookie before returning the user's browser to the React frontend.

<br>

---

<br>

## Requirements

This demo app requires the following prerequisites:

### Go

1. Visit [Go Downloads](https://go.dev/doc/install)
2. Download and install Go 1.24 or later
3. Verify the installation by opening a terminal or command prompt and running:
```bash
go version # Should show go1.24 or higher
```

### Node.js and NPM

1. Visit [NPM Downloads](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
2. Download and install the appropriate version for your OS
3. Verify the installation by opening a terminal or command prompt and running:
```bash
node --version # Should show v20.x.x or higher
npm --version  # Should show v9.6.x or higher
```

<br>

---

<br>

## Getting Started

You can start up the demo application in a few simple steps.

### 1) Sign up for a Wristband account

First, make sure you sign up for a Wristband account at [https://wristband.dev](https://wristband.dev).

### 2) Provision the Go demo application in the Wristband Dashboard

After your Wristband account is set up, log in to the Wristband dashboard. Once you land on the home page of the dashboard, click the "Add Application" button. Make sure you choose the following options:

- Step 1: Try a Demo
- Step 2: Subject to Authenticate - Humans
- Step 3: Application Framework - Go Backend, React Frontend

You can also follow the [Demo App Guide](https://docs.wristband.dev/docs/setting-up-a-demo-app) for more information.

### 3) Apply your Wristband configuration values

After completing demo app creation, you will be prompted with values that you should use to create environment variables for the Go server. You should see:

- `APPLICATION_VANITY_DOMAIN`
- `CLIENT_ID`
- `CLIENT_SECRET`

Copy those values, then create an environment variable file at the root of the project: `.env`. Once created, paste the copied values into this file.

### 4) Install Dependencies

**Go Dependencies:**
```bash
# From the go-demo-app directory
go mod download
```

**Frontend Dependencies:**
```bash
# Install React frontend dependencies
cd frontend
npm install
cd ..
```

### 5) Build the Frontend

```bash
cd frontend
npm run build
cd ..
```

This will write to the `dist/` directory used by `main.go` to serve the frontend.

### 6) Run the Application

```bash
# From the go-demo-app root directory
go run .
```

<br>

---

<br>

## How to interact with the demo app

The Go server starts on port 6001 and serves both the API endpoints and the built React frontend. During development, you can run the Vite dev server on port 5173, which is configured with a proxy to forward all `/api/*` requests to the Go backend at `http://localhost:6001/api/*`. This allows the frontend to make clean API calls using relative URLs like `/api/auth/session` while keeping the backend services separate and maintainable. The Go server also includes CORS middleware to allow cross-origin requests from the React frontend during development.

**Authentication Flow**

1. Visit `http://localhost:6001`
2. Click "Login" to initiate auth flow
3. Authenticate with Wristband
4. Return to application with session established
5. Access protected routes and API endpoints

### Home Page

The home page of the app can be accessed at `http://localhost:6001` (production) or `http://localhost:5173` (development). When the user is not authenticated, they will only see a Login button that will take them to the Application-level Login/Tenant Discovery page.

### Signup Users

You can sign up your first customer on the Signup Page at the following location:

- `https://{application_vanity_domain}/signup`, where `{application_vanity_domain}` should be replaced with the value of the "Application Vanity Domain" value of the application (found in the Wristband Dashboard).

This signup page is hosted by Wristband. Completing the signup form will provision both a new tenant with the specified tenant name and a new user that is assigned to that tenant.

### Application-level Login (Tenant Discovery)

Users of this app can access the Application-level Login Page at the following location:

- `https://{application_vanity_domain}/login`, where `{application_vanity_domain}` should be replaced with the value of the "Application Vanity Domain" value of the application.

This login page is hosted by Wristband. Here, the user will be prompted to enter either their email or their tenant's domain name, redirecting them to the Tenant-level Login Page for their specific tenant.

### Tenant-level Login

If users wish to directly access the Tenant-level Login Page without going through the Application-level Login Page, they can do so at:

- `http://localhost:6001/api/auth/login?tenant_name={tenant_name}`, where `{tenant_name}` should be replaced with the desired tenant's name.

This login page is hosted by Wristband. Here, the user will be prompted to enter their credentials to login to the application.

### Architecture

The application in this repository utilizes the Backend for Frontend (BFF) pattern, where Go is the backend for the React frontend. The server is responsible for:

- Storing the client ID and secret.
- Handling the OAuth2 authorization code flow redirections to and from Wristband during user login.
- Creating the application session cookie to be sent back to the browser upon successful login. The application session cookie contains the access and refresh tokens as well as some basic user info.
- Refreshing the access token if the access token is expired.
- Orchestrating all API calls from the frontend to Wristband.
- Destroying the application session cookie and revoking the refresh token when a user logs out.

API calls made from React to Go pass along the application session cookie with every request. The server has authentication middleware for all protected routes responsible for:

- Validating the session and refreshing the access token (if necessary)

Wristband hosts all onboarding workflow pages (signup, login, etc), and the Go server will redirect to Wristband in order to show users those pages.

<br>

---

<br>

## Development

For active development with hot-reloading:
```bash
# Terminal 1: Start Vite dev server
cd frontend
npm run dev

# Terminal 2: Start Go backend
go run .
```

Visit `http://localhost:5173` for development with hot module replacement.

<br>

---

## Wristband Go Auth SDK

This demo app is leveraging the [Wristband go-auth SDK](https://github.com/wristband-dev/go-auth) for all authentication interaction in the Go server. Refer to that GitHub repository for more information.

<br>

## Wristband React Client Auth SDK

This demo app is leveraging the [Wristband react-client-auth SDK](https://github.com/wristband-dev/react-client-auth) for any authenticated session interaction in the React frontend. Refer to that GitHub repository for more information.

<br/>

## Questions

Reach out to the Wristband team at <support@wristband.dev> for any questions regarding this demo app.

<br/>
