package main

import (
	"fmt"
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

func TestGeneratePDF(t *testing.T) {
	metadata := `{"email":{"from":"edakghar@gmail.com","to":"pras.p.in@gmail.com"},
	"data":{"name":"Usha Patel","addressLine1":"ketaki","addressLine2":"maneklal","addressLine3":"Ghatkopar",
	"city":"mumbai","pinCode":400086,"mobileNumber":9821284567,"policyNumber":1234567890},
	"status":{"mailSent":true,"pdfCreated":true}}`

	pm, err := marshallReq(metadata)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}

	_, errgen := generatePDF(pm)
	if errgen != nil {
		t.Errorf("error generating pdf %v", err)
	}
}

func TestMarshallPolData(t *testing.T) {
	metadata := `{"email":{"from":"edakghar@gmail.com","to":"pras.p.in@gmail.com"},"data":{"name":"Usha Patel","addressLine1":"ketaki","addressLine2":"maneklal","addressLine3":"Ghatkopar","city":"mumbai","pinCode":400086,"mobileNumber":9821284567,"policyNumber":1234567890},
	"status":{"mailSent":true,"pdfCreated":true}}`

	pm, err := marshallReq(metadata)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}
	fmt.Println(pm)
}
