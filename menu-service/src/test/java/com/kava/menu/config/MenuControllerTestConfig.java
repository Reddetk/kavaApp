package com.kava.menu.config;

import com.kava.menu.service.PersonalizedMenuService;
import org.mockito.Mockito;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;

@TestConfiguration
public class MenuControllerTestConfig {

    @Bean
    @Primary
    public PersonalizedMenuService personalizedMenuService() {
        return Mockito.mock(PersonalizedMenuService.class);
    }
}
