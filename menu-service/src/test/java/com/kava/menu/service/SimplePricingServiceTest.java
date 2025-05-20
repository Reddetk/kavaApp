package com.kava.menu.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.math.BigDecimal;

import static org.junit.jupiter.api.Assertions.assertEquals;

class SimplePricingServiceTest {

    private SimplePricingService pricingService;

    @BeforeEach
    void setUp() {
        pricingService = new SimplePricingService();
    }

    @Test
    @DisplayName("Calculate discounted price with 10% discount")
    void calculateDiscountedPrice_With10PercentDiscount() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal discountPercent = new BigDecimal("10.0");
        BigDecimal expected = new BigDecimal("90.00");

        // Act
        BigDecimal result = pricingService.calculateDiscountedPrice(basePrice, discountPercent);

        // Assert
        assertEquals(expected, result);
    }

    @Test
    @DisplayName("Calculate discounted price with 0% discount")
    void calculateDiscountedPrice_WithZeroDiscount() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal discountPercent = BigDecimal.ZERO;
        BigDecimal expected = new BigDecimal("100.00");

        // Act
        BigDecimal result = pricingService.calculateDiscountedPrice(basePrice, discountPercent);

        // Assert
        assertEquals(expected, result);
    }

    @ParameterizedTest
    @CsvSource({
        "100.00, 1.5, 97.50",  // High elasticity (>1) - apply discount
        "100.00, 2.0, 95.00",  // Higher elasticity - more discount
        "100.00, 0.5, 101.50", // Low elasticity (<1) - apply premium
        "100.00, 0.0, 103.00"  // Very low elasticity - more premium
    })
    @DisplayName("Calculate dynamic price based on elasticity")
    void calculateDynamicPrice(String basePriceStr, String elasticityStr, String expectedStr) {
        // Arrange
        BigDecimal basePrice = new BigDecimal(basePriceStr);
        BigDecimal elasticity = new BigDecimal(elasticityStr);
        BigDecimal expected = new BigDecimal(expectedStr);

        // Act
        BigDecimal result = pricingService.calculateDynamicPrice(basePrice, elasticity);

        // Assert
        assertEquals(expected, result);
    }
}
