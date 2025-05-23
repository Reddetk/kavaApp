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
public class ProductDemandMetricsId implements Serializable {
    
    private UUID productId;
    private UUID segmentId;
}
