package com.kava.menu.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class MenuRequestDTO {
    private UUID segmentId;
    private BigDecimal latitude;
    private BigDecimal longitude;
    private String regionCode;
}
