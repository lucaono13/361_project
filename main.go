package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lucaono13/361_project/controllers"
	"github.com/lucaono13/361_project/structure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// var sess = session.Must(session.NewSession(&aws.Config{
// 	Region: aws.String("us-west-2"),
// }))
// var db = dynamodb.New(sess)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func signIn(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// var body structure.User

	// jerr := json.NewDecoder(req.Body).Decode(&body)
	// if jerr != nil {
	// 	return events.APIGatewayProxyResponse{StatusCode: 100}, nil
	// }

	// b := new(structure.User)
	// err := json.Unmarshal([]byte(req.Body), b)

	// if err != nil {
	// 	return clientError(http.StatusUnprocessableEntity)
	// }
	fmt.Println(req.Body)
	// b := req.Body
	// err := controllers.Login(b["email"], b["password"])
	err := controllers.Login(req.QueryStringParameters["email"], req.QueryStringParameters["password"])

	if err == "Invalid Login" {
		return clientError(400)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Log in successful",
	}, nil
}

func register(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	b := structure.User{}
	err := json.Unmarshal([]byte(req.Body), b)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	ferr := controllers.SignUp(&b)
	// fmt.Println(req)
	// err := controllers.SignUp(req.QueryStringParameters["email"], req.QueryStringParameters["password"])

	if ferr != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Hi there")
	fmt.Println(req)
	switch req.HTTPMethod {
	case "GET":
		return signIn(req)
	case "POST":
		return register(req)
	default:
		fmt.Println("WTF")
		return clientError(http.StatusMethodNotAllowed)
	}
}
