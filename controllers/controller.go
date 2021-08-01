package controllers

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/lucaono13/361_project/structure"
	"golang.org/x/crypto/bcrypt"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("us-west-2"),
}))
var db = dynamodb.New(sess)

func HashPass(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func FindUser(email string) (*structure.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("profiles361_"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	usr := new(structure.User)
	err = dynamodbattribute.UnmarshalMap(result.Item, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func CreateUser(user structure.User) events.APIGatewayProxyResponse {
	usr, ferr := FindUser(user.Email)
	if ferr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error looking for user",
		}
	}
	if usr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 409,
			Body:       "User already exists",
		}
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("profiles361_"),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
			"password": {
				S: aws.String(HashPass(user.Password)),
			},
			"bio": {
				S: aws.String(""),
			},
		},
	}

	_, err := db.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error creating user",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       "User Created",
	}

}

func SignIn(user structure.User) events.APIGatewayProxyResponse {

	usr, err := FindUser(user.Email)

	if err != nil {
		errorLogger.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error searching for users",
		}
	}

	check := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
	if check != nil && check == bcrypt.ErrMismatchedHashAndPassword {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid Credentials",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Sign in successful",
	}

}

func UpdateBio(newBio structure.BioUpdate) events.APIGatewayProxyResponse {
	usr, ferr := FindUser(newBio.Email)
	if ferr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error looking for user",
		}
	}
	if usr == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "User doesn't exist",
		}
	}

	// nBio := new(BioUpdate)
	// nBio.Bio = newBio
	// nBio.Email = usr.Email

	expr, err := dynamodbattribute.MarshalMap(newBio)
	if err != nil {
		log.Fatalf("Got error marshalling info: %s", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("profiles361_"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(usr.Email),
			},
		},
		UpdateExpression:          aws.String("set bio = :bio"),
		ConditionExpression:       aws.String("email = :email"),
		ExpressionAttributeValues: expr,
		ReturnValues:              aws.String("ALL_NEW"),
	}

	_, uerr := db.UpdateItem(input)
	if uerr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error updating User bio",
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Successfully updated bio",
	}

}

// func userExists() events.APIGatewayProxyResponse {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 409,
// 		Body:       "User already exists",
// 	}
// }
