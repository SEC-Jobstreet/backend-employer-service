package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/SEC-Jobstreet/backend-employer-service/models"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

type Auth struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
}

type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
}

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func NewAuth(config utils.Config) *Auth {
	a := &Auth{
		cognitoRegion:     config.CognitoRegion,
		cognitoUserPoolID: config.CognitoUserPoolID,
	}

	a.jwkURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.cognitoRegion, a.cognitoUserPoolID)
	err := a.CacheJWK()
	if err != nil {
		logrus.Fatal(err)
	}

	return a
}

func (m *Auth) CacheJWK() error {
	req, err := http.NewRequest("GET", m.jwkURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jwk := new(JWK)
	err = json.Unmarshal(body, jwk)
	if err != nil {
		return err
	}

	m.jwk = jwk
	return nil
}

func (m *Auth) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key := convertKey(m.jwk.Keys[1].E, m.jwk.Keys[1].N)
		return key, nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func (m *Auth) JWK() *JWK {
	return m.jwk
}

func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}

func AuthMiddleware(config utils.Config, accessibleRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		accessToken := fields[1]

		// Validate token cognito
		auth := NewAuth(config)
		err := auth.CacheJWK()
		if err != nil {
			err := fmt.Errorf("AuthMiddleware - Error cacheJWK, error = %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		token, err := auth.ParseJWT(accessToken)
		if err != nil {
			err := fmt.Errorf("AuthMiddleware - Error ParseJWT, error = %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		if !token.Valid {
			err := fmt.Errorf("AuthMiddleware - Invalid token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		// Parse token to struct user and store into context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid token")
			return
		}
		var currentUser models.AuthClaim
		errDecode := mapstructure.Decode(claims, &currentUser)
		if errDecode != nil {
			err := fmt.Errorf("AuthMiddleware - Cannot decode claims %v to current user, error %v", utils.LogFull(claims), errDecode)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		if len(accessibleRoles) != 0 && !HasPermission(currentUser.CognitoGroups, accessibleRoles) {
			err := fmt.Errorf("permission denied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		ctx.Set(utils.AuthorizationPayloadKey, currentUser)
		ctx.Next()
	}
}

func HasPermission(userRoles []string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		for _, userRole := range userRoles {
			if userRole == role {
				return true
			}
		}
	}
	return false
}
