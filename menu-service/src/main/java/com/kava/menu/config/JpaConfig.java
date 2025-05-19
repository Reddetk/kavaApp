package com.kava.menu.config;

import org.hibernate.boot.model.naming.CamelCaseToUnderscoresNamingStrategy;
import org.hibernate.boot.model.naming.PhysicalNamingStrategy;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * JPA configuration for consistent database schema naming
 */
@Configuration
public class JpaConfig {

    /**
     * Configure Hibernate to use a consistent naming strategy for database objects
     * This helps prevent mismatches between entity definitions and database schema
     */
    @Bean
    public PhysicalNamingStrategy physicalNamingStrategy() {
        return new CamelCaseToUnderscoresNamingStrategy();
    }
}
