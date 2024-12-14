package account

import (
	"encoding/json"
	"fmt"
	"github.com/juankair/go_api_boilerplate/internal/entity"
	"github.com/juankair/go_api_boilerplate/pkg/log"
	"github.com/juankair/go_api_boilerplate/pkg/pagination"
	"github.com/juankair/go_api_boilerplate/pkg/response"
	"github.com/uptrace/bunrouter"
	"net/http"
)

func RegisterHandler(router *bunrouter.Group, service Service, log log.Logger) {
	res := resource{service, log}

	authGroup := router.NewGroup("/account")

	authGroup.GET("/list", res.get)
	authGroup.POST("/create", res.create)
	authGroup.PUT("/edit/:id", res.update)
	authGroup.PUT("/status/:id", res.toggleStatus)
	authGroup.DELETE("/delete/:id", res.delete)
	authGroup.DELETE("/delete/bulk", res.deleteBulk)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Fail Count", nil)
		return err
	}
	pages := pagination.NewFromRequest(req.Request, count)
	account, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		fmt.Println(err)
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Fail Paging", nil)
		return err
	}
	pages.Items = account

	response.RespondWithJSON(w, http.StatusOK, true, "Account Get Success", pages)
	return nil
}

func (r resource) create(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData FormAccountRequest

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	account, err := r.service.Create(ctx, requestData)
	if err != nil {
		return err
	}

	responseData := FormAccountRequest{
		AccountId:    account.AccountId,
		RoleId:       account.RoleId,
		EmployeeCode: account.EmployeeCode,
		FullName:     account.FullName,
		Email:        account.Email,
		PhoneNumber:  account.PhoneNumber,
		IsSuperAdmin: account.IsSuperAdmin,
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Account has been successfully created", responseData)
	return nil
}

func (r resource) update(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData FormAccountRequest

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	account, err := r.service.Update(ctx, req.Param("id"), requestData)
	if err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Update Failed", nil)
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Account has been successfully updated", account)
	return nil
}

func (r resource) toggleStatus(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	account, err := r.service.ToggleStatus(ctx, req.Param("id"))
	if err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Update Failed", nil)
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Account Status has been successfully changed", account)
	return nil
}

func (r resource) delete(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	account, err := r.service.Delete(ctx, req.Param("id"))
	if err != nil {
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Account has been successfully deleted", account)
	return nil
}

func (r resource) deleteBulk(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData []entity.Account

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	for _, value := range requestData {
		_, err := r.service.Delete(ctx, value.AccountId)
		if err != nil {
			return err
		}
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Account has been successfully deleted", requestData)
	return nil
}
