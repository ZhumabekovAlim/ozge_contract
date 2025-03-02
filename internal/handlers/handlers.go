package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"ozge/internal/models"
	"ozge/internal/services"
	"path/filepath"
	"strconv"
)

// TOO Handler
type TOOHandler struct {
	Service *services.TOOService
}

func (h *TOOHandler) CreateTOO(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(30 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusInternalServerError)
		return
	}

	// Retrieve form values
	name := r.FormValue("name")
	bin := r.FormValue("bin")
	bankDetails := r.FormValue("bank_details")
	email := r.FormValue("email")
	signer := r.FormValue("signer")
	iin := r.FormValue("iin")
	companyCode := r.FormValue("company_code")
	additionalInformation := r.FormValue("additional_information")

	// Create a TOO object
	too := models.TOO{
		Name:                  name,
		BIN:                   bin,
		BankDetails:           bankDetails,
		Email:                 email,
		Signer:                signer,
		IIN:                   iin,
		Token:                 "",
		CompanyCode:           companyCode,
		AdditionalInformation: additionalInformation,
	}

	// Call the service layer to save the TOO
	createdTOO, err := h.Service.CreateTOO(r.Context(), too)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTOO)
}

func (h *TOOHandler) UpdateUserContract(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	bin := r.FormValue("bin")
	company_code := r.FormValue("company_code")
	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"user_contract": "Подписанный_договор_пользователя.pdf",
	}

	// Save files
	savedFiles := map[string]string{}
	for formField, fileName := range fileFieldNames {
		file, _, err := r.FormFile(formField)
		if err == nil {
			defer file.Close()

			filePath := fmt.Sprintf("uploads/TOO/%s/%s", bin, fileName)
			err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create directory", http.StatusInternalServerError)
				return
			}

			out, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				http.Error(w, "Unable to save file content", http.StatusInternalServerError)
				return
			}
			savedFiles[formField] = filePath
		}
	}

	idInt, err := strconv.ParseInt(id, 10, 64)

	// Create a TOO object
	too := models.TOO{
		ID:           int(idInt),
		BIN:          bin,
		CompanyCode:  company_code,
		UserContract: savedFiles["user_contract"],
	}

	// Call the service layer to save the TOO
	createdTOO, err := h.Service.UpdateContractTOO(r.Context(), too)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdTOO)
	if err != nil {
		return
	}
}

// IP Handler
type IPHandler struct {
	Service *services.IPService
}

func (h *IPHandler) CreateIP(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		fmt.Println("123213")
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	name := r.FormValue("name")
	bin := r.FormValue("bin")
	bankDetails := r.FormValue("bank_details")
	email := r.FormValue("email")
	signer := r.FormValue("signer")
	iin := r.FormValue("iin")
	companyCode := r.FormValue("company_code")
	additionalInformation := r.FormValue("additional_information")

	// Create an IP object
	ip := models.IP{
		Name:                  name,
		BIN:                   bin,
		BankDetails:           bankDetails,
		Email:                 email,
		Signer:                signer,
		IIN:                   iin,
		Token:                 "",
		CompanyCode:           companyCode,
		AdditionalInformation: additionalInformation,
	}

	// Call the service layer to save the IP
	createdIP, err := h.Service.CreateIP(r.Context(), ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created IP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdIP)
}

func (h *IPHandler) UpdateUserContract(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	iin := r.FormValue("iin")
	company_code := r.FormValue("company_code")
	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"user_contract": "Подписанный_договор_пользователя.pdf",
	}

	// Save files
	savedFiles := map[string]string{}
	for formField, fileName := range fileFieldNames {
		file, _, err := r.FormFile(formField)
		if err == nil {
			defer file.Close()

			filePath := fmt.Sprintf("uploads/IP/%s/%s", iin, fileName)
			err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create directory", http.StatusInternalServerError)
				return
			}

			out, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				http.Error(w, "Unable to save file content", http.StatusInternalServerError)
				return
			}
			savedFiles[formField] = filePath
		}
	}

	idInt, err := strconv.ParseInt(id, 10, 64)

	// Create a TOO object
	ip := models.IP{
		ID:           int(idInt),
		IIN:          iin,
		CompanyCode:  company_code,
		UserContract: savedFiles["user_contract"],
	}

	// Call the service layer to save the TOO
	createdIP, err := h.Service.UpdateContractIP(r.Context(), ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdIP)
	if err != nil {
		return
	}
}

// Individual Handler
type IndividualHandler struct {
	Service *services.IndividualService
}

func (h *IndividualHandler) CreateIndividual(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB

	if err != nil {

		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}
	fmt.Println(r)
	// Retrieve form values
	fullName := r.FormValue("full_name")
	iin := r.FormValue("iin")
	email := r.FormValue("email")
	// File field mapping to table headers
	companyCode := r.FormValue("company_code")
	additional_information := r.FormValue("additional_information")

	// Create an Individual object
	individual := models.Individual{
		FullName:              fullName,
		IIN:                   iin,
		Email:                 email,
		CompanyCode:           companyCode,
		Token:                 "",
		AdditionalInformation: additional_information,
	}

	// Call the service layer to save the individual
	createdIndividual, err := h.Service.CreateIndividual(r.Context(), individual)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created individual
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdIndividual)
}

func (h *IndividualHandler) UpdateUserContract(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	iin := r.FormValue("iin")
	company_code := r.FormValue("company_code")
	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"user_contract": "Подписанный_договор_пользователя.pdf",
	}

	// Save files
	savedFiles := map[string]string{}
	for formField, fileName := range fileFieldNames {
		file, _, err := r.FormFile(formField)
		if err == nil {
			defer file.Close()

			filePath := fmt.Sprintf("uploads/Individual/%s/%s", iin, fileName)
			err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create directory", http.StatusInternalServerError)
				fmt.Println("1")
				return
			}

			out, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				fmt.Println("2")
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				http.Error(w, "Unable to save file content", http.StatusInternalServerError)
				fmt.Println("3")
				return
			}
			savedFiles[formField] = filePath
		}
	}

	fmt.Println("user_contarct:", savedFiles["user_contract"])

	idInt, err := strconv.ParseInt(id, 10, 64)

	// Create a TOO object
	individual := models.Individual{
		ID:           int(idInt),
		IIN:          iin,
		CompanyCode:  company_code,
		UserContract: savedFiles["user_contract"],
	}

	// Call the service layer to save the TOO
	createdIndividual, err := h.Service.UpdateContractIndividual(r.Context(), individual)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdIndividual)
	if err != nil {
		return
	}
}

// Search TOO by IIN
func (h *TOOHandler) SearchTOOs(w http.ResponseWriter, r *http.Request) {
	iin := r.URL.Query().Get(":iin")
	if iin == "" {
		http.Error(w, "Не указан параметр 'iin'", http.StatusBadRequest)
		return
	}

	pass := r.URL.Query().Get(":pass")
	if pass == "" {
		http.Error(w, "Не указан параметр 'pass'", http.StatusBadRequest)
		return
	}

	toos, err := h.Service.SearchTOOsByBIN(r.Context(), iin, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toos)
}

// Search IP by IIN
func (h *IPHandler) SearchIPs(w http.ResponseWriter, r *http.Request) {
	iin := r.URL.Query().Get(":iin")
	if iin == "" {
		http.Error(w, "Не указан параметр 'iin'", http.StatusBadRequest)
		return
	}

	pass := r.URL.Query().Get(":pass")
	if pass == "" {
		http.Error(w, "Не указан параметр 'bin'", http.StatusBadRequest)
		return
	}

	ips, err := h.Service.SearchIPsByIIN(r.Context(), iin, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ips)
}

// Search Individual by IIN

// SearchIndividuals ищет пользователя по IIN и возвращает JSON + PDF (если есть)
func (h *IndividualHandler) SearchIndividuals(w http.ResponseWriter, r *http.Request) {

	iin := r.URL.Query().Get(":iin")
	if iin == "" {
		http.Error(w, `{"error": "Не указан параметр 'iin'"}`, http.StatusBadRequest)
		return
	}

	pass := r.URL.Query().Get(":pass")
	if pass == "" {
		http.Error(w, "Не указан параметр 'pass' ", http.StatusBadRequest)
		return
	}

	// Получаем данные из сервиса
	individuals, err := h.Service.SearchIndividualsByIIN(r.Context(), iin, pass)
	fmt.Println("ind: ", individuals)
	if err != nil || len(individuals) == 0 {
		http.Error(w, `{"error": "Пользователь не найден"}`, http.StatusNotFound)
		return
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(individuals)
}

type CompanyDataHandler struct {
	Service *services.CompanyDataService
}

func (h *CompanyDataHandler) GetAllDataByIIN(w http.ResponseWriter, r *http.Request) {
	iin := r.URL.Query().Get(":iin")
	if iin == "" {
		http.Error(w, `{"error": "Не указан параметр 'iin'"}`, http.StatusBadRequest)
		return
	}

	pass := r.URL.Query().Get(":pass")
	if pass == "" {
		http.Error(w, "Не указан параметр 'pass' ", http.StatusBadRequest)
		return
	}

	// Получаем данные из сервиса
	data, err := h.Service.GetAllDataByIIN(r.Context(), iin, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Отправляем данные в JSON-формате
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *TOOHandler) SearchTOOsByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")

	too, err := h.Service.SearchTOOByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "TOO not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(too)
}

func (h *IPHandler) SearchIPsByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")

	ip, err := h.Service.SearchIPByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "IP not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ip)
}

func (h *IndividualHandler) SearchIndividualsByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get(":token")

	individual, err := h.Service.SearchIndividualByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Individual not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(individual)
}

func (h *TOOHandler) SearchTOOsByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	too, err := h.Service.SearchTOOsByID(r.Context(), id)
	if err != nil {
		http.Error(w, "TOO not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(too)
}

func (h *IPHandler) SearchIPsByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	ip, err := h.Service.SearchIPByID(r.Context(), id)
	if err != nil {
		http.Error(w, "IP not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ip)
}

func (h *IndividualHandler) SearchIndividualsByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	individual, err := h.Service.SearchIndividualByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Individual not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(individual)
}

func (h *TOOHandler) UpdateUserContractStatus(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get(":id")
	if id == "" {
		http.Error(w, "Не указан параметр 'id'", http.StatusBadRequest)
		return
	}

	// Call the service layer to save the TOO
	err := h.Service.UpdateUserContractStatus(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *IPHandler) UpdateUserContractStatus(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get(":id")
	if id == "" {
		http.Error(w, "Не указан параметр 'id'", http.StatusBadRequest)
		return
	}

	// Call the service layer to save the TOO
	err := h.Service.UpdateUserContractStatus(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *IndividualHandler) UpdateUserContractStatus(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get(":id")
	if id == "" {
		http.Error(w, "Не указан параметр 'id'", http.StatusBadRequest)
		return
	}

	// Call the service layer to save the TOO
	err := h.Service.UpdateUserContractStatus(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created TOO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

type CompanyHandler struct {
	Service *services.CompanyService
}

// Создание компании
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var company models.Company

	// Декодируем JSON-запрос в структуру
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем компанию через сервис
	createdCompany, err := h.Service.Create(r.Context(), company)
	if err != nil {
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ с созданной компанией
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCompany)
}
