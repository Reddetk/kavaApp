package com.kava.menu.repository;

import com.kava.menu.model.PriceHistory;
import com.kava.menu.model.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.PageRequest;

@Repository
public interface PriceHistoryRepository extends JpaRepository<PriceHistory, Long> {

    List<PriceHistory> findByProduct(Product product);

    @Query("SELECT ph FROM PriceHistory ph WHERE ph.product.id = :productId AND ph.changedAt >= :since ORDER BY ph.changedAt DESC")
    List<PriceHistory> findPriceHistoryForProductSince(UUID productId, LocalDateTime since);

    @Query("SELECT ph FROM PriceHistory ph WHERE ph.product.id = :productId ORDER BY ph.changedAt DESC")
    List<PriceHistory> findPriceHistoryForProductOrderByChangedAtDesc(UUID productId, Pageable pageable);

    default PriceHistory findLatestPriceChangeForProduct(UUID productId) {
        List<PriceHistory> history = findPriceHistoryForProductOrderByChangedAtDesc(
            productId,
            PageRequest.of(0, 1)
        );
        return history.isEmpty() ? null : history.get(0);
    }
}
