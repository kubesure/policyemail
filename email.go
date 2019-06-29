package main

import (
	"context"
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"log"
	"os"
	"strings"

	wpdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	e "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type polmetadata struct {
	Email email      `json:"email"`
	Data  policydata `json:"data"`
}

type policydata struct {
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	City         string `json:"city"`
	PinCode      int    `json:"pinCode"`
	MobileNumber int    `json:"mobileNumber"`
	PolicyNumber int    `json:"policyNumber"`
}

type email struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func init() {
	os.Setenv("WKHTMLTOPDF_PATH", os.Getenv("LAMBDA_TASK_ROOT"))
}

func handler(ctx context.Context, event e.S3Event) (string, error) {
	for _, record := range event.Records {
		log.Println("bucket" + record.S3.Bucket.Name)
		log.Println("object " + record.S3.Object.Key)
		processEvent(record)
	}
	return fmt.Sprintf("object processed "), nil
}

func processEvent(record e.S3EventRecord) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(record.S3.Bucket.Name),
		Key:    aws.String(record.S3.Object.Key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		log.Println("error getting obejct " + err.Error())
		return err
	}
	defer result.Body.Close()
	bodyBytes, err := io.ReadAll(result.Body)
	metaData := string(bodyBytes)
	log.Println(metaData)
	pm, err := marshallReq(metaData)
	if err != nil {
		return err
	}

	generatePDF(pm)
	if err != nil {
		return err
	}

	return nil
}

func marshallReq(data string) (*polmetadata, error) {
	var pd polmetadata
	err := json.Unmarshal([]byte(data), &pd)
	if err != nil {

		return nil, err
	}
	return &pd, nil
}

func generatePDF(metadata *polmetadata) ([]byte, error) {
	html := "<!doctype html><html><head><title>WKHTMLTOPDF TEST</title></head><body>HELLO PDF</body></html>"
	pdfg, err := wpdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	pdfg.AddPage(wpdf.NewPageReader(strings.NewReader(html)))
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	return pdfg.Bytes(), nil
}

func main() {
	//l.Start(handler)
}
