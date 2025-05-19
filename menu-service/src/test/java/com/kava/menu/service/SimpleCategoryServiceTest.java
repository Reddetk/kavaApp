package com.kava.menu.service;

import com.kava.menu.model.Category;
import com.kava.menu.repository.CategoryRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.*;

class SimpleCategoryServiceTest {

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
        MockitoAnnotations.openMocks(this);

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
    void getAllCategories_ShouldReturnAllCategories() {
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
}
