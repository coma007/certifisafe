package utils

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func ReadIDfromUrl(w http.ResponseWriter, r *http.Request) (int, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return 0, err
	}
	return id, err
}

func ReadCertificateIDFromUrl(w http.ResponseWriter, r *http.Request) (big.Int, error) {
	params := mux.Vars(r)
	id, err := StringToBigInt(params["id"])
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return big.Int{}, err
	}
	return id, err
}

func ReadVerificationCodeFromUrl(w http.ResponseWriter, r *http.Request) string {
	params := mux.Vars(r)
	code := params["verificationCode"]
	return code
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
