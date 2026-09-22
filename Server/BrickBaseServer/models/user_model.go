package models

import "go.mongodb.org/mongo-driver/v2/bson"

// User rappresenta il documenti principale dell'utente all'interno della collezione "users" in MongoDB
// ID è l'identificatore primario univoco generato da MongoDB
// UserID è l'identificativo univoco dell'utente
// FirstName è il nome dell'utente
// LastName è il cognome dell'utente
// Email è l'indirizzo email dell'utente
// Role definisce i permessi dell'utente nell'applicazione
// Token memorizza l'ultimo JWT Access Token generato per la sessione corrente dell'utente
// RefreshToken memorizza il token di ripristino utile per richiedere un nuovo Access Token senza re-autenticarsi
// Password contiene l'hash della password dell'utente
type User struct {
	ID           bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID       string        `json:"user_id" bson:"user_id"`
	FirstName    string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName     string        `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email        string        `json:"email" bson:"email" validate:"required,email"`
	Role         string        `json:"role" bson:"role" validate:"oneof=ADMIN USER"`
	Token        string        `json:"token" bson:"token"`
	RefreshToken string        `json:"refresh_token" bson:"refresh_token"`
	Password     string        `json:"password" bson:"password" validate:"required,min=6"`
}

// UserLogin è un DTO (Data Transfer Object) utilizzato per decodificare il blody JSON nelle richeste di login
// Email è l'indirizzo e-mail fornito dall'utente durante l'autenticazione
// Password è la password in chiaro inviata dall'utente per effettuare il login
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse è un DTO di risposta utilizzato per restituire i dati dell'utente al client dopo l'autenticazione
// UserId è l'identificativo unico dell'utente inviato al client
// FirstName è il nome dell'utente
// LastName è il cognome dell'utente
// Email è l'indirizzo e-mail dell'utente
// Role specifica il ruolo assegnato ("ADMIN" o "USER")
// Token è il JWT Access Token da utilizzare nelle intestazioni Authorization delle richieste successive
// RefreshToken è il token per rinnovare l'Access Token quando scade
type UserResponse struct {
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
