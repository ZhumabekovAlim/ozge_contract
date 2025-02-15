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
)

// TOO Handler
type TOOHandler struct {
	Service *services.TOOService
}

func (h *TOOHandler) CreateTOO(w http.ResponseWriter, r *http.Request) {
	// Parse the form-data
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Retrieve form values
	name := r.FormValue("name")
	bin := r.FormValue("bin")
	ceo_name := r.FormValue("ceo_name")
	bankDetails := r.FormValue("bank_details")
	legalAddress := r.FormValue("legal_address")
	actualAddress := r.FormValue("actual_address")
	contactDetails := r.FormValue("contact_details")
	email := r.FormValue("email")
	companyCode := r.FormValue("company_code")

	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"registration_file":  "Справка_о_регистрации.pdf",
		"ceo_order_file":     "Приказ_о_назначении.pdf",
		"ceo_id_file":        "Удостоверение_руководителя.pdf",
		"representative_poa": "Доверенность_представителя.pdf",
		"representative_id":  "Удостоверение_представителя.pdf",
		"egov_file":          "Адресная_справка.pdf",
		"company_card":       "Карточка_предприятия.pdf",
		"user_contract":      "Подписанный_договор_пользователя.pdf",
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

	// Create a TOO object
	too := models.TOO{
		Name:              name,
		BIN:               bin,
		RegistrationFile:  savedFiles["registration_file"],
		CEOName:           ceo_name,
		CEOOrderFile:      savedFiles["ceo_order_file"],
		CEOIDFile:         savedFiles["ceo_id_file"],
		RepresentativePOA: savedFiles["representative_poa"],
		RepresentativeID:  savedFiles["representative_id"],
		BankDetails:       bankDetails,
		LegalAddress:      legalAddress,
		ActualAddress:     actualAddress,
		ContactDetails:    contactDetails,
		Email:             email,
		EgovFile:          savedFiles["egov_file"],
		CompanyCard:       savedFiles["company_card"],
		CompanyCode:       companyCode,
		UserContract:      savedFiles["user_contract"],
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
	iin := r.FormValue("iin")
	bankDetails := r.FormValue("bank_details")
	legalAddress := r.FormValue("legal_address")
	actualAddress := r.FormValue("actual_address")
	contactDetails := r.FormValue("contact_details")
	email := r.FormValue("email")
	companyCode := r.FormValue("company_code")

	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"registration_file":  "Талон_о_регистрации.pdf",
		"representative_poa": "Доверенность_представителя.pdf",
		"representative_id":  "Удостоверение_личности_представителя_по_доверенности.pdf",
		"company_card":       "Карточка_предприятия.pdf",
		"user_contract":      "Подписанный_договор_пользователя.pdf",
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

	// Create an IP object
	ip := models.IP{
		Name:              name,
		IIN:               iin,
		RegistrationFile:  savedFiles["registration_file"],
		RepresentativePOA: savedFiles["representative_poa"],
		RepresentativeID:  savedFiles["representative_id"],
		BankDetails:       bankDetails,
		LegalAddress:      legalAddress,
		ActualAddress:     actualAddress,
		ContactDetails:    contactDetails,
		Email:             email,
		CompanyCard:       savedFiles["company_card"],
		CompanyCode:       companyCode,
		UserContract:      savedFiles["user_contract"],
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
	bankDetails := r.FormValue("bank_details")
	legalAddress := r.FormValue("legal_address")
	actualAddress := r.FormValue("actual_address")
	contactDetails := r.FormValue("contact_details")
	email := r.FormValue("email")
	// File field mapping to table headers
	fileFieldNames := map[string]string{
		"id_file":       "Удостоверение_личности_или_паспорт.pdf",
		"user_contract": "Подписанный_договор_пользователя.pdf",
	}
	companyCode := r.FormValue("company_code")

	fmt.Println("Full Name:", fullName)
	fmt.Println("IIN:", iin)
	fmt.Println("Email:", email)
	fmt.Println("Legal Address:", legalAddress)
	fmt.Println("Contact Details:", contactDetails)
	fmt.Println("Bank Details:", bankDetails)
	fmt.Println("file:", fileFieldNames["id_file"])

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

	// Create an Individual object
	individual := models.Individual{
		FullName:       fullName,
		IIN:            iin,
		IDFile:         savedFiles["id_file"],
		BankDetails:    bankDetails,
		LegalAddress:   legalAddress,
		ActualAddress:  actualAddress,
		ContactDetails: contactDetails,
		Email:          email,
		CompanyCode:    companyCode,
		UserContract:   savedFiles["user_contract"],
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

// Search TOO by BIN
func (h *TOOHandler) SearchTOOs(w http.ResponseWriter, r *http.Request) {
	bin := r.URL.Query().Get(":bin")
	if bin == "" {
		http.Error(w, "Не указан параметр 'bin'", http.StatusBadRequest)
		return
	}

	toos, err := h.Service.SearchTOOsByBIN(r.Context(), bin)
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

	ips, err := h.Service.SearchIPsByIIN(r.Context(), iin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ips)
}

// Search Individual by IIN

func (h *IndividualHandler) SearchIndividuals(w http.ResponseWriter, r *http.Request) {
	iin := r.URL.Query().Get(":iin")
	if iin == "" {
		http.Error(w, "Не указан параметр 'iin'", http.StatusBadRequest)
		return
	}

	individuals, err := h.Service.SearchIndividualsByIIN(r.Context(), iin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(individuals)
}
