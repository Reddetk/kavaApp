package com.kava.menu.repository;

import com.kava.menu.model.PersonalizedMenu;
import com.kava.menu.model.Segment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface PersonalizedMenuRepository extends JpaRepository<PersonalizedMenu, UUID> {
    
    List<PersonalizedMenu> findBySegment(Segment segment);
    
    @Query("SELECT pm FROM PersonalizedMenu pm WHERE pm.segment.id = :segmentId ORDER BY pm.generatedAt DESC")
    List<PersonalizedMenu> findLatestMenusBySegment(UUID segmentId);
    
    @Query("SELECT pm FROM PersonalizedMenu pm WHERE pm.segment.id = :segmentId AND pm.generatedAt >= :since ORDER BY pm.generatedAt DESC")
    List<PersonalizedMenu> findMenusBySegmentSince(UUID segmentId, LocalDateTime since);
    
    @Query("SELECT pm FROM PersonalizedMenu pm WHERE pm.segment.id = :segmentId ORDER BY pm.generatedAt DESC LIMIT 1")
    Optional<PersonalizedMenu> findLatestMenuForSegment(UUID segmentId);
}
