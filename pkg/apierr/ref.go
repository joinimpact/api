package apierr

// Ref gets an errors reference.
func Ref(errInterface interface{}) string {
	err, ok := errInterface.(RefCompatible)
	if !ok {
		// Fallback error
		return "generic.server_error"
	}

	return err.Ref()
}
