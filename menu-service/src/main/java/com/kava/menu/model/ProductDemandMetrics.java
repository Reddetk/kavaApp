package com.kava.menu.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "product_demand_metrics")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductDemandMetrics {
    
    @EmbeddedId
    private ProductDemandMetricsId id;
    
    @ManyToOne
    @MapsId("productId")
    @JoinColumn(name = "product_id")
    private Product product;
    
    @ManyToOne
    @MapsId("segmentId")
    @JoinColumn(name = "segment_id")
    private Segment segment;
    
    @Column(name = "lift_factor")
    private BigDecimal liftFactor;
    
    @Column(name = "redemption_rate")
    private BigDecimal redemptionRate;
    
    @Column(name = "price_elasticity")
    private BigDecimal priceElasticity;
    
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
    
    @PrePersist
    @PreUpdate
    protected void onUpdate() {
        updatedAt = LocalDateTime.now();
    }
}
