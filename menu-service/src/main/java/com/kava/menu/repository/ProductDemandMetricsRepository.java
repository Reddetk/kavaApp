package com.kava.menu.repository;

import com.kava.menu.model.Product;
import com.kava.menu.model.ProductDemandMetrics;
import com.kava.menu.model.ProductDemandMetricsId;
import com.kava.menu.model.Segment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface ProductDemandMetricsRepository extends JpaRepository<ProductDemandMetrics, ProductDemandMetricsId> {
    
    List<ProductDemandMetrics> findByProduct(Product product);
    
    List<ProductDemandMetrics> findBySegment(Segment segment);
    
    @Query("SELECT pdm FROM ProductDemandMetrics pdm WHERE pdm.segment.id = :segmentId ORDER BY pdm.liftFactor DESC")
    List<ProductDemandMetrics> findTopProductsByLiftFactorForSegment(UUID segmentId);
    
    @Query("SELECT pdm FROM ProductDemandMetrics pdm WHERE pdm.product.id = :productId AND pdm.segment.id = :segmentId")
    ProductDemandMetrics findByProductAndSegment(UUID productId, UUID segmentId);
}
