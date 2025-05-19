package com.kava.menu.repository;

import com.kava.menu.model.PriceHistory;
import com.kava.menu.model.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;

@Repository
public interface PriceHistoryRepository extends JpaRepository<PriceHistory, Long> {
    
    List<PriceHistory> findByProduct(Product product);
    
    @Query("SELECT ph FROM PriceHistory ph WHERE ph.product.id = :productId AND ph.changedAt >= :since ORDER BY ph.changedAt DESC")
    List<PriceHistory> findPriceHistoryForProductSince(Long productId, LocalDateTime since);
    
    @Query("SELECT ph FROM PriceHistory ph WHERE ph.product.id = :productId ORDER BY ph.changedAt DESC LIMIT 1")
    PriceHistory findLatestPriceChangeForProduct(Long productId);
}
