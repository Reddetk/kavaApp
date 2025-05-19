package com.kava.menu.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import org.springframework.context.support.PropertySourcesPlaceholderConfigurer;

/**
 * Configuration class to load environment variables from .env file
 */
@Configuration
@PropertySource(value = "file:.env", ignoreResourceNotFound = true)
public class EnvConfig {

    /**
     * Bean to enable resolution of ${...} placeholders in @Value annotations
     */
    @Bean
    public static PropertySourcesPlaceholderConfigurer propertySourcesPlaceholderConfigurer() {
        return new PropertySourcesPlaceholderConfigurer();
    }
}
