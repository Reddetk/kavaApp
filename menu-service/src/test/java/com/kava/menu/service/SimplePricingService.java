package com.kava.menu.service;

import java.math.BigDecimal;
import java.math.RoundingMode;

public class SimplePricingService {
    
    /**
     * Calculate a discounted price based on a discount percentage
     * 
     * @param basePrice The original price
     * @param discountPercent The discount percentage (e.g., 10.0 for 10%)
     * @return The discounted price
     */
    public BigDecimal calculateDiscountedPrice(BigDecimal basePrice, BigDecimal discountPercent) {
        BigDecimal discountMultiplier = BigDecimal.ONE.subtract(
                discountPercent.divide(BigDecimal.valueOf(100), 4, RoundingMode.HALF_UP));
        
        return basePrice.multiply(discountMultiplier)
                .setScale(2, RoundingMode.HALF_UP);
    }
    
    /**
     * Calculate a dynamic price based on price elasticity
     * 
     * @param basePrice The original price
     * @param priceElasticity The price elasticity factor
     * @return The dynamically adjusted price
     */
    public BigDecimal calculateDynamicPrice(BigDecimal basePrice, BigDecimal priceElasticity) {
        // If elasticity is high (>1), consider a small discount to drive volume
        // If elasticity is low (<1), consider a small premium as demand is less sensitive to price
        
        BigDecimal adjustmentFactor;
        
        if (priceElasticity.compareTo(BigDecimal.ONE) > 0) {
            // High elasticity - apply discount
            adjustmentFactor = BigDecimal.ONE.subtract(
                    priceElasticity.subtract(BigDecimal.ONE)
                            .multiply(BigDecimal.valueOf(0.05))  // 5% per elasticity point above 1
                            .min(BigDecimal.valueOf(0.15)));     // Max 15% discount
        } else {
            // Low elasticity - apply premium
            adjustmentFactor = BigDecimal.ONE.add(
                    BigDecimal.ONE.subtract(priceElasticity)
                            .multiply(BigDecimal.valueOf(0.03))  // 3% per elasticity point below 1
                            .min(BigDecimal.valueOf(0.10)));     // Max 10% premium
        }
        
        return basePrice.multiply(adjustmentFactor)
                .setScale(2, RoundingMode.HALF_UP);
    }
}
