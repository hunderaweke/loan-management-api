
# Project Setup and Running Instructions

[Postman Documentation](https://documenter.getpostman.com/view/37328786/2sAXjGducb)

This document provides instructions for setting up and running the project. The project is structured into several directories for better organization, including API controllers, middleware, routers, and more.

## Project Structure

- **api**: Contains the API layer of the application.
  - `controllers/`: Defines the request handlers for various endpoints.
    - `loan_controllers.go`: Handles loan-related requests.
    - `log_controllers.go`: Handles system log-related requests.
    - `user_controllers.go`: Handles user-related requests.
  - `middlewares/`: Contains middleware components.
    - `admin_middleware.go`: Middleware for admin authentication and authorization.
    - `jwt_middleware.go`: Middleware for JWT token validation.
  - `routers/`: Defines the routing of API endpoints.
    - `loan_routers.go`: Routes for loan-related endpoints.
    - `main_router.go`: Main router that integrates all sub-routers.
    - `user_routers.go`: Routes for user-related endpoints.

- **cmd**: Contains the entry point for the application.
  - `main.go`: The main application entry point.

- **config**: Configuration files and loading mechanisms.
  - `config.go`: Configuration handling logic.
  - `config.yaml`: Configuration file.
  - `config.yaml.example`: Example configuration file.

- **database**: Database connection and setup.
  - `mongo.go`: MongoDB connection setup.

- **go.mod**: Go module definition file.

- **go.sum**: Go module checksum file.

- **internal**: Contains internal domain models, repositories, and use cases.
  - `domain/`: Defines domain models.
    - `loan.go`: Loan domain model.
    - `logs.go`: System logs domain model.
    - `user.go`: User domain model.
  - `repositories/`: Defines data access layer.
    - `loan_repository.go`: Loan repository implementation.
    - `log_repository.go`: System log repository implementation.
    - `user_repository.go`: User repository implementation.
  - `usecases/`: Contains business logic and use cases.
    - `loan_usecases.go`: Loan-related business logic.
    - `log_usecases.go`: System log-related business logic.
    - `user_usecases.go`: User-related business logic.

- **pkg**: External package utilities.
  - `infrastructures/`: Utility functions and helpers.
    - `jwt_token.go`: JWT token handling.
    - `password_handlers.go`: Password hashing and validation.
    - `send_email.go`: Email sending functions.
    - `token_handlers.go`: Token handling functions.

## Setup Instructions

### 1. Clone the Repository

First, clone the repository to your local machine:

```sh
git clone https://github.com/yourusername/your-repository.git
cd your-repository
```

### 2. Install Dependencies

Ensure you have Go installed on your machine. Install the required Go dependencies:

```sh
go mod download
```

### 3. Configuration

Copy the example configuration file to a new file and modify it according to your environment:

```sh
cp config.yaml.example config.yaml
```

Edit `config.yaml` to include your specific configuration details.

### 4. Database Setup

Ensure you have MongoDB installed and running. Update the MongoDB connection settings in `config.yaml` if necessary.

### 5. Run the Application

To start the application, use the following command:

```sh
go run cmd/main.go
```

### 6. Testing

Run tests to ensure everything is working correctly:

```sh
go test ./...
```

## Endpoints

### Loan Endpoints

- **Create Loan**: `POST /loans`
- **View Loan Status**: `GET /loans/{id}`
- **View All Loans**: `GET /admin/loans`
- **Approve/Reject Loan**: `PATCH /admin/loans/{id}/status`
- **Delete Loan**: `DELETE /admin/loans/{id}`

### User Endpoints

- **Register User**: `POST /users/register`
- **Verify Email**: `GET /users/verify`
- **Login**: `POST /users/login`
- **Get User Profile**: `GET /users/{id}`
- **Forget Password**: `POST /users/forget-password`
- **Reset Password**: `POST /users/reset-password`
- **Refresh Access Token**: `POST /users/refresh-token`

### Log Endpoints

- **View System Logs**: `GET /admin/logs`

## Troubleshooting

- Ensure all environment variables and configurations are set correctly.
- Verify MongoDB is running and accessible.
- Check for any errors in the application logs for debugging.

For further assistance, consult the project documentation or reach out to the maintainers.

---

This documentation should help you set up and run your project smoothly. Adjust paths and commands based on your specific setup and environment.
