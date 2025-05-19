package com.kava.menu.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class MenuResponseDTO {
    private UUID id;
    private UUID segmentId;
    private String segmentName;
    private LocalDateTime generatedAt;
    private List<MenuItemDTO> items;
}
