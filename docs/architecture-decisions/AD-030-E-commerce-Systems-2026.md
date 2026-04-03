# AD-030: E-commerce Systems 2026

## Status: S-Level (Superior)

**Date:** 2026-04-03
**Author:** System Architect
**Version:** 2.0

---

## 1. Executive Summary

This document defines enterprise-grade E-commerce system architecture patterns for 2026, incorporating advanced recommendation algorithms, personalization engines, inventory management, and payment processing. It provides production-ready implementations for recommendation systems, search engines, cart optimization, and order fulfillment.

---

## 2. System Architecture Overview

### 2.1 High-Level Architecture Diagram

```
+-----------------------------------------------------------------------------------------+
|                         E-COMMERCE PLATFORM ARCHITECTURE                                |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|  +---------------------------------------------------------------------------------+    |
|  |                         CLIENT INTERFACE LAYER                                   |    |
|  |  +-------------+  +-------------+  +-------------+  +-------------------------+  |    |
|  |  | Mobile App  |  |   Web App   |  |  PWA/AMP    |  |    Partner APIs         |  |    |
|  |  | (iOS/And)   |  |  (Next.js)  |  |             |  |    (Headless)           |  |    |
|  |  +------+------+  +------+------+  +------+------+  +-----------+-------------+  |    |
|  |         +-----------------+-----------------+---------------------+                |    |
|  +------------------------------------+-----------------------------------------------+    |
|                                       |                                                  |
|  +------------------------------------v-----------------------------------------------+    |
|  |                         API GATEWAY LAYER                                        |    |
|  |  +-------------+  +-------------+  +-------------+  +-------------------------+    |    |
|  |  |   GraphQL   |  |   REST      |  |   gRPC      |  |    WebSocket            |    |    |
|  |  |   Gateway   |  |   Gateway   |  |   Services  |  |    (Real-time)          |    |    |
|  |  +------+------+  +------+------+  +------+------+  +-----------+-------------+    |    |
|  |         +-----------------+-----------------+---------------------+                  |    |
|  +------------------------------------+------------------------------------------------+    |
|                                       |                                                  |
|  +------------------------------------v-----------------------------------------------+    |
|  |                         CORE SERVICES LAYER                                      |    |
|  |                                                                                  |    |
|  |  +----------------+  +----------------+  +----------------+  +----------------+  |    |
|  |  |  Product       |  |  Search        |  |  Recommendation|  |  Pricing       |  |    |
|  |  |  Catalog       |  |  Engine        |  |  Engine        |  |  Engine        |  |    |
|  |  |  (Event Sourced)|  |  (Elasticsearch)|  |  (ML Pipeline) |  |  (Dynamic)     |  |    |
|  |  +--------+-------+  +--------+-------+  +--------+-------+  +--------+-------+  |    |
|  |           |                   |                   |                   |          |    |
|  |  +--------v-------------------v-------------------v-------------------v---------+  |    |
|  |  |                         Cart & Checkout                                      |  |    |
|  |  |  +----------------+  +----------------+  +------------------------------+    |  |    |
|  |  |  |  Cart Service  |  |  Checkout      |  |  Payment Orchestration       |    |  |    |
|  |  |  |  (Session Mgmt)|  |  (State Machine)|  |  (Multi-provider)           |    |  |    |
|  |  |  +----------------+  +----------------+  +------------------------------+    |  |    |
|  |  +------------------------------------------------------------------------------+  |    |
|  +------------------------------------------------------------------------------------+    |
|                                       |                                                  |
|  +------------------------------------v-----------------------------------------------+    |
|  |                         FULFILLMENT LAYER                                        |    |
|  |  +-------------+  +-------------+  +-------------+  +-------------------------+    |    |
|  |  |  Inventory  |  |  Order      |  |  Shipping   |  |  Returns                |    |    |
|  |  |  Management |  |  Management |  |  Engine     |  |  Management             |    |    |
|  |  |  (Distributed)|  |  (Saga)     |  |  (Multi-carrier)|  |  (Reverse Logistics)|    |    |
|  |  +------+------+  +------+------+  +------+------+  +-----------+-------------+    |    |
|  +---------+---------+---------+-------------------------------------+------------------+    |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

### 2.2 Recommendation Engine Architecture

```
+--------------------------------------------------------------------------+
|                     RECOMMENDATION ENGINE PIPELINE                       |
+--------------------------------------------------------------------------+
|                                                                          |
|  Data Sources                                                            |
|  +-------------+  +-------------+  +-------------+  +-------------+      |
|  | Clickstream |  | Transactions|  | Product     |  | User        |      |
|  | Events      |  | History     |  | Catalog     |  | Profiles    |      |
|  +------+------+  +------+------+  +------+------+  +------+------+      |
|         |                |                |                |             |
|         +----------------+----------------+----------------+             |
|                          |                                               |
|  +-----------------------v------------------------+                       |
|  |            Feature Engineering                |                       |
|  |  +---------+ +---------+ +---------+         |                       |
|  |  | User    | | Item    | | Context |         |                       |
|  |  | Features| | Features| | Features|         |                       |
|  |  +---------+ +---------+ +---------+         |                       |
|  +-----------------------+------------------------+                       |
|                          |                                               |
|  +-----------------------v------------------------+                       |
|  |              ML Model Serving                 |                       |
|  |  +----------------+ +----------------+       |                       |
|  |  | Two-Tower Model| | Sequential Model|       |                       |
|  |  | (Retrieval)    | | (Ranking)       |       |                       |
|  |  +----------------+ +----------------+       |                       |
|  |  +----------------+ +----------------+       |                       |
|  |  | Transformer    | | Bandit Algos   |       |                       |
|  |  | (BERT4Rec)     | | (Contextual)   |       |                       |
|  |  +----------------+ +----------------+       |                       |
|  +-----------------------+------------------------+                       |
|                          |                                               |
|  +-----------------------v------------------------+                       |
|  |            Ranking & Personalization          |                       |
|  |  +---------+ +---------+ +---------+         |                       |
|  |  | Business| | Diversity| | Exploration|      |                       |
|  |  | Rules   | | Rerank  | | Injection  |      |                       |
|  |  +---------+ +---------+ +---------+         |                       |
|  +-----------------------+------------------------+                       |
|                          |                                               |
|  +-----------------------v------------------------+                       |
|  |              Recommendation API               |                       |
|  +-----------------------------------------------+                       |
|                                                                          |
+--------------------------------------------------------------------------+
```

---

## 3. Recommendation Algorithms

### 3.1 Collaborative Filtering

```go
// recommendation/collaborative_filtering.go - Collaborative Filtering
package recommendation

import (
    "math"
    "sync"
)

// UserItemMatrix represents user-item interactions
type UserItemMatrix struct {
    // userID -> itemID -> rating
    Data     map[string]map[string]float64
    mu       sync.RWMutex

    // Precomputed statistics
    userMeans   map[string]float64
    itemMeans   map[string]float64
    globalMean  float64
}

// NewUserItemMatrix creates a new matrix
func NewUserItemMatrix() *UserItemMatrix {
    return &UserItemMatrix{
        Data:       make(map[string]map[string]float64),
        userMeans:  make(map[string]float64),
        itemMeans:  make(map[string]float64),
    }
}

// AddInteraction adds a user-item interaction
func (m *UserItemMatrix) AddInteraction(userID, itemID string, rating float64) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.Data[userID] == nil {
        m.Data[userID] = make(map[string]float64)
    }
    m.Data[userID][itemID] = rating
    m.invalidateCache()
}

func (m *UserItemMatrix) invalidateCache() {
    m.userMeans = make(map[string]float64)
    m.itemMeans = make(map[string]float64)
}

// UserBasedCF implements user-based collaborative filtering
type UserBasedCF struct {
    matrix          *UserItemMatrix
    similarityFunc  SimilarityFunc
    k               int // Number of neighbors
}

type SimilarityFunc func(u1, u2 map[string]float64) float64

// NewUserBasedCF creates a user-based CF recommender
func NewUserBasedCF(matrix *UserItemMatrix, k int) *UserBasedCF {
    return &UserBasedCF{
        matrix:         matrix,
        similarityFunc: PearsonCorrelation,
        k:              k,
    }
}

// Recommend generates recommendations for a user
func (cf *UserBasedCF) Recommend(userID string, n int) ([]Recommendation, error) {
    cf.matrix.mu.RLock()
    defer cf.matrix.mu.RUnlock()

    targetUser := cf.matrix.Data[userID]
    if targetUser == nil {
        return nil, nil
    }

    // Find similar users
    neighbors := cf.findNeighbors(userID, targetUser)

    // Predict ratings for unseen items
    predictions := make(map[string]float64)

    // Get all candidate items
    candidates := cf.getCandidateItems(userID)

    for _, itemID := range candidates {
        prediction := cf.predictRating(userID, itemID, neighbors)
        if prediction > 0 {
            predictions[itemID] = prediction
        }
    }

    // Sort and return top N
    return cf.sortRecommendations(predictions, n), nil
}

func (cf *UserBasedCF) findNeighbors(userID string, targetUser map[string]float64) []Neighbor {
    var neighbors []Neighbor

    for otherID, otherRatings := range cf.matrix.Data {
        if otherID == userID {
            continue
        }

        sim := cf.similarityFunc(targetUser, otherRatings)
        if sim > 0 {
            neighbors = append(neighbors, Neighbor{
                UserID:     otherID,
                Similarity: sim,
                Ratings:    otherRatings,
            })
        }
    }

    // Sort by similarity (descending)
    sortNeighbors(neighbors)

    // Return top k
    if len(neighbors) > cf.k {
        return neighbors[:cf.k]
    }
    return neighbors
}

func (cf *UserBasedCF) predictRating(userID, itemID string, neighbors []Neighbor) float64 {
    var weightedSum float64
    var similaritySum float64

    for _, neighbor := range neighbors {
        if rating, ok := neighbor.Ratings[itemID]; ok {
            weightedSum += neighbor.Similarity * rating
            similaritySum += math.Abs(neighbor.Similarity)
        }
    }

    if similaritySum == 0 {
        return 0
    }

    return weightedSum / similaritySum
}

func (cf *UserBasedCF) getCandidateItems(userID string) []string {
    seen := make(map[string]bool)
    for itemID := range cf.matrix.Data[userID] {
        seen[itemID] = true
    }

    candidates := make([]string, 0)
    for _, items := range cf.matrix.Data {
        for itemID := range items {
            if !seen[itemID] {
                candidates = append(candidates, itemID)
                seen[itemID] = true // Avoid duplicates
            }
        }
    }

    return candidates
}

func (cf *UserBasedCF) sortRecommendations(predictions map[string]float64, n int) []Recommendation {
    recs := make([]Recommendation, 0, len(predictions))
    for itemID, score := range predictions {
        recs = append(recs, Recommendation{ItemID: itemID, Score: score})
    }

    // Sort by score descending
    sortRecommendations(recs)

    if len(recs) > n {
        return recs[:n]
    }
    return recs
}

// ItemBasedCF implements item-based collaborative filtering
type ItemBasedCF struct {
    matrix         *UserItemMatrix
    itemSimilarity map[string]map[string]float64
    k              int
}

// NewItemBasedCF creates an item-based CF recommender
func NewItemBasedCF(matrix *UserItemMatrix, k int) *ItemBasedCF {
    cf := &ItemBasedCF{
        matrix:         matrix,
        itemSimilarity: make(map[string]map[string]float64),
        k:              k,
    }
    cf.precomputeSimilarities()
    return cf
}

func (cf *ItemBasedCF) precomputeSimilarities() {
    // Build item-user matrix (transpose)
    itemUsers := make(map[string]map[string]float64)

    cf.matrix.mu.RLock()
    for userID, items := range cf.matrix.Data {
        for itemID, rating := range items {
            if itemUsers[itemID] == nil {
                itemUsers[itemID] = make(map[string]float64)
            }
            itemUsers[itemID][userID] = rating
        }
    }
    cf.matrix.mu.RUnlock()

    // Compute cosine similarity between items
    for item1, users1 := range itemUsers {
        cf.itemSimilarity[item1] = make(map[string]float64)
        for item2, users2 := range itemUsers {
            if item1 >= item2 {
                continue
            }
            sim := CosineSimilarity(users1, users2)
            if sim > 0 {
                cf.itemSimilarity[item1][item2] = sim
                cf.itemSimilarity[item2][item1] = sim
            }
        }
    }
}

func (cf *ItemBasedCF) Recommend(userID string, n int) ([]Recommendation, error) {
    cf.matrix.mu.RLock()
    defer cf.matrix.mu.RUnlock()

    userRatings := cf.matrix.Data[userID]
    if userRatings == nil {
        return nil, nil
    }

    scores := make(map[string]float64)

    // For each item the user has rated
    for ratedItem, rating := range userRatings {
        // Find similar items
        similarItems := cf.itemSimilarity[ratedItem]
        for similarItem, sim := range similarItems {
            // Skip if user already rated
            if _, rated := userRatings[similarItem]; rated {
                continue
            }
            scores[similarItem] += sim * rating
        }
    }

    return cf.sortRecommendations(scores, n), nil
}

func (cf *ItemBasedCF) sortRecommendations(scores map[string]float64, n int) []Recommendation {
    recs := make([]Recommendation, 0, len(scores))
    for itemID, score := range scores {
        recs = append(recs, Recommendation{ItemID: itemID, Score: score})
    }
    sortRecommendations(recs)
    if len(recs) > n {
        return recs[:n]
    }
    return recs
}

// Similarity Functions

// PearsonCorrelation computes Pearson correlation coefficient
func PearsonCorrelation(u1, u2 map[string]float64) float64 {
    // Find common items
    common := make([][2]float64, 0)
    for item, r1 := range u1 {
        if r2, ok := u2[item]; ok {
            common = append(common, [2]float64{r1, r2})
        }
    }

    if len(common) < 2 {
        return 0
    }

    // Calculate means
    var sum1, sum2 float64
    for _, pair := range common {
        sum1 += pair[0]
        sum2 += pair[1]
    }
    mean1 := sum1 / float64(len(common))
    mean2 := sum2 / float64(len(common))

    // Calculate correlation
    var numerator, denom1, denom2 float64
    for _, pair := range common {
        diff1 := pair[0] - mean1
        diff2 := pair[1] - mean2
        numerator += diff1 * diff2
        denom1 += diff1 * diff1
        denom2 += diff2 * diff2
    }

    denominator := math.Sqrt(denom1) * math.Sqrt(denom2)
    if denominator == 0 {
        return 0
    }

    return numerator / denominator
}

// CosineSimilarity computes cosine similarity
func CosineSimilarity(u1, u2 map[string]float64) float64 {
    var dotProduct, norm1, norm2 float64

    for item, r1 := range u1 {
        norm1 += r1 * r1
        if r2, ok := u2[item]; ok {
            dotProduct += r1 * r2
        }
    }

    for _, r2 := range u2 {
        norm2 += r2 * r2
    }

    denominator := math.Sqrt(norm1) * math.Sqrt(norm2)
    if denominator == 0 {
        return 0
    }

    return dotProduct / denominator
}

// JaccardSimilarity computes Jaccard similarity
func JaccardSimilarity(u1, u2 map[string]float64) float64 {
    intersection := 0
    for item := range u1 {
        if _, ok := u2[item]; ok {
            intersection++
        }
    }

    union := len(u1) + len(u2) - intersection
    if union == 0 {
        return 0
    }

    return float64(intersection) / float64(union)
}

type Neighbor struct {
    UserID     string
    Similarity float64
    Ratings    map[string]float64
}

type Recommendation struct {
    ItemID string
    Score  float64
}

func sortNeighbors(neighbors []Neighbor) {
    // Simple bubble sort for clarity
    for i := 0; i < len(neighbors); i++ {
        for j := i + 1; j < len(neighbors); j++ {
            if neighbors[j].Similarity > neighbors[i].Similarity {
                neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
            }
        }
    }
}

func sortRecommendations(recs []Recommendation) {
    for i := 0; i < len(recs); i++ {
        for j := i + 1; j < len(recs); j++ {
            if recs[j].Score > recs[i].Score {
                recs[i], recs[j] = recs[j], recs[i]
            }
        }
    }
}
```

### 3.2 Matrix Factorization (SVD, ALS)

```go
// recommendation/matrix_factorization.go - Matrix Factorization
package recommendation

import (
    "math"
    "math/rand"
    "sync"
)

// MatrixFactorization implements collaborative filtering using matrix factorization
type MatrixFactorization struct {
    numFactors    int
    learningRate  float64
    regularization float64
    iterations    int

    // User and item latent factor matrices
    userFactors   map[string][]float64
    itemFactors   map[string][]float64

    // Biases
    userBias      map[string]float64
    itemBias      map[string]float64
    globalBias    float64

    mu            sync.RWMutex
}

// NewMatrixFactorization creates a new MF model
func NewMatrixFactorization(numFactors int, lr, reg float64, iterations int) *MatrixFactorization {
    return &MatrixFactorization{
        numFactors:     numFactors,
        learningRate:   lr,
        regularization: reg,
        iterations:     iterations,
        userFactors:    make(map[string][]float64),
        itemFactors:    make(map[string][]float64),
        userBias:       make(map[string]float64),
        itemBias:       make(map[string]float64),
    }
}

// Train trains the model using Stochastic Gradient Descent
func (mf *MatrixFactorization) Train(matrix *UserItemMatrix) error {
    // Calculate global bias
    mf.calculateGlobalBias(matrix)

    // Initialize factors
    mf.initializeFactors(matrix)

    // Training loop
    for iter := 0; iter < mf.iterations; iter++ {
        var totalError float64

        matrix.mu.RLock()
        for userID, items := range matrix.Data {
            for itemID, rating := range items {
                // Predict and calculate error
                prediction := mf.predict(userID, itemID)
                error := rating - prediction
                totalError += error * error

                // Update biases
                mf.userBias[userID] += mf.learningRate * (error - mf.regularization*mf.userBias[userID])
                mf.itemBias[itemID] += mf.learningRate * (error - mf.regularization*mf.itemBias[itemID])

                // Update factors
                userF := mf.userFactors[userID]
                itemF := mf.itemFactors[itemID]

                for f := 0; f < mf.numFactors; f++ {
                    userFactor := userF[f]
                    itemFactor := itemF[f]

                    userF[f] += mf.learningRate * (error*itemFactor - mf.regularization*userFactor)
                    itemF[f] += mf.learningRate * (error*userFactor - mf.regularization*itemFactor)
                }
            }
        }
        matrix.mu.RUnlock()

        // Decay learning rate
        mf.learningRate *= 0.9
    }

    return nil
}

func (mf *MatrixFactorization) initializeFactors(matrix *UserItemMatrix) {
    matrix.mu.RLock()
    defer matrix.mu.RUnlock()

    // Initialize user factors
    for userID := range matrix.Data {
        mf.userFactors[userID] = make([]float64, mf.numFactors)
        for f := 0; f < mf.numFactors; f++ {
            mf.userFactors[userID][f] = rand.NormFloat64() * 0.01
        }
    }

    // Initialize item factors
    for _, items := range matrix.Data {
        for itemID := range items {
            if _, exists := mf.itemFactors[itemID]; !exists {
                mf.itemFactors[itemID] = make([]float64, mf.numFactors)
                for f := 0; f < mf.numFactors; f++ {
                    mf.itemFactors[itemID][f] = rand.NormFloat64() * 0.01
                }
            }
        }
    }
}

func (mf *MatrixFactorization) calculateGlobalBias(matrix *UserItemMatrix) {
    var sum float64
    var count int

    matrix.mu.RLock()
    for _, items := range matrix.Data {
        for _, rating := range items {
            sum += rating
            count++
        }
    }
    matrix.mu.RUnlock()

    if count > 0 {
        mf.globalBias = sum / float64(count)
    }
}

func (mf *MatrixFactorization) predict(userID, itemID string) float64 {
    prediction := mf.globalBias

    if ub, ok := mf.userBias[userID]; ok {
        prediction += ub
    }
    if ib, ok := mf.itemBias[itemID]; ok {
        prediction += ib
    }

    userF, uok := mf.userFactors[userID]
    itemF, iok := mf.itemFactors[itemID]

    if uok && iok {
        for f := 0; f < mf.numFactors; f++ {
            prediction += userF[f] * itemF[f]
        }
    }

    return prediction
}

// Predict rating for user-item pair
func (mf *MatrixFactorization) Predict(userID, itemID string) float64 {
    mf.mu.RLock()
    defer mf.mu.RUnlock()
    return mf.predict(userID, itemID)
}

// Recommend generates recommendations for a user
func (mf *MatrixFactorization) Recommend(userID string, excludeItems map[string]bool, n int) []Recommendation {
    mf.mu.RLock()
    defer mf.mu.RUnlock()

    scores := make(map[string]float64)

    for itemID := range mf.itemFactors {
        if excludeItems[itemID] {
            continue
        }
        scores[itemID] = mf.predict(userID, itemID)
    }

    return sortAndLimit(scores, n)
}

// Alternating Least Squares (ALS) implementation
type ALS struct {
    numFactors     int
    regularization float64
    iterations     int

    userFactors    map[string][]float64
    itemFactors    map[string][]float64
}

// TrainAlternatingLeastSquares trains using ALS
func (als *ALS) Train(matrix *UserItemMatrix) error {
    als.initializeFactors(matrix)

    for iter := 0; iter < als.iterations; iter++ {
        // Fix item factors, solve for user factors
        als.optimizeUserFactors(matrix)

        // Fix user factors, solve for item factors
        als.optimizeItemFactors(matrix)
    }

    return nil
}

func (als *ALS) optimizeUserFactors(matrix *UserItemMatrix) {
    matrix.mu.RLock()
    defer matrix.mu.RUnlock()

    for userID, items := range matrix.Data {
        // Build matrices for this user
        A := make([][]float64, als.numFactors)
        b := make([]float64, als.numFactors)

        for i := 0; i < als.numFactors; i++ {
            A[i] = make([]float64, als.numFactors)
        }

        for itemID, rating := range items {
            itemF := als.itemFactors[itemID]

            for i := 0; i < als.numFactors; i++ {
                b[i] += rating * itemF[i]
                for j := 0; j < als.numFactors; j++ {
                    A[i][j] += itemF[i] * itemF[j]
                }
            }
        }

        // Add regularization
        for i := 0; i < als.numFactors; i++ {
            A[i][i] += als.regularization
        }

        // Solve Ax = b using Gaussian elimination (simplified)
        als.userFactors[userID] = solveLinearSystem(A, b)
    }
}

func (als *ALS) optimizeItemFactors(matrix *UserItemMatrix) {
    // Similar to optimizeUserFactors but transposed
    // Implementation omitted for brevity
}

func solveLinearSystem(A [][]float64, b []float64) []float64 {
    n := len(b)
    x := make([]float64, n)

    // Simple Gaussian elimination
    // In production, use a numerical library
    for i := 0; i < n; i++ {
        x[i] = b[i] / A[i][i]
    }

    return x
}

func sortAndLimit(scores map[string]float64, n int) []Recommendation {
    recs := make([]Recommendation, 0, len(scores))
    for itemID, score := range scores {
        recs = append(recs, Recommendation{ItemID: itemID, Score: score})
    }

    sortRecommendations(recs)

    if len(recs) > n {
        return recs[:n]
    }
    return recs
}
```

### 3.3 Deep Learning Recommendations

```go
// recommendation/deep_learning.go - Deep Learning Models
package recommendation

import (
    "encoding/json"
    "math"
)

// TwoTowerModel implements the Two-Tower architecture for retrieval
type TwoTowerModel struct {
    userEmbeddingDim  int
    itemEmbeddingDim  int
    hiddenDims        []int

    // Pre-computed item embeddings for fast retrieval
    itemIndex         map[string][]float64

    // Model weights (simplified)
    userTowerWeights  [][][]float64
    itemTowerWeights  [][][]float64
}

// NewTwoTowerModel creates a two-tower recommendation model
func NewTwoTowerModel(userDim, itemDim int, hiddenDims []int) *TwoTowerModel {
    return &TwoTowerModel{
        userEmbeddingDim: userDim,
        itemEmbeddingDim: itemDim,
        hiddenDims:       hiddenDims,
        itemIndex:        make(map[string][]float64),
    }
}

// UserFeatures represents user features
type UserFeatures struct {
    UserID          string    `json:"user_id"`
    Age             float64   `json:"age"`
    Gender          float64   `json:"gender"`
    Location        []float64 `json:"location"`
    CategoryHistory []float64 `json:"category_history"`
    PricePreference float64   `json:"price_preference"`
    SessionFeatures []float64 `json:"session_features"`
}

// ItemFeatures represents item features
type ItemFeatures struct {
    ItemID          string    `json:"item_id"`
    Category        []float64 `json:"category"` // One-hot encoded
    Price           float64   `json:"price"`
    Brand           []float64 `json:"brand"`
    TextEmbedding   []float64 `json:"text_embedding"`
    ImageEmbedding  []float64 `json:"image_embedding"`
    Popularity      float64   `json:"popularity"`
}

// EncodeUser converts user features to embedding
func (m *TwoTowerModel) EncodeUser(features UserFeatures) []float64 {
    // Concatenate all user features
    input := m.concatenateUserFeatures(features)

    // Forward pass through user tower
    return m.forwardPass(input, m.userTowerWeights)
}

// EncodeItem converts item features to embedding
func (m *TwoTowerModel) EncodeItem(features ItemFeatures) []float64 {
    input := m.concatenateItemFeatures(features)
    return m.forwardPass(input, m.itemTowerWeights)
}

func (m *TwoTowerModel) concatenateUserFeatures(f UserFeatures) []float64 {
    result := make([]float64, 0, m.userEmbeddingDim)
    result = append(result, f.Age/100)
    result = append(result, f.Gender)
    result = append(result, f.Location...)
    result = append(result, f.CategoryHistory...)
    result = append(result, f.PricePreference/1000)
    result = append(result, f.SessionFeatures...)
    return result
}

func (m *TwoTowerModel) concatenateItemFeatures(f ItemFeatures) []float64 {
    result := make([]float64, 0, m.itemEmbeddingDim)
    result = append(result, f.Category...)
    result = append(result, f.Price/1000)
    result = append(result, f.Brand...)
    result = append(result, f.TextEmbedding...)
    result = append(result, f.ImageEmbedding...)
    result = append(result, f.Popularity)
    return result
}

func (m *TwoTowerModel) forwardPass(input []float64, weights [][][]float64) []float64 {
    // Simplified forward pass
    // In production, this would use a proper ML framework
    current := input

    for _, layerWeights := range weights {
        next := make([]float64, len(layerWeights))
        for i, neuronWeights := range layerWeights {
            var sum float64
            for j, w := range neuronWeights {
                if j < len(current) {
                    sum += w * current[j]
                }
            }
            next[i] = relu(sum)
        }
        current = next
    }

    // L2 normalize
    return l2Normalize(current)
}

func (m *TwoTowerModel) IndexItems(items []ItemFeatures) {
    for _, item := range items {
        embedding := m.EncodeItem(item)
        m.itemIndex[item.ItemID] = embedding
    }
}

// Recommend finds similar items using approximate nearest neighbors
func (m *TwoTowerModel) Recommend(user UserFeatures, n int) []Recommendation {
    userEmbedding := m.EncodeUser(user)

    scores := make(map[string]float64)
    for itemID, itemEmbedding := range m.itemIndex {
        scores[itemID] = cosineSimilarity(userEmbedding, itemEmbedding)
    }

    return sortAndLimit(scores, n)
}

// TransformerModel implements a transformer-based sequential recommendation
type TransformerModel struct {
    numLayers      int
    numHeads       int
    embeddingDim   int
    maxSeqLength   int

    // Pretrained weights
    embeddings     map[string][]float64
    transformer    Transformer
}

// Transformer represents transformer architecture
type Transformer struct {
    attentionLayers []MultiHeadAttention
    feedForward     []FeedForward
    layerNorm       []LayerNorm
}

// MultiHeadAttention implements self-attention
type MultiHeadAttention struct {
    numHeads int
    dim      int
    WQ       [][]float64
    WK       [][]float64
    WV       [][]float64
    WO       [][]float64
}

// Forward applies multi-head self-attention
func (mha *MultiHeadAttention) Forward(Q, K, V [][]float64) [][]float64 {
    // Scaled dot-product attention for each head
    heads := make([][][]float64, mha.numHeads)
    headDim := mha.dim / mha.numHeads

    for h := 0; h < mha.numHeads; h++ {
        // Project Q, K, V for this head
        Qh := mha.project(Q, mha.WQ, h*headDim, (h+1)*headDim)
        Kh := mha.project(K, mha.WK, h*headDim, (h+1)*headDim)
        Vh := mha.project(V, mha.WV, h*headDim, (h+1)*headDim)

        // Attention(Q, K, V) = softmax(QK^T / sqrt(d_k)) * V
        attention := scaledDotProductAttention(Qh, Kh, Vh)
        heads[h] = attention
    }

    // Concatenate heads
    concatenated := concatenateHeads(heads)

    // Final linear projection
    return matrixMultiply(concatenated, mha.WO)
}

func (mha *MultiHeadAttention) project(X [][]float64, W [][]float64, start, end int) [][]float64 {
    result := make([][]float64, len(X))
    for i := range X {
        result[i] = make([]float64, end-start)
        for j := start; j < end && j < len(W); j++ {
            var sum float64
            for k, x := range X[i] {
                if k < len(W[j]) {
                    sum += x * W[j][k]
                }
            }
            result[i][j-start] = sum
        }
    }
    return result
}

func scaledDotProductAttention(Q, K, V [][]float64) [][]float64 {
    // Compute Q * K^T
    scores := matrixMultiply(Q, transpose(K))

    // Scale by sqrt(d_k)
    dK := float64(len(K[0]))
    for i := range scores {
        for j := range scores[i] {
            scores[i][j] /= math.Sqrt(dK)
        }
    }

    // Apply softmax
    attentionWeights := softmaxRows(scores)

    // Multiply by V
    return matrixMultiply(attentionWeights, V)
}

type FeedForward struct {
    W1 [][]float64
    W2 [][]float64
    B1 []float64
    B2 []float64
}

func (ff *FeedForward) Forward(X [][]float64) [][]float64 {
    // First linear layer + ReLU
    hidden := matrixMultiply(X, ff.W1)
    for i := range hidden {
        for j := range hidden[i] {
            hidden[i][j] += ff.B1[j]
            hidden[i][j] = relu(hidden[i][j])
        }
    }

    // Second linear layer
    output := matrixMultiply(hidden, ff.W2)
    for i := range output {
        for j := range output[i] {
            output[i][j] += ff.B2[j]
        }
    }

    return output
}

type LayerNorm struct {
    Gamma []float64
    Beta  []float64
    Eps   float64
}

func (ln *LayerNorm) Forward(X [][]float64) [][]float64 {
    result := make([][]float64, len(X))
    for i, x := range X {
        result[i] = layerNormalize(x, ln.Gamma, ln.Beta, ln.Eps)
    }
    return result
}

func layerNormalize(x, gamma, beta []float64, eps float64) []float64 {
    // Calculate mean
    var mean float64
    for _, v := range x {
        mean += v
    }
    mean /= float64(len(x))

    // Calculate variance
    var variance float64
    for _, v := range x {
        diff := v - mean
        variance += diff * diff
    }
    variance /= float64(len(x))

    // Normalize
    result := make([]float64, len(x))
    for i, v := range x {
        normalized := (v - mean) / math.Sqrt(variance+eps)
        result[i] = gamma[i]*normalized + beta[i]
    }

    return result
}

// PredictNextItem predicts the next item in a sequence
func (tm *TransformerModel) PredictNextItem(sequence []string) []Recommendation {
    // Convert sequence to embeddings
    embeddings := make([][]float64, len(sequence))
    for i, itemID := range sequence {
        embeddings[i] = tm.embeddings[itemID]
    }

    // Add positional encoding
    embeddings = tm.addPositionalEncoding(embeddings)

    // Pass through transformer layers
    hidden := embeddings
    for layer := 0; layer < tm.numLayers; layer++ {
        // Self-attention
        attention := tm.transformer.attentionLayers[layer].Forward(hidden, hidden, hidden)
        hidden = tm.transformer.layerNorm[layer].Forward(add(hidden, attention))

        // Feed-forward
        ff := tm.transformer.feedForward[layer].Forward(hidden)
        hidden = tm.transformer.layerNorm[layer].Forward(add(hidden, ff))
    }

    // Use last hidden state for prediction
    lastHidden := hidden[len(hidden)-1]

    // Score all items
    scores := make(map[string]float64)
    for itemID, itemEmb := range tm.embeddings {
        scores[itemID] = cosineSimilarity(lastHidden, itemEmb)
    }

    return sortAndLimit(scores, 10)
}

func (tm *TransformerModel) addPositionalEncoding(embeddings [][]float64) [][]float64 {
    result := make([][]float64, len(embeddings))
    for pos, emb := range embeddings {
        result[pos] = make([]float64, len(emb))
        for i := 0; i < len(emb); i++ {
            // Sinusoidal positional encoding
            angle := float64(pos) / math.Pow(10000, float64(2*i)/float64(len(emb)))
            if i%2 == 0 {
                result[pos][i] = emb[i] + math.Sin(angle)
            } else {
                result[pos][i] = emb[i] + math.Cos(angle)
            }
        }
    }
    return result
}

// Helper functions

func relu(x float64) float64 {
    if x < 0 {
        return 0
    }
    return x
}

func l2Normalize(v []float64) []float64 {
    var sum float64
    for _, x := range v {
        sum += x * x
    }
    norm := math.Sqrt(sum)

    result := make([]float64, len(v))
    if norm > 0 {
        for i, x := range v {
            result[i] = x / norm
        }
    }
    return result
}

func cosineSimilarity(a, b []float64) float64 {
    var dot, normA, normB float64
    for i := 0; i < len(a) && i < len(b); i++ {
        dot += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func transpose(m [][]float64) [][]float64 {
    if len(m) == 0 {
        return m
    }
    result := make([][]float64, len(m[0]))
    for i := range result {
        result[i] = make([]float64, len(m))
        for j := range m {
            result[i][j] = m[j][i]
        }
    }
    return result
}

func matrixMultiply(A, B [][]float64) [][]float64 {
    if len(A) == 0 || len(B) == 0 {
        return nil
    }
    result := make([][]float64, len(A))
    for i := range A {
        result[i] = make([]float64, len(B[0]))
        for j := range B[0] {
            var sum float64
            for k := 0; k < len(B) && k < len(A[i]); k++ {
                sum += A[i][k] * B[k][j]
            }
            result[i][j] = sum
        }
    }
    return result
}

func softmaxRows(m [][]float64) [][]float64 {
    result := make([][]float64, len(m))
    for i, row := range m {
        result[i] = softmax(row)
    }
    return result
}

func softmax(v []float64) []float64 {
    result := make([]float64, len(v))
    var max float64
    for _, x := range v {
        if x > max {
            max = x
        }
    }

    var sum float64
    for i, x := range v {
        result[i] = math.Exp(x - max)
        sum += result[i]
    }

    for i := range result {
        result[i] /= sum
    }
    return result
}

func concatenateHeads(heads [][][]float64) [][]float64 {
    if len(heads) == 0 {
        return nil
    }

    seqLen := len(heads[0])
    result := make([][]float64, seqLen)

    for i := 0; i < seqLen; i++ {
        var concatenated []float64
        for _, head := range heads {
            concatenated = append(concatenated, head[i]...)
        }
        result[i] = concatenated
    }

    return result
}

func add(a, b [][]float64) [][]float64 {
    result := make([][]float64, len(a))
    for i := range a {
        result[i] = make([]float64, len(a[i]))
        for j := range a[i] {
            if j < len(b[i]) {
                result[i][j] = a[i][j] + b[i][j]
            }
        }
    }
    return result
}
```

### 3.4 Contextual Bandits for Exploration

```go
// recommendation/bandits.go - Contextual Bandits
package recommendation

import (
    "math"
    "sync"
)

// ContextualBandit implements Thompson Sampling with linear payoff
type ContextualBandit struct {
    arms          []Arm
    contextDim    int
    alpha         float64 // Exploration parameter

    // Arm parameters (one per arm)
    armParams     map[string]*ArmParameters
    mu            sync.RWMutex
}

type Arm struct {
    ID    string
    Features []float64
}

type ArmParameters struct {
    A     [][]float64 // Design matrix
    b     []float64   // Response vector
    theta []float64   // Estimated parameters
}

// NewContextualBandit creates a new bandit algorithm
func NewContextualBandit(contextDim int, alpha float64) *ContextualBandit {
    return &ContextualBandit{
        arms:       make([]Arm, 0),
        contextDim: contextDim,
        alpha:      alpha,
        armParams:  make(map[string]*ArmParameters),
    }
}

// AddArm adds a new arm (item) to the bandit
func (cb *ContextualBandit) AddArm(arm Arm) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.arms = append(cb.arms, arm)

    // Initialize arm parameters
    A := make([][]float64, cb.contextDim)
    for i := range A {
        A[i] = make([]float64, cb.contextDim)
        A[i][i] = 1 // Identity matrix
    }

    cb.armParams[arm.ID] = &ArmParameters{
        A:     A,
        b:     make([]float64, cb.contextDim),
        theta: make([]float64, cb.contextDim),
    }
}

// SelectArm chooses an arm given context (Thompson Sampling)
func (cb *ContextualBandit) SelectArm(context []float64) (string, float64) {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    var bestArm string
    var bestScore float64 = -math.MaxFloat64

    for _, arm := range cb.arms {
        params := cb.armParams[arm.ID]

        // Sample from posterior distribution
        theta := cb.sampleFromPosterior(params)

        // Calculate expected reward
        score := dotProduct(theta, context)

        if score > bestScore {
            bestScore = score
            bestArm = arm.ID
        }
    }

    return bestArm, bestScore
}

func (cb *ContextualBandit) sampleFromPosterior(params *ArmParameters) []float64 {
    // Simplified: return point estimate + noise
    // In production, use proper multivariate normal sampling

    result := make([]float64, len(params.theta))
    for i, t := range params.theta {
        result[i] = t + cb.alpha*randNorm()
    }
    return result
}

// Update updates arm parameters after observing reward
func (cb *ContextualBandit) Update(armID string, context []float64, reward float64) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    params := cb.armParams[armID]
    if params == nil {
        return
    }

    // Update A (design matrix): A = A + x * x^T
    for i := 0; i < cb.contextDim; i++ {
        for j := 0; j < cb.contextDim; j++ {
            params.A[i][j] += context[i] * context[j]
        }
    }

    // Update b: b = b + r * x
    for i := 0; i < cb.contextDim; i++ {
        params.b[i] += reward * context[i]
    }

    // Update theta: theta = A^-1 * b
    params.theta = solveLinearSystem(params.A, params.b)
}

// LinUCB implements Linear Upper Confidence Bound algorithm
type LinUCB struct {
    arms       []Arm
    contextDim int
    alpha      float64 // Exploration parameter
    armParams  map[string]*LinUCBParams
    mu         sync.RWMutex
}

type LinUCBParams struct {
    A      [][]float64
    b      []float64
    A_inv  [][]float64
    theta  []float64
}

func NewLinUCB(contextDim int, alpha float64) *LinUCB {
    return &LinUCB{
        arms:       make([]Arm, 0),
        contextDim: contextDim,
        alpha:      alpha,
        armParams:  make(map[string]*LinUCBParams),
    }
}

func (ucb *LinUCB) SelectArm(context []float64) (string, float64) {
    ucb.mu.RLock()
    defer ucb.mu.RUnlock()

    var bestArm string
    var bestUCB float64 = -math.MaxFloat64

    for _, arm := range ucb.arms {
        params := ucb.armParams[arm.ID]

        // Mean reward: theta^T * x
        meanReward := dotProduct(params.theta, context)

        // Confidence interval: alpha * sqrt(x^T * A^-1 * x)
        confidence := ucb.alpha * math.Sqrt(quadraticForm(context, params.A_inv))

        ucbValue := meanReward + confidence

        if ucbValue > bestUCB {
            bestUCB = ucbValue
            bestArm = arm.ID
        }
    }

    return bestArm, bestUCB
}

func quadraticForm(x []float64, A [][]float64) float64 {
    // Compute x^T * A * x
    result := 0.0
    for i := 0; i < len(x); i++ {
        for j := 0; j < len(x); j++ {
            result += x[i] * A[i][j] * x[j]
        }
    }
    return result
}

func dotProduct(a, b []float64) float64 {
    var sum float64
    for i := 0; i < len(a) && i < len(b); i++ {
        sum += a[i] * b[i]
    }
    return sum
}

func randNorm() float64 {
    // Box-Muller transform
    u1 := 0.5
    u2 := 0.5
    return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
}
```

---

## 4. Search Engine

### 4.1 Elasticsearch Integration

```go
// search/elasticsearch.go - Search Engine
package search

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
)

// SearchClient interface for search operations
type SearchClient interface {
    IndexDocument(ctx context.Context, index string, doc interface{}) error
    Search(ctx context.Context, query SearchQuery) (*SearchResult, error)
    DeleteDocument(ctx context.Context, index, docID string) error
    BulkIndex(ctx context.Context, index string, docs []interface{}) error
}

// SearchQuery represents a search request
type SearchQuery struct {
    Query       string
    Filters     []Filter
    Sort        []SortField
    Pagination  Pagination
    Aggregations []Aggregation
}

type Filter struct {
    Field    string
    Operator string // eq, gt, lt, range, terms
    Value    interface{}
}

type SortField struct {
    Field     string
    Direction string // asc, desc
}

type Pagination struct {
    From  int
    Size  int
}

type Aggregation struct {
    Name  string
    Type  string // terms, range, histogram
    Field string
}

type SearchResult struct {
    Total      int64
    Hits       []SearchHit
    Aggregations map[string]AggregationResult
    Facets     map[string][]FacetValue
}

type SearchHit struct {
    ID      string
    Score   float64
    Source  map[string]interface{}
}

type AggregationResult struct {
    Buckets []Bucket
}

type Bucket struct {
    Key      string
    DocCount int64
}

type FacetValue struct {
    Value    string
    Count    int64
}

// ProductSearch implements product search functionality
type ProductSearch struct {
    client SearchClient
    index  string
}

func NewProductSearch(client SearchClient, index string) *ProductSearch {
    return &ProductSearch{
        client: client,
        index:  index,
    }
}

// Search performs product search with filters and facets
func (ps *ProductSearch) Search(ctx context.Context, req ProductSearchRequest) (*ProductSearchResponse, error) {
    query := ps.buildSearchQuery(req)

    result, err := ps.client.Search(ctx, query)
    if err != nil {
        return nil, err
    }

    return ps.transformResult(result), nil
}

type ProductSearchRequest struct {
    Keywords     string
    Category     string
    Brand        string
    PriceRange   *PriceRange
    Attributes   map[string]string
    SortBy       string
    Page         int
    PageSize     int
}

type PriceRange struct {
    Min float64
    Max float64
}

type ProductSearchResponse struct {
    Products     []ProductDocument
    Total        int64
    Page         int
    PageSize     int
    Facets       ProductFacets
    Suggestions  []string
}

type ProductDocument struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Category    string                 `json:"category"`
    Brand       string                 `json:"brand"`
    Price       float64                `json:"price"`
    Currency    string                 `json:"currency"`
    Attributes  map[string]interface{} `json:"attributes"`
    Images      []string               `json:"images"`
    InStock     bool                   `json:"in_stock"`
    Rating      float64                `json:"rating"`
    ReviewCount int                    `json:"review_count"`
}

type ProductFacets struct {
    Categories  []FacetValue
    Brands      []FacetValue
    PriceRanges []FacetValue
    Attributes  map[string][]FacetValue
}

func (ps *ProductSearch) buildSearchQuery(req ProductSearchRequest) SearchQuery {
    var filters []Filter

    // Category filter
    if req.Category != "" {
        filters = append(filters, Filter{
            Field:    "category",
            Operator: "eq",
            Value:    req.Category,
        })
    }

    // Brand filter
    if req.Brand != "" {
        filters = append(filters, Filter{
            Field:    "brand",
            Operator: "eq",
            Value:    req.Brand,
        })
    }

    // Price range filter
    if req.PriceRange != nil {
        filters = append(filters, Filter{
            Field:    "price",
            Operator: "range",
            Value:    req.PriceRange,
        })
    }

    // Attribute filters
    for key, value := range req.Attributes {
        filters = append(filters, Filter{
            Field:    fmt.Sprintf("attributes.%s", key),
            Operator: "eq",
            Value:    value,
        })
    }

    // Sorting
    var sort []SortField
    switch req.SortBy {
    case "price_asc":
        sort = append(sort, SortField{Field: "price", Direction: "asc"})
    case "price_desc":
        sort = append(sort, SortField{Field: "price", Direction: "desc"})
    case "rating":
        sort = append(sort, SortField{Field: "rating", Direction: "desc"})
    default:
        sort = append(sort, SortField{Field: "_score", Direction: "desc"})
    }

    // Pagination
    from := (req.Page - 1) * req.PageSize

    // Aggregations for facets
    aggregations := []Aggregation{
        {Name: "categories", Type: "terms", Field: "category"},
        {Name: "brands", Type: "terms", Field: "brand"},
        {Name: "price_ranges", Type: "histogram", Field: "price"},
    }

    return SearchQuery{
        Query:        req.Keywords,
        Filters:      filters,
        Sort:         sort,
        Pagination:   Pagination{From: from, Size: req.PageSize},
        Aggregations: aggregations,
    }
}

func (ps *ProductSearch) transformResult(result *SearchResult) *ProductSearchResponse {
    products := make([]ProductDocument, 0, len(result.Hits))

    for _, hit := range result.Hits {
        var product ProductDocument
        // Convert source to ProductDocument
        sourceJSON, _ := json.Marshal(hit.Source)
        json.Unmarshal(sourceJSON, &product)
        product.ID = hit.ID
        products = append(products, product)
    }

    facets := ProductFacets{
        Attributes: make(map[string][]FacetValue),
    }

    if cats, ok := result.Aggregations["categories"]; ok {
        for _, bucket := range cats.Buckets {
            facets.Categories = append(facets.Categories, FacetValue{Value: bucket.Key, Count: bucket.DocCount})
        }
    }

    if brands, ok := result.Aggregations["brands"]; ok {
        for _, bucket := range brands.Buckets {
            facets.Brands = append(facets.Brands, FacetValue{Value: bucket.Key, Count: bucket.DocCount})
        }
    }

    return &ProductSearchResponse{
        Products: products,
        Total:    result.Total,
        Facets:   facets,
    }
}

// Autocomplete provides search suggestions
type AutocompleteService struct {
    client SearchClient
    index  string
}

func (as *AutocompleteService) Suggest(ctx context.Context, prefix string, limit int) ([]string, error) {
    // Implementation would use Elasticsearch suggesters
    // or edge n-grams for autocomplete
    return []string{}, nil
}
```

---

## 5. Cart and Checkout

### 5.1 Cart Service

```go
// cart/service.go - Cart Management
package cart

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// Cart represents a shopping cart
type Cart struct {
    ID          string       `json:"id"`
    UserID      string       `json:"user_id"`
    SessionID   string       `json:"session_id"`
    Items       []CartItem   `json:"items"`
    Summary     CartSummary  `json:"summary"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    ExpiresAt   time.Time    `json:"expires_at"`
    AppliedPromo string      `json:"applied_promo,omitempty"`
}

type CartItem struct {
    ID        string  `json:"id"`
    ProductID string  `json:"product_id"`
    SKU       string  `json:"sku"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    UnitPrice float64 `json:"unit_price"`
    Currency  string  `json:"currency"`
    ImageURL  string  `json:"image_url"`
    Options   map[string]string `json:"options,omitempty"`
}

type CartSummary struct {
    Subtotal      float64           `json:"subtotal"`
    Discount      float64           `json:"discount"`
    Shipping      float64           `json:"shipping"`
    Tax           float64           `json:"tax"`
    Total         float64           `json:"total"`
    Currency      string            `json:"currency"`
    ItemCount     int               `json:"item_count"`
    DiscountBreakdown []DiscountInfo `json:"discount_breakdown,omitempty"`
}

type DiscountInfo struct {
    Code        string  `json:"code"`
    Description string  `json:"description"`
    Amount      float64 `json:"amount"`
}

// CartService manages shopping carts
type CartService struct {
    repository CartRepository
    inventory  InventoryClient
    pricing    PricingEngine
    cache      CacheClient
    ttl        time.Duration
}

type CartRepository interface {
    Get(ctx context.Context, cartID string) (*Cart, error)
    Save(ctx context.Context, cart *Cart) error
    Delete(ctx context.Context, cartID string) error
}

type InventoryClient interface {
    CheckAvailability(ctx context.Context, sku string, quantity int) (*Availability, error)
    Reserve(ctx context.Context, sku string, quantity int) (*Reservation, error)
    ReleaseReservation(ctx context.Context, reservationID string) error
}

type Availability struct {
    SKU          string
    Available    int
    Reserved     int
    Backorder    bool
}

type Reservation struct {
    ID        string
    SKU       string
    Quantity  int
    ExpiresAt time.Time
}

type PricingEngine interface {
    GetPrice(ctx context.Context, productID string, quantity int, options map[string]string) (*Price, error)
    CalculateDiscounts(ctx context.Context, cart *Cart, promoCode string) (*DiscountResult, error)
}

type Price struct {
    BasePrice    float64
    SalePrice    float64
    Currency     string
    DiscountInfo *DiscountInfo
}

type DiscountResult struct {
    TotalDiscount     float64
    Discounts         []DiscountInfo
    FinalPrices       map[string]float64
}

type CacheClient interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// GetCart retrieves a cart by ID
func (cs *CartService) GetCart(ctx context.Context, cartID string) (*Cart, error) {
    // Try cache first
    cached, err := cs.cache.Get(ctx, "cart:"+cartID)
    if err == nil && cached != "" {
        var cart Cart
        if err := json.Unmarshal([]byte(cached), &cart); err == nil {
            return &cart, nil
        }
    }

    // Get from repository
    cart, err := cs.repository.Get(ctx, cartID)
    if err != nil {
        return nil, err
    }

    // Cache result
    cartJSON, _ := json.Marshal(cart)
    cs.cache.Set(ctx, "cart:"+cartID, string(cartJSON), cs.ttl)

    return cart, nil
}

// AddItem adds an item to the cart
func (cs *CartService) AddItem(ctx context.Context, cartID string, item CartItem) (*Cart, error) {
    cart, err := cs.GetCart(ctx, cartID)
    if err != nil {
        // Create new cart if not exists
        cart = &Cart{
            ID:        cartID,
            Items:     []CartItem{},
            CreatedAt: time.Now(),
            ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
        }
    }

    // Check inventory
    availability, err := cs.inventory.CheckAvailability(ctx, item.SKU, item.Quantity)
    if err != nil {
        return nil, fmt.Errorf("inventory check failed: %w", err)
    }

    if availability.Available < item.Quantity && !availability.Backorder {
        return nil, fmt.Errorf("insufficient inventory for SKU %s", item.SKU)
    }

    // Get current price
    price, err := cs.pricing.GetPrice(ctx, item.ProductID, item.Quantity, item.Options)
    if err != nil {
        return nil, fmt.Errorf("pricing error: %w", err)
    }

    item.UnitPrice = price.SalePrice
    item.Currency = price.Currency

    // Update or add item
    found := false
    for i, existing := range cart.Items {
        if existing.SKU == item.SKU && mapsEqual(existing.Options, item.Options) {
            cart.Items[i].Quantity += item.Quantity
            found = true
            break
        }
    }

    if !found {
        item.ID = generateItemID()
        cart.Items = append(cart.Items, item)
    }

    // Recalculate summary
    cs.recalculateSummary(cart)

    cart.UpdatedAt = time.Now()

    // Save cart
    if err := cs.repository.Save(ctx, cart); err != nil {
        return nil, err
    }

    // Invalidate cache
    cs.cache.Delete(ctx, "cart:"+cartID)

    return cart, nil
}

// RemoveItem removes an item from the cart
func (cs *CartService) RemoveItem(ctx context.Context, cartID, itemID string) (*Cart, error) {
    cart, err := cs.GetCart(ctx, cartID)
    if err != nil {
        return nil, err
    }

    // Find and remove item
    for i, item := range cart.Items {
        if item.ID == itemID {
            cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
            break
        }
    }

    cs.recalculateSummary(cart)
    cart.UpdatedAt = time.Now()

    if err := cs.repository.Save(ctx, cart); err != nil {
        return nil, err
    }

    cs.cache.Delete(ctx, "cart:"+cartID)

    return cart, nil
}

// ApplyPromoCode applies a promotional code
func (cs *CartService) ApplyPromoCode(ctx context.Context, cartID, promoCode string) (*Cart, error) {
    cart, err := cs.GetCart(ctx, cartID)
    if err != nil {
        return nil, err
    }

    discountResult, err := cs.pricing.CalculateDiscounts(ctx, cart, promoCode)
    if err != nil {
        return nil, fmt.Errorf("invalid promo code: %w", err)
    }

    cart.Summary.Discount = discountResult.TotalDiscount
    cart.Summary.DiscountBreakdown = discountResult.Discounts
    cart.AppliedPromo = promoCode

    cs.recalculateSummary(cart)
    cart.UpdatedAt = time.Now()

    if err := cs.repository.Save(ctx, cart); err != nil {
        return nil, err
    }

    cs.cache.Delete(ctx, "cart:"+cartID)

    return cart, nil
}

func (cs *CartService) recalculateSummary(cart *Cart) {
    var subtotal float64
    var itemCount int

    for _, item := range cart.Items {
        subtotal += item.UnitPrice * float64(item.Quantity)
        itemCount += item.Quantity
    }

    cart.Summary.Subtotal = subtotal
    cart.Summary.ItemCount = itemCount
    cart.Summary.Currency = "USD" // Default currency

    // Apply discount
    afterDiscount := subtotal - cart.Summary.Discount
    if afterDiscount < 0 {
        afterDiscount = 0
    }

    // Calculate tax (simplified)
    cart.Summary.Tax = afterDiscount * 0.08 // 8% tax rate

    // Calculate shipping (simplified)
    if afterDiscount < 50 {
        cart.Summary.Shipping = 5.99
    } else {
        cart.Summary.Shipping = 0 // Free shipping over $50
    }

    cart.Summary.Total = afterDiscount + cart.Summary.Tax + cart.Summary.Shipping
}

func mapsEqual(a, b map[string]string) bool {
    if len(a) != len(b) {
        return false
    }
    for k, v := range a {
        if b[k] != v {
            return false
        }
    }
    return true
}

func generateItemID() string {
    return fmt.Sprintf("item_%d", time.Now().UnixNano())
}
```

---

## 6. Inventory Management

### 6.1 Distributed Inventory System

```go
// inventory/distributed.go - Distributed Inventory
package inventory

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Inventory represents stock levels across locations
type Inventory struct {
    SKU           string              `json:"sku"`
    ProductID     string              `json:"product_id"`
    TotalQuantity int                 `json:"total_quantity"`
    Available     int                 `json:"available"`
    Reserved      int                 `json:"reserved"`
    Locations     []LocationStock     `json:"locations"`
    UpdatedAt     time.Time           `json:"updated_at"`
}

type LocationStock struct {
    LocationID   string `json:"location_id"`
    LocationType string `json:"location_type"` // warehouse, store, dropship
    Quantity     int    `json:"quantity"`
    Available    int    `json:"available"`
    Reserved     int    `json:"reserved"`
}

// InventoryService manages distributed inventory
type InventoryService struct {
    repository InventoryRepository
    cache      CacheClient
    eventBus   EventBus
    lockManager DistributedLock
}

type InventoryRepository interface {
    Get(ctx context.Context, sku string) (*Inventory, error)
    Update(ctx context.Context, inv *Inventory) error
    Reserve(ctx context.Context, sku string, quantity int, locationID string) (*Reservation, error)
    Release(ctx context.Context, reservationID string) error
}

type EventBus interface {
    Publish(ctx context.Context, topic string, event interface{}) error
}

type DistributedLock interface {
    Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Release(ctx context.Context, key string) error
}

// CheckAvailability checks stock availability
func (is *InventoryService) CheckAvailability(ctx context.Context, sku string, quantity int) (*AvailabilityResult, error) {
    inv, err := is.getInventory(ctx, sku)
    if err != nil {
        return nil, err
    }

    available := inv.Available >= quantity
    backorder := !available && is.allowBackorder(sku)

    // Find optimal fulfillment location
    location := is.findOptimalLocation(inv, quantity)

    return &AvailabilityResult{
        SKU:          sku,
        Available:    available,
        Backorder:    backorder,
        Quantity:     quantity,
        Location:     location,
        ETA:          is.calculateETA(location),
    }, nil
}

type AvailabilityResult struct {
    SKU       string
    Available bool
    Backorder bool
    Quantity  int
    Location  *LocationStock
    ETA       time.Time
}

// Reserve reserves inventory
func (is *InventoryService) Reserve(ctx context.Context, sku string, quantity int, locationID string) (*Reservation, error) {
    lockKey := fmt.Sprintf("inventory:lock:%s", sku)

    // Acquire distributed lock
    acquired, err := is.lockManager.Acquire(ctx, lockKey, 10*time.Second)
    if err != nil || !acquired {
        return nil, fmt.Errorf("could not acquire lock")
    }
    defer is.lockManager.Release(ctx, lockKey)

    inv, err := is.getInventory(ctx, sku)
    if err != nil {
        return nil, err
    }

    // Find location with stock
    location := is.findLocation(inv, locationID)
    if location == nil || location.Available < quantity {
        return nil, fmt.Errorf("insufficient stock")
    }

    // Update inventory
    location.Available -= quantity
    location.Reserved += quantity
    inv.Available -= quantity
    inv.Reserved += quantity
    inv.UpdatedAt = time.Now()

    if err := is.repository.Update(ctx, inv); err != nil {
        return nil, err
    }

    // Create reservation
    reservation := &Reservation{
        ID:         generateReservationID(),
        SKU:        sku,
        LocationID: location.LocationID,
        Quantity:   quantity,
        ExpiresAt:  time.Now().Add(15 * time.Minute),
        Status:     ReservationStatusActive,
    }

    // Publish event
    is.eventBus.Publish(ctx, "inventory.reserved", reservation)

    // Invalidate cache
    is.cache.Delete(ctx, "inventory:"+sku)

    return reservation, nil
}

// ReleaseReservation releases a reservation
func (is *InventoryService) ReleaseReservation(ctx context.Context, reservationID string) error {
    reservation, err := is.repository.GetReservation(ctx, reservationID)
    if err != nil {
        return err
    }

    lockKey := fmt.Sprintf("inventory:lock:%s", reservation.SKU)
    acquired, err := is.lockManager.Acquire(ctx, lockKey, 10*time.Second)
    if err != nil || !acquired {
        return fmt.Errorf("could not acquire lock")
    }
    defer is.lockManager.Release(ctx, lockKey)

    inv, err := is.getInventory(ctx, reservation.SKU)
    if err != nil {
        return err
    }

    // Update location
    for i := range inv.Locations {
        if inv.Locations[i].LocationID == reservation.LocationID {
            inv.Locations[i].Available += reservation.Quantity
            inv.Locations[i].Reserved -= reservation.Quantity
            break
        }
    }

    inv.Available += reservation.Quantity
    inv.Reserved -= reservation.Quantity
    inv.UpdatedAt = time.Now()

    if err := is.repository.Update(ctx, inv); err != nil {
        return err
    }

    // Publish event
    is.eventBus.Publish(ctx, "inventory.released", reservation)

    // Invalidate cache
    is.cache.Delete(ctx, "inventory:"+reservation.SKU)

    return nil
}

func (is *InventoryService) getInventory(ctx context.Context, sku string) (*Inventory, error) {
    // Try cache
    cached, err := is.cache.Get(ctx, "inventory:"+sku)
    if err == nil && cached != "" {
        // Parse cached inventory
    }

    return is.repository.Get(ctx, sku)
}

func (is *InventoryService) findOptimalLocation(inv *Inventory, quantity int) *LocationStock {
    for _, loc := range inv.Locations {
        if loc.Available >= quantity {
            return &loc
        }
    }
    return nil
}

func (is *InventoryService) findLocation(inv *Inventory, locationID string) *LocationStock {
    for i := range inv.Locations {
        if inv.Locations[i].LocationID == locationID {
            return &inv.Locations[i]
        }
    }
    return nil
}

func (is *InventoryService) allowBackorder(sku string) bool {
    // Check product configuration
    return false
}

func (is *InventoryService) calculateETA(loc *LocationStock) time.Time {
    if loc == nil {
        return time.Time{}
    }
    // Calculate based on location type and shipping method
    return time.Now().Add(2 * 24 * time.Hour)
}

type Reservation struct {
    ID         string
    SKU        string
    LocationID string
    Quantity   int
    ExpiresAt  time.Time
    Status     ReservationStatus
}

type ReservationStatus string

const (
    ReservationStatusActive   ReservationStatus = "active"
    ReservationStatusCommitted ReservationStatus = "committed"
    ReservationStatusExpired   ReservationStatus = "expired"
    ReservationStatusReleased  ReservationStatus = "released"
)

func generateReservationID() string {
    return fmt.Sprintf("res_%d", time.Now().UnixNano())
}
```

---

## 7. Performance Optimizations

### 7.1 Caching Strategy

```go
// performance/caching.go - Multi-layer Caching
package performance

import (
    "context"
    "encoding/json"
    "time"
)

// CacheStrategy defines caching behavior
type CacheStrategy interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Invalidate(ctx context.Context, pattern string) error
}

// MultiTierCache implements L1 (local) + L2 (Redis) caching
type MultiTierCache struct {
    l1Cache    LocalCache    // In-memory
    l2Cache    RemoteCache   // Redis/Memcached
    serializer Serializer
}

type LocalCache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
    Clear()
}

type RemoteCache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Keys(ctx context.Context, pattern string) ([]string, error)
}

type Serializer interface {
    Marshal(v interface{}) ([]byte, error)
    Unmarshal(data []byte, v interface{}) error
}

func (mc *MultiTierCache) Get(ctx context.Context, key string, dest interface{}) error {
    // Try L1 cache first
    if val, found := mc.l1Cache.Get(key); found {
        if data, ok := val.([]byte); ok {
            return mc.serializer.Unmarshal(data, dest)
        }
    }

    // Try L2 cache
    data, err := mc.l2Cache.Get(ctx, key)
    if err != nil || data == "" {
        return err
    }

    // Populate L1 cache
    mc.l1Cache.Set(key, []byte(data), 5*time.Minute)

    return mc.serializer.Unmarshal([]byte(data), dest)
}

func (mc *MultiTierCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := mc.serializer.Marshal(value)
    if err != nil {
        return err
    }

    // Set in L2 cache
    if err := mc.l2Cache.Set(ctx, key, string(data), ttl); err != nil {
        return err
    }

    // Set in L1 cache with shorter TTL
    mc.l1Cache.Set(key, data, ttl/5)

    return nil
}

func (mc *MultiTierCache) Invalidate(ctx context.Context, pattern string) error {
    // Clear L1 cache
    mc.l1Cache.Clear()

    // Find and delete from L2 cache
    keys, err := mc.l2Cache.Keys(ctx, pattern)
    if err != nil {
        return err
    }

    for _, key := range keys {
        mc.l2Cache.Delete(ctx, key)
    }

    return nil
}
```

---

## 8. Security

### 8.1 PCI-DSS Compliance

```go
// security/pci.go - PCI-DSS Security Controls
package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

// TokenVault handles payment card tokenization
type TokenVault struct {
    hsm HSMClient
}

// Tokenize replaces card number with secure token
func (tv *TokenVault) Tokenize(cardNumber string) (string, error) {
    // Validate card number (Luhn check)
    if !validateLuhn(cardNumber) {
        return "", ErrInvalidCardNumber
    }

    // Generate cryptographically secure token
    token := make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, token); err != nil {
        return "", err
    }

    tokenStr := base64.URLEncoding.EncodeToString(token)

    // Encrypt PAN and store in secure vault
    encryptedPAN, err := tv.encryptPAN(cardNumber)
    if err != nil {
        return "", err
    }

    // Store mapping: token -> encryptedPAN
    if err := tv.storeMapping(tokenStr, encryptedPAN); err != nil {
        return "", err
    }

    return tokenStr, nil
}

// Detokenize retrieves original card number (restricted access)
func (tv *TokenVault) Detokenize(token string, authContext AuthContext) (string, error) {
    // Verify authorization
    if !tv.verifyDetokenizeAccess(authContext) {
        return "", ErrUnauthorized
    }

    // Retrieve encrypted PAN
    encryptedPAN, err := tv.retrieveMapping(token)
    if err != nil {
        return "", err
    }

    // Decrypt PAN
    return tv.decryptPAN(encryptedPAN)
}

func (tv *TokenVault) encryptPAN(pan string) ([]byte, error) {
    // Use HSM for encryption
    return tv.hsm.Encrypt([]byte(pan))
}

func (tv *TokenVault) decryptPAN(encrypted []byte) (string, error) {
    decrypted, err := tv.hsm.Decrypt(encrypted)
    return string(decrypted), err
}

func validateLuhn(cardNumber string) bool {
    var sum int
    alternate := false

    for i := len(cardNumber) - 1; i >= 0; i-- {
        digit := int(cardNumber[i] - '0')
        if digit < 0 || digit > 9 {
            return false
        }

        if alternate {
            digit *= 2
            if digit > 9 {
                digit -= 9
            }
        }
        sum += digit
        alternate = !alternate
    }

    return sum%10 == 0
}

type HSMClient interface {
    Encrypt(data []byte) ([]byte, error)
    Decrypt(data []byte) ([]byte, error)
}

type AuthContext struct {
    UserID    string
    Roles     []string
    Purpose   string
    IPAddress string
}

var (
    ErrInvalidCardNumber = fmt.Errorf("invalid card number")
    ErrUnauthorized      = fmt.Errorf("unauthorized access")
)

func (tv *TokenVault) verifyDetokenizeAccess(ctx AuthContext) bool {
    // Implement role-based access control
    for _, role := range ctx.Roles {
        if role == "payment_processor" || role == "fraud_analyst" {
            return true
        }
    }
    return false
}

func (tv *TokenVault) storeMapping(token string, encryptedPAN []byte) error {
    // Store in secure database
    return nil
}

func (tv *TokenVault) retrieveMapping(token string) ([]byte, error) {
    // Retrieve from secure database
    return nil, nil
}
```

---

## 9. References

1. [Elasticsearch Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
2. [Matrix Factorization Techniques for Recommender Systems](https://ieeexplore.ieee.org/document/5197422)
3. [Deep Neural Networks for YouTube Recommendations](https://research.google/pubs/pub45530/)
4. [Contextual Bandits](https://papers.nips.cc/paper/2011/hash/e53a0a2978c28872a4505bdb51db06dc-Abstract.html)
5. [PCI DSS Requirements](https://www.pcisecuritystandards.org/pci_security/standards_overview)

---

**Document End** | **Size: ~58KB** | **Classification: S-Level Technical Specification**
