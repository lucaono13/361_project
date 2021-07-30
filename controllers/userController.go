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

// type ErrorResponse struct {
// 	Err string
// }

// type error interface {
// 	Error() string
// }

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))


// Function created to encrypt the password to store in the DynamoDB
func Hash( password string ) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// Function to verify the user inputted password against the one in the database
func VerifyPass( userPass string, inputtedPass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPass))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Login/Password is incorrect")
		check = false
		return
	}



	return check
}

func FindUsername ( username string ) (*user, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("---"),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
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

func AddUser ( ur *User) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(),
		Item: map[string]*dynamodb.AttributeValue{
			"Email":{ 
				S: aws.String(ur.Email) 
			},
			"Password":{ 
				S: aws.String(ur.Password) 
			},
			"Username":{ 
				S: aws.String(ur.Username) 
			},
			"TokenHash":{ 
				S: aws.String(ur.TokenHash) 
			},
			"RefToken":{ 
				S: aws.String(ur.RefToken) 
			},
			"CreatedAt":{ 
				S: aws.String(ur.CreatedAt) 
			},
			"UpdatedAt":{ 
				S: aws.String(ur.UpdatedAt) 
			},
			
		},
		ConditionExpression: "Username <> :username AND Email <> :email",
		ExpressionAttributeValues: {
			":username" : { S: ur.Username},
			":email" : { S: ur.Email}
		},
	}

	_, err := db.PutItem(input)
	return err
}


// Function to register a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	user := &structure.User{}
	json.NewDecoder(r.Body).Decode(user)

	pass := Hash(user.Password)

	user.Password = pass
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	token, refToken, _ := handlers.GenerateAllTokens(*user.Username, *user.Email)
	user.Token = &token
	user.RefToken = &refToken

	err := AddUser(user)

	if err != nil {
		fmt.Println(err)
		return
	}

	newUser = append(newUser, user.Username)
	newUser = append(newUser, user.Email)
	newUser = append(newUser, use.CreatedAt)
	json.NewEcoder(w).Encode(newUser)

}

func UpdateTokens ( signedToken string, signedRefToken string, username string) {
	// var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

	// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":TokenHash": {
				S: aws.String(signedToken),
			},
			"RefToken": {
				S: aws.String(signedRefToken)
			},
		},
		TableName: aws.String(""),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
		ReturnValues: 		aws.String("UPDATED_NEW"),
		UpdateExpression: 	aws.String("set TokenHash = :TokenHash"),
	}

	_, err := db.UpdateItem(input)
	if err != nil {
		log.Fatalf("Error updating Item: ", err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := structure.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, err2 := FindUsername(user.Username)
	if err2 != nil {
		var resp = map[string]interface{}{"status":false, "message":"Couldn't find username"}
		return resp
	}

	resp := VerifyLogin(user.Username, user.Password)
	json.NewEncoder(w).Encode(resp)
}

func VerifyLogin( username string, password string) map[string]interface{} {
	user := &structure.User{}
	ur, err := FindUsername(username)
	
	if err != nil {
		var resp = map[string]interface{}{"status":false, "message":"Username not found."}
		return resp
	}
	expires := time.Now().Add(time.Minute * 100000).Unix()

	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
	if check != nil && check == bcrypt.ErrMismatchedHashAndPassword{
		var resp = map[string]interface{}{"status":false, "message" :"Invalid login. Please try again"}
		return resp
	}

	tok := &handlers.SignedDetails {
		Username: 	ur.Username,
		Email:		ur.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expires
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tok)

	tString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status":false, "message":"Successful Login!"}
	resp["token"] = tString
	resp["user"] = ur.Username

	return resp
}