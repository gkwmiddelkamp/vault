package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"time"
	"vault/internal"
)

const TokenCollection = "token"
const TokenLength = 50
const KeyLength = 32
const ObjectIdLength = 24

type Token struct {
	Id             primitive.ObjectID `bson:"_id"`
	Description    string             `bson:"description,omitempty"`
	Token          string             `bson:"token,omitempty"`
	secret         string             `bson:"secret,omitempty"`
	EnvironmentId  primitive.ObjectID `bson:"environmentId,omitempty"`
	CreatedAt      primitive.DateTime `bson:"createdAt,omitempty"`
	CreatedBy      string             `bson:"createdBy,omitempty"`
	LastModifiedAt primitive.DateTime `bson:"lastModifiedAt,omitempty"`
	LastModifiedBy string             `bson:"lastModifiedBy,omitempty"`
	ExpiresAt      primitive.DateTime `bson:"expiresAt,omitempty"`
	TokenType      TokenType          `bson:"tokenType,omitempty"`
}

func (t *Token) SetExpiresAt(time primitive.DateTime) {
	t.ExpiresAt = time
}

func NewToken(description string, environmentId primitive.ObjectID, createdBy string, tokenType TokenType, environmentSecret []byte) Token {
	token := Token{}
	token.Id = primitive.NewObjectID()
	objectIdHex := token.Id.Hex()

	token.Description = description
	token.EnvironmentId = environmentId
	token.CreatedBy = createdBy
	token.TokenType = tokenType
	token.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Create token and secret
	checkToken := internal.RandomString(TokenLength)
	key := internal.RandomString(KeyLength)
	userTokenSecret := string(key) + objectIdHex + string(checkToken) + string(environmentSecret)
	userTokenB64 := base64.StdEncoding.EncodeToString([]byte(userTokenSecret))
	token.secret = userTokenB64

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf(err.Error())
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf(err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf(err.Error())
	}

	toBeSealed := append([]byte(userTokenSecret), checkToken...)

	dbToken := gcm.Seal(nonce, nonce, toBeSealed, nil)

	dbTokenB64 := base64.StdEncoding.EncodeToString(dbToken)
	token.Token = dbTokenB64

	return token
}

func (t *Token) GetSecret() string {
	return t.secret
}

func (t *Token) GetEnvironmentId() primitive.ObjectID {
	return t.EnvironmentId
}
