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

func (i authInterceptor) lookup(tokenInput string, cfg *vault.InterceptorConfig) (vault.Token, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(tokenInput)
	if err != nil {
		return vault.Token{}, err
	}
	if len(decodedToken) != 106 {
		return vault.Token{}, errors.New("invalid token")
	}
	inputKey, inputObjectId, inputTokenToCheck :=
		decodedToken[:len(decodedToken)-(24+vault.TokenLength)],
		decodedToken[len(decodedToken)-(24+vault.TokenLength):len(decodedToken)-vault.TokenLength],
		decodedToken[len(decodedToken)-vault.TokenLength:]

	decodeCipher, err := aes.NewCipher(inputKey)
	if err != nil {
		return vault.Token{}, err
	}
	decodeGcm, err := cipher.NewGCM(decodeCipher)
	if err != nil {
		return vault.Token{}, err
	}

	tokenId, err := primitive.ObjectIDFromHex(string(inputObjectId))
	if err != nil {
		return vault.Token{}, err
	}
	var dbToken vault.Token
	err = cfg.GetDb().Connection.Collection(vault.TokenCollection).FindOne(context.Background(), bson.M{"_id": tokenId}).Decode(&dbToken)
	if err != nil {
		return vault.Token{}, err
	}

	if dbToken.ExpiresAt > 0 && dbToken.ExpiresAt.Time().Before(time.Now()) {
		return vault.Token{}, errors.New("token expired")
	}

	decodedDbToken, err := base64.StdEncoding.DecodeString(dbToken.Token)
	if err != nil {
		return vault.Token{}, err
	}

	decodeNonceSize := decodeGcm.NonceSize()
	if len(inputTokenToCheck) < decodeNonceSize {
		return vault.Token{}, errors.New("invalid nonce in token")
	}
	decodeNonce, decodeCipherText := decodedDbToken[:decodeNonceSize], decodedDbToken[decodeNonceSize:]
	plainText, err := decodeGcm.Open(nil, decodeNonce, decodeCipherText, nil)
	if err != nil {
		return vault.Token{}, err
	}
	decryptedDbToken := plainText[len(plainText)-vault.TokenLength:]

	if string(decryptedDbToken) == string(inputTokenToCheck) {
		cfg.GetHandlerConfig().SetEnvironmentId(dbToken.GetEnvironmentId())
		cfg.GetHandlerConfig().SetTokenType(dbToken.TokenType)
		return dbToken, nil
	}
	return vault.Token{}, errors.New("token is invalid")
}
