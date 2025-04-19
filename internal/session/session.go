package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session
type Session struct {
	ID        string                 `json:"id"`
	UserID    uint                   `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// SessionStore interface defines the methods that all session store implementations must provide
type SessionStore interface {
	// Create creates a new session
	Create(ctx context.Context, userID uint, data map[string]interface{}, expiration time.Duration) (*Session, error)

	// Get retrieves a session by ID
	Get(ctx context.Context, id string) (*Session, error)

	// Update updates a session
	Update(ctx context.Context, session *Session) error

	// Delete deletes a session
	Delete(ctx context.Context, id string) error

	// Cleanup removes expired sessions
	Cleanup(ctx context.Context) error
}

// RedisSessionStore implements the SessionStore interface using Redis
type RedisSessionStore struct {
	client *redis.Client
	prefix string
}

// NewRedisSessionStore creates a new Redis session store
func NewRedisSessionStore(host string, port int, password string, db int, prefix string) (*RedisSessionStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisSessionStore{
		client: client,
		prefix: prefix,
	}, nil
}

// Create creates a new session in Redis
func (s *RedisSessionStore) Create(ctx context.Context, userID uint, data map[string]interface{}, expiration time.Duration) (*Session, error) {
	id := uuid.New().String()
	now := time.Now()

	session := &Session{
		ID:        id,
		UserID:    userID,
		Data:      data,
		CreatedAt: now,
		ExpiresAt: now.Add(expiration),
	}

	// Serialize the session
	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	// Store the session in Redis
	key := fmt.Sprintf("%s:%s", s.prefix, id)
	if err := s.client.Set(ctx, key, sessionData, expiration).Err(); err != nil {
		return nil, err
	}

	// Store a user ID to session ID mapping for quick lookups
	userKey := fmt.Sprintf("%s:user:%d", s.prefix, userID)
	if err := s.client.SAdd(ctx, userKey, id).Err(); err != nil {
		return nil, err
	}

	// Set expiration on the user key
	if err := s.client.Expire(ctx, userKey, expiration).Err(); err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a session from Redis
func (s *RedisSessionStore) Get(ctx context.Context, id string) (*Session, error) {
	key := fmt.Sprintf("%s:%s", s.prefix, id)

	data, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Update updates a session in Redis
func (s *RedisSessionStore) Update(ctx context.Context, session *Session) error {
	// Serialize the session
	sessionData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Store the session in Redis
	key := fmt.Sprintf("%s:%s", s.prefix, session.ID)
	expiration := time.Until(session.ExpiresAt)

	if err := s.client.Set(ctx, key, sessionData, expiration).Err(); err != nil {
		return err
	}

	return nil
}

// Delete deletes a session from Redis
func (s *RedisSessionStore) Delete(ctx context.Context, id string) error {
	// Get the session to find the user ID
	session, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}

	// Delete the session
	key := fmt.Sprintf("%s:%s", s.prefix, id)
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	// Remove the session ID from the user's set
	userKey := fmt.Sprintf("%s:user:%d", s.prefix, session.UserID)
	if err := s.client.SRem(ctx, userKey, id).Err(); err != nil {
		return err
	}

	return nil
}

// Cleanup removes expired sessions from Redis
func (s *RedisSessionStore) Cleanup(ctx context.Context) error {
	// Redis automatically removes expired keys, so no cleanup is needed
	return nil
}

// GetUserSessions retrieves all sessions for a user
func (s *RedisSessionStore) GetUserSessions(ctx context.Context, userID uint) ([]*Session, error) {
	userKey := fmt.Sprintf("%s:user:%d", s.prefix, userID)

	// Get all session IDs for the user
	sessionIDs, err := s.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*Session, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		session, err := s.Get(ctx, id)
		if err != nil {
			continue
		}
		if session != nil {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

// DatabaseSessionStore implements the SessionStore interface using a database
type DatabaseSessionStore struct {
	db *gorm.DB
}

// NewDatabaseSessionStore creates a new database session store
func NewDatabaseSessionStore(db *gorm.DB) *DatabaseSessionStore {
	return &DatabaseSessionStore{
		db: db,
	}
}

// Create creates a new session in the database
func (s *DatabaseSessionStore) Create(ctx context.Context, userID uint, data map[string]interface{}, expiration time.Duration) (*Session, error) {
	id := uuid.New().String()
	now := time.Now()

	session := &Session{
		ID:        id,
		UserID:    userID,
		Data:      data,
		CreatedAt: now,
		ExpiresAt: now.Add(expiration),
	}

	// Store the session in the database
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a session from the database
func (s *DatabaseSessionStore) Get(ctx context.Context, id string) (*Session, error) {
	var session Session

	if err := s.db.WithContext(ctx).Where("id = ? AND expires_at > ?", id, time.Now()).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// Update updates a session in the database
func (s *DatabaseSessionStore) Update(ctx context.Context, session *Session) error {
	return s.db.WithContext(ctx).Save(session).Error
}

// Delete deletes a session from the database
func (s *DatabaseSessionStore) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&Session{}).Error
}

// Cleanup removes expired sessions from the database
func (s *DatabaseSessionStore) Cleanup(ctx context.Context) error {
	return s.db.WithContext(ctx).Where("expires_at <= ?", time.Now()).Delete(&Session{}).Error
}

// GetUserSessions retrieves all sessions for a user
func (s *DatabaseSessionStore) GetUserSessions(ctx context.Context, userID uint) ([]*Session, error) {
	var sessions []*Session

	if err := s.db.WithContext(ctx).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

// SessionStoreFactory creates a session store based on the provided configuration
func SessionStoreFactory(storeType string, config map[string]interface{}) (SessionStore, error) {
	switch storeType {
	case "redis":
		host := config["host"].(string)
		port := config["port"].(int)
		password := config["password"].(string)
		db := config["db"].(int)
		prefix := config["prefix"].(string)

		return NewRedisSessionStore(host, port, password, db, prefix)

	case "database":
		db := config["db"].(*gorm.DB)

		return NewDatabaseSessionStore(db), nil

	default:
		return nil, fmt.Errorf("unsupported session store type: %s", storeType)
	}
}
