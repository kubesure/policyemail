package main

import (
	"fmt"
	"testing"

	e "github.com/aws/aws-lambda-go/events"
)

const metaData = `{"email":{"from":"edakghar@gmail.com","to":"pras.p.in@gmail.com"},
"data":{"name":"Usha Patel","addressLine1":"ketaki","addressLine2":"maneklal","addressLine3":"Ghatkopar",
"city":"mumbai","pinCode":400086,"mobileNumber":9821284567,"policyNumber":1234567890},
"status":{"mailSent":true,"pdfCreated":true}}`

var pmetadata polmetadata

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

	pm, err := marshallReq(metaData)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}

	_, errgen := generatePDF(pm)
	if errgen != nil {
		t.Errorf("error generating pdf %v", err)
	}
}

func TestMarshallPolData(t *testing.T) {

	pm, err := marshallReq(metaData)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}
	fmt.Println(pm)
}

func TestGenerateHTML(t *testing.T) {
	pm, err := marshallReq(metaData)
	if err != nil {
		t.Errorf("marshall err %v", err)
	}

	html, err := generateHTML(pm)

	if err != nil {
		t.Errorf("Err while html generation %v", err)
	}

	if len(html) == 0 {
		t.Errorf("html not generated  %v", err)
	}
}

func TestCurrentDate(t *testing.T) {
	date, err := currentdate()

	if err != nil {
		t.Errorf("time formatting error")
	}

	if date != "2019-06-30" {
		t.Errorf("Incorrect date format")
	}
}
