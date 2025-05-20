package com.kava.menu.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.kava.menu.dto.MenuItemDTO;
import com.kava.menu.dto.MenuRequestDTO;
import com.kava.menu.dto.MenuResponseDTO;
import com.kava.menu.model.PersonalizedMenu;
import com.kava.menu.model.PersonalizedMenuItem;
import com.kava.menu.model.Product;
import com.kava.menu.model.Segment;
import com.kava.menu.service.PersonalizedMenuService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentMatchers;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.mockito.junit.jupiter.MockitoExtension;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.*;

import static org.hamcrest.Matchers.hasSize;
import static org.hamcrest.Matchers.is;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@ExtendWith(MockitoExtension.class)
class MenuControllerTest {

    private MockMvc mockMvc;

    private ObjectMapper objectMapper = new ObjectMapper();

    @Mock
    private PersonalizedMenuService menuService;

    @InjectMocks
    private MenuController menuController;

    private PersonalizedMenu menu;
    private UUID menuId;
    private UUID segmentId;
    private Segment segment;
    private Set<PersonalizedMenuItem> menuItems;

    @BeforeEach
    void setUp() {
        mockMvc = MockMvcBuilders.standaloneSetup(menuController).build();

        menuId = UUID.randomUUID();
        segmentId = UUID.randomUUID();

        segment = new Segment();
        segment.setId(segmentId);
        segment.setName("Premium");
        segment.setDescription("Premium customers");

        menu = new PersonalizedMenu();
        menu.setId(menuId);
        menu.setSegment(segment);
        menu.setGeneratedAt(LocalDateTime.now());

        menuItems = new HashSet<>();
        // We would normally add menu items here, but for simplicity in the test,
        // we'll mock the conversion to DTO in the controller
    }

    @Test
    @DisplayName("GET /api/menus should return all menus")
    void getAllMenus_ShouldReturnAllMenus() throws Exception {
        // Arrange
        when(menuService.getAllMenus()).thenReturn(Collections.singletonList(menu));

        // Act & Assert
        mockMvc.perform(get("/api/menus"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$", hasSize(1)))
                .andExpect(jsonPath("$[0].id", is(menuId.toString())))
                .andExpect(jsonPath("$[0].segmentId", is(segmentId.toString())))
                .andExpect(jsonPath("$[0].segmentName", is("Premium")));

        verify(menuService, times(1)).getAllMenus();
    }

    @Test
    @DisplayName("GET /api/menus/{id} should return menu when exists")
    void getMenuById_ShouldReturnMenu_WhenExists() throws Exception {
        // Arrange
        when(menuService.getMenuById(menuId)).thenReturn(Optional.of(menu));

        // Act & Assert
        mockMvc.perform(get("/api/menus/{id}", menuId))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(menuId.toString())))
                .andExpect(jsonPath("$.segmentId", is(segmentId.toString())))
                .andExpect(jsonPath("$.segmentName", is("Premium")));

        verify(menuService, times(1)).getMenuById(menuId);
    }

    @Test
    @DisplayName("GET /api/menus/{id} should return 404 when not exists")
    void getMenuById_ShouldReturn404_WhenNotExists() throws Exception {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(menuService.getMenuById(nonExistentId)).thenReturn(Optional.empty());

        // Act & Assert
        mockMvc.perform(get("/api/menus/{id}", nonExistentId))
                .andExpect(status().isNotFound());

        verify(menuService, times(1)).getMenuById(nonExistentId);
    }

    @Test
    @DisplayName("GET /api/menus/segment/{segmentId} should return menus for segment")
    void getMenusBySegment_ShouldReturnMenusForSegment() throws Exception {
        // Arrange
        when(menuService.getMenusBySegment(segmentId)).thenReturn(Collections.singletonList(menu));

        // Act & Assert
        mockMvc.perform(get("/api/menus/segment/{segmentId}", segmentId))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$", hasSize(1)))
                .andExpect(jsonPath("$[0].id", is(menuId.toString())))
                .andExpect(jsonPath("$[0].segmentId", is(segmentId.toString())))
                .andExpect(jsonPath("$[0].segmentName", is("Premium")));

        verify(menuService, times(1)).getMenusBySegment(segmentId);
    }

    @Test
    @DisplayName("GET /api/menus/segment/{segmentId}/latest should return latest menu for segment")
    void getLatestMenuForSegment_ShouldReturnLatestMenu() throws Exception {
        // Arrange
        when(menuService.getLatestMenuForSegment(segmentId)).thenReturn(Optional.of(menu));

        // Act & Assert
        mockMvc.perform(get("/api/menus/segment/{segmentId}/latest", segmentId))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(menuId.toString())))
                .andExpect(jsonPath("$.segmentId", is(segmentId.toString())))
                .andExpect(jsonPath("$.segmentName", is("Premium")));

        verify(menuService, times(1)).getLatestMenuForSegment(segmentId);
    }

    @Test
    @DisplayName("POST /api/menus/generate should generate and return new menu")
    void generateMenu_ShouldGenerateAndReturnNewMenu() throws Exception {
        // Arrange
        MenuRequestDTO request = new MenuRequestDTO();
        request.setSegmentId(segmentId);
        request.setLatitude(new BigDecimal("37.7749"));
        request.setLongitude(new BigDecimal("-122.4194"));
        request.setRegionCode("US-CA");

        when(menuService.generateMenuForSegment(segmentId)).thenReturn(Optional.of(menu));

        // Act & Assert
        mockMvc.perform(post("/api/menus/generate")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(menuId.toString())))
                .andExpect(jsonPath("$.segmentId", is(segmentId.toString())))
                .andExpect(jsonPath("$.segmentName", is("Premium")));

        verify(menuService, times(1)).generateMenuForSegment(segmentId);
    }

    @Test
    @DisplayName("DELETE /api/menus/{id} should return 204 when exists")
    void deleteMenu_ShouldReturn204_WhenExists() throws Exception {
        // Arrange
        when(menuService.deleteMenu(menuId)).thenReturn(true);

        // Act & Assert
        mockMvc.perform(delete("/api/menus/{id}", menuId))
                .andExpect(status().isNoContent());

        verify(menuService, times(1)).deleteMenu(menuId);
    }

    @Test
    @DisplayName("DELETE /api/menus/{id} should return 404 when not exists")
    void deleteMenu_ShouldReturn404_WhenNotExists() throws Exception {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(menuService.deleteMenu(nonExistentId)).thenReturn(false);

        // Act & Assert
        mockMvc.perform(delete("/api/menus/{id}", nonExistentId))
                .andExpect(status().isNotFound());

        verify(menuService, times(1)).deleteMenu(nonExistentId);
    }
}
