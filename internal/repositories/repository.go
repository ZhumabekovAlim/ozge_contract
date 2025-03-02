package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"ozge/internal/models"
)

// TOO Repository
type TOORepository struct {
	Db *sql.DB
}

func (r *TOORepository) CreateTOO(ctx context.Context, too models.TOO) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO TOO (name, bin, bank_details, email, signer, iin, company_code, additional_information) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		too.Name, too.BIN, too.BankDetails, too.Email, too.Signer, too.IIN, too.CompanyCode, too.AdditionalInformation,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *TOORepository) GetCreatedAt(ctx context.Context, id int) (string, error) {
	var created string
	err := r.Db.QueryRowContext(ctx, "SELECT created_at FROM TOO WHERE id = ?", id).Scan(&created)
	return created, err
}

func (r *TOORepository) UpdateContractTOO(ctx context.Context, too models.TOO) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE TOO
	SET user_contract = ? , status = 2 
		WHERE id = ?`,
		too.UserContract, too.ID,
	)
	if err != nil {
		return err
	}

	return err
}

// IP Repository
type IPRepository struct {
	Db *sql.DB
}

func (r *IPRepository) CreateIP(ctx context.Context, ip models.IP) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO IP (name, bin, bank_details, email, signer, iin, company_code, additional_information) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		ip.Name, ip.BIN, ip.BankDetails, ip.Email, ip.Signer, ip.IIN, ip.CompanyCode, ip.AdditionalInformation,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *IPRepository) GetCreatedAt(ctx context.Context, id int) (string, error) {
	var created string
	err := r.Db.QueryRowContext(ctx, "SELECT created_at FROM IP WHERE id = ?", id).Scan(&created)
	return created, err
}

func (r *IPRepository) UpdateContractIP(ctx context.Context, ip models.IP) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE IP
	SET user_contract = ? , status = 2 
		WHERE id = ?`,
		ip.UserContract, ip.ID,
	)
	if err != nil {
		return err
	}

	return err
}

// Individual Repository
type IndividualRepository struct {
	Db *sql.DB
}

func (r *IndividualRepository) CreateIndividual(ctx context.Context, individual models.Individual) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO Individual (full_name, iin, email, company_code, additional_information) 
		VALUES (?, ?, ?, ?, ?)`,
		individual.FullName, individual.IIN, individual.Email, individual.CompanyCode, individual.AdditionalInformation,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

func (r *IndividualRepository) GetCreatedAt(ctx context.Context, id int) (string, error) {
	var created string
	err := r.Db.QueryRowContext(ctx, "SELECT created_at FROM Individual WHERE id = ?", id).Scan(&created)
	return created, err
}

func (r *IndividualRepository) UpdateContractIndividual(ctx context.Context, individual models.Individual) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE Individual
		SET user_contract = ? , status = 2  
		WHERE id = ?`,
		individual.UserContract, individual.ID,
	)
	if err != nil {
		return err
	}

	return err
}

// For TOO (search by BIN)
func (r *TOORepository) GetTOOsByBIN(ctx context.Context, iin, pass string) ([]models.TOO, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		log.Fatal(err)
	}

	query := `
SELECT t.id,
       t.name,
       t.bin,
       t.bank_details,
       t.email,
       t.signer,
       t.iin,
       t.company_code,
       t.additional_information,
       t.user_contract,
       t.status,
       t.created_at,
       t.updated_at,
       d.id,
       d.full_name,
       d.iin,
       d.phone_number,
       d.contract_id,
       d.reason,
       d.company_name,
       d.bin,
       d.signer,
       d.contract_path,
       d.created_at,
       d.updated_at
FROM TOO t
         LEFT JOIN discard d ON t.id = d.contract_id
         JOIN companies c ON CAST(SUBSTRING_INDEX(t.company_code, '.', 1) AS UNSIGNED) = c.id
WHERE t.iin = ? AND c.password = ?
ORDER BY t.created_at DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, iin, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toos []models.TOO

	for rows.Next() {
		var t models.TOO
		var d models.Discard
		err = rows.Scan(
			&t.ID, &t.Name, &t.BIN, &t.BankDetails, &t.Email, &t.Signer, &t.IIN, &t.CompanyCode,
			&t.AdditionalInformation, &t.UserContract, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&d.ID, &d.FullName, &d.IIN, &d.PhoneNumber, &d.ContractID, &d.Reason, &d.CompanyName,
			&d.BIN, &d.Signer, &d.ContractPath, &d.CreatedAt, &d.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Проверяем, есть ли данные в `discard`
		if d.ID != 0 {
			t.Discard = &d
		} else {
			t.Discard = nil
		}

		toos = append(toos, t)
	}

	return toos, rows.Err()
}

// For IP (search by IIN)
func (r *IPRepository) GetIPsByIIN(ctx context.Context, iin, pass string) ([]models.IP, error) {
	// Хэшируем пароль перед выполнением запроса
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	SELECT 
	    ip.id, ip.name, ip.bin, ip.bank_details, ip.email, ip.signer, ip.iin, ip.company_code, ip.additional_information, 
	    ip.user_contract, ip.status, ip.created_at, ip.updated_at,
	    d.id, d.full_name, d.iin, d.phone_number, d.contract_id, d.reason, d.company_name, d.bin, 
	    d.signer, d.contract_path, d.created_at, d.updated_at
	FROM IP ip
	LEFT JOIN discard d ON ip.id = d.contract_id
	JOIN companies c ON CAST(SUBSTRING_INDEX(ip.company_code, '.', 1) AS UNSIGNED) = c.id
	WHERE ip.iin = ? AND c.password = ?
	ORDER BY ip.created_at DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, iin, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []models.IP

	for rows.Next() {
		var ip models.IP
		var d models.Discard
		err = rows.Scan(
			&ip.ID, &ip.Name, &ip.BIN, &ip.BankDetails, &ip.Email, &ip.Signer, &ip.IIN, &ip.CompanyCode,
			&ip.AdditionalInformation, &ip.UserContract, &ip.Status, &ip.CreatedAt, &ip.UpdatedAt,
			&d.ID, &d.FullName, &d.IIN, &d.PhoneNumber, &d.ContractID, &d.Reason, &d.CompanyName,
			&d.BIN, &d.Signer, &d.ContractPath, &d.CreatedAt, &d.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Проверяем, есть ли данные в `discard`
		if d.ID != 0 {
			ip.Discard = &d
		} else {
			ip.Discard = nil
		}

		ips = append(ips, ip)
	}

	return ips, rows.Err()
}

// For Individual (search by IIN)
func (r *IndividualRepository) GetIndividualsByIIN(ctx context.Context, iin, pass string) ([]models.Individual, error) {
	// Хэшируем пароль перед выполнением запроса
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	SELECT 
	    ind.id, ind.full_name, ind.iin, COALESCE(ind.email, '') as email, ind.company_code, 
	    COALESCE(ind.user_contract, '') as user_contract, COALESCE(ind.additional_information, '') as additional_information, 
	    ind.status, ind.created_at, ind.updated_at,
	    d.id, d.full_name, d.iin, d.phone_number, d.contract_id, d.reason, d.company_name, d.bin, 
	    d.signer, d.contract_path, d.created_at, d.updated_at
	FROM Individual ind
	LEFT JOIN discard d ON ind.id = d.contract_id
	JOIN companies c ON CAST(SUBSTRING_INDEX(ind.company_code, '.', 1) AS UNSIGNED) = c.id
	WHERE ind.iin = ? AND c.password = ?
	ORDER BY ind.created_at DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, iin, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var individuals []models.Individual

	for rows.Next() {
		var ind models.Individual
		var d models.Discard
		err = rows.Scan(
			&ind.ID, &ind.FullName, &ind.IIN, &ind.Email, &ind.CompanyCode,
			&ind.UserContract, &ind.AdditionalInformation, &ind.Status, &ind.CreatedAt, &ind.UpdatedAt,
			&d.ID, &d.FullName, &d.IIN, &d.PhoneNumber, &d.ContractID, &d.Reason, &d.CompanyName,
			&d.BIN, &d.Signer, &d.ContractPath, &d.CreatedAt, &d.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Проверяем, есть ли данные в `discard`
		if d.ID != 0 {
			ind.Discard = &d
		} else {
			ind.Discard = nil
		}

		individuals = append(individuals, ind)
	}

	return individuals, rows.Err()
}

func (r *TOORepository) UpdateToken(ctx context.Context, id int, token string) error {
	_, err := r.Db.ExecContext(ctx, `UPDATE TOO SET token = ? WHERE id = ?`, token, id)
	return err
}

func (r *IPRepository) UpdateToken(ctx context.Context, id int, token string) error {
	_, err := r.Db.ExecContext(ctx, `UPDATE IP SET token = ? WHERE id = ?`, token, id)
	return err
}

func (r *IndividualRepository) UpdateToken(ctx context.Context, id int, token string) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE Individual SET token = ? WHERE id = ?`,
		token, id,
	)
	return err
}

func (r *TOORepository) FindByToken(ctx context.Context, token string) (models.TOO, error) {
	var too models.TOO

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code,   COALESCE(user_contract, '') as user_contract, COALESCE(additional_information, '') as additional_information, status, created_at, updated_at
		FROM TOO WHERE token = ?`, token).
		Scan(&too.ID, &too.Name, &too.BIN, &too.BankDetails, &too.Email, &too.Signer, &too.IIN, &too.CompanyCode, &too.UserContract, &too.AdditionalInformation, &too.Status, &too.CreatedAt, &too.UpdatedAt)

	if err != nil {
		return models.TOO{}, err
	}

	return too, nil
}

func (r *IPRepository) FindByToken(ctx context.Context, token string) (models.IP, error) {
	var ip models.IP

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code,   COALESCE(user_contract, '') as user_contract, COALESCE(additional_information, '') as additional_information, status, created_at, updated_at
		FROM IP WHERE token = ?`, token).
		Scan(&ip.ID, &ip.Name, &ip.BIN, &ip.BankDetails, &ip.Email, &ip.Signer, &ip.IIN, &ip.CompanyCode, &ip.UserContract, &ip.AdditionalInformation, &ip.Status, &ip.CreatedAt, &ip.UpdatedAt)

	if err != nil {
		return models.IP{}, err
	}

	return ip, nil
}

func (r *IndividualRepository) FindByToken(ctx context.Context, token string) (models.Individual, error) {
	var individual models.Individual

	fmt.Println(token)

	err := r.Db.QueryRowContext(ctx, `
		SELECT 
			id, 
			full_name, 
			iin, 
			COALESCE(email, '') as email,
			company_code,
			COALESCE(user_contract, '') as user_contract,
			COALESCE(additional_information, '') as additional_information,
			status,
			created_at,
			updated_at
		FROM Individual WHERE token = ?`, token).
		Scan(&individual.ID,
			&individual.FullName,
			&individual.IIN,
			&individual.Email,
			&individual.CompanyCode,
			&individual.UserContract,
			&individual.AdditionalInformation,
			&individual.Status,
			&individual.CreatedAt,
			&individual.UpdatedAt)

	if err != nil {
		return models.Individual{}, err
	}

	fmt.Println(individual)

	return individual, nil
}

func (r *TOORepository) FindByID(ctx context.Context, id string) (models.TOO, error) {
	var too models.TOO

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code,   COALESCE(user_contract, '') as user_contract, COALESCE(additional_information, '') as additional_information, token, status, created_at, updated_at
		FROM TOO WHERE id = ?`, id).
		Scan(&too.ID, &too.Name, &too.BIN, &too.BankDetails, &too.Email, &too.Signer, &too.IIN, &too.CompanyCode, &too.UserContract, &too.AdditionalInformation, &too.Token, &too.Status, &too.CreatedAt, &too.UpdatedAt)

	if err != nil {
		return models.TOO{}, err
	}

	return too, nil
}

func (r *IPRepository) FindByID(ctx context.Context, id string) (models.IP, error) {
	var ip models.IP

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code,   COALESCE(user_contract, '') as user_contract, COALESCE(additional_information, '') as additional_information, token, status, created_at, updated_at
		FROM IP WHERE id = ?`, id).
		Scan(&ip.ID, &ip.Name, &ip.BIN, &ip.BankDetails, &ip.Email, &ip.Signer, &ip.IIN, &ip.CompanyCode, &ip.UserContract, &ip.AdditionalInformation, &ip.Token, &ip.Status, &ip.CreatedAt, &ip.UpdatedAt)

	if err != nil {
		return models.IP{}, err
	}

	return ip, nil
}

func (r *IndividualRepository) FindByID(ctx context.Context, id string) (models.Individual, error) {
	var individual models.Individual

	err := r.Db.QueryRowContext(ctx, `
		SELECT 
			id, 
			full_name, 
			iin, 
			COALESCE(email, '') as email,
			company_code,
			COALESCE(user_contract, '') as user_contract,
			COALESCE(additional_information, '') as additional_information,
			token,
			status,
			created_at,
			updated_at
		FROM Individual WHERE id = ?`, id).
		Scan(&individual.ID,
			&individual.FullName,
			&individual.IIN,
			&individual.Email,
			&individual.CompanyCode,
			&individual.UserContract,
			&individual.AdditionalInformation,
			&individual.Token,
			&individual.Status,
			&individual.CreatedAt,
			&individual.UpdatedAt)

	if err != nil {
		return models.Individual{}, err
	}

	fmt.Println(individual)

	return individual, nil
}

func (r *TOORepository) UpdateUserContractStatus(ctx context.Context, id string) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE TOO
			SET  status = 3  
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *IPRepository) UpdateUserContractStatus(ctx context.Context, id string) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE IP
			SET  status = 3  
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	return err
}

func (r *IndividualRepository) UpdateUserContractStatus(ctx context.Context, id string) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE Individual
			SET  status = 3  
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	return err
}

type DiscardRepository struct {
	Db *sql.DB
}

type CompanyRepository struct {
	Db *sql.DB
}

// Создание компании с сохранением пароля в хэшированном виде
func (r *CompanyRepository) Create(ctx context.Context, company models.Company) (uint, error) {
	query := `
	INSERT INTO companies (company_name, password) 
	VALUES (?, ?)
	RETURNING id`

	var id uint
	err := r.Db.QueryRowContext(ctx, query, company.CompanyName, company.Password).Scan(&id)
	return id, err
}
