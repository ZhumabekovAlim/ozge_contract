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
	ContractName          string   `json:"contract_name,omitempty"`
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
	ContractName          string   `json:"contract_name,omitempty"`
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
	ContractName          string   `json:"contract_name,omitempty"`
	CreatedAt             string   `json:"created_at,omitempty"`
	UpdatedAt             string   `json:"updated_at,omitempty"`
	Discard               *Discard `json:"discard,omitempty"`
}

type Discard struct {
	ID           int    `json:"id,omitempty"`
	FullName     string `json:"full_name,omitempty"`
	IIN          string `json:"iin,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	ContractID   int    `json:"contract_id,omitempty"`
	Reason       string `json:"reason,omitempty"`
	CompanyName  string `json:"company_name,omitempty"`
	BIN          string `json:"bin,omitempty"`
	Signer       string `json:"signer,omitempty"`
	ContractPath string `json:"contract_path,omitempty"`
	Token        string `json:"token,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type Company struct {
	ID          uint   `json:"id"`
	CompanyName string `json:"company_name"`
	Password    string `json:"password"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
