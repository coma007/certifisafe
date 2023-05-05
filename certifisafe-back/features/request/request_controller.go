package request

import (
	"certifisafe-back/features/auth"
	certificate2 "certifisafe-back/features/certificate"
	"certifisafe-back/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type RequestController struct {
	service            RequestService
	certificateService certificate2.CertificateService
	authService        auth.AuthService
}

func NewRequestController(service RequestService, certificateService certificate2.CertificateService, authService auth.AuthService) *RequestController {
	return &RequestController{service: service, certificateService: certificateService, authService: authService}
}

func (c *RequestController) CreateRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req NewRequestDTO
	err := utils.ReadRequestBody(w, r, &req)
	if err != nil {
		return
	}

	request, err := c.service.CreateRequest(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, err, request, http.StatusCreated)
}

func (c *RequestController) GetRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadIDfromUrl(w, ps)
	if err != nil {
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

	utils.ReturnResponse(w, err, request, http.StatusOK)
}

func (controller *RequestController) GetAllRequestsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))

	requests, err := controller.service.GetAllRequestsByUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, err, requests, http.StatusOK)
}

func (c *RequestController) DeleteRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadIDfromUrl(w, ps)
	if err != nil {
		return
	}

	err = c.service.DeleteRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) AcceptRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := utils.ReadIDfromUrl(w, ps)
	if err != nil {
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
	id, err := utils.ReadIDfromUrl(w, ps)
	if err != nil {
		return
	}

	err = controller.service.DeclineRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) GenerateCertificates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// dummy data
	rootDTO := &NewRequestDTO{
		ParentSerial:    nil,
		CertificateName: "root",
		CertificateType: "ROOT",
		SubjectId:       1,
	}

	root, err := controller.certificateService.CreateCertificate(rootDTO.ParentSerial, rootDTO.CertificateName, rootDTO.CertificateType, rootDTO.SubjectId)
	if err != nil {
		panic(err)
	}
	//rootSerial := new(big.Int)
	//rootSerial.SetString(root.Serial, 10)
	//rootCreated, err := controller.service.GetCertificate(*rootSerial)
	if err != nil {
		panic(err)
	}

	parentSerial := uint(*root.Serial)
	intermediateDTO := &NewRequestDTO{
		ParentSerial:    &parentSerial,
		CertificateName: "intermediate",
		CertificateType: "INTERMEDIATE",
		SubjectId:       1,
	}
	intermidiate, err := controller.certificateService.CreateCertificate(intermediateDTO.ParentSerial, intermediateDTO.CertificateName, intermediateDTO.CertificateType, intermediateDTO.SubjectId)

	intermediateSerial := uint(*intermidiate.Serial)
	leafDTO := &NewRequestDTO{
		ParentSerial:    &intermediateSerial,
		CertificateName: "end",
		CertificateType: "END",
		SubjectId:       1,
	}

	leaf, err := controller.certificateService.CreateCertificate(leafDTO.ParentSerial, leafDTO.CertificateName, leafDTO.CertificateType, leafDTO.SubjectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.ReturnResponse(w, err, []certificate2.CertificateDTO{root, intermidiate, leaf}, http.StatusCreated)
}
