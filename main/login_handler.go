package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ybakhan/tax-calculator/common"
)

// handleLogin handles login api call go doc
//
//	@Summary		login to taxes api
//	@Description	returns api key for calling taxes api
//	@Tags			taxes
//	@Produce		json
//	@Param			user	body		user	true	"username password"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	taxServerError
//	@Failure		401		{object}	taxServerResponse
//	@Failure		405		{object}	taxServerError
//	@Failure		500		{object}	taxServerError
//	@Router			/login [post]
func (s *taxServer) handleLogin(w http.ResponseWriter, r *http.Request) (*handlerResponse, error) {
	if r.Method != "POST" {
		return &handlerResponse{Status: http.StatusMethodNotAllowed}, fmt.Errorf("method not supported %s", r.Method)
	}

	var user user
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return &handlerResponse{Status: http.StatusBadRequest}, err
	}

	if user.Username != "admin" || user.Password != "password" {
		return &handlerResponse{http.StatusUnauthorized, taxServerResponse{"invalid username password"}}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * time.Duration(s.ApiTokenConfig.ExpirationMinutes)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.ApiTokenConfig.Secret))
	if err != nil {
		return &handlerResponse{Status: http.StatusInternalServerError}, err
	}

	return &handlerResponse{http.StatusOK, &loginResponse{tokenString}}, nil
}

func (s *taxServer) validateApiKey(next requestHandler) requestHandler {
	return func(w http.ResponseWriter, r *http.Request) (*handlerResponse, error) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			s.Logger.Log("requestID", common.GetRequestID(r.Context()), "message", "Authorization header missing")
			return &handlerResponse{http.StatusUnauthorized, taxServerResponse{"Authorization header missing"}}, nil
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.ApiTokenConfig.Secret), nil
		})

		if err != nil {
			return &handlerResponse{Status: http.StatusUnauthorized}, err
		}

		if !token.Valid {
			return &handlerResponse{http.StatusUnauthorized, taxServerResponse{"Invalid token"}}, nil
		}

		// Token is valid, continue with the API logic
		return next(w, r)
	}
}
