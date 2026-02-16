// Package api provides Data Transfer Objects for API request/response handling.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// This file defines all request and response structures for the REST API.
// DTOs decouple the API contract from internal domain models.
//
// =============================================================================
// AGENT DTOs:
// =============================================================================
//
// type AgentListResponse struct {
//     Agents     []AgentSummary `json:"agents"`
//     Pagination PaginationMeta `json:"pagination"`
// }
//
// type AgentSummary struct {
//     ID              string  `json:"id"`
//     Name            string  `json:"name"`
//     PersonalityType string  `json:"personalityType"` // e.g., "introvert", "leader"
//     CurrentMood     string  `json:"currentMood"`     // e.g., "happy", "anxious"
//     MoodIntensity   float64 `json:"moodIntensity"`   // 0.0 - 1.0
//     IsActive        bool    `json:"isActive"`
// }
//
// type AgentDetailResponse struct {
//     ID          string              `json:"id"`
//     Name        string              `json:"name"`
//     Personality PersonalityDTO      `json:"personality"`
//     CurrentMood MoodDTO             `json:"currentMood"`
//     Goals       []GoalDTO           `json:"goals"`
//     Stats       AgentStatsDTO       `json:"stats"`
//     CreatedAt   time.Time           `json:"createdAt"`
// }
//
// type PersonalityDTO struct {
//     Openness          float64 `json:"openness"`          // Big Five trait
//     Conscientiousness float64 `json:"conscientiousness"`
//     Extraversion      float64 `json:"extraversion"`
//     Agreeableness     float64 `json:"agreeableness"`
//     Neuroticism       float64 `json:"neuroticism"`
//     CoreValues        []string `json:"coreValues"`       // e.g., ["honesty", "creativity"]
//     Quirks            []string `json:"quirks"`           // Unique behavioral patterns
// }
//
// =============================================================================
// MEMORY DTOs:
// =============================================================================
//
// type MemoryResponse struct {
//     Memories []MemoryEntryDTO `json:"memories"`
//     Summary  string           `json:"summary"` // AI-generated memory summary
// }
//
// type MemoryEntryDTO struct {
//     ID            string    `json:"id"`
//     Type          string    `json:"type"`          // "episodic" | "semantic"
//     Content       string    `json:"content"`
//     EmotionalTag  string    `json:"emotionalTag"`  // Emotion during memory formation
//     Importance    float64   `json:"importance"`    // 0.0 - 1.0
//     Timestamp     time.Time `json:"timestamp"`
//     RelatedAgents []string  `json:"relatedAgents"` // IDs of involved agents
// }
//
// =============================================================================
// RELATIONSHIP DTOs:
// =============================================================================
//
// type RelationshipGraphResponse struct {
//     Nodes []GraphNodeDTO `json:"nodes"`
//     Edges []GraphEdgeDTO `json:"edges"`
// }
//
// type GraphNodeDTO struct {
//     ID    string `json:"id"`
//     Label string `json:"label"`
//     Type  string `json:"type"`  // Agent personality type for coloring
//     Size  int    `json:"size"`  // Based on relationship count
// }
//
// type GraphEdgeDTO struct {
//     Source   string  `json:"source"`
//     Target   string  `json:"target"`
//     Type     string  `json:"type"`     // "friend", "rival", "neutral", "romantic"
//     Strength float64 `json:"strength"` // -1.0 (hostile) to 1.0 (close)
//     Label    string  `json:"label"`    // Optional description
// }
//
// =============================================================================
// EVENT DTOs:
// =============================================================================
//
// type EventDTO struct {
//     ID             string    `json:"id"`
//     Type           string    `json:"type"`           // "global", "agent", "interaction"
//     Category       string    `json:"category"`       // "disaster", "celebration", etc.
//     Description    string    `json:"description"`
//     AffectedAgents []string  `json:"affectedAgents"`
//     Timestamp      time.Time `json:"timestamp"`
//     Status         string    `json:"status"`         // "pending", "active", "completed"
// }
//
// type InjectEventRequest struct {
//     Type           string   `json:"type" binding:"required"`
//     Description    string   `json:"description" binding:"required"`
//     AffectedAgents []string `json:"affectedAgents"` // empty = all agents
// }
//
// =============================================================================
// WORLD DTOs:
// =============================================================================
//
// type WorldStatusResponse struct {
//     CurrentTick     int64   `json:"currentTick"`
//     SimulationSpeed float64 `json:"simulationSpeed"` // 1.0 = real-time
//     IsPaused        bool    `json:"isPaused"`
//     ActiveAgents    int     `json:"activeAgents"`
//     TotalEvents     int     `json:"totalEvents"`
//     Uptime          string  `json:"uptime"`
// }
//
// type WorldStatisticsResponse struct {
//     MoodDistribution      map[string]int `json:"moodDistribution"`
//     RelationshipStats     RelationshipStatsDTO `json:"relationshipStats"`
//     ActivityMetrics       ActivityMetricsDTO   `json:"activityMetrics"`
//     TopInteractingAgents  []AgentSummary       `json:"topInteractingAgents"`
// }

package api
