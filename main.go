package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/lucaono13/361_project/controllers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

// var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)


func signIn(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	
	err := Login(req.Headers["Email"], req.Headers["Password"])
	// status, err := SignUp(req.Headers["Email"], req.Headers["Pass"])

	if err == "Invalid Login" {
		return clientError(400)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 	200,
		Body:			"Log in successful"
	}
}

func register(req events.APIGatewayProxyRequest) (event.APIGatewayProxyResponse, error) {
	err := SignUp(req.Headers["Email"], req.Headers["Pass"])

	if err != nil {
		return serverError(err)
	}
	
	return events.APIGatewayProxyResponse {
		StatusCode: 	200
	}

}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode:		http:StatusInternalServerError,
		Body:			http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error){
	return events.APIGatewayProxyResponse{
		StatusCode:		status,
		Body:			http.StatusText(status),
	}, nil
}



func main() {
	lambda.Start(router)
}

func router( req events.APIGatewayProxyRequest ) ( events.APIGatewayProxyResponse, error ) {
	switch req.HTTPMethod {
	case "GET":
		return signIn(req)
	case "POST":
		return register(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}