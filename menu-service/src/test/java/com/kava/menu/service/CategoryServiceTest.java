package com.kava.menu.service;

import com.kava.menu.model.Category;
import com.kava.menu.repository.CategoryRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CategoryServiceTest {

    @Mock
    private CategoryRepository categoryRepository;

    @InjectMocks
    private CategoryService categoryService;

    private Category category1;
    private Category category2;
    private UUID categoryId1;
    private UUID categoryId2;

    @BeforeEach
    void setUp() {
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
    @DisplayName("Get all categories should return list of categories")
    void getAllCategories_ShouldReturnListOfCategories() {
        // Arrange
        when(categoryRepository.findAll()).thenReturn(Arrays.asList(category1, category2));

        // Act
        List<Category> result = categoryService.getAllCategories();

        // Assert
        assertEquals(2, result.size());
        assertEquals("Beverages", result.get(0).getName());
        assertEquals("Food", result.get(1).getName());
        verify(categoryRepository, times(1)).findAll();
    }

    @Test
    @DisplayName("Get category by ID should return category when exists")
    void getCategoryById_ShouldReturnCategory_WhenExists() {
        // Arrange
        when(categoryRepository.findById(categoryId1)).thenReturn(Optional.of(category1));

        // Act
        Optional<Category> result = categoryService.getCategoryById(categoryId1);

        // Assert
        assertTrue(result.isPresent());
        assertEquals("Beverages", result.get().getName());
        verify(categoryRepository, times(1)).findById(categoryId1);
    }

    @Test
    @DisplayName("Get category by ID should return empty when not exists")
    void getCategoryById_ShouldReturnEmpty_WhenNotExists() {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(categoryRepository.findById(nonExistentId)).thenReturn(Optional.empty());

        // Act
        Optional<Category> result = categoryService.getCategoryById(nonExistentId);

        // Assert
        assertFalse(result.isPresent());
        verify(categoryRepository, times(1)).findById(nonExistentId);
    }

    @Test
    @DisplayName("Create category should return saved category")
    void createCategory_ShouldReturnSavedCategory() {
        // Arrange
        Category newCategory = new Category();
        newCategory.setName("Desserts");
        newCategory.setDescription("Sweet treats");

        when(categoryRepository.save(any(Category.class))).thenReturn(newCategory);

        // Act
        Category result = categoryService.createCategory(newCategory);

        // Assert
        assertEquals("Desserts", result.getName());
        assertEquals("Sweet treats", result.getDescription());
        verify(categoryRepository, times(1)).save(newCategory);
    }

    @Test
    @DisplayName("Update category should return updated category when exists")
    void updateCategory_ShouldReturnUpdatedCategory_WhenExists() {
        // Arrange
        Category updatedDetails = new Category();
        updatedDetails.setName("Updated Beverages");
        updatedDetails.setDescription("Updated description");

        when(categoryRepository.findById(categoryId1)).thenReturn(Optional.of(category1));
        when(categoryRepository.save(any(Category.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Optional<Category> result = categoryService.updateCategory(categoryId1, updatedDetails);

        // Assert
        assertTrue(result.isPresent());
        assertEquals("Updated Beverages", result.get().getName());
        assertEquals("Updated description", result.get().getDescription());
        verify(categoryRepository, times(1)).findById(categoryId1);
        verify(categoryRepository, times(1)).save(any(Category.class));
    }

    @Test
    @DisplayName("Update category should return empty when not exists")
    void updateCategory_ShouldReturnEmpty_WhenNotExists() {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        Category updatedDetails = new Category();
        updatedDetails.setName("Updated Beverages");

        when(categoryRepository.findById(nonExistentId)).thenReturn(Optional.empty());

        // Act
        Optional<Category> result = categoryService.updateCategory(nonExistentId, updatedDetails);

        // Assert
        assertFalse(result.isPresent());
        verify(categoryRepository, times(1)).findById(nonExistentId);
        verify(categoryRepository, never()).save(any(Category.class));
    }

    @Test
    @DisplayName("Delete category should return true when exists")
    void deleteCategory_ShouldReturnTrue_WhenExists() {
        // Arrange
        when(categoryRepository.findById(categoryId1)).thenReturn(Optional.of(category1));
        doNothing().when(categoryRepository).delete(any(Category.class));

        // Act
        boolean result = categoryService.deleteCategory(categoryId1);

        // Assert
        assertTrue(result);
        verify(categoryRepository, times(1)).findById(categoryId1);
        verify(categoryRepository, times(1)).delete(category1);
    }

    @Test
    @DisplayName("Delete category should return false when not exists")
    void deleteCategory_ShouldReturnFalse_WhenNotExists() {
        // Arrange
        UUID nonExistentId = UUID.randomUUID();
        when(categoryRepository.findById(nonExistentId)).thenReturn(Optional.empty());

        // Act
        boolean result = categoryService.deleteCategory(nonExistentId);

        // Assert
        assertFalse(result);
        verify(categoryRepository, times(1)).findById(nonExistentId);
        verify(categoryRepository, never()).delete(any(Category.class));
    }
}
