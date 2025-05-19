package com.kava.menu.service;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.math.BigDecimal;
import java.math.RoundingMode;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class PricingServiceOptimalPriceTest {

    private PricingService pricingService;

    @BeforeEach
    void setUp() {
        pricingService = new PricingService();
    }

    @Test
    @DisplayName("Calculate optimal price with unit elasticity (e=1)")
    void calculateOptimalPrice_WithUnitElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = BigDecimal.ONE;

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        // For unit elasticity, the price should be adjusted slightly
        assertTrue(result.compareTo(basePrice) >= 0);
    }

    @Test
    @DisplayName("Calculate optimal price with high elasticity (e=2)")
    void calculateOptimalPrice_WithHighElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = new BigDecimal("2.0");
        BigDecimal expected = new BigDecimal("200.00");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        assertEquals(expected, result);
    }

    @Test
    @DisplayName("Calculate optimal price with low elasticity (e=0.5)")
    void calculateOptimalPrice_WithLowElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = new BigDecimal("0.5");
        BigDecimal expected = new BigDecimal("120.00");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        assertEquals(expected, result);
    }

    @Test
    @DisplayName("Calculate optimal price with very high elasticity (e=5)")
    void calculateOptimalPrice_WithVeryHighElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = new BigDecimal("5.0");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        // For very high elasticity, the price should be higher than base price
        // but constrained by the max markup limit
        assertTrue(result.compareTo(basePrice) > 0);
        assertTrue(result.compareTo(basePrice.multiply(new BigDecimal("1.5"))) <= 0);
    }

    @Test
    @DisplayName("Calculate optimal price with zero elasticity (e=0)")
    void calculateOptimalPrice_WithZeroElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = BigDecimal.ZERO;
        BigDecimal expected = new BigDecimal("120.00");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        assertEquals(expected, result);
    }

    @Test
    @DisplayName("Calculate optimal price with negative elasticity (e=-0.5)")
    void calculateOptimalPrice_WithNegativeElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = new BigDecimal("-0.5");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        // For negative elasticity (which is unusual), the price should be higher
        assertTrue(result.compareTo(basePrice) > 0);
    }

    @Test
    @DisplayName("Calculate optimal price with very small elasticity (e=0.01)")
    void calculateOptimalPrice_WithVerySmallElasticity() {
        // Arrange
        BigDecimal basePrice = new BigDecimal("100.00");
        BigDecimal elasticity = new BigDecimal("0.01");

        // Act
        BigDecimal result = pricingService.calculateOptimalPrice(basePrice, elasticity);

        // Assert
        // For very small elasticity, the price should be higher than base price
        assertTrue(result.compareTo(basePrice) > 0);
    }
}
