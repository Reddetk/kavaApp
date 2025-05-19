# Menu Service Source Code

This directory contains the Java source code for the Menu Service application.

## Package Structure

### `com.kava.menu`

- **MenuServiceApplication.java**: Main application entry point with Spring Boot configuration

### `com.kava.menu.controller`

Controllers handle HTTP requests and define the REST API endpoints.

- **CategoryController**: Manages product categories
- **ProductController**: Manages product catalog
- **MenuController**: Handles personalized menu generation and retrieval
- **PromotionController**: Handles promotional offers and discounts
- **GeoPromotionController**: Manages location-based promotions

### `com.kava.menu.service`

Services implement the business logic of the application.

- **CategoryService**: Business logic for category management
- **ProductService**: Business logic for product management
- **PricingService**: Handles price calculations and dynamic pricing
- **PersonalizedMenuService**: Generates personalized menus based on user segments
- **PromotionService**: Manages promotional offers
- **GeoPromotionService**: Handles location-based promotions

### `com.kava.menu.repository`

Repositories provide data access interfaces for database operations.

- **CategoryRepository**: Data access for categories
- **ProductRepository**: Data access for products
- **PriceHistoryRepository**: Tracks product price changes
- **PersonalizedMenuRepository**: Stores generated menus
- **PromotionRepository**: Data access for promotions
- **GeoPromotionRepository**: Data access for location-based promotions
- **SegmentRepository**: Manages customer segments

### `com.kava.menu.model`

Models represent the domain entities of the application.

- **Category**: Represents product categories
- **Product**: Represents menu items
- **PriceHistory**: Tracks product price changes over time
- **PersonalizedMenu**: Represents a menu tailored to a specific segment
- **PersonalizedMenuItem**: Items within a personalized menu
- **Promotion**: Represents promotional offers
- **GeoPromotion**: Location-based promotions
- **Segment**: Customer segments for personalization

### `com.kava.menu.dto`

Data Transfer Objects (DTOs) are used for transferring data between layers.

- **MenuRequestDTO**: Request parameters for menu generation
- **MenuResponseDTO**: Response format for menu data
- **MenuItemDTO**: Representation of menu items for API responses
- **ProductDTO**: Simplified product representation for API responses
- **CategoryDTO**: Simplified category representation for API responses

### `com.kava.menu.config`

Configuration classes for the application.

- **EnvConfig**: Environment variable configuration
- **RedisConfig**: Redis cache configuration
- **SwaggerConfig**: OpenAPI documentation configuration
- **JpaConfig**: Database and JPA configuration
- **EnvironmentService**: Service to access environment variables

### `com.kava.menu.exception`

Custom exceptions and exception handling.

- **ResourceNotFoundException**: Thrown when a requested resource is not found
- **GlobalExceptionHandler**: Handles exceptions and returns appropriate HTTP responses

## Key Interfaces

### `CategoryRepository`

```java
public interface CategoryRepository extends JpaRepository<Category, UUID> {
    List<Category> findByNameContainingIgnoreCase(String name);
}
```

### `ProductRepository`

```java
public interface ProductRepository extends JpaRepository<Product, UUID> {
    List<Product> findByCategoryId(UUID categoryId);
    List<Product> findByIsActiveTrue();
}
```

### `PersonalizedMenuRepository`

```java
public interface PersonalizedMenuRepository extends JpaRepository<PersonalizedMenu, UUID> {
    List<PersonalizedMenu> findBySegmentId(UUID segmentId);
    Optional<PersonalizedMenu> findLatestMenuForSegment(UUID segmentId);
}
```

## Design Patterns

The application implements several design patterns:

1. **Repository Pattern**: For data access abstraction
2. **Service Layer Pattern**: For business logic encapsulation
3. **DTO Pattern**: For data transfer between layers
4. **Dependency Injection**: For loose coupling between components
5. **Builder Pattern**: For complex object construction (e.g., in DTOs)

## Best Practices

1. **Separation of Concerns**: Each class has a single responsibility
2. **Immutable DTOs**: DTOs are designed to be immutable
3. **Validation**: Input validation at controller level
4. **Exception Handling**: Centralized exception handling
5. **Logging**: Consistent logging throughout the application
6. **Testing**: Unit and integration tests for all components
