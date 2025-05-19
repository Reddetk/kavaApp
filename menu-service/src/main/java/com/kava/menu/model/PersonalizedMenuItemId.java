package com.kava.menu.model;

import javax.persistence.*;
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
