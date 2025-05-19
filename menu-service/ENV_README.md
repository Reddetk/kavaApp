# Environment Configuration for Menu Service

This document explains how to use the `.env` file for configuring the Menu Service application.

## Overview

The Menu Service uses environment variables for configuration to:
- Keep sensitive information out of the codebase
- Allow different configurations for different environments
- Make the application more portable and easier to deploy

## Using the .env File

1. Copy the `.env` file to the root of your project
2. Modify the values as needed for your environment
3. The application will automatically load these values at startup

## Available Configuration Options

### Database Configuration
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_NAME` - Database name (default: menu)
- `DB_USERNAME` - Database username (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)

### Redis Configuration
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (default: empty)

### Server Configuration
- `SERVER_PORT` - Server port (default: 8080)
- `ACTIVE_PROFILE` - Active Spring profile (default: dev)

### External Service URLs
- `ANALYTICS_SERVICE_URL` - URL of the analytics service
- `USER_SERVICE_URL` - URL of the user service
- `DISCOUNT_ENGINE_URL` - URL of the discount engine

### Logging Configuration
- `LOG_LEVEL` - Log level (default: INFO)
- `LOG_FILE_PATH` - Path to log file (default: logs/menu-service.log)

### Performance Tuning
- `CACHE_TTL_SECONDS` - Cache time-to-live in seconds (default: 300)
- `CONNECTION_POOL_SIZE` - Database connection pool size (default: 10)
- `CONNECTION_TIMEOUT_MS` - Database connection timeout in milliseconds (default: 30000)

### Feature Flags
- `ENABLE_REDIS_CACHE` - Enable Redis caching (default: true)
- `ENABLE_GEO_TARGETING` - Enable geo-targeting features (default: true)
- `ENABLE_DYNAMIC_PRICING` - Enable dynamic pricing features (default: true)

### Test Configuration
- `TEST_DB_URL` - Test database URL
- `TEST_DB_USERNAME` - Test database username
- `TEST_DB_PASSWORD` - Test database password

### Java Development Kit Configuration
- `JAVA_HOME` - Path to JDK installation
- `PATH` - System PATH including JDK bin directory

## Accessing Environment Variables in Code

You can access these environment variables in your code using the `EnvironmentService`:

```java
@Service
public class YourService {
    
    private final EnvironmentService environmentService;
    
    @Autowired
    public YourService(EnvironmentService environmentService) {
        this.environmentService = environmentService;
    }
    
    public void someMethod() {
        if (environmentService.isDevelopment()) {
            // Development-specific code
        }
        
        String analyticsUrl = environmentService.getAnalyticsServiceUrl();
        // Use the URL to make API calls
    }
}
```

## Overriding Environment Variables

You can override the values in the `.env` file by:
1. Setting actual environment variables on your system
2. Passing JVM arguments when starting the application:
   ```
   java -jar menu-service.jar --spring.datasource.url=jdbc:postgresql://new-host:5432/menu
   ```
3. Using Spring Boot's application-{profile}.yaml files for different environments
