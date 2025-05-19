package com.kava.menu.model;

import jakarta.persistence.Embeddable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.UUID;

@Embeddable
@Data
@NoArgsConstructor
@AllArgsConstructor
public class PersonalizedMenuItemId implements Serializable {
    
    private UUID menuId;
    private UUID productId;
}
