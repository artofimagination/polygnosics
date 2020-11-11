package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Store key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")), []byte(os.Getenv("SESSION_ENCRYPTION")))
)

func EncryptUserAndOrigin(userid uuid.UUID, origin string) (*[]byte, error) {
	data := make(map[string]interface{})
	data["userid"] = userid.String()
	data["origin"] = origin

	binary, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(binary, 4)
	if err != nil {
		return nil, err
	}

	return &hashedPassword, nil
}

func matchingUserAndOrigin(userid uuid.UUID, origin string, cookieData *[]byte) (bool, error) {
	data := make(map[string]interface{})
	data["userid"] = userid.String()
	data["origin"] = origin

	binary, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword(*cookieData, binary); err != nil {
		return false, nil
	}

	return true, nil
}

func IsAuthenticated(userID uuid.UUID, session *sessions.Session, r *http.Request) (bool, error) {
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false, nil
	}

	binary, ok := session.Values["cookie_key"].([]byte)
	if !ok {
		return false, errors.New("Failed to decode cookie key")
	}

	match, err := matchingUserAndOrigin(userID, r.RemoteAddr, &binary)
	if err != nil {
		return false, err
	}

	return match, nil
}
