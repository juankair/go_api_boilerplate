package keperluan

import (
	"encoding/json"
	"fmt"
	"github.com/juankair/docs_sign_be/internal/entity"
	"github.com/juankair/docs_sign_be/pkg/log"
	"github.com/juankair/docs_sign_be/pkg/pagination"
	"github.com/juankair/docs_sign_be/pkg/response"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strconv"
)

func RegisterHandler(router *bunrouter.Group, service Service, log log.Logger) {
	res := resource{service, log}

	authGroup := router.NewGroup("/keperluan")

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
	keperluan, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		fmt.Println(err)
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Fail Paging", nil)
		return err
	}
	pages.Items = keperluan

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan Get Success", pages)
	return nil
}

func (r resource) create(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData FormKeperluanRequest

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	keperluan, err := r.service.Create(ctx, requestData)
	if err != nil {
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan has been successfully created", keperluan)
	return nil
}

func (r resource) update(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData FormKeperluanRequest

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	idStr := req.Param("id")
	id, err := strconv.Atoi(idStr)
	keperluan, err := r.service.Update(ctx, id, requestData)
	if err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Update Failed", nil)
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan has been successfully updated", keperluan)
	return nil
}

func (r resource) toggleStatus(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	idStr := req.Param("id")
	id, err := strconv.Atoi(idStr)

	keperluan, err := r.service.ToggleStatus(ctx, id)
	if err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Update Failed", nil)
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan Status has been successfully changed", keperluan)
	return nil
}

func (r resource) delete(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	idStr := req.Param("id")
	id, err := strconv.Atoi(idStr)

	keperluan, err := r.service.Delete(ctx, id)
	if err != nil {
		return err
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan has been successfully deleted", keperluan)
	return nil
}

func (r resource) deleteBulk(w http.ResponseWriter, req bunrouter.Request) error {
	var requestData []entity.Keperluan

	ctx := req.Context()

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		response.RespondWithJSON(w, http.StatusBadRequest, false, "Invalid Value", nil)
		return nil
	}

	for _, value := range requestData {
		_, err := r.service.Delete(ctx, value.ID)
		if err != nil {
			return err
		}
	}

	response.RespondWithJSON(w, http.StatusOK, true, "Keperluan has been successfully deleted", requestData)
	return nil
}
