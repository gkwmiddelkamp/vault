package interceptors

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
	"vault/server/services"
	"vault/vault"
)

type authInterceptor struct {
}

const Header = "X-API-KEY"
const UnauthorizedText = "Unauthorized"

var AuthInterceptor vault.Interceptor = authInterceptor{}

func (i authInterceptor) Before(w http.ResponseWriter, r *http.Request, cfg *vault.InterceptorConfig) vault.Result {
	if cfg.GetMinimalAuth() > vault.Public {
		header := r.Header.Get(Header)
		if header == "" {
			http.Error(w, UnauthorizedText, http.StatusUnauthorized)
			return vault.Done()
		}

		token, err := i.lookup(header, cfg)
		if err != nil {
			log.Println("Invalid token: " + err.Error())
			http.Error(w, UnauthorizedText, http.StatusUnauthorized)
			return vault.Done()
		}
		if token.TokenType < cfg.GetMinimalAuth() {
			log.Println("Unauthorized")
			http.Error(w, UnauthorizedText, http.StatusUnauthorized)
			return vault.Done()
		}
	}

	return vault.NotDone()
}

func (i authInterceptor) lookup(tokenInput string, cfg *vault.InterceptorConfig) (*vault.Token, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(tokenInput)
	if err != nil {
		return nil, err
	}
	if len(decodedToken) != 138 {
		return nil, errors.New("invalid token")
	}

	inputKey, inputObjectId, inputTokenToCheck, inputEnvironmentSecret :=
		decodedToken[:len(decodedToken)-(vault.ObjectIdLength+vault.TokenLength+vault.KeyLength)],
		decodedToken[len(decodedToken)-(vault.ObjectIdLength+vault.TokenLength+vault.KeyLength):len(decodedToken)-(vault.TokenLength+vault.KeyLength)],
		decodedToken[len(decodedToken)-(vault.TokenLength+vault.KeyLength):len(decodedToken)-vault.KeyLength],
		decodedToken[len(decodedToken)-(vault.KeyLength):]

	decodeCipher, err := aes.NewCipher(inputKey)
	if err != nil {
		return nil, err
	}
	decodeGcm, err := cipher.NewGCM(decodeCipher)
	if err != nil {
		return nil, err
	}

	tokenId, err := primitive.ObjectIDFromHex(string(inputObjectId))
	if err != nil {
		return nil, err
	}
	var dbToken vault.Token
	err = cfg.GetDb().Connection.Collection(vault.TokenCollection).FindOne(context.Background(), bson.M{"_id": tokenId}).Decode(&dbToken)
	if err != nil {
		return nil, err
	}

	if dbToken.ExpiresAt > 0 && dbToken.ExpiresAt.Time().Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	decodedDbToken, err := base64.StdEncoding.DecodeString(dbToken.Token)
	if err != nil {
		return nil, err
	}

	decodeNonceSize := decodeGcm.NonceSize()
	if len(inputTokenToCheck) < decodeNonceSize {
		return nil, errors.New("invalid nonce in token")
	}
	decodeNonce, decodeCipherText := decodedDbToken[:decodeNonceSize], decodedDbToken[decodeNonceSize:]
	plainText, err := decodeGcm.Open(nil, decodeNonce, decodeCipherText, nil)
	if err != nil {
		return nil, err
	}
	decryptedDbToken := plainText[len(plainText)-vault.TokenLength:]

	environmentService := services.NewEnvironmentService(cfg.GetDb())

	environment, err := environmentService.GetEnvironment(dbToken.GetEnvironmentId().Hex())
	if err != nil || !environment.Active {
		return nil, errors.New("invalid environment: " + err.Error())
	}

	if string(decryptedDbToken) == string(inputTokenToCheck) {
		handlerConfig := cfg.GetHandlerConfig()
		handlerConfig.SetEnvironmentId(dbToken.GetEnvironmentId())
		handlerConfig.SetTokenType(dbToken.TokenType)
		handlerConfig.SetEnvironmentSecret(inputEnvironmentSecret)
		handlerConfig.SetTokenDescription(dbToken.Description)
		return &dbToken, nil
	}
	return nil, errors.New("token is invalid")
}
