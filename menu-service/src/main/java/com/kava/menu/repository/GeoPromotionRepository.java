package com.kava.menu.repository;

import com.kava.menu.model.GeoPromotion;
import com.kava.menu.model.Promotion;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

@Repository
public interface GeoPromotionRepository extends JpaRepository<GeoPromotion, UUID> {

    List<GeoPromotion> findByPromotion(Promotion promotion);

    List<GeoPromotion> findByRegionCode(String regionCode);

    List<GeoPromotion> findByCity(String city);

    @Query(value = "SELECT * FROM geo_promotion g WHERE " +
            "6371 * acos(cos(radians(:latitude)) * cos(radians(g.latitude)) * cos(radians(g.longitude) - radians(:longitude)) + sin(radians(:latitude)) * sin(radians(g.latitude))) <= g.radius_km", nativeQuery = true)
    List<GeoPromotion> findPromotionsNearLocation(BigDecimal latitude, BigDecimal longitude);
}
