package main

import (
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-fthealth"
)

func (hc healthcheck) health() func(w http.ResponseWriter, r *http.Request) {
	return fthealth.HandlerParallel("Dependent services healthcheck", "Checks if all the dependent services are reachable and healthy.", hc.cmsMetadataNotifierReachable(), hc.berthaSpreadsheetAvailable())
}

func (hc healthcheck) gtg(w http.ResponseWriter, r *http.Request) {
	healthChecks := []func() error{hc.checkCmsMetadataNotifierHealth, hc.checkBerthaSpreadsheetHealth}

	for _, hCheck := range healthChecks {
		if err := hCheck(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
}

func (hc healthcheck) cmsMetadataNotifierReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Video metadata will not reach UPP stack. Videos will not be annotated.",
		Name:             "CMS Metadata Notifier Reachable",
		PanicGuide:       "<coco runbook>",
		Severity:         1,
		TechnicalSummary: "CMS Metadata Notifier is not reachable/healthy",
		Checker:          hc.checkCmsMetadataNotifierHealth,
	}
}

func (hc healthcheck) checkCmsMetadataNotifierHealth() error {
	req, err := http.NewRequest("GET", hc.config.cmsMetadataNotifierAddr+"/__health", nil)
	if err != nil {
		return err
	}
	req.Host = hc.config.cmsMetadataNotifierHost

	resp, err := hc.client.Do(req)
	if err != nil {
		return err
	}
	defer cleanupResp(resp)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Unhealthy status code received: [%d]", resp.StatusCode)
	}
	return nil
}

func (hc healthcheck) berthaSpreadsheetAvailable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Video metadata might not be generated.",
		Name:             "Bertha Spreadsheet Available",
		PanicGuide:       "<coco runbook>",
		Severity:         1,
		TechnicalSummary: "Spreadsheet containing the mappings is not available through Bertha.",
		Checker:          hc.checkBerthaSpreadsheetHealth,
	}
}

func (hc healthcheck) checkBerthaSpreadsheetHealth() error {
	req, err := http.NewRequest("GET", hc.config.berthaAddr, nil)
	if err != nil {
		return err
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return err
	}
	defer cleanupResp(resp)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid status code received: [%d]", resp.StatusCode)
	}
	return nil
}
