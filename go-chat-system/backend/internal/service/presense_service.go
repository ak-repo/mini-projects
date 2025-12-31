package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type PresenceRepository interface {
	SetUserOnline(userID string, ttl time.Duration) error
	SetUserOffline(userID string) error
	IsUserOnline(userID string) (bool, error)
	SetLastSeen(userID string) error
	GetLastSeen(userID string) (*time.Time, error)
}

type PresenceService struct {
	presenceRepo PresenceRepository
}

func NewPresenceService(presenceRepo PresenceRepository) *PresenceService {
	return &PresenceService{presenceRepo: presenceRepo}
}

func (s *PresenceService) SetOnline(userID string) error {
	return s.presenceRepo.SetUserOnline(userID, 30*time.Second)
}

func (s *PresenceService) SetOffline(userID string) error {
	if err := s.presenceRepo.SetLastSeen(userID); err != nil {
		return err
	}
	return s.presenceRepo.SetUserOffline(userID)
}

func (s *PresenceService) IsOnline(userID string) (bool, error) {
	return s.presenceRepo.IsUserOnline(userID)
}

func (s *PresenceService) GetLastSeen(userID string) (*time.Time, error) {
	return s.presenceRepo.GetLastSeen(userID)
}

// internal/service/signaling_service.go
type SignalingService struct {
	cacheRepo CacheRepository
}

type CacheRepository interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

func NewSignalingService(cacheRepo CacheRepository) *SignalingService {
	return &SignalingService{cacheRepo: cacheRepo}
}

func (s *SignalingService) StoreOffer(callID, offer string) error {
	return s.cacheRepo.Set("call:offer:"+callID, offer, 5*time.Minute)
}

func (s *SignalingService) GetOffer(callID string) (string, error) {
	return s.cacheRepo.Get("call:offer:" + callID)
}

func (s *SignalingService) StoreAnswer(callID, answer string) error {
	return s.cacheRepo.Set("call:answer:"+callID, answer, 5*time.Minute)
}

func (s *SignalingService) GetAnswer(callID string) (string, error) {
	return s.cacheRepo.Get("call:answer:" + callID)
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
