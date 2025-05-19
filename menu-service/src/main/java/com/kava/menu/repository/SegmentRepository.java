package com.kava.menu.repository;

import com.kava.menu.model.Segment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface SegmentRepository extends JpaRepository<Segment, UUID> {
}
