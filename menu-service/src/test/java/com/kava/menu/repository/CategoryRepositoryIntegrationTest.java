package com.kava.menu.repository;

import com.kava.menu.model.Category;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.ContextConfiguration;
import com.kava.menu.MenuServiceApplication;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

@DataJpaTest
@ActiveProfiles("test")
@ContextConfiguration(classes = MenuServiceApplication.class)
class CategoryRepositoryIntegrationTest {

    @Autowired
    private CategoryRepository categoryRepository;

    @Test
    @DisplayName("Save category should persist category")
    void save_ShouldPersistCategory() {
        // Arrange
        Category category = new Category();
        category.setId(UUID.randomUUID());
        category.setName("Beverages");
        category.setDescription("All types of drinks");
        category.setCreatedAt(LocalDateTime.now());

        // Act
        Category savedCategory = categoryRepository.save(category);

        // Assert
        assertNotNull(savedCategory.getId());
        assertEquals("Beverages", savedCategory.getName());
        assertEquals("All types of drinks", savedCategory.getDescription());
        assertNotNull(savedCategory.getCreatedAt());
    }

    @Test
    @DisplayName("Find all categories should return all categories")
    void findAll_ShouldReturnAllCategories() {
        // Arrange
        Category category1 = new Category();
        category1.setId(UUID.randomUUID());
        category1.setName("Beverages");
        category1.setDescription("All types of drinks");
        category1.setCreatedAt(LocalDateTime.now());

        Category category2 = new Category();
        category2.setId(UUID.randomUUID());
        category2.setName("Food");
        category2.setDescription("All types of food");
        category2.setCreatedAt(LocalDateTime.now());

        categoryRepository.save(category1);
        categoryRepository.save(category2);

        // Act
        List<Category> categories = categoryRepository.findAll();

        // Assert
        assertEquals(2, categories.size());
        assertTrue(categories.stream().anyMatch(c -> c.getName().equals("Beverages")));
        assertTrue(categories.stream().anyMatch(c -> c.getName().equals("Food")));
    }

    @Test
    @DisplayName("Find by ID should return category when exists")
    void findById_ShouldReturnCategory_WhenExists() {
        // Arrange
        Category category = new Category();
        category.setId(UUID.randomUUID());
        category.setName("Beverages");
        category.setDescription("All types of drinks");
        category.setCreatedAt(LocalDateTime.now());

        Category savedCategory = categoryRepository.save(category);

        // Act
        Optional<Category> foundCategory = categoryRepository.findById(savedCategory.getId());

        // Assert
        assertTrue(foundCategory.isPresent());
        assertEquals("Beverages", foundCategory.get().getName());
    }

    @Test
    @DisplayName("Find by ID should return empty when not exists")
    void findById_ShouldReturnEmpty_WhenNotExists() {
        // Act
        Optional<Category> foundCategory = categoryRepository.findById(UUID.randomUUID());

        // Assert
        assertFalse(foundCategory.isPresent());
    }

    @Test
    @DisplayName("Delete should remove category")
    void delete_ShouldRemoveCategory() {
        // Arrange
        Category category = new Category();
        category.setId(UUID.randomUUID());
        category.setName("Beverages");
        category.setDescription("All types of drinks");
        category.setCreatedAt(LocalDateTime.now());

        Category savedCategory = categoryRepository.save(category);

        // Act
        categoryRepository.delete(savedCategory);
        Optional<Category> foundCategory = categoryRepository.findById(savedCategory.getId());

        // Assert
        assertFalse(foundCategory.isPresent());
    }
}
