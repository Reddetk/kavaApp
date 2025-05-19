package com.kava.menu.service;

import com.kava.menu.model.Category;
import com.kava.menu.model.PriceHistory;
import com.kava.menu.model.Product;
import com.kava.menu.repository.CategoryRepository;
import com.kava.menu.repository.PriceHistoryRepository;
import com.kava.menu.repository.ProductRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ProductServiceTest {

    @Mock
    private ProductRepository productRepository;

    @Mock
    private CategoryRepository categoryRepository;

    @Mock
    private PriceHistoryRepository priceHistoryRepository;

    @InjectMocks
    private ProductService productService;

    @Captor
    private ArgumentCaptor<PriceHistory> priceHistoryCaptor;

    private Product product1;
    private Product product2;
    private Category category;
    private UUID productId1;
    private UUID productId2;
    private UUID categoryId;

    @BeforeEach
    void setUp() {
        categoryId = UUID.randomUUID();
        productId1 = UUID.randomUUID();
        productId2 = UUID.randomUUID();

        category = new Category();
        category.setId(categoryId);
        category.setName("Beverages");
        category.setDescription("All types of drinks");

        product1 = new Product();
        product1.setId(productId1);
        product1.setName("Coffee");
        product1.setDescription("Hot coffee");
        product1.setBasePrice(new BigDecimal("3.50"));
        product1.setCategory(category);
        product1.setIsActive(true);
        product1.setCreatedAt(LocalDateTime.now());
        product1.setUpdatedAt(LocalDateTime.now());

        product2 = new Product();
        product2.setId(productId2);
        product2.setName("Tea");
        product2.setDescription("Hot tea");
        product2.setBasePrice(new BigDecimal("2.50"));
        product2.setCategory(category);
        product2.setIsActive(true);
        product2.setCreatedAt(LocalDateTime.now());
        product2.setUpdatedAt(LocalDateTime.now());
    }

    @Test
    @DisplayName("Get all products should return list of products")
    void getAllProducts_ShouldReturnListOfProducts() {
        // Arrange
        when(productRepository.findAll()).thenReturn(Arrays.asList(product1, product2));

        // Act
        List<Product> result = productService.getAllProducts();

        // Assert
        assertEquals(2, result.size());
        assertEquals("Coffee", result.get(0).getName());
        assertEquals("Tea", result.get(1).getName());
        verify(productRepository, times(1)).findAll();
    }

    @Test
    @DisplayName("Get active products should return list of active products")
    void getActiveProducts_ShouldReturnListOfActiveProducts() {
        // Arrange
        when(productRepository.findByIsActiveTrue()).thenReturn(Arrays.asList(product1, product2));

        // Act
        List<Product> result = productService.getActiveProducts();

        // Assert
        assertEquals(2, result.size());
        assertEquals("Coffee", result.get(0).getName());
        assertEquals("Tea", result.get(1).getName());
        verify(productRepository, times(1)).findByIsActiveTrue();
    }

    @Test
    @DisplayName("Get product by ID should return product when exists")
    void getProductById_ShouldReturnProduct_WhenExists() {
        // Arrange
        when(productRepository.findById(productId1)).thenReturn(Optional.of(product1));

        // Act
        Optional<Product> result = productService.getProductById(productId1);

        // Assert
        assertTrue(result.isPresent());
        assertEquals("Coffee", result.get().getName());
        verify(productRepository, times(1)).findById(productId1);
    }

    @Test
    @DisplayName("Get products by category should return list of products in category")
    void getProductsByCategory_ShouldReturnListOfProductsInCategory() {
        // Arrange
        when(categoryRepository.findById(categoryId)).thenReturn(Optional.of(category));
        when(productRepository.findByCategory(category)).thenReturn(Arrays.asList(product1, product2));

        // Act
        List<Product> result = productService.getProductsByCategory(categoryId);

        // Assert
        assertEquals(2, result.size());
        assertEquals("Coffee", result.get(0).getName());
        assertEquals("Tea", result.get(1).getName());
        verify(categoryRepository, times(1)).findById(categoryId);
        verify(productRepository, times(1)).findByCategory(category);
    }

    @Test
    @DisplayName("Get products by category should return empty list when category not exists")
    void getProductsByCategory_ShouldReturnEmptyList_WhenCategoryNotExists() {
        // Arrange
        UUID nonExistentCategoryId = UUID.randomUUID();
        when(categoryRepository.findById(nonExistentCategoryId)).thenReturn(Optional.empty());

        // Act
        List<Product> result = productService.getProductsByCategory(nonExistentCategoryId);

        // Assert
        assertTrue(result.isEmpty());
        verify(categoryRepository, times(1)).findById(nonExistentCategoryId);
        verify(productRepository, never()).findByCategory(any(Category.class));
    }

    @Test
    @DisplayName("Create product should return saved product")
    void createProduct_ShouldReturnSavedProduct() {
        // Arrange
        Product newProduct = new Product();
        newProduct.setName("Juice");
        newProduct.setDescription("Fresh juice");
        newProduct.setBasePrice(new BigDecimal("4.00"));
        newProduct.setCategory(category);

        when(productRepository.save(any(Product.class))).thenReturn(newProduct);

        // Act
        Product result = productService.createProduct(newProduct);

        // Assert
        assertEquals("Juice", result.getName());
        assertEquals("Fresh juice", result.getDescription());
        assertEquals(new BigDecimal("4.00"), result.getBasePrice());
        verify(productRepository, times(1)).save(newProduct);
    }

    @Test
    @DisplayName("Update product price should record price history and update price")
    void updateProductPrice_ShouldRecordPriceHistoryAndUpdatePrice() {
        // Arrange
        BigDecimal newPrice = new BigDecimal("4.00");
        String reason = "Price increase";

        when(productRepository.findById(productId1)).thenReturn(Optional.of(product1));
        when(priceHistoryRepository.save(any(PriceHistory.class))).thenAnswer(invocation -> invocation.getArgument(0));
        when(productRepository.save(any(Product.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Optional<Product> result = productService.updateProductPrice(productId1, newPrice, reason);

        // Assert
        assertTrue(result.isPresent());
        assertEquals(newPrice, result.get().getBasePrice());
        
        verify(productRepository, times(1)).findById(productId1);
        verify(priceHistoryRepository, times(1)).save(priceHistoryCaptor.capture());
        verify(productRepository, times(1)).save(product1);
        
        PriceHistory capturedPriceHistory = priceHistoryCaptor.getValue();
        assertEquals(product1, capturedPriceHistory.getProduct());
        assertEquals(new BigDecimal("3.50"), capturedPriceHistory.getOldPrice());
        assertEquals(newPrice, capturedPriceHistory.getNewPrice());
        assertEquals(reason, capturedPriceHistory.getChangeReason());
    }

    @Test
    @DisplayName("Delete product should return true when exists")
    void deleteProduct_ShouldReturnTrue_WhenExists() {
        // Arrange
        when(productRepository.findById(productId1)).thenReturn(Optional.of(product1));
        doNothing().when(productRepository).delete(any(Product.class));

        // Act
        boolean result = productService.deleteProduct(productId1);

        // Assert
        assertTrue(result);
        verify(productRepository, times(1)).findById(productId1);
        verify(productRepository, times(1)).delete(product1);
    }

    @Test
    @DisplayName("Deactivate product should set isActive to false")
    void deactivateProduct_ShouldSetIsActiveToFalse() {
        // Arrange
        when(productRepository.findById(productId1)).thenReturn(Optional.of(product1));
        when(productRepository.save(any(Product.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Optional<Product> result = productService.deactivateProduct(productId1);

        // Assert
        assertTrue(result.isPresent());
        assertFalse(result.get().getIsActive());
        verify(productRepository, times(1)).findById(productId1);
        verify(productRepository, times(1)).save(product1);
    }

    @Test
    @DisplayName("Activate product should set isActive to true")
    void activateProduct_ShouldSetIsActiveToTrue() {
        // Arrange
        product1.setIsActive(false);
        when(productRepository.findById(productId1)).thenReturn(Optional.of(product1));
        when(productRepository.save(any(Product.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Optional<Product> result = productService.activateProduct(productId1);

        // Assert
        assertTrue(result.isPresent());
        assertTrue(result.get().getIsActive());
        verify(productRepository, times(1)).findById(productId1);
        verify(productRepository, times(1)).save(product1);
    }
}
