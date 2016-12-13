package main

import "testing"

func TestProcessMapping_EntriesInUnexpectedFormat_ErrorReturned(t *testing.T) {
	var testMappings = []map[string]string{
		map[string]string{
			"streamurl": "/stream/sectionsId/MQ==-U2VjdGlvbnM=", //missing brightcovesearchterm field
		},
		map[string]string{
			"brightcovesearchterm": "tag:section:world", //missing streamurl field
		},
		map[string]string{
			"brightcovesearchterm": "tag:section:world",
			"streamurl":            "MQ==-U2VjdGlvbnM=", //streamurl in unexpected format: no '/' (slash) char to be found
		},
		map[string]string{
			"brightcovesearchterm": "tag:section:world",
			"streamurl":            "/stream/sectionsId/MQ==-U2VjdGlvbnM=/", //streamurl in unexpected format: last '/' should not be the last character
		},
		map[string]string{
			"brightcovesearchterm": "tag:section:world",
			"streamurl":            "/stream/sectionsId/MQ==U2VjdGlvbnM=", //TME ID in unexpected format: missing '-' (dash)
		},
	}

	for _, m := range testMappings {
		if _, err := processMapping(m); err == nil {
			t.Errorf("Expected failure. Testcase: [%+v]", m)
		}
	}
}

func TestProcessMapping_HappyScenarios_BrightcoveSearchTermAndTMEIDAndTaxonomyAreCorrect(t *testing.T) {
	var testCases = []struct {
		mapping  map[string]string
		bcTag    string
		tmeID    string
		taxonomy string
	}{
		{
			mapping: map[string]string{
				"streamurl":            "/stream/sectionsId/MQ==-U2VjdGlvbnM=",
				"brightcovesearchterm": "tag:section:world",
			},
			bcTag:    "section:world",
			tmeID:    "MQ==-U2VjdGlvbnM=",
			taxonomy: "Sections",
		},
		{
			mapping: map[string]string{
				"streamurl":            "/stream/authorsId/Q0ItMDAwMDkyMw==-QXV0aG9ycw==",
				"brightcovesearchterm": "john authers"},
			bcTag:    "john authers",
			tmeID:    "Q0ItMDAwMDkyMw==-QXV0aG9ycw==",
			taxonomy: "Authors",
		},
		{
			mapping: map[string]string{
				"streamurl":            "/stream/regionsId/TnN0ZWluX0dMX0xhdGluQW1lcmljYQ==-R0w=",
				"brightcovesearchterm": "latin-america",
			},
			bcTag:    "latin-america",
			tmeID:    "TnN0ZWluX0dMX0xhdGluQW1lcmljYQ==-R0w=",
			taxonomy: "GL",
		},
	}

	for _, tc := range testCases {
		m, err := processMapping(tc.mapping)
		if err != nil {
			t.Errorf("Expected success. Testcase: [%+v]", tc)
		}
		if m.key != tc.bcTag {
			t.Errorf("Expected: [%s]. Actual: [%s]. Testcase: [%+v]", tc.bcTag, m.key, tc)
		}
		if m.value.ID != tc.tmeID {
			t.Errorf("Expected: [%s]. Actual: [%s]. Testcase: [%+v]", tc.tmeID, m.value.ID, tc)
		}
		if m.value.Taxonomy != tc.taxonomy {
			t.Errorf("Expected: [%s]. Actual: [%s]. Testcase: [%+v]", tc.taxonomy, m.value.Taxonomy, tc)
		}
	}
}

func TestDecodeTaxonomy_InvalidTMEIDs_ErrorReturned(t *testing.T) {
	var testTMEIDs = []string{"MQ==U2VjdGlvbnM=", "MQ==-U2VjdGlvbnM=-", "MQ==-U2VjdGl?vbnM="}

	for _, tmeID := range testTMEIDs {
		if _, err := decodeTaxonomy(tmeID); err == nil {
			t.Errorf("Expected failure. Testcase: [%s]", tmeID)
		}
	}
}

func TestDecodeTaconomy_HappyScenario_DecodedTaxonomyIsCorrect(t *testing.T) {
	validTMEID := "MQ==-U2VjdGlvbnM="
	expected := "Sections"

	actual, err := decodeTaxonomy(validTMEID)

	if err != nil {
		t.Error("Expected success.")
	}
	if actual != expected {
		t.Errorf("Expected: [%s]. Actual: [%s]", expected, actual)
	}
}
