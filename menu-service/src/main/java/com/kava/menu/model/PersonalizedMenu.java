package com.kava.menu.model;

import javax.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Entity
@Table(name = "personalized_menus")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class PersonalizedMenu {
    
    @Id
    private UUID id;
    
    @ManyToOne
    @JoinColumn(name = "segment_id")
    private Segment segment;
    
    @Column(name = "generated_at")
    private LocalDateTime generatedAt;
    
    @OneToMany(mappedBy = "menu", cascade = CascadeType.ALL, orphanRemoval = true)
    private Set<PersonalizedMenuItem> menuItems = new HashSet<>();
    
    @PrePersist
    protected void onCreate() {
        if (id == null) {
            id = UUID.randomUUID();
        }
        if (generatedAt == null) {
            generatedAt = LocalDateTime.now();
        }
    }
}
