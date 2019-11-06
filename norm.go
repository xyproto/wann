package wann

// NormalizationInfo contains if and how the score function should be normalized
type NormalizationInfo struct {
	shouldNormalize bool
	mul, add        float64
}

// NewNormalizationInfo returns a new struct, containing if and how the score function should be normalized
func NewNormalizationInfo(enable bool) *NormalizationInfo {
	return &NormalizationInfo{enable, 0.0, 1.0}
}

// Enable signifies that normalization is enabled when this struct is used
func (norm *NormalizationInfo) Enable() {
	norm.shouldNormalize = true
}

// Disable signifies that normalization is disabled when this struct is used
func (norm *NormalizationInfo) Disable() {
	norm.shouldNormalize = false
}

// Get retrieves the multiplication and addition numbers that can be used for normalization
func (norm *NormalizationInfo) Get() (float64, float64) {
	return norm.mul, norm.add
}

// Set sets the multiplication and addition numbers that can be used for normalization
func (norm *NormalizationInfo) Set(mul, add float64) {
	norm.mul = mul
	norm.add = add
}
