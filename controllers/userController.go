package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/lucaono13/361_project/structure"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// var validate = validator.New()
// var mySession = session.Must(session.NewSession())
var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("us-west-2"),
}))
var db = dynamodb.New(sess)

// var db = dynamodb.New(session.NewSession(), aws.NewConfig().WithRegion("us-west-2"))

// Function created to encrypt the password to store in the DynamoDB
func Hash(password string) string {
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

	ur := new(structure.User)
	err = dynamodbattribute.UnmarshalMap(result.Item, ur)
	if err != nil {
		return nil, err
	}

	return ur, nil

}

// func AddUser(ur *structure.User) error {
func AddUser(email string, password string) error {
	ur, ferr := FindUser(email)
	if ferr != nil {
		return ferr
	}
	if ur == nil {
		return errors.New("email already in use")
	}
	fmt.Println(email, password)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("profiles361_"),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"password": {
				S: aws.String(password),
			},
		},
	}
	_, err := db.PutItem(input)
	return err
}

// Function to register a new user
// func SignUp(email string, pass string) error {
func SignUp(b *structure.User) error {
	// user := &structure.User{}

	password := Hash(b.Password)

	// user.Email = email
	// user.Password = pass
	fmt.Println(b.Email, password, b.Password)
	err := AddUser(b.Email, password)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err

}

// func Login(w http.ResponseWriter, r *http.Request) {
func Login(email string, password string) string {
	user := structure.User{}
	user.Email = email
	user.Password = password
	resp := VerifyLogin(user)
	return resp
}

func VerifyLogin(user structure.User) string {
	ur, err := FindUser(user.Email)

	if err != nil {
		return "Invalid Login"
	}

	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
	if check != nil && check == bcrypt.ErrMismatchedHashAndPassword {
		return "Invalid Login"
	}

	return "Login Successful"
}
