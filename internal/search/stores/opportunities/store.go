package opportunities

// Store represents a storage of opportunities in the Elasticsearch database.
type Store interface {
	// Save saves an opportunity by ID in the Elasticsearch store.
	Save(opportunityID int64) error
}
