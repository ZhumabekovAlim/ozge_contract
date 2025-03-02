package models

type TOO struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	BIN                   string   `json:"bin"`
	BankDetails           string   `json:"bank_details"`
	Email                 string   `json:"email"`
	Signer                string   `json:"signer"`
	IIN                   string   `json:"iin"`
	CompanyCode           string   `json:"company_code"`
	UserContract          string   `json:"user_contract,omitempty"`
	Token                 string   `json:"token,omitempty"`
	Status                int      `json:"status"`
	AdditionalInformation string   `json:"additional_information,omitempty"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
	Discard               *Discard `json:"discard,omitempty"`
}

type IP struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	BIN                   string   `json:"bin"`
	BankDetails           string   `json:"bank_details"`
	Email                 string   `json:"email"`
	Signer                string   `json:"signer"`
	IIN                   string   `json:"iin"`
	CompanyCode           string   `json:"company_code"`
	UserContract          string   `json:"user_contract,omitempty"`
	Token                 string   `json:"token,omitempty"`
	Status                int      `json:"status"`
	AdditionalInformation string   `json:"additional_information,omitempty"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
	Discard               *Discard `json:"discard,omitempty"`
}

type Individual struct {
	ID                    int      `json:"id"`
	FullName              string   `json:"full_name"`
	IIN                   string   `json:"iin"`
	Email                 string   `json:"email,omitempty"`
	CompanyCode           string   `json:"company_code"`
	UserContract          string   `json:"user_contract,omitempty"`
	Token                 string   `json:"token,omitempty"`
	Status                int      `json:"status"`
	AdditionalInformation string   `json:"additional_information,omitempty"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
	Discard               *Discard `json:"discard,omitempty"`
}

type Discard struct {
	ID           int    `json:"id"`
	FullName     string `json:"full_name"`
	IIN          string `json:"iin,omitempty"`
	PhoneNumber  string `json:"phone_number"`
	ContractID   int    `json:"contract_id"`
	Reason       string `json:"reason,omitempty"`
	CompanyName  string `json:"company_name"`
	BIN          string `json:"bin,omitempty"`
	Signer       string `json:"signer"`
	ContractPath string `json:"contract_path,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type Company struct {
	ID          uint   `json:"id"`
	CompanyName string `json:"company_name"`
	Password    string `json:"password"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
