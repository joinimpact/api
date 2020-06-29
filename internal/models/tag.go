package models

// Tag categories
const (
	TagCategoryInterest = iota
	TagCategoryJobType  = iota
)

// Tag represents a single tag.
type Tag struct {
	Model
	Searchable bool   `json:"-"`
	Name       string `json:"name" gorm:"unique,not null,index:name"`
	Category   int    `json:"category"`
}

// TagRepository represents a repository of tags.
type TagRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*Tag, error)
	// FindByCategory finds entities by Category.
	FindByCategory(category int) ([]Tag, error)
	// FindByName finds a single entity by name.
	FindByName(name string) (*Tag, error)
	// SearchTags searches for tags with a query string.
	SearchTags(query string, limit int) ([]Tag, error)
	// Create creates a new entity.
	Create(tag Tag) error
	// Update updates an entity with the ID in the provided entity.
	Update(tag Tag) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
