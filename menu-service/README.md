# Menu Service

## Overview
The Menu Service is a microservice for dynamic personalized menu generation within the Kava application ecosystem. It provides RESTful APIs for managing menu items, categories, products, and promotions, with support for personalization based on user segments and geolocation.

## Purpose
The primary goal of this service is to enhance the customer experience by delivering personalized menu options and promotions that are relevant to each user's preferences, location, and behavior patterns.

## Features
- **Personalized Menu Generation**: Creates customized menus based on user segments
- **Dynamic Pricing**: Applies intelligent pricing strategies based on demand elasticity
- **Geo-targeted Promotions**: Delivers location-specific offers and recommendations
- **Price History Tracking**: Maintains historical record of price changes
- **Caching Strategy**: Implements Redis-based caching for performance optimization
- **API Documentation**: Provides OpenAPI documentation for seamless integration

## Architecture

The service follows a layered architecture:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Implement business logic
- **Repositories**: Manage data access
- **Models**: Represent domain entities

### Technology Stack

- **Framework**: Spring Boot 2.7.x
- **Database**: PostgreSQL
- **Cache**: Redis
- **Documentation**: SpringDoc OpenAPI
- **Build Tool**: Gradle

### Algorithms and Methods

- **Elasticity of Demand**: For price adjustments based on discount response
- **Redis Caching**: For performance optimization
- **Geo-filtering**: Location-based filtering using Haversine formula
- **OpenAPI Generation**: Using springdoc-openapi for API documentation

## Project Structure
```
menu-service/
│
├── src/
│   └── main/
│       ├── java/com/kava/menu/
│       │   ├── controller/       # REST API endpoints
│       │   ├── model/            # Domain models
│       │   ├── repository/       # Data access interfaces
│       │   ├── service/          # Business logic
│       │   ├── dto/              # Data transfer objects
│       │   ├── config/           # Application configuration
│       │   └── MenuServiceApplication.java
│       │
│       └── resources/
│           └── application.yaml  # Configuration (DB, Redis, cache)
│
├── build.gradle                  # Dependency management
├── gradlew                       # Gradle wrapper script
└── README.md                     # This file
```

## Core Components

### Models
- **Category**: Organizes products into logical groups
- **Product**: Represents menu items with pricing information
- **PriceHistory**: Tracks product price changes over time
- **Promotion**: Defines discount offers and their validity periods
- **GeoPromotion**: Stores location data for geo-targeting
- **Segment**: Defines customer segments for personalization
- **PersonalizedMenu**: Represents a menu tailored to a specific segment
- **PersonalizedMenuItem**: Items within a personalized menu

### Services
- **CategoryService**: Business logic for category management
- **ProductService**: Business logic for product management
- **PricingService**: Handles price calculations and dynamic pricing
- **PersonalizedMenuService**: Generates personalized menus based on user segments
- **PromotionService**: Manages promotional offers
- **GeoPromotionService**: Handles location-based promotions

### Controllers
- **CategoryController**: Manages product categories
- **MenuController**: Handles personalized menu generation and retrieval
- **ProductController**: Manages product catalog
- **PromotionController**: Handles promotional offers and discounts
- **GeoPromotionController**: Manages location-based promotions

## API Endpoints

### Categories
```
GET    /api/categories
GET    /api/categories/{id}
POST   /api/categories
PUT    /api/categories/{id}
DELETE /api/categories/{id}
```

### Products
```
GET    /api/products
GET    /api/products/active
GET    /api/products/{id}
GET    /api/products/category/{categoryId}
POST   /api/products
PUT    /api/products/{id}
PUT    /api/products/{id}/price
PUT    /api/products/{id}/activate
PUT    /api/products/{id}/deactivate
DELETE /api/products/{id}
```

### Menus
```
GET    /api/menus
GET    /api/menus/{id}
GET    /api/menus/segment/{segmentId}
GET    /api/menus/segment/{segmentId}/latest
POST   /api/menus/generate
DELETE /api/menus/{id}
```

### Promotions
```
GET    /api/promotions
GET    /api/promotions/active
GET    /api/promotions/{id}
POST   /api/promotions
PUT    /api/promotions/{id}
PUT    /api/promotions/{id}/activate
PUT    /api/promotions/{id}/deactivate
DELETE /api/promotions/{id}
```

### Geo-Promotions
```
GET    /api/geo-promotions
GET    /api/geo-promotions/{id}
GET    /api/geo-promotions/promotion/{promotionId}
GET    /api/geo-promotions/region/{regionCode}
GET    /api/geo-promotions/city/{city}
GET    /api/geo-promotions/near
POST   /api/geo-promotions
PUT    /api/geo-promotions/{id}
DELETE /api/geo-promotions/{id}
```

## Integration Points

The Menu Service integrates with:

1. **Analytics Service**: For menu view events and popularity metrics
2. **User Service**: For user segment information and preferences
3. **Discount Engine**: For personalized discounts and promotion eligibility

## Database Schema

The service uses the following database tables:

- `categories`: Product categories
- `products`: Menu items
- `price_history`: Historical price changes
- `promotions`: Promotional offers
- `geo_promotions`: Location-based promotions
- `segments`: Customer segments
- `personalized_menus`: Generated menus
- `personalized_menu_items`: Items within personalized menus

## Key Features

### Dynamic Pricing

The service implements dynamic pricing based on:
- Price elasticity
- User segment
- Time of day
- Historical purchase data

### Personalization

Menus are personalized based on:
- User segment (demographics, behavior)
- Location
- Purchase history
- Time of day

### Geo-Targeting

The service supports location-based promotions:
- Region/state-specific promotions
- City-specific promotions
- Radius-based promotions (using Haversine formula)

## Configuration

The service is configured via environment variables or application properties:

### Database Configuration
```properties
spring.datasource.url=jdbc:postgresql://${DB_HOST:localhost}:${DB_PORT:5432}/${DB_NAME:menu}
spring.datasource.username=${DB_USERNAME:postgres}
spring.datasource.password=${DB_PASSWORD:postgres}
```

### Redis Configuration
```properties
spring.redis.host=${REDIS_HOST:localhost}
spring.redis.port=${REDIS_PORT:6379}
spring.redis.password=${REDIS_PASSWORD:}
```

### Server Configuration
```properties
server.port=${SERVER_PORT:8080}
spring.profiles.active=${ACTIVE_PROFILE:dev}
```

## Getting Started

### Prerequisites
- JDK 8 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- Gradle 7 or higher

### Installation
1. Clone the repository
2. Navigate to the project directory
3. Run `./gradlew build` to build the project

### Running the Service
```bash
./gradlew bootRun
```

### API Documentation
Once the service is running, access the OpenAPI documentation at:
```
http://localhost:8080/swagger-ui.html
```

## Development

### Testing
```bash
./gradlew test
```

### Key Dependencies
- Spring Boot Web
- Spring Data JPA
- Spring Data Redis
- PostgreSQL Driver
- SpringDoc OpenAPI
- Lombok

### Future Enhancements
- Machine learning integration for predictive pricing
- Enhanced personalization based on dietary preferences
- A/B testing framework for menu layouts
- Advanced analytics for revenue optimization
