package models

type TOO struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	BIN                   string `json:"bin"`
	RegistrationFile      string `json:"registration_file"`
	CEOName               string `json:"ceo_name"`
	CEOOrderFile          string `json:"ceo_order_file"`
	CEOIDFile             string `json:"ceo_id_file"`
	RepresentativePOA     string `json:"representative_poa"`
	RepresentativeID      string `json:"representative_id"`
	BankDetails           string `json:"bank_details"`
	LegalAddress          string `json:"legal_address"`
	ActualAddress         string `json:"actual_address"`
	ContactDetails        string `json:"contact_details"`
	Email                 string `json:"email"`
	EgovFile              string `json:"egov_file"`
	CompanyCard           string `json:"company_card"`
	CompanyCode           string `json:"company_code"`
	UserContract          string `json:"user_contract,omitempty"`
	Token                 string `json:"token,omitempty"`
	Status                int    `json:"status"`
	AdditionalInformation string `json:"additional_information,omitempty"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
}

type IP struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	IIN                   string `json:"iin"`
	RegistrationFile      string `json:"registration_file"`
	RepresentativePOA     string `json:"representative_poa"`
	RepresentativeID      string `json:"representative_id"`
	BankDetails           string `json:"bank_details"`
	LegalAddress          string `json:"legal_address"`
	ActualAddress         string `json:"actual_address"`
	ContactDetails        string `json:"contact_details"`
	Email                 string `json:"email"`
	CompanyCard           string `json:"company_card"`
	CompanyCode           string `json:"company_code"`
	UserContract          string `json:"user_contract,omitempty"`
	Token                 string `json:"token,omitempty"`
	Status                int    `json:"status"`
	AdditionalInformation string `json:"additional_information,omitempty"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
}

type Individual struct {
	ID                    int    `json:"id"`
	FullName              string `json:"full_name"`
	IIN                   string `json:"iin"`
	IDFile                string `json:"id_file,omitempty"`
	BankDetails           string `json:"bank_details,omitempty"`
	LegalAddress          string `json:"legal_address,omitempty"`
	ActualAddress         string `json:"actual_address,omitempty"`
	ContactDetails        string `json:"contact_details,omitempty"`
	Email                 string `json:"email,omitempty"`
	CompanyCode           string `json:"company_code"`
	UserContract          string `json:"user_contract,omitempty"`
	Token                 string `json:"token,omitempty"`
	Status                int    `json:"status"`
	AdditionalInformation string `json:"additional_information,omitempty"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
}
