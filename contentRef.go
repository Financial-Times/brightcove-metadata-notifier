package main

// contentRef models the data as it comes from the metadata publishing event
type contentRef struct {
	TagHolder      tags `xml:"tags"`
	PrimarySection term `xml:"primarySection,omitempty"`
}

type tags struct {
	Tags []tag `xml:"tag"`
}

type tag struct {
	Term     term     `xml:"term"`
	TagScore tagScore `xml:"score"`
}

type term struct {
	CanonicalName string `xml:"canonicalName,omitempty"`
	Taxonomy      string `xml:"taxonomy,attr"`
	ID            string `xml:"id,attr"`
}

type tagScore struct {
	Confidence int `xml:"confidence,attr"`
	Relevance  int `xml:"relevance,attr"`
}
