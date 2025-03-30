package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Функция для получения JWKS и поиска нужного ключа по kid
// func getVKPublicKey(token *jwt.Token) (interface{}, error) {
// 	// Проверяем, что алгоритм подписи - RS256
// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 		return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
// 	}

// 	// Получаем kid из заголовка токена
// 	kid, ok := token.Header["kid"].(string)
// 	if !ok {
// 		return nil, errors.New("kid header not found")
// 	}

// 	// Загружаем JWKS из VK (в реальном приложении это лучше кэшировать)
// 	jwksURL := "https://id.vk.com/.well-known/jwks"
// 	resp, err := http.Get(jwksURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var jwks struct {
// 		Keys []struct {
// 			Kid string `json:"kid"`
// 			N   string `json:"n"`
// 			E   string `json:"e"`
// 			Kty string `json:"kty"`
// 		} `json:"keys"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
// 		return nil, err
// 	}

// 	// Ищем ключ с соответствующим kid
// 	for _, key := range jwks.Keys {
// 		if key.Kid == kid {
// 			// Декодируем modulus (n) и exponent (e) из base64
// 			modulus, err := base64.RawURLEncoding.DecodeString(key.N)
// 			if err != nil {
// 				return nil, err
// 			}
// 			exponent, err := base64.RawURLEncoding.DecodeString(key.E)
// 			if err != nil {
// 				return nil, err
// 			}

// 			// Создаем RSA публичный ключ
// 			pubKey := &rsa.PublicKey{
// 				N: big.NewInt(0).SetBytes(modulus),
// 				E: int(big.NewInt(0).SetBytes(exponent).Int64()),
// 			}
// 			return pubKey, nil
// 		}
// 	}

// 	return nil, errors.New("public key not found for kid: " + kid)
// }

// ExtractVKIDCredentials извлекает user_id из id_token VK
func ExtractVKIDCredentials(tokenString string) (int64, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// PEM-форматированный публичный ключ
		publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvsvJlhFX9Ju/pvCz1frB
DgJs592VjdwQuRAmnlJAItyHkoiDIOEocPzgcUBTbDf1plDcTyO2RCkUt0pz0WK6
6HNhpJyIfARjaWHeUlv4TpuHXAJJsBKklkU2gf1cjID+40sWWYjtq5dAkXnSJUVA
UR+sq0lJ7GmTdJtAr8hzESqGEcSP15PTs7VUdHZ1nkC2XgkuR8KmKAUb388ji1Q4
n02rJNOPQgd9r0ac4N2v/yTAFPXumO78N25bpcuWf5vcL9e8THk/U2zt7wf+aAWL
748e0pREqNluTBJNZfmhC79Xx6GHtwqHyyduiqfPmejmiujNM/rqnA4e30Tg86Yn
cNZ6vLJyF72Eva1wXchukH/aLispbY+EqNPxxn4zzCWaLKHG87gaCxpVv9Tm0jSD
2es22NjrUbtb+2pAGnXbyDp2eGUqw0RrTQFZqt/VcmmSCE45FlcZMT28otrwG1ZB
kZAb5Js3wLEch3ZfYL8sjhyNRPBmJBrAvzrd8qa3rdUjkC9sKyjGAaHu2MNmFl1Y
JFQ3J54tGpkGgJjD7Kz3w0K6OiPDlVCNQN5sqXm24fCw85Pbi8SJiaLTp/CImrs1
Z3nHW5q8hljA7OGmqfOP0nZS/5zW9GHPyepsI1rW6CympYLJ15WeNzePxYS5KEX9
EncmkSD9b45ge95hJeJZteUCAwEAAQ==
-----END PUBLIC KEY-----`

		// Декодируем PEM
		block, _ := pem.Decode([]byte(publicKeyPEM))
		if block == nil || block.Type != "PUBLIC KEY" {
			return nil, errors.New("failed to decode PEM block containing public key")
		}

		// Парсим публичный ключ
		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		// Проверяем, что ключ является *rsa.PublicKey
		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("key is not of type *rsa.PublicKey")
		}

		return rsaPubKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Проверяем, что токен валиден
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Дополнительные проверки токена
	issuer := claims["iis"]
	if issuer != "VK" {
		return 0, errors.New("invalid issuer")
	}

	// Извлекаем время истечения
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}
	if exp.Before(time.Now()) {
		return 0, errors.New("token expired")
	}

	// Извлекаем user_id из sub
	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid sub format")
	}

	return int64(userID), nil
}

// Вспомогательная функция для проверки audience
func contains(slice jwt.ClaimStrings, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
