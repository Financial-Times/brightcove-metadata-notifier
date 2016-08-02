package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Financial-Times/transactionid-utils-go"
)

type video struct {
	UUID string
	Tags []string
}

type nativeCmsMetadataPublicationEvent struct {
	Value string `json:"value"`
	UUID  string `json:"uuid"`
}

var defaultTagScore = tagScore{Confidence: 90, Relevance: 90}

func (mm metadataMapper) handleNotification(w http.ResponseWriter, r *http.Request) {
	tid := transactionidutils.GetTransactionIDFromRequest(r)
	infoLogger.Printf("Received video. tid=[%s]", tid);
	var v video
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		handleServerErr(w, fmt.Sprintf("tid=[%s]. Cannot decode video metadata: [%v]", tid, err))
		return
	}
	if v.UUID == "" {
		handleClientErr(w, fmt.Sprintf("tid=[%s]. Missing uuid: [%#v]", tid, v))
		return
	}
	if len(v.Tags) == 0 {
		infoLogger.Printf("tid=[%s]. Video with uuid [%s] has no tags. No metadata will be generated.", tid, v.UUID)
		return
	}

	ev, err := mm.createMetadataPublishEventMsg(v, tid)
	if err != nil {
		handleServerErr(w, fmt.Sprintf("tid=[%s]. %v", tid, err))
		return
	}
	m, err := json.Marshal(*ev)
	if err != nil {
		handleServerErr(w, fmt.Sprintf("tid=[%s]. JSON Marshalling: [%v]", tid, err))
		return
	}
	if err = mm.sendMetadata(m, tid); err != nil {
		handleServerErr(w, fmt.Sprintf("tid=[%s]. %v", tid, err))
		return
	}
	infoLogger.Printf("Sent metadata event for video=[%s] tid=[%s]", v.UUID, tid);
}

func (mm metadataMapper) createMetadataPublishEventMsg(v video, tid string) (*nativeCmsMetadataPublicationEvent, error) {
	marshalled, err := xml.Marshal(buildContentRef(v.UUID, mm.getAnnotations(v.Tags, tid), tid))
	if err != nil {
		return nil, fmt.Errorf("tid=[%s]. XML Marshalling: [%v]", tid, err)
	}
	return &nativeCmsMetadataPublicationEvent{
		Value: base64.StdEncoding.EncodeToString(marshalled),
		UUID:  v.UUID,
	}, nil
}

func buildContentRef(uuid string, terms []term, tid string) contentRef {
	var tagz []tag

	for _, term := range terms {
		tagz = append(tagz, tag{term, defaultTagScore})
	}

	return contentRef{
		TagHolder: tags{tagz},
	}
}

func (mm metadataMapper) getAnnotations(tags []string, tid string) []term {
	var annotations []term
	for _, tag := range tags {
		t, present := mm.mappings[tag]
		if !present {
			infoLogger.Printf("tid=[%s]. Brightcove tag [%s] has no TME mapping.", tid, tag)
			continue
		}
		annotations = append(annotations, t)
	}
	return annotations
}

func (mm metadataMapper) sendMetadata(metadata []byte, tid string) error {
	req, err := http.NewRequest("POST", mm.config.cmsMetadataNotifierAddr+"/notify", bytes.NewReader(metadata))
	if err != nil {
		return fmt.Errorf("Creating request: [%v]", err)
	}
	req.Header.Add("X-Origin-System-Id", "brightcove")
	req.Header.Add("X-Request-Id", tid)
	req.Header.Add("Content-type", "application/json")
	req.Host = mm.config.cmsMetadataNotifierHost
	if mm.config.cmsMetadataNotifierAuth != "" {
		req.Header.Add("Authorization", mm.config.cmsMetadataNotifierAuth)
	}
	resp, err := mm.client.Do(req)
	if err != nil {
		return fmt.Errorf("Sending metadata to notifier: [%v]", err)
	}
	defer cleanupResp(resp)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Sending metadata to notifier: unexpected status code: [%d]", resp.StatusCode)
	}
	return nil
}

func handleServerErr(w http.ResponseWriter, errMsg string) {
	warnLogger.Printf(errMsg)
	w.WriteHeader(http.StatusInternalServerError)
}

func handleClientErr(w http.ResponseWriter, errMsg string) {
	warnLogger.Printf(errMsg)
	w.WriteHeader(http.StatusBadRequest)
}

func cleanupResp(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		warnLogger.Printf("[%v]", err)
	}
	err = resp.Body.Close()
	if err != nil {
		warnLogger.Printf("[%v]", err)
	}
}
