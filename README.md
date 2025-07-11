
# 🛠️ Project Overview

This project is my first step into learning Go (Golang). I created this boilerplate as a starting point to help me build Go projects more easily in the future. It’s designed to be clean, reusable, and flexible so I don’t have to start from scratch every time I create a new project.

🔧 Key Features:

- **Flexible database support**: Easily switch between GORM and native SQL based on `.env` configuration.
- **Clean architecture**: Organized into controllers, services, repositories, and helpers for better maintainability and testability.
- **Dynamic filtering & pagination**: Supports powerful API query filtering (e.g., `?name[like]=john`) and paginated results.
- **PostgreSQL and MySQL support**: Fully compatible with both databases, including intelligent SQL placeholder formatting.
- **Migration-friendly**: Automatically generates `CREATE TABLE` statements and database triggers from Go structs.
- **Supports GORM and native SQL**: Choose your database interaction mode via the `USE_GORM` flag in `.env`.
- **Auto migration & auto seeding**: Automatically runs migrations and seeds based on registered entities.
- **CLI commands via Cobra**: Includes `migrate`, `seed`, and `dropdb` commands for streamlined database management.
- **Entity auto-registration**: Simply register your models once; they’re used for migration, seeding, and table dropping automatically.

🎯 Purpose & Vision:

This boilerplate is more than just a learning exercise—it's a tool I plan to evolve as I grow with Go. My goal is to write clean, maintainable Go code from the beginning. As I build more applications, I’ll continue refining this project to improve structure, developer experience, and performance.

---

## Getting Started
1. **Clone the repository**:
```bash
git clone https://github.com/ahmadsaubani/go-rest.git
```

2. **How To Run Development**:
```sh
Requirements:
- go > 1.20.x
- postgre
```

```sh
# Salin file .env default
cp .env.example .env

# Jalankan aplikasi dengan auto-reload (rekomendasi saat dev)
gowatch

# Atau manual
go run main.go

# Migration 
go run main.go migrate

# Seeder
go run main.go seed

# Drop Table
go run main.go dropdb

Semua entitas database didefinisikan dalam file:
/src/entities/registered.go


var RegisteredEntities = []any{
	users.User{},
	auth.AccessToken{},
	auth.RefreshToken{},
	// Tambahkan entitas lain di sini...
}


```

3. **List Endpoint**:
```sh
GET    /api/v1/ping             
POST   /api/v1/user/register    
POST   /api/v1/user/login       
GET    /api/v1/user/profile     
GET    /api/v1/users            
POST   /api/v1/token/refresh     
POST   /api/v1/user/logout       
```

4. **Filter Usage**:
```sh
Example :
1. /api/v1/users?email[like]=%john%&age[moreThan]=18&order_by=id,desc&page=1&per_page=10
```

#### STRUCTURE PROJECT
```sh
myapp/
├── src/
│   ├── config/
│   │  └── database/
│   │     └── database.go
│   ├── controllers/
│   │   └── api/
│   │       ├── v1/
│   │       │   ├── auth/
│   │       │   │   ├── login_controller.go
│   │       │   │   └── register_controller.go
│   │       │   └── user/
│   │       │       └── user_controller.go
│   ├── entities/
│   │   ├── auth/
│   │   │   ├── access_token.go
│   │   │   └── register_controller.go
│   │   ├── users/
│   │   │   └── user.go
│   ├── helpers/
│   │   ├── debug.go
│   │   └── response.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   ├── routes/
│   │   └── routes.go
│   ├── seeders/
│   │   └── user_seeders/
│   │       └── user_seeder.go
│   ├── repositories/
│   │   ├── auth_repositories/
│   │   │   └── auth_repository_interface.go
│   │   │   └── auth_repository.go
│   ├── services/
│   │   ├── auth_services/
│   │   │   └── auth_service_interface.go
│   │   │   └── auth_service.go
│   ├── utils/
│   │   └── loggers/
│   │       └── logger.go
│   └── storage/
│       └── logs/
├── .env
├── go.sum
├── go.mod
└── main.go


```
## Folder Structure

### **1. `src/`**
The core source code of the application. This folder contains controllers, routes, entities, helpers, middleware, seeders, services, and utilities.
- **`config/`**
  - **`database/`**
    - **`database.go`**: Handles the configuration and initialization of the database connection. It contains the database connection setup and the DB instance initialization.
- **`controllers/`**
  - **`api/`**: Contains the API controllers for handling incoming requests.
    - **`v1/`**: Version 1 of the API, organizing the controllers into subfolders for specific functionality.
      - **`auth/`**: Authentication-related controllers.
        - **`login_controller.go`**: Handles the login logic and authentication requests.
        - **`register_controller.go`**: Manages user registration logic and new user account creation.
      - **`user/`**: User-related controllers.
        - **`user_controller.go`**: Handles user-related functionalities, such as fetching or updating user data.

- **`entities/`**
  - **`auth/`**: Auth-related entities or models.
    - **`access_token.go`**: Defines the model and logic for managing access tokens used in the app.
    - **`refresh_token.go`**: Defines the model and logic for managing access tokens used in the app.
  - **`users/`**: User-related entities.
    - **`user.go`**: Defines the user model and handles related database operations.

- **`helpers/`**
  - **`debug.go`**: Contains utility functions for debugging, such as print-based tools or structured debugging helpers.
  - **`response.go`**: Provides helper functions to format and send standardized responses (success/failure).

- **`middleware/`**
  - **`auth_middleware.go`**: Handles authentication checks, including JWT token verification, to protect routes requiring authenticated access.

- **`routes/`**
  - **`routes.go`**: Manages routing of HTTP requests, maps controllers to routes, and defines the main HTTP request handling logic for the application.

- **`seeders/`**
  - **`user_seeders/`**
    - **`user_seeder.go`**: A seeder file to populate the database with initial or test user data, useful for development or testing.

- **`services/`**
  - **`auth_services/`**
    - **`auth_service.go`**: Implements the business logic for authentication, including token generation and refreshing tokens.

- **`utils/`**
  - **`loggers/`**
    - **`logger.go`**: Implements centralized logging functionality for the application, which might include logging levels, output formats, and log file storage.

- **`storage/`**
  - **`logs/`**: Stores log files generated during application runtime for debugging or monitoring purposes.

---

### **2. Root Files**

- **`.env`**: Stores environment variables used by the application, such as database credentials, API keys, and the JWT secret.
- **`go.mod`**: The Go module file that defines the dependencies for the project.
- **`go.sum`**: Ensures the integrity of the dependencies in `go.mod`.
- **`main.go`**: The entry point of the application, which initializes the application and starts the HTTP server.

---

## Summary

The project follows a modular structure where:

- **Controllers** are responsible for handling HTTP requests and routing logic.
- **Entities** define the data models used throughout the application.
- **Helpers** provide utility functions for common tasks such as formatting responses.
- **Middleware** enforces authentication and authorization checks before accessing certain routes.
- **Seeders** populate the database with initial data for testing or development.
- **Services** implement the core business logic, particularly for authentication.
- **Utilities** provide lower-level functionalities like logging.
- **Storage** holds files such as logs that are generated by the application during runtime.

---