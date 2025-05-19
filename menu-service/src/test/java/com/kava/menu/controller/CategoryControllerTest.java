package com.kava.menu.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.kava.menu.model.Category;
import com.kava.menu.service.CategoryService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.mockito.junit.jupiter.MockitoExtension;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.Optional;
import java.util.UUID;

import static org.hamcrest.Matchers.hasSize;
import static org.hamcrest.Matchers.is;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@ExtendWith(MockitoExtension.class)
class CategoryControllerTest {

    private MockMvc mockMvc;

    private ObjectMapper objectMapper = new ObjectMapper();

    @Mock
    private CategoryService categoryService;

    @InjectMocks
    private CategoryController categoryController;

    private Category category1;
    private Category category2;
    private UUID categoryId1;
    private UUID categoryId2;

    @BeforeEach
    void setUp() {
        mockMvc = MockMvcBuilders.standaloneSetup(categoryController).build();

        categoryId1 = UUID.randomUUID();
        categoryId2 = UUID.randomUUID();

        category1 = new Category();
        category1.setId(categoryId1);
        category1.setName("Beverages");
        category1.setDescription("All types of drinks");
        category1.setCreatedAt(LocalDateTime.now());

        category2 = new Category();
        category2.setId(categoryId2);
        category2.setName("Food");
        category2.setDescription("All types of food");
        category2.setCreatedAt(LocalDateTime.now());
    }

    @Test
    @DisplayName("GET /api/categories should return all categories")
    void getAllCategories_ShouldReturnAllCategories() throws Exception {
        // Arrange
        when(categoryService.getAllCategories()).thenReturn(Arrays.asList(category1, category2));

        // Act & Assert
        mockMvc.perform(get("/api/categories"))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$", hasSize(2)))
                .andExpect(jsonPath("$[0].id", is(categoryId1.toString())))
                .andExpect(jsonPath("$[0].name", is("Beverages")))
                .andExpect(jsonPath("$[1].id", is(categoryId2.toString())))
                .andExpect(jsonPath("$[1].name", is("Food")));

        verify(categoryService, times(1)).getAllCategories();
    }

    @Test
    @DisplayName("GET /api/categories/{id} should return category when exists")
    void getCategoryById_ShouldReturnCategory_WhenExists() throws Exception {
        // Arrange
        when(categoryService.getCategoryById(categoryId1)).thenReturn(Optional.of(category1));

        // Act & Assert
        mockMvc.perform(get("/api/categories/{id}", categoryId1))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(categoryId1.toString())))
                .andExpect(jsonPath("$.name", is("Beverages")))
                .andExpect(jsonPath("$.description", is("All types of drinks")));

        verify(categoryService, times(1)).getCategoryById(categoryId1);
    }

    @Test
    @DisplayName("GET /api/categories/{id} should return 404 when not exists")
    void getCategoryById_ShouldReturn404_WhenNotExists() throws Exception {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(categoryService.getCategoryById(nonExistentId)).thenReturn(Optional.empty());

        // Act & Assert
        mockMvc.perform(get("/api/categories/{id}", nonExistentId))
                .andExpect(status().isNotFound());

        verify(categoryService, times(1)).getCategoryById(nonExistentId);
    }

    @Test
    @DisplayName("POST /api/categories should create and return new category")
    void createCategory_ShouldCreateAndReturnNewCategory() throws Exception {
        // Arrange
        Category newCategory = new Category();
        newCategory.setName("Desserts");
        newCategory.setDescription("Sweet treats");

        Category savedCategory = new Category();
        savedCategory.setId(UUID.randomUUID());
        savedCategory.setName("Desserts");
        savedCategory.setDescription("Sweet treats");
        savedCategory.setCreatedAt(LocalDateTime.now());

        when(categoryService.createCategory(any(Category.class))).thenReturn(savedCategory);

        // Act & Assert
        mockMvc.perform(post("/api/categories")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(newCategory)))
                .andExpect(status().isCreated())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(savedCategory.getId().toString())))
                .andExpect(jsonPath("$.name", is("Desserts")))
                .andExpect(jsonPath("$.description", is("Sweet treats")));

        verify(categoryService, times(1)).createCategory(any(Category.class));
    }

    @Test
    @DisplayName("PUT /api/categories/{id} should update and return category when exists")
    void updateCategory_ShouldUpdateAndReturnCategory_WhenExists() throws Exception {
        // Arrange
        Category updatedCategory = new Category();
        updatedCategory.setName("Updated Beverages");
        updatedCategory.setDescription("Updated description");

        Category savedCategory = new Category();
        savedCategory.setId(categoryId1);
        savedCategory.setName("Updated Beverages");
        savedCategory.setDescription("Updated description");
        savedCategory.setCreatedAt(LocalDateTime.now());

        when(categoryService.updateCategory(eq(categoryId1), any(Category.class))).thenReturn(Optional.of(savedCategory));

        // Act & Assert
        mockMvc.perform(put("/api/categories/{id}", categoryId1)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(updatedCategory)))
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(categoryId1.toString())))
                .andExpect(jsonPath("$.name", is("Updated Beverages")))
                .andExpect(jsonPath("$.description", is("Updated description")));

        verify(categoryService, times(1)).updateCategory(eq(categoryId1), any(Category.class));
    }

    @Test
    @DisplayName("PUT /api/categories/{id} should return 404 when not exists")
    void updateCategory_ShouldReturn404_WhenNotExists() throws Exception {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        Category updatedCategory = new Category();
        updatedCategory.setName("Updated Beverages");
        updatedCategory.setDescription("Updated description");

        when(categoryService.updateCategory(eq(nonExistentId), any(Category.class))).thenReturn(Optional.empty());

        // Act & Assert
        mockMvc.perform(put("/api/categories/{id}", nonExistentId)
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(updatedCategory)))
                .andExpect(status().isNotFound());

        verify(categoryService, times(1)).updateCategory(eq(nonExistentId), any(Category.class));
    }

    @Test
    @DisplayName("DELETE /api/categories/{id} should return 204 when exists")
    void deleteCategory_ShouldReturn204_WhenExists() throws Exception {
        // Arrange
        when(categoryService.deleteCategory(categoryId1)).thenReturn(true);

        // Act & Assert
        mockMvc.perform(delete("/api/categories/{id}", categoryId1))
                .andExpect(status().isNoContent());

        verify(categoryService, times(1)).deleteCategory(categoryId1);
    }

    @Test
    @DisplayName("DELETE /api/categories/{id} should return 404 when not exists")
    void deleteCategory_ShouldReturn404_WhenNotExists() throws Exception {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(categoryService.deleteCategory(nonExistentId)).thenReturn(false);

        // Act & Assert
        mockMvc.perform(delete("/api/categories/{id}", nonExistentId))
                .andExpect(status().isNotFound());

        verify(categoryService, times(1)).deleteCategory(nonExistentId);
    }
}
