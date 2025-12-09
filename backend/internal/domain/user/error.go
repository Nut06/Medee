package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user: user not found")
	ErrSkillNotFound     = errors.New("user: skill not found in profile")
	ErrSkillAlreadyAdded = errors.New("user: skill already added to profile")
	ErrSkillRequired     = errors.New("user: skill id or name is required")
)
