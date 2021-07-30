package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	user := req.QueryStringParameters["username"]
	if !isbnRegexp.MatchString(user) {
		return clientError(http.StatusBadRequest)
	}

	// ur, err :=
}

func main() {
	lambda.Start(show)
}
