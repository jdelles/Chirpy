# Chirpy

A Twitter-like social media API built with Go, featuring user authentication, chirp (tweet) creation, and admin functionality.

## Features

- **User Management**: Create accounts, login, and update user profiles
- **Authentication**: JWT-based authentication with refresh tokens
- **Chirps**: Create, read, and delete short messages (max 140 characters)
- **Profanity Filter**: Automatic filtering of inappropriate content
- **Admin Panel**: Metrics tracking and development utilities
- **Database**: PostgreSQL with automated migrations
- **API Documentation**: RESTful endpoints with proper HTTP status codes

## Tech Stack

- **Language**: Go 1.24+
- **Database**: PostgreSQL
- **SQL Generator**: [sqlc](https://sqlc.dev/)
- **Migration Tool**: [Goose](https://github.com/pressly/goose)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Environment**: dotenv for configuration

## API Endpoints

### Health & Metrics
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - View server metrics (HTML)
- `POST /admin/reset` - Reset database (dev only)

### Authentication
- `POST /api/users` - Create new user account
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token
- `PUT /api/users` - Update user profile

### Chirps
- `POST /api/chirps` - Create new chirp (authenticated)
- `GET /api/chirps` - Get all chirps (supports sorting and filtering)
- `GET /api/chirps/{chirpID}` - Get specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp (owner only)

### Webhooks
- `POST /api/polka/webhooks` - Handle Polka payment webhooks

## Project Structure

```
├── main.go                 # Application entry point
├── .env                    # Environment variables
├── go.mod                  # Go module definition
├── sqlc.yaml              # SQL code generation config
├── handlers/              # HTTP request handlers
├── internal/
│   ├── auth/             # Authentication utilities
│   └── database/         # Generated database code
└── sql/
    ├── queries/          # SQL query definitions
    └── schema/           # Database migration files
```

## Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Chirpy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file with:
   ```env
   DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
   JWT_SECRET=your-secret-key
   POLKA_KEY=your-polka-api-key
   PLATFORM=dev
   ```

4. **Set up database**
   ```bash
   # Run migrations
   goose -dir sql/schema postgres $DB_URL up
   
   # Generate database code
   sqlc generate
   ```

5. **Run the server**
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## Development

### Database Changes

1. Create new migration file in `sql/schema`
2. Add queries in `sql/queries`
3. Run `sqlc generate` to update Go code
4. Apply migrations with `goose up`

### Testing

Run tests with:
```bash
go test ./...
```

### Code Generation

This project uses `sqlc` for type-safe database operations. After modifying SQL files, regenerate code:
```bash
sqlc generate
```

## Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Authentication**: Secure token-based auth
- **Input Validation**: Request parameter validation
- **SQL Injection Protection**: Parameterized queries
- **Profanity Filter**: Automatic content moderation

## Configuration

Environment variables:
- `DB_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT signing
- `POLKA_KEY`: API key for Polka webhook verification
- `PLATFORM`: Set to "dev" for development features

## License

This project is part of the Boot.dev Go course curriculum.