package request

import (
	"certifisafe-back/features/auth"
	certificate2 "certifisafe-back/features/certificate"
	"certifisafe-back/utils"
	validation "github.com/go-ozzo/ozzo-validation"
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

func (controller *RequestController) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req NewRequestDTO
	err := utils.ReadRequestBody(w, r, &req)
	if err != nil {
		return
	}

	err = req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.authService.CheckRecaptcha(req.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	user := controller.authService.GetUserFromToken(token)
	request, err := controller.service.CreateRequest(&req, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, err, request, http.StatusCreated)
}

func (controller *RequestController) GetRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDfromUrl(w, r)
	if err != nil {
		return
	}

	request, err := controller.service.GetRequest(id)
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

func (controller *RequestController) GetAllRequestsByUserSigning(w http.ResponseWriter, r *http.Request) {
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))

	requests, err := controller.service.GetAllRequestsByUserSigning(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, err, requests, http.StatusOK)
}

func (controller *RequestController) GetAllRequestsByUser(w http.ResponseWriter, r *http.Request) {
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))

	//TODO check if error is returned

	requests, err := controller.service.GetAllRequestsByUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnResponse(w, err, requests, http.StatusOK)
}

func (controller *RequestController) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDfromUrl(w, r)
	if err != nil {
		return
	}

	err = controller.service.DeleteRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDfromUrl(w, r)
	if err != nil {
		return
	}

	_, err = controller.service.AcceptRequest(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDfromUrl(w, r)
	if err != nil {
		return
	}

	reason := struct {
		Reason string
	}{}

	err = utils.ReadRequestBody(w, r, &reason)
	if err != nil {
		http.Error(w, "invalid reason", http.StatusBadRequest)
		return
	}

	err = validation.ValidateStruct(&reason,
		validation.Field(&reason.Reason, validation.Required, validation.Length(3, 50)),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.service.DeclineRequest(id, reason.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *RequestController) GenerateCertificates(w http.ResponseWriter, r *http.Request) {
	// dummy data
	user := controller.authService.GetUserFromToken(r.Header.Get("Authorization"))
	rootDTO := &NewRequestDTO{
		ParentSerial:    nil,
		CertificateName: "some root name",
		CertificateType: "ROOT",
	}

	root, err := controller.certificateService.CreateCertificate(rootDTO.ParentSerial, rootDTO.CertificateName, certificate2.StringToType(rootDTO.CertificateType), user.ID)
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
		CertificateName: "localhost",
		CertificateType: "INTERMEDIATE",
	}
	intermidiate, err := controller.certificateService.CreateCertificate(intermediateDTO.ParentSerial, intermediateDTO.CertificateName, certificate2.StringToType(intermediateDTO.CertificateType), user.ID)

	intermediateSerial := uint(*intermidiate.Serial)
	leafDTO := &NewRequestDTO{
		ParentSerial:    &intermediateSerial,
		CertificateName: "end",
		CertificateType: "END",
	}

	leaf, err := controller.certificateService.CreateCertificate(leafDTO.ParentSerial, leafDTO.CertificateName, certificate2.StringToType(leafDTO.CertificateType), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.ReturnResponse(w, err, []certificate2.CertificateDTO{root, intermidiate, leaf}, http.StatusCreated)
}
