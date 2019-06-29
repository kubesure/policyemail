package main

import (
	"testing"

	e "github.com/aws/aws-lambda-go/events"
)

func TestS3GetObject(t *testing.T) {
	record := e.S3EventRecord{}
	record.S3.Bucket = e.S3Bucket{Name: "kubesure-policyissued"}
	record.S3.Object = e.S3Object{Key: "unprocessed_1234567890.json"}
	err := processEvent(record)
	if err != nil {
		t.Errorf("No file returned %v", err)
	}
}
