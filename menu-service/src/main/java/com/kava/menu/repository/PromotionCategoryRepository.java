package com.kava.menu.repository;

import com.kava.menu.model.PromotionCategory;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface PromotionCategoryRepository extends JpaRepository<PromotionCategory, UUID> {
    
    Optional<PromotionCategory> findByName(String name);
}
