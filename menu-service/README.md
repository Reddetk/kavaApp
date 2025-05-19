# Menu Service

## Overview
Menu Service is a microservice responsible for dynamically generating personalized menus and promotional offers based on customer metrics and preference analysis. It's a core component of the Kava application ecosystem, providing tailored product recommendations and pricing strategies.

## Purpose
The primary goal of this service is to enhance the customer experience by delivering personalized menu options and promotions that are relevant to each user's preferences, location, and behavior patterns.

## Features
- 🛍️ **Personalized Menu Generation**: Creates customized menus based on user segments
- 🔖 **Dynamic Discount Management**: Applies intelligent pricing strategies based on demand elasticity
- 📍 **Geo-targeted Promotions**: Delivers location-specific offers and recommendations
- 🧠 **Caching Strategy**: Implements Redis-based caching with stale-while-revalidate approach
- 📘 **API Documentation**: Provides OpenAPI documentation for seamless integration

## Architecture

### Input Data
- 📊 Recommendations from analytics-service (associated products, lift factors)
- 💰 Product margin data
- 📍 User geolocation (for local promotions)
- 📦 Categories and products from database

### Output Data
- 🛍️ Personalized menu based on user segment
- 🔖 List of products with dynamic discounts
- 📘 OpenAPI documentation for integration with discount-engine and pwa-client

### Algorithms and Methods
- 📈 **Elasticity of Demand**: For price adjustments based on discount response
- 🧠 **Redis Caching**: Using stale-while-revalidate strategy for performance
- 🔁 **Geo-filtering**: Location-based filtering at query level
- 🧾 **OpenAPI Generation**: Using springdoc-openapi for API documentation

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
- **Product**: Represents menu items with pricing information
- **Category**: Organizes products into logical groups
- **Promotion**: Defines discount offers and their validity periods
- **GeoLocation**: Stores location data for geo-targeting
- **PersonalizedMenu**: Aggregates recommended products and promotions
- **PricingInfo**: Contains dynamic pricing calculations

### Services
- **MenuService**: Coordinates menu generation
- **PricingService**: Calculates dynamic prices based on elasticity
- **PromotionService**: Manages active promotions
- **GeoTargetingService**: Handles location-based filtering
- **RecommendationAdapterService**: Interfaces with analytics for recommendations
- **MenuAssembler**: Combines various components into a cohesive menu

### Controllers
- **MenuController**: Exposes endpoints for personalized menus
- **AdminCatalogController**: Provides administrative functions for catalog management

## Integration Points
| Service | Connection Type | Description |
|---------|----------------|-------------|
| analytics-service | REST | Receives recommendations and lift factors |
| user-service | REST | Obtains customer segment information |
| discount-engine | REST | Shares promotional product lists and discount response metrics |
| pwa-client | REST/WS | Delivers personalized menus to the frontend |
| Redis | Cache | Implements caching strategy |

## Metrics and Models
| Metric | Usage in Menu Service |
|--------|----------------------|
| Customer Segmentation | Determines menu categories by segment |
| Price Elasticity of Demand | Adjusts discounts based on customer sensitivity |
| GeoAffinity | Offers promotions by proximity |
| Redemption Rate | Analyzes menu success through discount-engine feedback |

## Getting Started

### Prerequisites
- Java 17 or higher
- Gradle
- Redis server

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

### Key Dependencies
- Spring Boot Web
- Spring Data Redis
- SpringDoc OpenAPI
- Lombok

### Future Enhancements
- Integration with PostgreSQL for persistent storage
- Implementation of A/B testing for menu layouts
- Enhanced analytics integration for real-time recommendation updates
