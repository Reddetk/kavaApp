# Analytics Service

The Analytics Service is a microservice that provides data analysis capabilities for the KavaApp platform. It uses various algorithms to analyze transaction data, generate product recommendations, and provide insights for business decisions.

## Features

- **Association Rule Mining**: Uses the Apriori algorithm to discover relationships between products in transaction data.
- **Product Recommendations**: Generates personalized product recommendations based on association rules.
- **ABC Analysis**: Categorizes products into A, B, and C segments based on their contribution to revenue.

## Architecture

The service follows a clean architecture pattern with the following layers:

- **Domain**: Contains business entities and repository interfaces.
- **Application**: Contains business logic and use cases.
- **Infrastructure**: Contains implementations of repositories and external services.
- **Interfaces**: Contains HTTP handlers and other interfaces to the outside world.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL database

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_USER=postgres
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=kavaapp_analytics
```

### Running the Service

```bash
go run cmd/main.go
```

## API Endpoints

- `GET /api/v1/recommendations?user_id=123`: Get product recommendations for a user.
- `POST /api/v1/analyze`: Analyze transaction data and generate insights.
- `GET /api/v1/abc-analysis`: Get ABC analysis results for products.

## Dependencies

- [go-apriori](https://github.com/eMAGTechLabs/go-apriori): Implementation of the Apriori algorithm for association rule mining.
- [logrus](https://github.com/sirupsen/logrus): Structured logger for Go.
- [godotenv](https://github.com/joho/godotenv): Load environment variables from .env files.
- [yaml.v3](https://gopkg.in/yaml.v3): YAML support for Go.

## License

This project is licensed under the MIT License - see the LICENSE file for details.