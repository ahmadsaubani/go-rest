
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
│   ├── services/
│   │   └── auth_services/
│   │       └── auth_service.go
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

# Project Overview

This Go-based project is designed to provide an API with authentication and user management functionalities. The project is structured to include components for handling database connections, seeding data, authentication logic, middleware, and utilities like logging. Below is the description of the project's folder structure and its contents.

---

## Folder Structure

### **1. `config/`**
Contains configuration files for your application.

- **`database/`**
  - **`database.go`**: Handles the configuration and initialization of the database connection. It contains the database connection setup and the DB instance initialization.

---

### **2. `src/`**
The core source code of the application. This folder contains controllers, routes, entities, helpers, middleware, seeders, services, and utilities.

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

### **3. Root Files**

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

## Getting Started

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ahmadsaubani/go-rest.git



# _How To Run_

#### Development
```sh
Requirements:
- go > 1.20.x
- postgre
```

```sh
How to run :
- cp .env.example .env
- run gowatch or go main.go
```
