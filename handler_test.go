package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
}

func TestHandleNotification_RequestBodyValidation_CorrectStatusCodeReturned(t *testing.T) {
	var testCases = []struct {
		req        string
		respStatus int
	}{
		{
			`{"scenario" : invalidJson"}`,
			500,
		},
		{
			`{"scenario" : "no uuid"}`,
			400,
		},
		{
			`{"scenario" : "tags list empty", "uuid" : "1a78d8e7-473d-4e9f-ae2e-7f20a45e31fc", "tags" : []}`,
			200,
		},
	}

	mm := metadataMapper{}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "test-url", bytes.NewReader([]byte(tc.req)))
		if err != nil {
			t.Fatalf("[%v]", err)
		}
		mm.handleNotification(w, req)
		if w.Code != tc.respStatus {
			t.Errorf("Expected status code: [%d]. Actual: [%d]", tc.respStatus, w.Code)
		}
	}
}

func TestCreateMetadataPublishEvent_CreatedEventMatchesExpectedEvent(t *testing.T) {
	video := video{
		UUID: "1234",
		Tags: []string{"Emerging-Markets", "Commodities"},
	}
	expected := nativeCmsMetadataPublicationEvent{
		UUID:  "1234",
		Value: "PGNvbnRlbnRSZWY+PHRhZ3M+PHRhZz48dGVybSB0YXhvbm9teT0iU2VjdGlvbnMiIGlkPSJNVEEyLVUyVmpkR2x2Ym5NPSI+PC90ZXJtPjxzY29yZSBjb25maWRlbmNlPSI5MCIgcmVsZXZhbmNlPSI5MCI+PC9zY29yZT48L3RhZz48dGFnPjx0ZXJtIHRheG9ub215PSJTZWN0aW9ucyIgaWQ9Ik1UQTEtVTJWamRHbHZibk09Ij48L3Rlcm0+PHNjb3JlIGNvbmZpZGVuY2U9IjkwIiByZWxldmFuY2U9IjkwIj48L3Njb3JlPjwvdGFnPjwvdGFncz48cHJpbWFyeVNlY3Rpb24gdGF4b25vbXk9IiIgaWQ9IiI+PC9wcmltYXJ5U2VjdGlvbj48L2NvbnRlbnRSZWY+",
	}
	mm := metadataMapper{
		mappings: map[string]term{
			"Emerging-Markets": term{
				ID:       "MTA2-U2VjdGlvbnM=",
				Taxonomy: "Sections",
			},
			"Commodities": term{
				ID:       "MTA1-U2VjdGlvbnM=",
				Taxonomy: "Sections",
			},
		},
	}

	actual, err := mm.createMetadataPublishEventMsg(video, "unit-test")
	if err != nil {
		t.Errorf("Expected no error. Found: [%v]", err)
	}
	if expected != *actual {
		t.Errorf("Expected: [%v]. Actual: [%v]", expected, *actual)
	}
}

func TestSendMetadata_RequestHeadersAreSet(t *testing.T) {

}
