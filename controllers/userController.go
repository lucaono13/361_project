package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lucaono13/361_project/structure"
	"github.com/lucaono13/361_project/handlers"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)

var validate = validator.New()

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))


// Function created to encrypt the password to store in the DynamoDB
func Hash( password string ) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func FindUser( email string ) (*user, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("---"),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nul {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	ur := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil

}

func AddUser( ur *User) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(),
		Item: map[string]*dynamodb.AttributeValue{
			"Email":{ 
				S: aws.String(ur.Email) 
			},
			"Password":{ 
				S: aws.String(ur.Password) 
			},
			"CreatedAt":{ 
				S: aws.String(ur.CreatedAt) 
			},
			"UpdatedAt":{ 
				S: aws.String(ur.UpdatedAt) 
			},
			
		},
		ConditionExpression: "Email <> :email",
		ExpressionAttributeValues: {
			":email" : { S: ur.Email}
		},
	}

	_, err := db.PutItem(input)
	return err
}


// Function to register a new user
func SignUp(email string, pass string) (error) {
	user := &structure.User{}

	pass := Hash(user.Password)

	user.Email = email
	user.Password = pass
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	err := AddUser(user)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err


}



// func Login(w http.ResponseWriter, r *http.Request) {
func Login(email string, password string) error {
	user := structure.User{}
	user.Email = email
	user.Password = password
	resp := VerifyLogin(user)
	return resp
}

func VerifyLogin( user structure.User) map[string]interface{} {
	ur, err := FindUser(user.Email)

	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
	if check != nil && check == bcrypt.ErrMismatchedHashAndPassword{
		return "Invalid Login"
	}

	return "Login Successful"
}