package customer

// Provider constants for customer providers
const (
	ProviderLine = "LINE"
)

// IsValidProvider checks if the given provider is valid
func IsValidProvider(provider string) bool {
	switch provider {
	case ProviderLine:
		return true
	default:
		return false
	}
}
