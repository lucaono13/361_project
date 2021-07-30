package handlers

import (
	"time"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lucaono13/361_project/structure"
)

// type Token struct {

// }

func getSecret() string {
	secretName := "UserAuthSecret"
	region := "us-west-2"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
                                  aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

				case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

				case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

				case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

				case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return ("error")
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return secretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret
	}
	
}



func GenerateAllTokens(user string, email string) (signedToken string, signedRefToken string, err error) {
	claims := &structure.SignedDetails{
		Email:    email,
		Username: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &structure.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		}
	}

	secretKey = getSecret()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func ValidateToken(signedToken string) (claims *structure.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&structure.SignedDetails{},
		func(token *jwt.Token) (interface{}, err) {
			return []byte(getSecret()), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*structure.SignedDetails)
	if !ok {
		msg = fmt.Sprintf("invalid token")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("expired token")
		msg = err.Error()
	}

	return claims, msg
}

func UpdateTokens ( signedToken string, signedRefToken string, username string) {
	var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

	// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{}
	}
}