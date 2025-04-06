package session

import (
	"github.com/Rasikrr/core/enum"
)

type Session struct {
	userID string
	email  string
	role   enum.AccountRole
	claims map[string]any
}

func NewSession(id, email string, role enum.AccountRole, claims map[string]any) *Session {
	return &Session{
		userID: id,
		email:  email,
		role:   role,
		claims: claims,
	}
}

func (s *Session) SetUserID(id string) {
	s.userID = id
}

func (s *Session) SetEmail(email string) {
	s.email = email
}

func (s *Session) SetRole(role enum.AccountRole) {
	s.role = role
}

func (s *Session) Email() string {
	return s.email
}

func (s *Session) UserID() string {
	return s.userID
}

func (s *Session) Claims() map[string]any {
	return s.claims
}

func (s *Session) AccountRole() enum.AccountRole {
	return s.role
}

func (s *Session) SetClaim(k string, v any) {
	if s.claims == nil {
		s.claims = make(map[string]any)
	}
	s.claims[k] = v
}

func (s *Session) SetClaims(claims map[string]any) {
	s.claims = claims
}
