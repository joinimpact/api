package apierr

// RefCompatible references an error with a Ref method.
type RefCompatible interface {
	Ref() string
}
