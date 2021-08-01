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

func signIn(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]
	usr.Password = req.QueryStringParameters["pass"]

	err := controllers.SignIn(*usr)

	return err
}

func register(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]
	usr.Password = req.QueryStringParameters["pass"]

	err := controllers.CreateUser(*usr)
	return err
}

func updateBio(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// usr := new(structure.User)
	b := new(structure.BioUpdate)
	b.Email = req.QueryStringParameters["email"]
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Wrong content type. Ensure that content type is application/json",
		}
	}

	err := json.Unmarshal([]byte(req.Body), b)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error reading body",
		}
	}
	if b.Bio == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad request",
		}
	}

	return controllers.UpdateBio(b)
}

func getInfo(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	usr := new(structure.User)
	usr.Email = req.QueryStringParameters["email"]

	user, err := controllers.FindUser(usr.Email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Bad request",
		}
	}
	nbu := new(structure.BioUpdate)
	nbu.Bio = user.Bio
	nbu.Email = user.Email
	stringBody, _ := json.Marshal(nbu)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(stringBody),
	}
}

func serverError(err error) events.APIGatewayProxyResponse {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}
}

func clientError(status int) events.APIGatewayProxyResponse {
	// log.Println("Hello why")
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}
}

func main() {
	lambda.Start(router)
}

func router(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// fmt.Println("Hi there")
	// log.Println(req.HTTPMethod)
	htype := req.Headers["type"]
	switch htype {
	case "signin":
		return signIn(req)
	case "create":
		return register(req)
	case "updateBio":
		return updateBio(req)
	case "getInfo":
		return getInfo(req)
	default:
		// fmt.Println("WTF")
		return clientError(http.StatusMethodNotAllowed)
	}
}
