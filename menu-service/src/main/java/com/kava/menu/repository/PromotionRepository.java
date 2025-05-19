package com.kava.menu.repository;

import com.kava.menu.model.Promotion;
import com.kava.menu.model.PromotionCategory;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Repository
public interface PromotionRepository extends JpaRepository<Promotion, UUID> {
    
    List<Promotion> findByPromotionCategory(PromotionCategory category);
    
    @Query("SELECT p FROM Promotion p WHERE p.isActive = true AND p.startDate <= :now AND p.endDate >= :now")
    List<Promotion> findActivePromotions(LocalDateTime now);
    
    @Query("SELECT p FROM Promotion p JOIN p.products prod WHERE prod.id = :productId AND p.isActive = true AND p.startDate <= :now AND p.endDate >= :now")
    List<Promotion> findActivePromotionsForProduct(UUID productId, LocalDateTime now);
}
