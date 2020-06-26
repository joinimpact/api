package organizations

import (
	"io"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a provider of Organization services.
type Service interface {
	// GetOrganizationProfile retrieves a single user's profile.
	// GetOrganizationProfile(organizationID int64) (*OrganizationProfile, error)
	// UpdateOrganizationProfile updates a user's profile.
	// UpdateOrganizationProfile(userID int64, profile OrganizationProfile) error

	// CreateOrganization creates a new organization and returns the ID on success.
	CreateOrganization(organization models.Organization) (int64, error)
	// GetOrganizationTags gets all of a user's tags.
	GetOrganizationTags(organizationID int64) ([]models.Tag, error)
	// AddOrganizationTags adds tags to a user by tag name.
	AddOrganizationTags(organizationID int64, tags []string) (int, error)
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
	UploadProfilePicture(organizationID int64, fileReader io.Reader) error
	// RemoveOrganizationTag removes a tag from a user by id.
	RemoveOrganizationTag(organizationID, tagID int64) error
}

// service represents the internal implementation of the organizations Service.
type service struct {
	organizationRepository           models.OrganizationRepository
	organizationMembershipRepository models.OrganizationMembershipRepository
	organizationTagRepository        models.OrganizationTagRepository
	tagRepository                    models.TagRepository
	config                           *config.Config
	logger                           *zerolog.Logger
	snowflakeService                 snowflakes.SnowflakeService
	cdnClient                        *cdn.Client
}

// NewService creates and returns a new Users service with the provifded dependencies.
func NewService(organizationRepository models.OrganizationRepository, organizationMembershipRepository models.OrganizationMembershipRepository, organizationTagRepository models.OrganizationTagRepository,
	tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService) Service {
	return &service{
		organizationRepository,
		organizationMembershipRepository,
		organizationTagRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		cdn.NewCDNClient(config),
	}
}
