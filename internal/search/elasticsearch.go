package search

import (
	"fmt"

	"github.com/elastic/go-elasticsearch"
)

// NewElasticsearch creates and returns an elasticsearch.Client from a hostname and port.
func NewElasticsearch(hostname, port string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", hostname, port),
		},
	}

	return elasticsearch.NewClient(cfg)
}
