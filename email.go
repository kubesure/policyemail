package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	io "io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	wpdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	e "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event e.S3Event) (string, error) {

	for _, record := range event.Records {
		log.Println("bucket" + record.S3.Bucket.Name)
		log.Println("object " + record.S3.Object.Key)
		err := processEvent(record)
		if err != nil {
			log.Println(err)
			return "processing error ", err
		}
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
		log.Println("error getting object " + err.Error())
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

	pdfbytes, perr := generatePDF(pm)
	if perr != nil {
		log.Println("error while pdf generation")
		log.Println(perr)
		return err
	}

	//newKey := strings.Replace(record.S3.Object.Key, ".json", ".pdf", -1)
	newKey := strconv.Itoa(pm.Data.PolicyNumber) + ".pdf"
	log.Println(newKey)

	pinput := s3.PutObjectInput{
		Bucket: aws.String(record.S3.Bucket.Name),
		Key:    aws.String(newKey),
		Body:   bytes.NewReader(pdfbytes),
	}

	presult, putErr := svc.PutObject(&pinput)

	if putErr != nil {
		return putErr
	}
	log.Println(presult)

	return nil
}

func generatePDF(metadata *polmetadata) ([]byte, error) {

	html, herr := generateHTML(metadata)
	if herr != nil {
		return nil, herr
	}
	log.Printf(html)

	log.Println("os.getenv(WKHTMLTOPDF_PATH)--" + os.Getenv("WKHTMLTOPDF_PATH"))
	pdfg, err := wpdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	pdfg.AddPage(wpdf.NewPageReader(strings.NewReader(html)))
	pdfg.Dpi.Set(350)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	err = pdfg.Create()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Write buffer contents to file on disk
	// err = pdfg.WriteFile("./simplesample.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return pdfg.Bytes(), nil
}

func generateHTML(metaData *polmetadata) (string, error) {
	fmap := template.FuncMap{
		"currentdate": currentdate,
	}

	t, errp := template.New("esyhealth-pdf.html").Funcs(fmap).ParseFiles("esyhealth-pdf.html")
	if errp != nil {
		return "", errp
	}
	buff := new(bytes.Buffer)
	err := t.Execute(buff, metaData)
	if err != nil {
		return "", err
	}
	log.Println(buff.String())
	return buff.String(), nil
}

func marshallReq(data string) (*polmetadata, error) {
	var pd polmetadata
	err := json.Unmarshal([]byte(data), &pd)
	if err != nil {

		return nil, err
	}
	return &pd, nil
}

func currentdate() (string, error) {
	const layoutISO = "2006-01-02"
	const custom = "Mon Jan _2 15:04:05 2006"
	currentDate := time.Now().Format(custom)
	return currentDate, nil
}
