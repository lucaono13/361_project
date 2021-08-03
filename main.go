package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/lucaono13/361_project/controllers"
	"github.com/lucaono13/361_project/structure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func signIn(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]
	usr.Password = req.QueryStringParameters["pass"]

	err := controllers.SignIn(*usr)

	return err, nil
}

func register(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]
	usr.Password = req.QueryStringParameters["pass"]

	err := controllers.CreateUser(*usr)
	return err, nil
}

func updateBio(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// usr := new(structure.User)
	b := new(structure.BioUpdate)
	b.Email = req.QueryStringParameters["email"]
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Wrong content type. Ensure that content type is application/json",
		}, nil
	}

	b.Bio = req.Body
	if b.Bio == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad request",
		}, nil
	}

	return controllers.UpdateBio(b), nil
}

func getInfo(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]

	user, err := controllers.FindUser(usr.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad request",
		}, nil
	}
	nbu := new(structure.BioUpdate)
	nbu.Bio = user.Bio
	nbu.Email = user.Email
	stringBody, _ := json.Marshal(nbu)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(stringBody),
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
	log.Println("Hello why")
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	qtype := req.QueryStringParameters["type"]
	switch req.HTTPMethod {
	case "GET":
		switch qtype {
		case "signin":
			return signIn(req)
		case "getInfo":
			return getInfo(req)
		default:
			return clientError(http.StatusMethodNotAllowed)
		}
	case "POST":
		return register(req)
	case "PUT":
		return updateBio(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}
