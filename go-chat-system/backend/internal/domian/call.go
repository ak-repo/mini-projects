package domain

import "time"

type CallType string

const (
	CallTypeAudio CallType = "audio"
	CallTypeVideo CallType = "video"
)

type CallStatus string

const (
	CallStatusInitiated CallStatus = "initiated"
	CallStatusRinging   CallStatus = "ringing"
	CallStatusActive    CallStatus = "active"
	CallStatusEnded     CallStatus = "ended"
	CallStatusMissed    CallStatus = "missed"
	CallStatusRejected  CallStatus = "rejected"
)

type Call struct {
	ID             string
	ConversationID string
	CallerID       string
	CallType       CallType
	Status         CallStatus
	StartedAt      *time.Time
	EndedAt        *time.Time
	CreatedAt      time.Time
}

type CallParticipant struct {
	CallID   string
	UserID   string
	JoinedAt *time.Time
	LeftAt   *time.Time
}
