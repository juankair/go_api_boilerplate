package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juankair/docs_sign_be/pkg/log"
	"github.com/juankair/docs_sign_be/pkg/response"
	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(router *bunrouter.Router, service Service, log log.Logger) {

	authGroup := router.NewGroup("/auth")

	authGroup.POST("/login", login(service, log))
	authGroup.POST("/activation", activation(service, log))
	authGroup.POST("/activation-confirmation", activationConfirmation(service, log))
}

func login(service Service, logger log.Logger) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		ctx := req.Context()

		if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
			response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
			return nil
		}

		data, err := service.Login(ctx, requestData.Email, requestData.Password)
		if err != nil {
			response.RespondWithJSON(w, http.StatusForbidden, false, fmt.Sprint(err), nil)
			return nil
		}

		response.RespondWithJSON(w, http.StatusOK, true, "Login Success", data)
		return nil
	}
}

func activation(service Service, logger log.Logger) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var requestData struct {
			Email string `json:"email"`
		}

		ctx := req.Context()

		if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
			response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
			return nil
		}

		data, err := service.CheckActivation(ctx, requestData.Email)
		if err != nil {
			response.RespondWithJSON(w, http.StatusForbidden, false, fmt.Sprint(err), nil)
			return nil
		}

		response.RespondWithJSON(w, http.StatusOK, true, "Activation Allowed", data)
		return nil
	}
}

func activationConfirmation(service Service, logger log.Logger) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var requestData struct {
			Id       string `json:"id"`
			Password string `json:"password"`
		}

		ctx := req.Context()

		if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
			response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
			return nil
		}

		account, err := service.Get(ctx, requestData.Id)

		if account.IsActive == 1 {
			response.RespondWithJSON(w, http.StatusForbidden, false, "Akun sudah dikonfirmasi", nil)
			return nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
		if err != nil {
			response.RespondWithJSON(w, http.StatusInternalServerError, false, "Error hashing password", nil)
			return nil
		}

		activationConfirmation, err := service.ActivationConfirmation(ctx, account.Email, string(hashedPassword))
		if err != nil {
			response.RespondWithJSON(w, http.StatusBadRequest, false, "Update Failed", nil)
			return err
		}

		response.RespondWithJSON(w, http.StatusOK, true, "Activation Successed", activationConfirmation)
		return nil
	}
}
