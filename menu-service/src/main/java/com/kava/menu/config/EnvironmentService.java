package com.kava.menu.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

/**
 * Service to access environment variables in the application
 */
@Service
public class EnvironmentService {

    @Value("${spring.profiles.active:dev}")
    private String activeProfile;

    @Value("${services.analytics.url}")
    private String analyticsServiceUrl;

    @Value("${services.user.url}")
    private String userServiceUrl;

    @Value("${services.discount-engine.url}")
    private String discountEngineUrl;

    @Value("${features.geo-targeting:true}")
    private boolean geoTargetingEnabled;

    @Value("${features.dynamic-pricing:true}")
    private boolean dynamicPricingEnabled;

    /**
     * Check if the application is running in development mode
     * @return true if the active profile is dev
     */
    public boolean isDevelopment() {
        return "dev".equals(activeProfile);
    }

    /**
     * Check if the application is running in production mode
     * @return true if the active profile is prod
     */
    public boolean isProduction() {
        return "prod".equals(activeProfile);
    }

    /**
     * Check if the application is running in test mode
     * @return true if the active profile is test
     */
    public boolean isTest() {
        return "test".equals(activeProfile);
    }

    /**
     * Get the URL of the analytics service
     * @return the URL of the analytics service
     */
    public String getAnalyticsServiceUrl() {
        return analyticsServiceUrl;
    }

    /**
     * Get the URL of the user service
     * @return the URL of the user service
     */
    public String getUserServiceUrl() {
        return userServiceUrl;
    }

    /**
     * Get the URL of the discount engine
     * @return the URL of the discount engine
     */
    public String getDiscountEngineUrl() {
        return discountEngineUrl;
    }

    /**
     * Check if geo-targeting is enabled
     * @return true if geo-targeting is enabled
     */
    public boolean isGeoTargetingEnabled() {
        return geoTargetingEnabled;
    }

    /**
     * Check if dynamic pricing is enabled
     * @return true if dynamic pricing is enabled
     */
    public boolean isDynamicPricingEnabled() {
        return dynamicPricingEnabled;
    }
}
