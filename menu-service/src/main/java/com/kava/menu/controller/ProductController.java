package com.kava.menu.controller;

import com.kava.menu.model.Product;
import com.kava.menu.service.ProductService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
@Tag(name = "Products", description = "Product management endpoints")
public class ProductController {
    
    private final ProductService productService;
    
    @GetMapping
    @Operation(summary = "Get all products", description = "Retrieves a list of all products")
    public ResponseEntity<List<Product>> getAllProducts() {
        return ResponseEntity.ok(productService.getAllProducts());
    }
    
    @GetMapping("/active")
    @Operation(summary = "Get active products", description = "Retrieves a list of all active products")
    public ResponseEntity<List<Product>> getActiveProducts() {
        return ResponseEntity.ok(productService.getActiveProducts());
    }
    
    @GetMapping("/{id}")
    @Operation(summary = "Get product by ID", description = "Retrieves a product by its UUID")
    public ResponseEntity<Product> getProductById(@PathVariable UUID id) {
        return productService.getProductById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @GetMapping("/category/{categoryId}")
    @Operation(summary = "Get products by category", description = "Retrieves all products in a specific category")
    public ResponseEntity<List<Product>> getProductsByCategory(@PathVariable UUID categoryId) {
        return ResponseEntity.ok(productService.getProductsByCategory(categoryId));
    }
    
    @GetMapping("/category/{categoryId}/active")
    @Operation(summary = "Get active products by category", description = "Retrieves active products in a specific category")
    public ResponseEntity<List<Product>> getActiveProductsByCategory(@PathVariable UUID categoryId) {
        return ResponseEntity.ok(productService.getActiveProductsByCategory(categoryId));
    }
    
    @PostMapping
    @Operation(summary = "Create a new product", description = "Creates a new product")
    public ResponseEntity<Product> createProduct(@RequestBody Product product) {
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(productService.createProduct(product));
    }
    
    @PutMapping("/{id}")
    @Operation(summary = "Update a product", description = "Updates an existing product by its UUID")
    public ResponseEntity<Product> updateProduct(@PathVariable UUID id, @RequestBody Product product) {
        return productService.updateProduct(id, product)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PatchMapping("/{id}/price")
    @Operation(summary = "Update product price", description = "Updates only the price of an existing product")
    public ResponseEntity<Product> updateProductPrice(
            @PathVariable UUID id,
            @RequestParam BigDecimal price,
            @RequestParam(required = false, defaultValue = "Manual price update") String reason) {
        return productService.updateProductPrice(id, price, reason)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PatchMapping("/{id}/deactivate")
    @Operation(summary = "Deactivate a product", description = "Marks a product as inactive")
    public ResponseEntity<Product> deactivateProduct(@PathVariable UUID id) {
        return productService.deactivateProduct(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @PatchMapping("/{id}/activate")
    @Operation(summary = "Activate a product", description = "Marks a product as active")
    public ResponseEntity<Product> activateProduct(@PathVariable UUID id) {
        return productService.activateProduct(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
    
    @DeleteMapping("/{id}")
    @Operation(summary = "Delete a product", description = "Deletes a product by its UUID")
    public ResponseEntity<Void> deleteProduct(@PathVariable UUID id) {
        boolean deleted = productService.deleteProduct(id);
        return deleted ? ResponseEntity.noContent().build() : ResponseEntity.notFound().build();
    }
}
