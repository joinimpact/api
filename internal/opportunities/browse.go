package opportunities

// Section represents a section from the browse section.
type Section struct {
	Name          string            `json:"name"`
	Tag           string            `json:"tag,omitempty"`
	Opportunities []OpportunityView `json:"opportunities"`
}
