package com.kava.menu.controller;

import com.kava.menu.dto.MenuItemDTO;
import com.kava.menu.dto.MenuRequestDTO;
import com.kava.menu.dto.MenuResponseDTO;
import com.kava.menu.model.PersonalizedMenu;
import com.kava.menu.model.PersonalizedMenuItem;
import com.kava.menu.service.PersonalizedMenuService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/api/menus")
@RequiredArgsConstructor
@Tag(name = "Menus", description = "Personalized menu endpoints")
public class MenuController {
    
    private final PersonalizedMenuService menuService;
    
    @GetMapping
    @Operation(summary = "Get all menus", description = "Retrieves a list of all personalized menus")
    public ResponseEntity<List<MenuResponseDTO>> getAllMenus() {
        List<MenuResponseDTO> menus = menuService.getAllMenus().stream()
                .map(this::convertToDTO)
                .collect(Collectors.toList());
        return ResponseEntity.ok(menus);
    }
    
    @GetMapping("/{id}")
    @Operation(summary = "Get menu by ID", description = "Retrieves a personalized menu by its UUID")
    public ResponseEntity<MenuResponseDTO> getMenuById(@PathVariable UUID id) {
        return menuService.getMenuById(id)
                .map(this::convertToDTO)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/segment/{segmentId}")
    @Operation(summary = "Get menus by segment", description = "Retrieves all menus for a specific segment")
    public ResponseEntity<List<MenuResponseDTO>> getMenusBySegment(@PathVariable UUID segmentId) {
        List<MenuResponseDTO> menus = menuService.getMenusBySegment(segmentId).stream()
                .map(this::convertToDTO)
                .collect(Collectors.toList());
        return ResponseEntity.ok(menus);
    }
    
    @GetMapping("/segment/{segmentId}/latest")
    @Operation(summary = "Get latest menu for segment", description = "Retrieves the most recent menu for a segment")
    public ResponseEntity<MenuResponseDTO> getLatestMenuForSegment(@PathVariable UUID segmentId) {
        return menuService.getLatestMenuForSegment(segmentId)
                .map(this::convertToDTO)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PostMapping("/generate")
    @Operation(summary = "Generate a personalized menu", description = "Generates a new personalized menu for a segment")
    public ResponseEntity<MenuResponseDTO> generateMenu(@RequestBody MenuRequestDTO request) {
        return menuService.generateMenuForSegment(request.getSegmentId())
                .map(this::convertToDTO)
                .map(dto -> ResponseEntity.status(HttpStatus.CREATED).body(dto))
                .orElse(ResponseEntity.badRequest().build());
    }
    
    @DeleteMapping("/{id}")
    @Operation(summary = "Delete a menu", description = "Deletes a personalized menu by its UUID")
    public ResponseEntity<Void> deleteMenu(@PathVariable UUID id) {
        boolean deleted = menuService.deleteMenu(id);
        return deleted ? ResponseEntity.noContent().build() : ResponseEntity.notFound().build();
    }
    
    /**
     * Converts a PersonalizedMenu entity to a MenuResponseDTO
     */
    private MenuResponseDTO convertToDTO(PersonalizedMenu menu) {
        MenuResponseDTO dto = new MenuResponseDTO();
        dto.setId(menu.getId());
        dto.setSegmentId(menu.getSegment().getId());
        dto.setSegmentName(menu.getSegment().getName());
        dto.setGeneratedAt(menu.getGeneratedAt());
        
        List<MenuItemDTO> items = menu.getMenuItems().stream()
                .map(this::convertToItemDTO)
                .collect(Collectors.toList());
        
        dto.setItems(items);
        return dto;
    }
    
    /**
     * Converts a PersonalizedMenuItem entity to a MenuItemDTO
     */
    private MenuItemDTO convertToItemDTO(PersonalizedMenuItem item) {
        MenuItemDTO dto = new MenuItemDTO();
        dto.setProductId(item.getProduct().getId());
        dto.setProductName(item.getProduct().getName());
        dto.setProductDescription(item.getProduct().getDescription());
        dto.setBasePrice(item.getProduct().getBasePrice());
        dto.setFinalPrice(item.getFinalPrice());
        dto.setDiscountApplied(item.getDiscountApplied());
        
        if (item.getPromotion() != null) {
            dto.setPromotionId(item.getPromotion().getId());
            dto.setPromotionName(item.getPromotion().getName());
        }
        
        return dto;
    }
}
