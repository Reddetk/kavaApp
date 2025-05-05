package services

import (
	"errors"
	"math"
	"math/rand"
	"time"
	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

type KMeansSegmentation struct {
	numClusters int
}

func NewKMeansSegmentation(numClusters int) *KMeansSegmentation {
	return &KMeansSegmentation{
		numClusters: numClusters,
	}
}

func (s *KMeansSegmentation) PerformRFMClustering(users []entities.UserMetrics) ([]entities.Segment, error) {
	if len(users) == 0 {
		return nil, errors.New("no user metrics provided")
	}

	points := make([][]float64, len(users))
	for i, u := range users {
		points[i] = []float64{
			float64(u.Recency),
			float64(u.Frequency),
			u.Monetary,
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := initializeCentroids(points, s.numClusters, r)

	assignments := make([]int, len(points))
	for iter := 0; iter < 100; iter++ {
		for i, p := range points {
			assignments[i] = closestCentroid(p, centroids)
		}

		newCentroids := updateCentroids(points, assignments, s.numClusters)
		if converged(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}

	segments := make([]entities.Segment, s.numClusters)
	for i := 0; i < s.numClusters; i++ {
		segments[i] = entities.Segment{
			ID:           uuid.New(),
			Name:         generateSegmentName(i),
			Type:         "RFM",
			Algorithm:    "KMeans",
			CentroidData: map[string]interface{}{"values": centroids[i]},
			CreatedAt:    time.Now(),
		}
	}

	return segments, nil
}

func (s *KMeansSegmentation) PerformBehaviorClustering(transactions []entities.Transaction) ([]entities.Segment, error) {
	if len(transactions) == 0 {
		return nil, errors.New("no transactions provided")
	}

	// Aggregate transactions by user and calculate behavior metrics
	userBehaviors := make(map[uuid.UUID][]float64)
	for _, t := range transactions {
		if _, exists := userBehaviors[t.UserID]; !exists {
			userBehaviors[t.UserID] = make([]float64, 3)
		}
		// Example behavior metrics:
		// 1. Average transaction amount
		userBehaviors[t.UserID][0] += t.Amount
		// 2. Transaction frequency (count)
		userBehaviors[t.UserID][1]++
		// 3. Time of day preference (0-24 hours)
		userBehaviors[t.UserID][2] += float64(t.Timestamp.Hour())
	}

	// Convert to points for clustering
	points := make([][]float64, 0, len(userBehaviors))
	for _, metrics := range userBehaviors {
		// Normalize metrics
		if metrics[1] > 0 { // If there are transactions
			metrics[0] /= metrics[1] // Average amount
			metrics[2] /= metrics[1] // Average hour
		}
		points = append(points, metrics)
	}

	// Initialize random seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Perform k-means clustering
	centroids := initializeCentroids(points, s.numClusters, r)
	assignments := make([]int, len(points))

	// Iterate until convergence
	for iter := 0; iter < 100; iter++ {
		for i, p := range points {
			assignments[i] = closestCentroid(p, centroids)
		}

		newCentroids := updateCentroids(points, assignments, s.numClusters)
		if converged(centroids, newCentroids) {
			break
		}
		centroids = newCentroids
	}

	// Create segments
	segments := make([]entities.Segment, s.numClusters)
	for i := 0; i < s.numClusters; i++ {
		segments[i] = entities.Segment{
			ID:        uuid.New(),
			Name:      "Behavior-" + generateSegmentName(i),
			Type:      "Behavioral",
			Algorithm: "KMeans",
			CentroidData: map[string]interface{}{
				"values": centroids[i],
				"metrics": []string{
					"avg_transaction_amount",
					"transaction_frequency",
					"preferred_time",
				},
			},
			CreatedAt: time.Now(),
		}
	}

	return segments, nil
}

func (s *KMeansSegmentation) AssignUserToSegment(userID uuid.UUID, metrics entities.UserMetrics) (entities.Segment, error) {
	// Convert user metrics to point
	point := []float64{
		float64(metrics.Recency),
		float64(metrics.Frequency),
		metrics.Monetary,
	}

	// Create mock centroids for testing
	centroids := make([][]float64, s.numClusters)
	for i := 0; i < s.numClusters; i++ {
		centroids[i] = make([]float64, 3)
	}

	// Find closest centroid
	closestIdx := closestCentroid(point, centroids)

	// Create segment for closest centroid
	segment := entities.Segment{
		ID:           uuid.New(),
		Name:         generateSegmentName(closestIdx),
		Type:         "RFM",
		Algorithm:    "KMeans",
		CentroidData: map[string]interface{}{"values": centroids[closestIdx]},
		CreatedAt:    time.Now(),
	}

	return segment, nil
}

func generateSegmentName(index int) string {
	return string(rune('A' + index))
}

func initializeCentroids(points [][]float64, k int, r *rand.Rand) [][]float64 {
	centroids := make([][]float64, k)
	perm := r.Perm(len(points))
	for i := 0; i < k; i++ {
		if i < len(perm) {
			centroids[i] = append([]float64{}, points[perm[i]]...)
		} else {
			centroids[i] = make([]float64, len(points[0]))
		}
	}
	return centroids
}

func closestCentroid(point []float64, centroids [][]float64) int {
	minDist := math.MaxFloat64
	minIdx := 0
	for i, c := range centroids {
		dist := euclideanDistance(point, c)
		if dist < minDist {
			minDist = dist
			minIdx = i
		}
	}
	return minIdx
}

func updateCentroids(points [][]float64, assignments []int, k int) [][]float64 {
	centroids := make([][]float64, k)
	counts := make([]int, k)

	for i := range centroids {
		centroids[i] = make([]float64, len(points[0]))
	}

	for i, p := range points {
		cluster := assignments[i]
		for j := range p {
			centroids[cluster][j] += p[j]
		}
		counts[cluster]++
	}

	for i := range centroids {
		if counts[i] > 0 {
			for j := range centroids[i] {
				centroids[i][j] /= float64(counts[i])
			}
		}
	}
	return centroids
}

func euclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func converged(old, new [][]float64) bool {
	const epsilon = 1e-6
	for i := range old {
		for j := range old[i] {
			if math.Abs(old[i][j]-new[i][j]) > epsilon {
				return false
			}
		}
	}
	return true
}
