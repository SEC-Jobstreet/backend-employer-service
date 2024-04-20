package middleware

import (
	"os"

	"github.com/SEC-Jobstreet/backend-employer-service/api/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type CognitoClient interface {
	SignUp(employer dto.Employer) (error, string)
	ConfirmSignUp(email string, code string) (error, string)
	Login(email string, password string) (error, string, string)
}

type awsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientID   string
}

func createSecretHash(username string, clientId string, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognitoClient(cognitoRegion string, cognitoAppClientID string) CognitoClient {
	conf := &aws.Config{Region: aws.String("ap-southeast-1")}

	sess, err := session.NewSession(conf)

	client := cognito.New(sess)

	if err != nil {
		panic(err)
	}

	return &awsCognitoClient{
		cognitoClient: client,
		appClientID:   cognitoAppClientID,
	}
}

func (c *awsCognitoClient) SignUp(employer dto.Employer) (error, string) {
	clientId := os.Getenv("COGNITO_CLIENT_ID")
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")

	secretHash := createSecretHash(employer.Email, clientId, clientSecret)

	user := &cognito.SignUpInput{
		ClientId:   aws.String(c.appClientID),
		Username:   aws.String(employer.Email),
		Password:   aws.String(employer.Password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(employer.Email),
			},
			{
				Name:  aws.String("given_name"),
				Value: aws.String(employer.FirstName),
			},
			{
				Name:  aws.String("family_name"),
				Value: aws.String(employer.LastName),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(employer.Phone),
			},
			{
				Name:  aws.String("address"),
				Value: aws.String(employer.Address),
			},
		},
	}

	result, err := c.cognitoClient.SignUp(user)

	if err != nil {
		return err, ""
	}

	return nil, result.String()
}

func (c *awsCognitoClient) ConfirmSignUp(email string, code string) (error, string) {
	clientId := os.Getenv("COGNITO_CLIENT_ID")
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")

	secretHash := createSecretHash(email, clientId, clientSecret)

	confirmSignupInput := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.appClientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
		SecretHash:       aws.String(secretHash),
	}

	result, err := c.cognitoClient.ConfirmSignUp(confirmSignupInput)

	if err != nil {
		return err, ""
	}

	return nil, result.String()
}

func (c *awsCognitoClient) Login(email string, password string) (error, string, string) {
	clientId := os.Getenv("COGNITO_CLIENT_ID")
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")

	secretHash := createSecretHash(email, clientId, clientSecret)

	signInInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(email),
			"PASSWORD":    aws.String(password),
			"SECRET_HASH": aws.String(secretHash),
		},
		ClientId: aws.String(c.appClientID),
	}

	output, err := c.cognitoClient.InitiateAuth(signInInput)
	if err != nil {
		return err, "", ""
	}

	// The JWT tokens are in output.AuthenticationResult
	accessToken := *output.AuthenticationResult.AccessToken
	refreshToken := *output.AuthenticationResult.RefreshToken

	return nil, accessToken, refreshToken
}
