package com.kava.menu.service;

import com.kava.menu.model.Category;
import com.kava.menu.model.PriceHistory;
import com.kava.menu.model.Product;
import com.kava.menu.repository.CategoryRepository;
import com.kava.menu.repository.PriceHistoryRepository;
import com.kava.menu.repository.ProductRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ProductService {
    
    private final ProductRepository productRepository;
    private final CategoryRepository categoryRepository;
    private final PriceHistoryRepository priceHistoryRepository;
    
    public List<Product> getAllProducts() {
        return productRepository.findAll();
    }
    
    public List<Product> getActiveProducts() {
        return productRepository.findByIsActiveTrue();
    }
    
    public Optional<Product> getProductById(UUID id) {
        return productRepository.findById(id);
    }
    
    public List<Product> getProductsByCategory(UUID categoryId) {
        return categoryRepository.findById(categoryId)
                .map(productRepository::findByCategory)
                .orElse(List.of());
    }
    
    public List<Product> getActiveProductsByCategory(UUID categoryId) {
        return categoryRepository.findById(categoryId)
                .map(productRepository::findByCategoryAndIsActiveTrue)
                .orElse(List.of());
    }
    
    @Transactional
    public Product createProduct(Product product) {
        return productRepository.save(product);
    }
    
    @Transactional
    public Optional<Product> updateProduct(UUID id, Product productDetails) {
        return productRepository.findById(id)
                .map(existingProduct -> {
                    // Check if price has changed
                    if (existingProduct.getBasePrice().compareTo(productDetails.getBasePrice()) != 0) {
                        // Record price history
                        PriceHistory priceHistory = new PriceHistory();
                        priceHistory.setProduct(existingProduct);
                        priceHistory.setOldPrice(existingProduct.getBasePrice());
                        priceHistory.setNewPrice(productDetails.getBasePrice());
                        priceHistory.setChangeReason("Manual update");
                        priceHistoryRepository.save(priceHistory);
                    }
                    
                    // Update product details
                    existingProduct.setName(productDetails.getName());
                    existingProduct.setDescription(productDetails.getDescription());
                    existingProduct.setBasePrice(productDetails.getBasePrice());
                    existingProduct.setCategory(productDetails.getCategory());
                    existingProduct.setIsActive(productDetails.getIsActive());
                    
                    return productRepository.save(existingProduct);
                });
    }
    
    @Transactional
    public Optional<Product> updateProductPrice(UUID id, BigDecimal newPrice, String reason) {
        return productRepository.findById(id)
                .map(existingProduct -> {
                    // Record price history
                    PriceHistory priceHistory = new PriceHistory();
                    priceHistory.setProduct(existingProduct);
                    priceHistory.setOldPrice(existingProduct.getBasePrice());
                    priceHistory.setNewPrice(newPrice);
                    priceHistory.setChangeReason(reason);
                    priceHistoryRepository.save(priceHistory);
                    
                    // Update product price
                    existingProduct.setBasePrice(newPrice);
                    return productRepository.save(existingProduct);
                });
    }
    
    @Transactional
    public boolean deleteProduct(UUID id) {
        return productRepository.findById(id)
                .map(product -> {
                    productRepository.delete(product);
                    return true;
                })
                .orElse(false);
    }
    
    @Transactional
    public Optional<Product> deactivateProduct(UUID id) {
        return productRepository.findById(id)
                .map(existingProduct -> {
                    existingProduct.setIsActive(false);
                    return productRepository.save(existingProduct);
                });
    }
    
    @Transactional
    public Optional<Product> activateProduct(UUID id) {
        return productRepository.findById(id)
                .map(existingProduct -> {
                    existingProduct.setIsActive(true);
                    return productRepository.save(existingProduct);
                });
    }
}
