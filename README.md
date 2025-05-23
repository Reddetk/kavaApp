# KavaApp: Dynamic Pricing & Personalization Platform

## Overview

KavaApp is a sophisticated platform designed for optimizing customer engagement through dynamic pricing and personalized discount systems. Leveraging a microservices architecture, the platform aims to enhance customer retention and drive repeat purchases by offering targeted, profitable promotions.

## Core Objectives

The platform's primary objectives are to:

- **Customer Segmentation**: Accurately classify customers into distinct segments upon registration and dynamically update their segment based on subsequent purchase behavior.
- **Personalized Discount Calculation**: Compute individualized discounts that maximally incentivize repeat transactions.
- **Profitability Balancing**: Strategically balance demand stimulation with customer retention goals, ensuring promotional offers are both attractive to the customer and profitable for the business.

## Microservices Architecture

The KavaApp platform is composed of several interconnected microservices, each responsible for a specific domain or function. This architecture promotes scalability, resilience, and maintainability.

### 1. User Service

**Purpose**: Manages user profiles and updates customer segments based on behavioral analytics.

**Architecture**: Adheres to Clean Architecture principles with distinct layers for Domain, Application, Infrastructure, and Interfaces.

**Key Features:**
- Comprehensive user profile management.
- Advanced customer segmentation (RFM, Behavioral, K-means clustering).
- Customer Lifetime Value (CLV) calculation and tracking.
- User retention prediction and survival analysis (Cox proportional hazards model).
- State transition analysis using Markov chains.
- PostgreSQL for persistent data storage.
- Kafka consumer for real-time transaction event processing.
- HTTP handlers for external API interactions.

**Technology Stack:**
- **Language**: Golang
- **Database**: PostgreSQL
- **Messaging**: Kafka
- **Frameworks/Libraries**: gRPC, Gin-gonic, `gonum/stat`, `github.com/sajari/regression`

### 2. Analytics Service

**Purpose**: Provides data analysis capabilities, including association rule mining, product recommendations, and product performance analysis.

**Architecture**: Follows a clean architecture pattern with layers for Domain, Application, Infrastructure, and Interfaces.

**Key Features:**
- Association Rule Mining using the Apriori algorithm.
- Generation of personalized product recommendations based on discovered rules.
- ABC Analysis for product categorization based on revenue contribution.
- Supports A/B testing of discount strategies.
- Identifies product associations and purchase patterns.

**Technology Stack:**
- **Language**: Golang
- **Database**: ClickHouse (for columnar data storage), PostgreSQL (for transactional data)
- **Messaging**: Kafka (for real-time processing)
- **Frameworks/Libraries**: `github.com/eMAGTechLabs/go-apriori`, `github.com/sirupsen/logrus`, `github.com/joho/godotenv`, `gopkg.in/yaml.v3`

### 3. Menu Service

**Purpose**: Facilitates dynamic and personalized menu generation.

**Architecture**: Employs a layered architecture comprising Controllers, Services, Repositories, and Models.

**Key Features:**
- Dynamic generation of personalized menus based on user segments.
- Implementation of dynamic pricing strategies informed by demand elasticity.
- Delivery of geo-targeted promotions based on location data.
- Tracking of product price history.
- Performance optimization through Redis-based caching.
- OpenAPI documentation for API endpoints.
- Management of menu items, categories, products, and promotions.

**Technology Stack:**
- **Language**: Java
- **Framework**: Spring Boot 2.7.x
- **Database**: PostgreSQL
- **Cache**: Redis
- **Documentation**: SpringDoc OpenAPI
- **Build Tool**: Gradle
- **Algorithms/Methods**: Elasticity of Demand, Redis Caching, Geo-filtering (Haversine formula)

### 4. Discount Engine

**Purpose**: Calculates and manages personalized discounts.

**Key Features:**
- Calculation of personalized discounts using game theory and XGBoost.
- Integration with payment gateways.
- Confirmation buffer utilizing Redis.

**Technology Stack:**
- **Language**: Golang
- **Database**: PostgreSQL
- **Cache**: Redis
- **Libraries**: XGBoost

### 5. API Gateway

**Purpose**: Serves as the single entry point for external clients, handling request routing, authentication, and aggregation.

**Key Features:**
- Centralized request processing.
- JWT authentication and authorization.
- Request aggregation.
- Rate limiting.
- Static content caching.
- Health check monitoring.

**Technology Stack:**
- **Language**: Golang
- **Frameworks**: Gin-gonic, KrakenD

### 6. KavaApp PWA

**Purpose**: The Progressive Web Application frontend providing the user interface.

**Key Features:**
- Responsive design for various devices.
- Offline-first capabilities.
- Personal user account management.
- Dynamic menu display.
- Interactive analytics visualization.

**Technology Stack:**
- **Framework**: React
- **State Management**: Redux, Redux Thunk
- **WebAssembly**: WASM modules on Go
- **Real-time Updates**: WebSocket

## Core Modules and Operational Logic

### 1. Customer Segmentation Module

**Objective**: Assign customers to specific segments using clustering algorithms (K-Means, DBSCAN).

**Input Data**: Recency (R), Frequency (F), Average Order Value (AOV), Order Composition, Customer Lifetime Value (CLV).

**Output**: Customer segment, updated dynamically after each purchase transaction.

### 2. Personalized Discount Calculation Module

**Objective**: Propose discounts that enhance customer engagement without compromising profitability.

**Methodologies**: Game Theory (modeling customer response to discounts), Gradient Boosting (determining optimal discount levels).

**Output**: Dynamically calculated discount percentage (ranging from 0% to a defined maximum).

### 3. Preference Analysis Module

**Objective**: Identify products that a customer is most likely to purchase.

**Methodologies**: Apriori algorithm (analysis of frequently co-purchased items), Correlation Analysis (identifying preferences of similar customer segments).

**Output**: Identification of products suitable for targeted discounts.

### 4. Profitability Control Module

**Objective**: Prevent the application of unprofitable discounts.

**Input Data**: Product margin data, cost of goods sold.

**Output**: Adjustment of proposed discounts to ensure business profitability.

## Integrated Operational Flow

1.  A customer completes a purchase transaction, triggering an update to their segment.
2.  The system analyzes the customer's RFM metrics, purchase preferences, and propensity for discounts.
3.  A personalized discount is calculated, and relevant products are recommended.
4.  The profitability control module verifies that the proposed discount is economically viable for the business.

**Outcome**: The customer receives a tailored offer, significantly increasing their engagement and incentive for future purchases.

## Technical Implementation Details

### Data Storage
- **Transactional/Relational Data**: PostgreSQL (User profiles, purchase history, segments, menu items, promotions, etc.)
- **Columnar Data**: ClickHouse (for analytics, A/B testing, behavioral data)
- **Caching/Temporary Data**: Redis (session caching, confirmation buffer, geo-location caching)

### Messaging and Event Processing
- **Asynchronous Updates**: Kafka (for real-time processing of transactions and triggering segment updates)

### Machine Learning and Analytics Libraries
- `gonum/stat`: Statistical methods for data analysis.
- `github.com/sajari/regression`: Regression analysis for predictive modeling.
- `github.com/eMAGTechLabs/go-apriori`: Apriori algorithm implementation.
- XGBoost: Gradient Boosting for discount optimization.

## Key System Characteristics

- **Architecture**: Fully microservice-oriented.
- **Performance**: High throughput and low latency leveraging Kafka and Redis.
- **Optimization**: Discount strategies informed by machine learning models.
- **Flexibility**: Adaptive APIs implemented in Golang and Java.
- **Scalability**: Designed for horizontal scaling of individual services.
- **Frontend**: Scalable and responsive PWA with offline capabilities.

## Setup & Deployment

### Prerequisites
- Go 1.23+
- Java Development Kit (JDK) [Specific version based on Menu Service requirements, e.g., JDK 11+ or 17+]
- PostgreSQL 12+
- ClickHouse
- Redis 7+
- Kafka
- Docker and Kubernetes (for containerized deployment)

### Configuration
Each service typically uses a combination of environment variables and configuration files (e.g., `.env`, `config.yaml`, `application.yaml`). Refer to the individual service directories for specific configuration details.

### Running Services
Refer to the `README.md` within each service directory (`analitics-service`, `user-service`, `menu-service`, etc.) for specific instructions on building and running individual components.

## Monitoring and Observability

The platform incorporates tools for monitoring and tracing:
- **Metrics**: Prometheus
- **Visualization**: Grafana
- **Distributed Tracing**: Jaeger
- **Error Tracking**: Sentry

## Integration Points

The services integrate with each other and potentially external systems:
- **Menu Service** integrates with **Analytics Service** (for menu view events) and **User Service** (for segment information).
- **Discount Engine** integrates with payment systems.
- **Analytics Service** processes events from Kafka.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
