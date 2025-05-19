package com.kava.menu.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class MenuItemDTO {
    private UUID productId;
    private String productName;
    private String productDescription;
    private BigDecimal basePrice;
    private BigDecimal finalPrice;
    private boolean discountApplied;
    private UUID promotionId;
    private String promotionName;
}
