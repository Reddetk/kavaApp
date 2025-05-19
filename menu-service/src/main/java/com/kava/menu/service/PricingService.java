package com.kava.menu.service;

import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;

@Service
public class PricingService {
    
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
    
    /**
     * Calculate optimal price based on elasticity to maximize revenue
     * 
     * @param basePrice The original price
     * @param priceElasticity The price elasticity of demand
     * @return The optimal price
     */
    public BigDecimal calculateOptimalPrice(BigDecimal basePrice, BigDecimal priceElasticity) {
        // For unit elastic demand (e=1), any price gives same revenue
        // For inelastic demand (e<1), higher price gives more revenue
        // For elastic demand (e>1), lower price gives more revenue
        
        // Optimal markup formula: m = e/(e-1) where e is elasticity and e>1
        // For e<1, we use a different approach
        
        if (priceElasticity.compareTo(BigDecimal.ONE) <= 0) {
            // For inelastic demand, increase price by a factor based on how inelastic it is
            BigDecimal inelasticityFactor = BigDecimal.ONE.subtract(priceElasticity)
                    .multiply(BigDecimal.valueOf(0.2))  // 20% of the inelasticity factor
                    .add(BigDecimal.ONE);               // Add 1 to get the multiplier
            
            return basePrice.multiply(inelasticityFactor)
                    .setScale(2, RoundingMode.HALF_UP);
        } else {
            // For elastic demand, use the optimal markup formula
            BigDecimal denominator = priceElasticity.subtract(BigDecimal.ONE);
            
            // Avoid division by zero or very small numbers
            if (denominator.abs().compareTo(BigDecimal.valueOf(0.1)) < 0) {
                denominator = BigDecimal.valueOf(0.1).multiply(
                        denominator.signum() > 0 ? BigDecimal.ONE : BigDecimal.valueOf(-1));
            }
            
            BigDecimal optimalMarkup = priceElasticity.divide(denominator, 4, RoundingMode.HALF_UP);
            
            // Constrain the markup to reasonable bounds
            optimalMarkup = optimalMarkup.max(BigDecimal.valueOf(0.8))  // Min 20% discount
                    .min(BigDecimal.valueOf(1.5));                      // Max 50% premium
            
            return basePrice.multiply(optimalMarkup)
                    .setScale(2, RoundingMode.HALF_UP);
        }
    }
}
