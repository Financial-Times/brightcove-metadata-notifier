package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//used for convenience
type mapping struct {
	key   string
	value term
}

func fetchMappings(mappingURL string) map[string]term {
	resp, err := http.Get(mappingURL)
	if err != nil {
		errorLogger.Panicf("Couldn't fetch mappings: [%#v]", err)
	}
	defer cleanupResp(resp)
	if resp.StatusCode != 200 {
		errorLogger.Panicf("Unhealthy status code received: [%v]", resp.StatusCode)
	}

	var entries []map[string]string
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		errorLogger.Panicf("Couldn't decode mappings: [%#v]", err)
	}

	infoLogger.Printf("Processing mappings...\n")
	mappings := make(map[string]term, 0)
	for _, entry := range entries {
		infoLogger.Printf("%v\n", entry)
		mapping, err := processMapping(entry)
		if err != nil {
			errorLogger.Println(err)
			continue
		}
		mappings[mapping.key] = mapping.value
	}
	return mappings
}

func processMapping(entry map[string]string) (*mapping, error) {
	bcTag, present := entry["brightcovesearchterm"]
	if !present {
		return nil, fmt.Errorf("Couldn't found brightcoveSearchTerm in mapping: [%+v]", entry)
	}

	streamURL, present := entry["streamurl"]
	if !present {
		return nil, fmt.Errorf("Couldn't found streamURL in mapping: [%+v]", entry)
	}

	i := strings.LastIndex(streamURL, "/")
	if i == -1 || i == len(streamURL)-1 {
		return nil, fmt.Errorf("Couldn't parse TME ID from streamURL: [%s]", streamURL)
	}
	termID := streamURL[i+1:]

	taxonomy, err := decodeTaxonomy(termID)
	if err != nil {
		return nil, err
	}
	return &mapping{
		key: bcTag,
		value: term{
			ID:       termID,
			Taxonomy: taxonomy,
		},
	}, nil
}

func decodeTaxonomy(termID string) (string, error) {
	i := strings.LastIndex(termID, "-")
	if i == -1 || i == len(termID)-1 {
		return "", fmt.Errorf("Couldn't decode taxonomy name from TME ID: [%s]", termID)
	}
	decoded, err := base64.StdEncoding.DecodeString(termID[i+1:])
	if err != nil {
		return "", fmt.Errorf("Couldn't decode taxonomy name from TME ID: [%s]. [%v]", termID, err)
	}
	return string(decoded), nil
}

func (mm metadataMapper) prettyPrintMappings() string {
	s := fmt.Sprintf("metadataMapper.mappings: [\n")
	for _, entry := range mm.mappings {
		s += fmt.Sprintf("\tCanonicalName: [%s], ID: [%s], Taxonomy: [%s]\n", entry.CanonicalName, entry.ID, entry.Taxonomy)
	}
	s += fmt.Sprintf("]\n")
	return s
}
