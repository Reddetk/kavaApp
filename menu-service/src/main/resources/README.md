# Menu Service Resources

This directory contains configuration files and resources for the Menu Service application.
   
## Configuration Files

### `application.yaml`

The main configuration file for the Spring Boot application. It includes:

```yaml
spring:
  profiles:
    active: ${ACTIVE_PROFILE:dev}
  datasource:
    url: jdbc:postgresql://${DB_HOST:localhost}:${DB_PORT:5432}/${DB_NAME:menu}
    username: ${DB_USERNAME:postgres}
    password: ${DB_PASSWORD:postgres}
    driver-class-name: org.postgresql.Driver
    hikari:
      connection-timeout: ${CONNECTION_TIMEOUT_MS:30000}
      maximum-pool-size: ${CONNECTION_POOL_SIZE:10}
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: true
    properties:
      hibernate:
        format_sql: true
        dialect: org.hibernate.dialect.PostgreSQLDialect
  data:
    redis:
      host: ${REDIS_HOST:localhost}
      port: ${REDIS_PORT:6379}
      password: ${REDIS_PASSWORD:}
      repositories:
        enabled: ${ENABLE_REDIS_CACHE:true}
  cache:
    redis:
      time-to-live: ${CACHE_TTL_SECONDS:300}000

server:
  port: ${SERVER_PORT:8080}

springdoc:
  api-docs:
    path: /api-docs
  swagger-ui:
    path: /swagger-ui.html
    operationsSorter: method
```

### `application-dev.yaml`

Development environment specific configuration:

```yaml
spring:
  jpa:
    show-sql: true
    hibernate:
      ddl-auto: update

logging:
  level:
    com.kava.menu: DEBUG
    org.hibernate.SQL: DEBUG
```

### `application-prod.yaml`

Production environment specific configuration:

```yaml
spring:
  jpa:
    show-sql: false
    hibernate:
      ddl-auto: validate

logging:
  level:
    com.kava.menu: INFO
    root: WARN
```

### `application-test.yaml`

Test environment specific configuration:

```yaml
spring:
  datasource:
    url: jdbc:h2:mem:testdb;DB_CLOSE_DELAY=-1;DB_CLOSE_ON_EXIT=FALSE
    username: sa
    password: 
    driver-class-name: org.h2.Driver
  jpa:
    hibernate:
      ddl-auto: create-drop
    show-sql: true
    properties:
      hibernate:
        format_sql: true
        dialect: org.hibernate.dialect.H2Dialect
  redis:
    repositories:
      enabled: false
  autoconfigure:
    exclude:
      - org.springframework.boot.autoconfigure.data.redis.RedisAutoConfiguration
```

## Database Migration

### `db/migration`

This directory contains Flyway database migration scripts:

- `V1__init_schema.sql`: Initial database schema creation
- `V2__add_price_history.sql`: Adds price history table
- `V3__add_geo_promotions.sql`: Adds geo-promotions functionality

## Static Resources

### `static`

This directory contains static resources served by the application:

- `images/`: Product and category images
- `docs/`: Additional documentation

### `templates`

This directory contains Thymeleaf templates (if applicable):

- `error/`: Error page templates

## Internationalization

### `messages`

This directory contains message bundles for internationalization:

- `messages.properties`: Default messages
- `messages_en.properties`: English messages
- `messages_es.properties`: Spanish messages

## Environment Variables

The application uses environment variables for configuration. These can be set in a `.env` file in the project root:

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=menu
DB_USERNAME=postgres
DB_PASSWORD=postgres

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server Configuration
SERVER_PORT=8080
ACTIVE_PROFILE=dev

# External Service URLs
ANALYTICS_SERVICE_URL=http://localhost:8081/api
USER_SERVICE_URL=http://localhost:8082/api
DISCOUNT_ENGINE_URL=http://localhost:8083/api
```

## Usage

The configuration files use Spring's property placeholder syntax `${VARIABLE:default_value}` to allow overriding values with environment variables.

For example, to change the database host:

```bash
export DB_HOST=new-database-host
./gradlew bootRun
```

Or when running with Java:

```bash
java -jar menu-service.jar --spring.datasource.url=jdbc:postgresql://new-host:5432/menu
```
