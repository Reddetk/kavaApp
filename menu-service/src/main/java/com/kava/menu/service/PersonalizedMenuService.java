package com.kava.menu.service;

import com.kava.menu.model.*;
import com.kava.menu.repository.PersonalizedMenuRepository;
import com.kava.menu.repository.ProductDemandMetricsRepository;
import com.kava.menu.repository.ProductRepository;
import com.kava.menu.repository.SegmentRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class PersonalizedMenuService {
    
    private final PersonalizedMenuRepository personalizedMenuRepository;
    private final SegmentRepository segmentRepository;
    private final ProductRepository productRepository;
    private final ProductDemandMetricsRepository productDemandMetricsRepository;
    private final PromotionService promotionService;
    private final PricingService pricingService;
    
    public List<PersonalizedMenu> getAllMenus() {
        return personalizedMenuRepository.findAll();
    }
    
    public Optional<PersonalizedMenu> getMenuById(UUID id) {
        return personalizedMenuRepository.findById(id);
    }
    
    public List<PersonalizedMenu> getMenusBySegment(UUID segmentId) {
        return segmentRepository.findById(segmentId)
                .map(personalizedMenuRepository::findBySegment)
                .orElse(List.of());
    }
    
    public Optional<PersonalizedMenu> getLatestMenuForSegment(UUID segmentId) {
        return personalizedMenuRepository.findLatestMenuForSegment(segmentId);
    }
    
    @Transactional
    public PersonalizedMenu createMenu(PersonalizedMenu menu) {
        return personalizedMenuRepository.save(menu);
    }
    
    @Transactional
    public Optional<PersonalizedMenu> generateMenuForSegment(UUID segmentId) {
        return segmentRepository.findById(segmentId)
                .map(segment -> {
                    // Create a new personalized menu
                    PersonalizedMenu menu = new PersonalizedMenu();
                    menu.setSegment(segment);
                    
                    // Get top products for this segment based on lift factor
                    List<ProductDemandMetrics> topProducts = 
                            productDemandMetricsRepository.findTopProductsByLiftFactorForSegment(segmentId);
                    
                    // Create menu items for each product
                    for (ProductDemandMetrics metrics : topProducts) {
                        Product product = metrics.getProduct();
                        
                        // Skip inactive products
                        if (!product.getIsActive()) {
                            continue;
                        }
                        
                        // Get active promotions for this product
                        List<Promotion> activePromotions = 
                                promotionService.getActivePromotionsForProduct(product.getId());
                        
                        // Calculate final price based on pricing strategy
                        BigDecimal finalPrice;
                        Promotion appliedPromotion = null;
                        boolean discountApplied = false;
                        
                        if (!activePromotions.isEmpty()) {
                            // Use the promotion with the highest discount
                            appliedPromotion = activePromotions.stream()
                                    .max((p1, p2) -> p1.getDiscountPercent().compareTo(p2.getDiscountPercent()))
                                    .orElse(null);
                            
                            if (appliedPromotion != null) {
                                // Apply promotion discount
                                finalPrice = pricingService.calculateDiscountedPrice(
                                        product.getBasePrice(), 
                                        appliedPromotion.getDiscountPercent());
                                discountApplied = true;
                            } else {
                                finalPrice = product.getBasePrice();
                            }
                        } else {
                            // Use dynamic pricing based on elasticity if no promotions
                            finalPrice = pricingService.calculateDynamicPrice(
                                    product.getBasePrice(),
                                    metrics.getPriceElasticity());
                        }
                        
                        // Create menu item
                        PersonalizedMenuItem menuItem = new PersonalizedMenuItem();
                        menuItem.setId(new PersonalizedMenuItemId(menu.getId(), product.getId()));
                        menuItem.setMenu(menu);
                        menuItem.setProduct(product);
                        menuItem.setFinalPrice(finalPrice);
                        menuItem.setDiscountApplied(discountApplied);
                        menuItem.setPromotion(appliedPromotion);
                        
                        menu.getMenuItems().add(menuItem);
                    }
                    
                    // Save the menu
                    return personalizedMenuRepository.save(menu);
                });
    }
    
    @Transactional
    public boolean deleteMenu(UUID id) {
        return personalizedMenuRepository.findById(id)
                .map(menu -> {
                    personalizedMenuRepository.delete(menu);
                    return true;
                })
                .orElse(false);
    }
    
    @Transactional
    public Optional<PersonalizedMenu> addItemToMenu(UUID menuId, UUID productId, BigDecimal finalPrice, 
                                                   boolean discountApplied, UUID promotionId) {
        Optional<PersonalizedMenu> menuOpt = personalizedMenuRepository.findById(menuId);
        Optional<Product> productOpt = productRepository.findById(productId);
        
        if (menuOpt.isPresent() && productOpt.isPresent()) {
            PersonalizedMenu menu = menuOpt.get();
            Product product = productOpt.get();
            
            // Create menu item
            PersonalizedMenuItem menuItem = new PersonalizedMenuItem();
            menuItem.setId(new PersonalizedMenuItemId(menuId, productId));
            menuItem.setMenu(menu);
            menuItem.setProduct(product);
            menuItem.setFinalPrice(finalPrice);
            menuItem.setDiscountApplied(discountApplied);
            
            if (promotionId != null) {
                promotionService.getPromotionById(promotionId)
                        .ifPresent(menuItem::setPromotion);
            }
            
            menu.getMenuItems().add(menuItem);
            return Optional.of(personalizedMenuRepository.save(menu));
        }
        
        return Optional.empty();
    }
    
    @Transactional
    public Optional<PersonalizedMenu> removeItemFromMenu(UUID menuId, UUID productId) {
        return personalizedMenuRepository.findById(menuId)
                .map(menu -> {
                    menu.getMenuItems().removeIf(item -> 
                            item.getId().getProductId().equals(productId));
                    return personalizedMenuRepository.save(menu);
                });
    }
}
