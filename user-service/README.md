# User Service Architecture

## Overview

User Service is a microservice designed to manage user profiles and update customer segments based on their behavior. It implements sophisticated analytics for user segmentation, retention prediction, and Customer Lifetime Value (CLV) calculation.

## Architectural Approach

The service follows Clean Architecture principles with these key layers:

- **Domain Layer**: Core business entities and business rules
- **Application Layer**: Use cases that orchestrate business logic
- **Infrastructure Layer**: Technical implementations of interfaces
- **Interfaces Layer**: External system communication (HTTP, Kafka)

## Project Structure

```
user-service/
├── cmd/
│   └── main.go                        # Application entry point
├── config/
│   └── config.yaml                    # Configuration file
├── internal/
│   ├── application/                   # Use Cases
│   │   ├── clv_service.go
│   │   ├── retention_service.go
│   │   ├── segmentation_service.go
│   │   └── user_service.go
│   ├── domain/                        # Business entities and logic
│   │   ├── entities/
│   │   │   ├── segment.go
│   │   │   ├── transaction.go
│   │   │   ├── user.go
│   │   │   └── user-metrics.go
│   │   ├── repositories/
│   │   │   ├── segment_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   ├── user_metrics_repository.go
│   │   │   └── user_repository.go
│   │   └── services/
│   │       ├── clv_service.go
│   │       ├── segmentation_service.go
│   │       ├── state_transition_service.go
│   │       └── survival_analysis_service.go
│   ├── infrastructure/
│   │   └── postgres/                  # Repository implementations
│   │       ├── transaction_repository.go
│   │       ├── segment_repository.go
│   │       ├── user_metrics_repository.go
│   │       └── user_repository.go
│   └── interfaces/                    # External interfaces
│       ├── http/
│       │   ├── handlers/
│       │   │   ├── segment_handler.go
│       │   │   └── user_handler.go
│       │   ├── consumer.go
│       │   └── router.go
├── pkg/                               # Helper packages
│   ├── logger/
│   └── utils/
├── test/                              # Tests
├── Dockerfile                         # Docker build instructions
├── go.mod                             # Go dependency manifest
└── README.md                          # Project documentation
```

## Core Components

### Domain Layer

The domain layer contains business entities and core business rules:

- **Entities**: User, Transaction, Segment, UserMetrics
- **Repositories**: Interfaces for data access
- **Services**: Business logic interfaces (segmentation, survival analysis, etc.)

### Application Layer

The application layer implements use cases:

- **UserService**: User profile management
- **SegmentationService**: User clustering and segment assignment
- **RetentionService**: Churn prediction and survival analysis
- **CLVService**: Customer Lifetime Value calculation

### Infrastructure Layer

The infrastructure layer provides concrete implementations:

- **Repository Implementations**: PostgreSQL-based data access
- **Algorithm Implementations**: KMeans clustering, Cox survival analysis

### Interfaces Layer

The interfaces layer handles external communications:

- **HTTP Handlers**: REST API endpoints
- **Kafka Consumers**: Event processing for transactions

## Key Features

- **User Segmentation**: RFM, behavioral, demographic, and promotional segmentation
- **Survival Analysis**: Predicting user churn and retention
- **State Transition Analysis**: Markov chains for user behavior modeling
- **CLV Calculation**: Dynamic lifetime value updates based on retention predictions

## API Endpoints

The service exposes RESTful endpoints:

- `/api/v1/users/`: User management
- `/api/v1/segments/`: Segment operations

## Event Processing

The service consumes transaction events via Kafka to:
- Update user metrics
- Trigger re-segmentation when necessary
- Recalculate retention probabilities and CLV

## Setup & Deployment

### Prerequisites
- Go 1.16+
- PostgreSQL 12+
- Kafka

### Configuration
Configure database connection, Kafka brokers, and other settings in `config/config.yaml`.

### Running the Service
```bash
# Build
go build -o user-service ./cmd/main.go

# Run
./user-service

# Using Docker
docker build -t user-service .
docker run -p 8080:8080 user-service
```

## Development Guidelines

1. Follow Clean Architecture principles
2. Keep business logic in domain and application layers
3. Infrastructure and interfaces should depend on inner layers, not vice versa
4. Use dependency injection for service composition

## Contribution

1. Fork this repository
2. Create a feature branch
3. Submit a pull request with detailed descriptions

## License

This project is licensed under the MIT License.
