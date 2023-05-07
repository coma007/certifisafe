package utils

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func ReadIDfromUrl(w http.ResponseWriter, ps httprouter.Params) (int, error) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return 0, err
	}
	return id, err
}

func ReadCertificateIDFromUrl(w http.ResponseWriter, ps httprouter.Params) (big.Int, error) {
	id, err := StringToBigInt(ps.ByName("id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return big.Int{}, err
	}
	return id, err
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, req interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	return err
}

func ReturnResponse(w http.ResponseWriter, err error, data interface{}, httpStatus int) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}

func AddFileToResponse(w http.ResponseWriter, publicPath string) bool {
	f, err := os.Open(publicPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return false
	}
	defer f.Close()
	io.Copy(w, f)
	return true
}
