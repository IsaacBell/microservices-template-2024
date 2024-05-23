package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

var (
	onceGetSearchClient sync.Once
	client              *elasticsearch.Client
)

type SearchInterface[T any] interface {
	CreateIndex(idx string) (string, error)
	Search(idx string, filters map[string]interface{}) ([]T, error)
	IndexDocument(indexName string, doc interface{}) error
}

type Search[T any] struct {
	SearchInterface[T]
	results []T
}

func getSearchClient() *elasticsearch.Client {
	onceGetSearchClient.Do(func() {
		res, err := elasticsearch.NewClient(elasticsearch.Config{
			CloudID: os.Getenv("ES_CloudID"),
			APIKey:  os.Getenv("ES_ApiKey"),
		})

		if err != nil {
			log.Fatalln("Error creating search client", err)
		}

		client = res
	})

	return client
}

func (s *Search[T]) CreateIndex(idx string) (string, error) {
	res, err := client.Indices.Create(idx)
	if err != nil {
		return "", err
	}
	if res.IsError() {
		return "", errors.New(res.Status())
	}
	return res.Status(), nil
}

func (s *Search[T]) Search(idx string, filters map[string]interface{}) ([]T, error) {
	var query strings.Builder
	query.WriteString(`{"query": {"bool": {"must": [`)

	var filterParts []string
	for field, value := range filters {
		filterParts = append(filterParts, fmt.Sprintf(`{"match": {"%s": "%v"}}`, field, value))
	}

	query.WriteString(strings.Join(filterParts, ","))
	query.WriteString(`]}}}`)

	res, err := client.Search(
		client.Search.WithIndex(idx),
		client.Search.WithBody(strings.NewReader(query.String())),
	)
	if err != nil {
		return nil, err
	}

	var results []T
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		return nil, err
	}

	return results, nil
}

// example usage:
// err := search.IndexDocument("users", &biz.User{...})
func (s *Search[T]) IndexDocument(indexName string, doc interface{}) error {
	client := getSearchClient()

	indexData := struct {
		Document interface{} `json:"document"`
	}{
		// why not this?
		// "go-elasticsearch",
		Document: doc,
	}

	data, err := json.Marshal(indexData)
	if err != nil {
		return err
	}

	res, err := client.Index(indexName, bytes.NewReader(data))
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.Status())
	}

	return nil
}

func (s *Search[T]) DeleteDocument(idx, id string) error {
	res, err := client.Delete(idx, id)
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.Status())
	}
	return nil
}

func (s *Search[T]) DeleteIndex(idx string) error {
	res, err := client.Indices.Delete([]string{idx})
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.Status())
	}
	return nil
}

func (s *Search[T]) UpdateDocument(idx, id, input string) error {
	// example input: `{ language: "Go" }`
	docStr := fmt.Sprintf(`{doc: %s}`, input)
	res, err := client.Update(idx, id, strings.NewReader(docStr))
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.Status())
	}
	return nil
}
