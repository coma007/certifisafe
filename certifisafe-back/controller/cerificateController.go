package controller

import (
	"certifisafe-back/domain"
	"certifisafe-back/model"
	"encoding/json"
	//"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type ICertificateService interface {
	UpdateCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type CertificateHandler struct {
	service domain.ICertificateService
}

func NewCertificateHandler(cs domain.ICertificateService) *CertificateHandler {
	return &CertificateHandler{service: cs}
}

// curl -X PATCH "localhost:8080/movies/1" -d '{ "title": "Beautiful film" }'
func (ch *CertificateHandler) UpdateCertificate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	var certificate model.Certificate
	err := json.NewDecoder(r.Body).Decode(&certificate)
	if err != nil {
		http.Error(w, "error when decoding json", http.StatusInternalServerError)
		return
	}

	err = ch.service.UpdateCertificate(id, certificate)

	if err != nil {
		// if errors.Is(err, service.ErrIDIsNotValid) ||
		// 	errors.Is(err, service.ErrTitleIsNotEmpty) {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// } else if errors.Is(err, service.ErrMovieNotFound) {
		// 	http.Error(w, err.Error(), http.StatusNotFound)
		// 	return
		// }
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Successfully Updated"))
}
