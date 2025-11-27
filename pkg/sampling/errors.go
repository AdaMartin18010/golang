package sampling

import "errors"

var (
	// ErrInvalidSampleRate 无效的采样率
	ErrInvalidSampleRate = errors.New("invalid sample rate: must be between 0.0 and 1.0")
)
