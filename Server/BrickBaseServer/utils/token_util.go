package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gianlucapalumbo0/BD2-TEAM-BRICKBASE/Server/BrickBaseServer/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SignedDetails rappresenta la struttura dei dati salvati dentro al token JWT
type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

// La funzione GenerateAllTokens crea e firma sia l'Access Token che il Refresh Token per un utente
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	// recupera le chiavi segrete di firma dalle variabili d'ambiente
	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "BrickBase",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// crea e firma l'Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "BrickBase",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}

	// crea e firma il Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil

}

// La funzione UpdateAllTokens salva o aggiorna i token JWT dell'utente e la data di modifica su MongoDB
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	if err != nil {
		return err
	}
	return nil
}

// La funzione GetAccessToken recupera il valore dell'Access Token dai cookie inviati dal client
func GetAccessToken(c *gin.Context) (string, error) {

	tokenString, err := c.Cookie("access_token")
	if err != nil {

		return "", err
	}

	return tokenString, nil

}

// La funzione ValidateToken decodifica e verifica la validità di un Access Token
func ValidateToken(tokenString string) (*SignedDetails, error) {

	// recupera la chiave segreta dell'Access Token dalle variabili d'ambiente
	secretKey := os.Getenv("SECRET_KEY")

	claims := &SignedDetails{}

	// verifica la firma del token usando la chiave segreta
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// La funzione GetUserIdFromContext estrae e valida l'ID utente dal contesto della richiesta Gin
func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")

	if !exists {
		return "", errors.New("userId does not exists in this context")
	}

	id, ok := userId.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil
}

// La funzione GetRoleFromContext estrae e valida il ruolo dell'utente dal contesto della richiesta Gin
func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("role does not exists in this context")
	}

	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("unable to retrieve role")
	}

	return memberRole, nil
}

// La funzione ValidateRefreshToken decodifica e verifica la validità di un Refresh Token
func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	// recupera la chiave segreta specifica per i Refresh Token dalle variabili d'ambiente
	secretRefreshKey := os.Getenv("SECRET_REFRESH_KEY")

	// verifica della firma usando la chiave segreta di refresh
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(secretRefreshKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
