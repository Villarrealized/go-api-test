package models

type ModelMissingError struct {
	Message string
}

func (e *ModelMissingError) Error() string { return e.Message }
