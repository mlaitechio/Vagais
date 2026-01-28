package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model with common fields
// Use string for ID for SQLite compatibility
// Remove type:uuid and default:gen_random_uuid()
type BaseModel struct {
	ID        string     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Organization represents a company or team
type Organization struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	Plan        string `json:"plan" gorm:"default:'free'"`
	Users       []User `json:"users,omitempty" gorm:"foreignKey:OrganizationID"`
}

// User represents a platform user
type User struct {
	BaseModel
	Email          string        `json:"email" gorm:"uniqueIndex;not null"`
	Username       string        `json:"username" gorm:"uniqueIndex;not null"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	PasswordHash   string        `json:"-" gorm:"not null"`
	Role           string        `json:"role" gorm:"default:'user'"`
	IsActive       bool          `json:"is_active" gorm:"default:true"`
	EmailVerified  bool          `json:"email_verified" gorm:"default:false"`
	Avatar         string        `json:"avatar"`
	OrganizationID *string       `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	Credits        int64         `json:"credits" gorm:"default:0"`
	LastLoginAt    *time.Time    `json:"last_login_at"`
	Preferences    JSON          `json:"preferences" gorm:"type:jsonb"`
	Agents         []Agent       `json:"agents,omitempty" gorm:"foreignKey:CreatorID"`
	Reviews        []Review      `json:"reviews,omitempty" gorm:"foreignKey:UserID"`
}

// Agent represents an AI agent
type Agent struct {
	BaseModel
	Name              string        `json:"name" gorm:"not null"`
	Description       string        `json:"description"`
	Slug              string        `json:"slug" gorm:"uniqueIndex;not null"`
	Version           string        `json:"version" gorm:"default:'1.0.0'"`
	Status            string        `json:"status" gorm:"default:'draft'"`
	Type              string        `json:"type" gorm:"default:'custom'"`
	Category          string        `json:"category"`
	Tags              JSON          `json:"tags" gorm:"type:jsonb"`
	Config            JSON          `json:"config" gorm:"type:jsonb"`
	LLMProvider       string        `json:"llm_provider"`
	LLMModel          string        `json:"llm_model"`
	EmbeddingProvider string        `json:"embedding_provider"`
	EmbeddingModel    string        `json:"embedding_model"`
	CreatorID         string        `json:"creator_id"`
	Creator           User          `json:"creator"`
	OrganizationID    *string       `json:"organization_id,omitempty"`
	Organization      *Organization `json:"organization,omitempty"`
	IsPublic          bool          `json:"is_public" gorm:"default:false"`
	IsEnabled         bool          `json:"is_enabled" gorm:"default:false"`
	Price             float64       `json:"price" gorm:"default:0"`
	Currency          string        `json:"currency" gorm:"default:'USD'"`
	PricingModel      string        `json:"pricing_model" gorm:"default:'free'"`
	Rating            float64       `json:"rating" gorm:"default:0"`
	ReviewCount       int           `json:"review_count" gorm:"default:0"`
	UsageCount        int64         `json:"usage_count" gorm:"default:0"`
	Downloads         int           `json:"downloads" gorm:"default:0"`
	Icon              string        `json:"icon"`
	Screenshots       JSON          `json:"screenshots" gorm:"type:jsonb"`
	Documentation     string        `json:"documentation"`
	Repository        string        `json:"repository"`
	VideoURL          string        `json:"video_url"`           // URL to demo/tutorial video
	HowItWorks        string        `json:"how_it_works"`        // Detailed explanation of how the agent works
	FilePath          string        `json:"file_path"`           // Path to the agent file/script
	ExecutablePath    string        `json:"executable_path"`     // Path to the executable or script
	Reviews           []Review      `json:"reviews,omitempty" gorm:"foreignKey:AgentID"`
	Executions        []Execution   `json:"executions,omitempty" gorm:"foreignKey:AgentID"`
}

// Review represents user reviews for agents
type Review struct {
	BaseModel
	AgentID    string `json:"agent_id"`
	Agent      Agent  `json:"agent"`
	UserID     string `json:"user_id"`
	User       User   `json:"user"`
	Rating     int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	IsVerified bool   `json:"is_verified" gorm:"default:false"`
	IsHelpful  int    `json:"is_helpful" gorm:"default:0"`
	Response   string `json:"response,omitempty"`
}

// Execution represents agent execution logs
type Execution struct {
	BaseModel
	AgentID        string        `json:"agent_id"`
	Agent          Agent         `json:"agent"`
	UserID         string        `json:"user_id"`
	User           User          `json:"user"`
	OrganizationID *string       `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	Status         string        `json:"status"`
	Input          JSON          `json:"input" gorm:"type:jsonb"`
	Output         JSON          `json:"output" gorm:"type:jsonb"`
	Error          string        `json:"error,omitempty"`
	Duration       int64         `json:"duration"`
	Cost           float64       `json:"cost" gorm:"default:0"`
	CreditsUsed    int64         `json:"credits_used" gorm:"default:0"`
	IPAddress      string        `json:"ip_address"`
	UserAgent      string        `json:"user_agent"`
	SessionID      string        `json:"session_id"`
}

// Webhook represents webhook configurations
type Webhook struct {
	BaseModel
	Name           string        `json:"name" gorm:"not null"`
	URL            string        `json:"url" gorm:"not null"`
	Events         []string      `json:"events" gorm:"type:jsonb"`
	Secret         string        `json:"secret,omitempty"`
	IsActive       bool          `json:"is_active" gorm:"default:true"`
	OrganizationID *string       `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	UserID         string        `json:"user_id"`
	User           User          `json:"user"`
	LastTriggered  *time.Time    `json:"last_triggered,omitempty"`
	FailureCount   int           `json:"failure_count" gorm:"default:0"`
	Headers        JSON          `json:"headers" gorm:"type:jsonb"`
}

// Notification represents user notifications
type Notification struct {
	BaseModel
	UserID         string        `json:"user_id"`
	User           User          `json:"user"`
	OrganizationID *string       `json:"organization_id,omitempty"`
	Organization   *Organization `json:"organization,omitempty"`
	Type           string        `json:"type"`
	Title          string        `json:"title"`
	Message        string        `json:"message"`
	Status         string        `json:"status" gorm:"default:'unread'"`
	Priority       string        `json:"priority" gorm:"default:'normal'"`
	Category       string        `json:"category"`
	ReadAt         *time.Time    `json:"read_at,omitempty"`
	Metadata       JSON          `json:"metadata" gorm:"type:jsonb"`
}

// LLMProvider represents an LLM provider configuration
type LLMProvider struct {
	BaseModel
	Name      string `json:"name" gorm:"not null"`
	Type      string `json:"type" gorm:"not null"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	RateLimit int    `json:"rate_limit" gorm:"default:1000"`
	MaxTokens int    `json:"max_tokens" gorm:"default:4096"`
	Config    JSON   `json:"config" gorm:"type:jsonb"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	BaseModel
	UserID    string    `json:"user_id"`
	User      User      `json:"user"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used" gorm:"default:false"`
}

// JSON type for storing JSON data (compatible with SQLite and Postgres)
type JSON []byte

// Value implements the driver.Valuer interface for database serialization
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	return []byte(j), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON([]byte("{}"))
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
		return nil
	case string:
		*j = JSON([]byte(v))
		return nil
	default:
		return fmt.Errorf("unsupported type for JSON: %T", value)
	}
}

// Helper to marshal a map to JSON
func MapToJSON(m map[string]interface{}) JSON {
	if m == nil {
		return JSON([]byte("{}"))
	}
	b, _ := json.Marshal(m)
	return JSON(b)
}

// BeforeCreate hook to set timestamps
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// IsFeatureEnabled checks if a feature is enabled based on configuration
func IsFeatureEnabled(feature string, config interface{}) bool {
	// This function can be implemented based on your configuration
	// For now, we'll assume all features are enabled by default
	// In production, this would check against your configuration
	return true
}
