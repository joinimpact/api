package opportunities

// OpportunityDocument represents an opportunity as an Elasticsearch NoSQL document.
type OpportunityDocument struct {
	ID          int64                    `json:"_id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Tags        []OpportunityTagDocument `json:"tags"`
}

// OpportunityTagDocument represents an opportunity's tag in the Elasticsearch database.
type OpportunityTagDocument struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category int    `json:"category"`
}
