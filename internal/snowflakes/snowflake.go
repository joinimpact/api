package snowflakes

import "github.com/bwmarrin/snowflake"

// Node is the snowflake node number. I'll leave it static for now, but it might change in the future.
const Node = 1

// SnowflakeService is a service meant to generate snowflake IDs.
type SnowflakeService interface {
	// GenerateID generates and returns a new Snowflake ID at the current timestamp.
	GenerateID() int64
}

// snowflakeService is the default version of the SnowflakeService.
type snowflakeService struct {
	node *snowflake.Node
}

// NewSnowflakeService creates and returns a new SnowflakeService.
func NewSnowflakeService() (SnowflakeService, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	return &snowflakeService{node}, nil
}

// GenerateID generates and returns a new snowflake ID.
func (s *snowflakeService) GenerateID() int64 {
	return s.node.Generate().Int64()
}
