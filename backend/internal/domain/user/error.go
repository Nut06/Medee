package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user: user not found")
	ErrSkillNotFound      = errors.New("user: skill not found in profile")
	ErrSkillAlreadyAdded  = errors.New("user: skill already added to profile")
	ErrSkillRequired      = errors.New("user: skill id or name is required")
	ErrInvalidRequestBody = errors.New("user: invalid request body")
	ErrExperienceNotFound = errors.New("user: experience not found")
	ErrEducationNotFound  = errors.New("user: education not found")
	ErrProjectNotFound    = errors.New("user: project not found")
	ErrFileUploadFailed   = errors.New("user: file upload failed")
	ErrInvalidFileType    = errors.New("user: invalid file type")
	ErrFileSizeTooLarge   = errors.New("user: file size too large")
)
