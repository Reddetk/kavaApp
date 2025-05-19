package com.kava.menu.repository;

import com.kava.menu.model.Category;
import com.kava.menu.model.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface ProductRepository extends JpaRepository<Product, UUID> {
    
    List<Product> findByCategory(Category category);
    
    List<Product> findByCategoryAndIsActiveTrue(Category category);
    
    List<Product> findByIsActiveTrue();
}
