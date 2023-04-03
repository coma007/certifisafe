package controller

import (
	"certifisafe-back/dto"
	"certifisafe-back/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type RequestController struct {
	service service.RequestService
}

func NewRequestController(service service.RequestService) *RequestController {
	return &RequestController{service: service}
}

func (c *RequestController) GetRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	request, err := c.service.GetRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request == nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (controller *RequestController) GetAllRequests(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requests, err := controller.service.GetAllRequests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (c *RequestController) CreateRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req dto.NewRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	request, err := c.service.CreateRequest(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func (controller *RequestController) AcceptRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	err = controller.service.AcceptRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) DeclineRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	err = controller.service.DeclineRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *RequestController) DeleteRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
