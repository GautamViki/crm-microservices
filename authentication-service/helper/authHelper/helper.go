package authhelper

import (
	"authentication-service/config"
	h "authentication-service/helper"
	"authentication-service/internals/dto"
	"authentication-service/models"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// Function to validate the token
func ValidateToken(tokenString string, db *gorm.DB) error {
	publicKey, err := ParseTokenUnvarified(tokenString)
	if err != nil {
		return fmt.Errorf("unable to fetch public key: %s", err.Error())
	}
	// Parse the public key from PEM format
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return errors.New("failed to parse public key")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	// Validate the token using the parsed public key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return parsedPublicKey, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("token invalid")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return errors.New("token is expired")
	}

	var user models.User
	db.First(&user, claims["aud"])
	if user.EntityId == 0 {
		return errors.New("user not found for token")
	}
	return nil
}

// Function to generate token
func GenerateToken(user models.User, redis *redis.Client) (dto.TokenResponse, error) {
	kid := h.GenerateUuidV4()
	jwkConfigKey, err := GenerateJwkConfigSecret(kid)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Parse the private key from PEM format
	block, _ := pem.Decode([]byte(jwkConfigKey.PrivateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		return dto.TokenResponse{}, errors.New("failed to parse private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// Create the token using RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud":       user.EntityId,
		"firstname": user.FirstName,
		"email":     user.Email,
		"exp":       time.Now().Add(time.Minute * 60).Unix(),
	})

	// Sign the token with the RSA private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	jwkKey := h.GenerateJwkRedisKey(user.EntityId)
	if err := redis.Set(context.Background(), jwkKey, jwkConfigKey.PublicKey, 5*time.Minute).Err(); err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		Token:     tokenString,
		TokenType: "Bearer",
		ExpiresIn: int64(time.Now().Add(time.Minute * 5).Unix()),
	}, nil
}

// Function to generate JWK config secrets
func GenerateJwkConfigSecret(kid string) (dto.IdsConfig, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return dto.IdsConfig{}, err
	}

	// Encode the private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Encode the public key to PEM format
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return dto.IdsConfig{}, err
	}
	publicKeyPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEM,
	})

	return dto.IdsConfig{
		PrivateKey: string(privateKeyPEM),
		PublicKey:  string(publicKeyPEMBytes),
	}, nil
}

func ParseTokenUnvarified(tokenString string) (string, error) {
	redis := config.ConnectRedis()
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("unable to fetch token claims")
	}
	jwkKey := h.GenerateJwkRedisKey(mapClaims["aud"])
	publicKey, err := redis.Get(context.Background(), jwkKey).Result()
	if err != nil {
		return "", err
	}
	return publicKey, nil
}
