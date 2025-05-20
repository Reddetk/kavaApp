package com.kava.menu.config;

import com.kava.menu.service.CategoryService;
import org.mockito.Mockito;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;

@TestConfiguration
public class CategoryControllerTestConfig {

    @Bean
    @Primary
    public CategoryService categoryService() {
        return Mockito.mock(CategoryService.class);
    }
}
