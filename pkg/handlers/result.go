package handlers

// EventHandlingResult is the outcome of a webhook event handler.
type EventHandlingResult struct {
	Success      bool
	ErrorMessage string
	Err          error
}

// Success returns a successful handling result.
func Success() EventHandlingResult {
	return EventHandlingResult{Success: true}
}

// Failure returns a failed handling result with an optional wrapped error.
func Failure(message string, err ...error) EventHandlingResult {
	res := EventHandlingResult{
		Success:      false,
		ErrorMessage: message,
	}
	if len(err) > 0 {
		res.Err = err[0]
	}
	return res
}

// IsSuccess reports whether handling completed without failure.
func (r EventHandlingResult) IsSuccess() bool {
	return r.Success
}
