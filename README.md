# Go Auth API

This is a simple API built with Golang that uses JWT and OAuth2 for user authentication. It uses MongoDB as the database and provides endpoints for user registration, login, and OAuth2 authentication with Google. The API also includes role management functionalities for admins.

## Features

- User Registration: Register a new user by sending a POST request to `/register` with the user’s email and password.
- User Login: Authenticate users by sending a POST request to `/login` with credentials.
- User Logout: Logout users by sending a POST request to `/logout`.
- OAuth2 with Google:
  - Login: Redirect users to Google’s OAuth consent screen via `/oauth/google`.
  - Callback: Handle Google OAuth callback and obtain user information via `/oauth/google/callback`.
- User Role Management:
  - Get User: Fetch user details by sending a GET request to `/user`. Authentication is required.
  - Admin Access: Access admin-specific functionalities by sending a GET request to `/admin`. Admin privileges required.
  - Edit User Role: Update a user’s role by sending a PUT request to `/admin/edit-role/:user_id`. Admin privileges required.

## Endpoints

### User Registration

- POST `/register`
  - Request Body: { "username": "yourusername", "password": "yourpassword" }
  - Description: Register a new user with an username and password.

### User Login

- POST `/login`
  - Request Body: { "username": "yourusername", "password": "yourpassword" }
  - Description: Authenticate a user and receive a JWT token.

### User Logout

- POST `/logout`
  - Description: Invalidate the user’s session (token).

### OAuth2 with Google

- GET /oauth/google
  - Description: Redirect to Google’s OAuth consent screen.
- GET /oauth/google/callback
  - Query Parameters: code (The authorization code from Google)
  - Description: Handle the OAuth callback, retrieve user info, and provide a JWT token.

### User Endpoints

- GET `/user`
  - Description: Get the current authenticated user’s details. Requires authentication.
- GET `/admin`
  - Description: Access admin dashboard. Requires admin privileges.
- PUT `/admin/edit-role/:user_id`
  - URL Parameter: user_id (ID of the user whose role is to be updated)
  - Request Body: { "role": "admin" }
  - Description: Update a user’s role. Requires admin privileges.

### Environment Variables

The application uses the following environment variables:

- `GOOGLE_REDIRECT_URL`: The URL where Google will redirect after authentication.
- `GOOGLE_CLIENT_ID`: Google OAuth Client ID.
- `GOOGLE_CLIENT_SECRET`: Google OAuth Client Secret.
- `JWT_SECRE`T: Secret key for signing JWT tokens.
- `MONGODB_URI`: MongoDB connection string.
- `MONGODB_DB`: MongoDB database name.
- `ALLOW_ORIGINS`: Comma-separated list of allowed origins for CORS.

### Example .env File

```.env
MONGODB_URI=your_mongodb_uri
MONGODB_DB=your_db_name
JWT_SECRET=your_jwt_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=your_google_redirect_url
ALLOW_ORIGINS=your_allow_origins
```

## Setup

1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-auth-api.git
cd go-auth-api
```

2. Install Dependencies

```bash
go mod tidy
```

3. Create .env File
   Create a .env file in the root directory and add your environment variables.
4. Run the Application

```bash
go run main.go
```

5. Access the API
The API will be running on http://localhost:3000. You can interact with the endpoints using tools like Postman or curl.

## Enjoy the API! 🚀