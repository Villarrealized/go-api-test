package models

type ModelMissingError struct {
	Message string
}

func (e *ModelMissingError) Error() string { return e.Message }

type ModelMissingRequiredFieldError struct {
	Message string
}

func (e *ModelMissingRequiredFieldError) Error() string { return e.Message }

type ModelRelationshipError struct {
	Message string
}

func (e *ModelRelationshipError) Error() string { return e.Message }

type UniqueViolationError struct {
	Message string
}

func (e *UniqueViolationError) Error() string { return e.Message }
