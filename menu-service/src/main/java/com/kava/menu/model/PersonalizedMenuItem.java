package com.kava.menu.model;

import javax.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Entity
@Table(name = "personalized_menu_items")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class PersonalizedMenuItem {
    
    @EmbeddedId
    private PersonalizedMenuItemId id;
    
    @ManyToOne
    @MapsId("menuId")
    @JoinColumn(name = "menu_id")
    private PersonalizedMenu menu;
    
    @ManyToOne
    @MapsId("productId")
    @JoinColumn(name = "product_id")
    private Product product;
    
    @Column(name = "final_price")
    private BigDecimal finalPrice;
    
    @Column(name = "discount_applied")
    private Boolean discountApplied = false;
    
    @ManyToOne
    @JoinColumn(name = "promotion_id")
    private Promotion promotion;
}
