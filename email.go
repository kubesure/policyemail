package main

import (
	"context"
	"fmt"
	io "io/ioutil"
	"log"

	e "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
	return nil
}

func main() {

	//l.Start(handler)
}
