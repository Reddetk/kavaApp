package com.kava.menu.service;

import com.kava.menu.model.GeoPromotion;
import com.kava.menu.model.Product;
import com.kava.menu.model.Promotion;
import com.kava.menu.repository.GeoPromotionRepository;
import com.kava.menu.repository.ProductRepository;
import com.kava.menu.repository.PromotionRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.UUID;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class PromotionService {
    
    private final PromotionRepository promotionRepository;
    private final ProductRepository productRepository;
    private final GeoPromotionRepository geoPromotionRepository;
    
    public List<Promotion> getAllPromotions() {
        return promotionRepository.findAll();
    }
    
    public List<Promotion> getActivePromotions() {
        return promotionRepository.findActivePromotions(LocalDateTime.now());
    }
    
    public Optional<Promotion> getPromotionById(UUID id) {
        return promotionRepository.findById(id);
    }
    
    public List<Promotion> getActivePromotionsForProduct(UUID productId) {
        return promotionRepository.findActivePromotionsForProduct(productId, LocalDateTime.now());
    }
    
    @Transactional
    public Promotion createPromotion(Promotion promotion) {
        return promotionRepository.save(promotion);
    }
    
    @Transactional
    public Optional<Promotion> updatePromotion(UUID id, Promotion promotionDetails) {
        return promotionRepository.findById(id)
                .map(existingPromotion -> {
                    existingPromotion.setName(promotionDetails.getName());
                    existingPromotion.setDescription(promotionDetails.getDescription());
                    existingPromotion.setDiscountPercent(promotionDetails.getDiscountPercent());
                    existingPromotion.setStartDate(promotionDetails.getStartDate());
                    existingPromotion.setEndDate(promotionDetails.getEndDate());
                    existingPromotion.setPromotionCategory(promotionDetails.getPromotionCategory());
                    existingPromotion.setIsActive(promotionDetails.getIsActive());
                    return promotionRepository.save(existingPromotion);
                });
    }
    
    @Transactional
    public boolean deletePromotion(UUID id) {
        return promotionRepository.findById(id)
                .map(promotion -> {
                    promotionRepository.delete(promotion);
                    return true;
                })
                .orElse(false);
    }
    
    @Transactional
    public Optional<Promotion> deactivatePromotion(UUID id) {
        return promotionRepository.findById(id)
                .map(existingPromotion -> {
                    existingPromotion.setIsActive(false);
                    return promotionRepository.save(existingPromotion);
                });
    }
    
    @Transactional
    public Optional<Promotion> addProductToPromotion(UUID promotionId, UUID productId) {
        Optional<Promotion> promotionOpt = promotionRepository.findById(promotionId);
        Optional<Product> productOpt = productRepository.findById(productId);
        
        if (promotionOpt.isPresent() && productOpt.isPresent()) {
            Promotion promotion = promotionOpt.get();
            Product product = productOpt.get();
            
            promotion.getProducts().add(product);
            return Optional.of(promotionRepository.save(promotion));
        }
        
        return Optional.empty();
    }
    
    @Transactional
    public Optional<Promotion> removeProductFromPromotion(UUID promotionId, UUID productId) {
        return promotionRepository.findById(promotionId)
                .map(promotion -> {
                    Set<Product> updatedProducts = promotion.getProducts().stream()
                            .filter(product -> !product.getId().equals(productId))
                            .collect(Collectors.toSet());
                    
                    promotion.setProducts(updatedProducts);
                    return promotionRepository.save(promotion);
                });
    }
    
    @Transactional
    public GeoPromotion createGeoPromotion(GeoPromotion geoPromotion) {
        return geoPromotionRepository.save(geoPromotion);
    }
    
    public List<GeoPromotion> getGeoPromotionsForPromotion(UUID promotionId) {
        return promotionRepository.findById(promotionId)
                .map(geoPromotionRepository::findByPromotion)
                .orElse(List.of());
    }
    
    public List<GeoPromotion> getGeoPromotionsByRegion(String regionCode) {
        return geoPromotionRepository.findByRegionCode(regionCode);
    }
    
    public List<GeoPromotion> getGeoPromotionsByCity(String city) {
        return geoPromotionRepository.findByCity(city);
    }
    
    public List<GeoPromotion> getGeoPromotionsNearLocation(BigDecimal latitude, BigDecimal longitude) {
        return geoPromotionRepository.findPromotionsNearLocation(latitude, longitude);
    }
    
    @Transactional
    public boolean deleteGeoPromotion(UUID id) {
        return geoPromotionRepository.findById(id)
                .map(geoPromotion -> {
                    geoPromotionRepository.delete(geoPromotion);
                    return true;
                })
                .orElse(false);
    }
}
