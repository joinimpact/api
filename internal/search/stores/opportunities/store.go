package opportunities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

const indexName = "opportunities"

// Store represents a storage of opportunities in the Elasticsearch database.
type Store interface {
	// Start starts the queue processor asynchronously.
	Start() error
	// Save saves an opportunity by ID in the Elasticsearch store.
	Save(opportunityID int64)
	// Search searches opportunities, and returns relevant a list of documents.
	Search(query Query) ([]OpportunityDocument, error)
	// Recommendations searches opportunities, and returns relevant a list of documents.
	Recommendations(query RecommendationQuery) ([]OpportunityDocument, error)
}

// store represents the internal implementation of the Store.
type store struct {
	client                            *elasticsearch.Client
	opportunityRepository             models.OpportunityRepository
	opportunityRequirementsRepository models.OpportunityRequirementsRepository
	opportunityLimitsRepository       models.OpportunityLimitsRepository
	opportunityTagRepository          models.OpportunityTagRepository
	tagRepository                     models.TagRepository
	organizationRepository            models.OrganizationRepository
	logger                            *zerolog.Logger
	config                            *config.Config
	saveQueue                         chan queueItem
}

// NewStore creates and returns a new Store with the provided dependencies.
func NewStore(
	client *elasticsearch.Client,
	opportunityRepository models.OpportunityRepository,
	opportunityRequirementsRepository models.OpportunityRequirementsRepository,
	opportunityLimitsRepository models.OpportunityLimitsRepository,
	opportunityTagRepository models.OpportunityTagRepository,
	tagRepository models.TagRepository,
	organizationRepository models.OrganizationRepository,
	logger *zerolog.Logger,
	config *config.Config,
) Store {
	saveQueue := make(chan queueItem, 4)

	return &store{
		client,
		opportunityRepository,
		opportunityRequirementsRepository,
		opportunityLimitsRepository,
		opportunityTagRepository,
		tagRepository,
		organizationRepository,
		logger,
		config,
		saveQueue,
	}
}

// Start starts the queue processor asynchronously.
func (s *store) Start() error {
	// Run the queueWatcher asynchronously with a goroutine.
	go s.queueWatcher()

	return nil
}

// queueWatcher watches the queue for events and processes them.
func (s *store) queueWatcher() {
	for {
		select {
		case item := <-s.saveQueue:
			err := s.save(item.opportunityID)
			if err != nil {
				s.logger.Error().Err(err).Msgf("Error saving opportunity %d", item.opportunityID)
			}
		}
	}
}

// save saves an opportunity by ID into the Elasticsearch store.
func (s *store) save(opportunityID int64) error {
	ctx := context.Background()

	document := &OpportunityDocument{}
	document.Requirements = &Requirements{}
	document.Limits = &Limits{}
	document.Organization = &OpportunityOrganizationDocument{}
	document.Tags = []OpportunityTagDocument{}

	opportunity, err := s.opportunityRepository.FindByID(ctx, opportunityID)
	if err != nil {
		return err
	}

	document.ID = opportunity.ID
	document.Organization.ID = opportunity.OrganizationID
	document.Title = opportunity.Title
	document.Description = opportunity.Description
	document.Public = opportunity.Public

	opportunityRequirements, err := s.opportunityRequirementsRepository.FindByOpportunityID(opportunity.ID)
	if err == nil {
		if opportunityRequirements.AgeLimitActive {
			document.Requirements.AgeLimit = AgeLimit{
				Active: true,
				From:   opportunityRequirements.AgeLimitFrom,
				To:     opportunityRequirements.AgeLimitTo,
			}
		}
		if opportunityRequirements.ExpectedHoursActive {
			document.Requirements.ExpectedHours = ExpectedHours{
				Active: true,
				Hours:  opportunityRequirements.ExpectedHours,
			}
		}
	}

	opportunityLimits, err := s.opportunityLimitsRepository.FindByOpportunityID(opportunity.ID)
	if err == nil {
		if opportunityLimits.VolunteersCapActive {
			document.Limits.VolunteersCap = VolunteersCap{
				Active: true,
				Cap:    opportunityLimits.VolunteersCap,
			}
		}
	}

	tags, err := s.opportunityTagRepository.FindByOpportunityID(opportunityID)
	if err == nil {
		for _, tag := range tags {
			tag, err := s.tagRepository.FindByID(tag.TagID)
			if err != nil {
				continue
			}

			document.Tags = append(document.Tags, OpportunityTagDocument{
				ID:       tag.ID,
				Name:     tag.Name,
				Category: tag.Category,
			})
		}
	}

	organization, err := s.organizationRepository.FindByID(opportunity.OrganizationID)
	if err == nil {
		document.Organization.Name = organization.Name
		if organization.LocationLatitude != 0 || organization.LocationLongitude != 0 {
			document.Location = &LocationDocument{}
			document.Location.Latitude = organization.LocationLatitude
			document.Location.Longitude = organization.LocationLongitude
		}
	}

	payloadMap := map[string]interface{}{
		"doc":           document,
		"doc_as_upsert": true,
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	res, err := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", document.ID),
		Body:       bytes.NewReader(payload),
	}.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

// Save adds an opportunity save event to the queue.
func (s *store) Save(opportunityID int64) {
	// Add the item to the queue.
	s.saveQueue <- queueItem{
		opportunityID: opportunityID,
	}
}

// Search searches opportunities, and returns relevant a list of documents.
func (s *store) Search(query Query) ([]OpportunityDocument, error) {
	documents := []OpportunityDocument{}

	queryReader := buildQuery(query)

	res, err := s.client.Search(
		s.client.Search.WithIndex(indexName),
		s.client.Search.WithBody(queryReader),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	for _, hit := range r.Hits.Hits {
		doc := OpportunityDocument{}

		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// Recommendations searches opportunities, and returns relevant a list of documents.
func (s *store) Recommendations(query RecommendationQuery) ([]OpportunityDocument, error) {
	documents := []OpportunityDocument{}

	queryReader := buildRecommendationQuery(query)

	res, err := s.client.Search(
		s.client.Search.WithIndex(indexName),
		s.client.Search.WithBody(queryReader),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	for _, hit := range r.Hits.Hits {
		doc := OpportunityDocument{}

		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}

		documents = append(documents, doc)
	}

	return documents, nil
}
