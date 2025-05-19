package com.kava.menu.model;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "geo_promotions")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class GeoPromotion {
    
    @Id
    private UUID id;
    
    @ManyToOne
    @JoinColumn(name = "promotion_id")
    private Promotion promotion;
    
    @Column(name = "region_code", nullable = false)
    private String regionCode;
    
    private String city;
    
    private BigDecimal latitude;
    
    private BigDecimal longitude;
    
    @Column(name = "radius_km")
    private BigDecimal radiusKm;
    
    @Column(name = "created_at")
    private LocalDateTime createdAt;
    
    @PrePersist
    protected void onCreate() {
        if (id == null) {
            id = UUID.randomUUID();
        }
        if (createdAt == null) {
            createdAt = LocalDateTime.now();
        }
    }
}
