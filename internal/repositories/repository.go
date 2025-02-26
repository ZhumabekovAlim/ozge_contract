package repositories

import (
	"context"
	"database/sql"
	"fmt"
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
func (r *TOORepository) GetTOOsByBIN(ctx context.Context, bin, code string) ([]models.TOO, error) {
	query := `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code, additional_information, user_contract, status, created_at, updated_at
		FROM TOO
		WHERE bin = ? AND company_code LIKE CONCAT('%', ?, '%') AND status = 2 OR status = 3
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, bin, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toos []models.TOO

	for rows.Next() {
		var t models.TOO
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.BIN,
			&t.BankDetails,
			&t.Email,
			&t.Signer,
			&t.IIN,
			&t.CompanyCode,
			&t.AdditionalInformation,
			&t.UserContract,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		toos = append(toos, t)
	}

	return toos, rows.Err()
}

// For IP (search by IIN)
func (r *IPRepository) GetIPsByIIN(ctx context.Context, iin, code string) ([]models.IP, error) {
	query := `
		SELECT id, name, bin, bank_details, email, signer, iin, company_code, additional_information,user_contract, status, created_at, updated_at
		FROM IP
		WHERE iin = ? AND company_code LIKE CONCAT('%', ?, '%') AND status = 2 OR status = 3
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, iin, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []models.IP
	for rows.Next() {
		var ip models.IP
		err = rows.Scan(
			&ip.ID,
			&ip.Name,
			&ip.BIN,
			&ip.BankDetails,
			&ip.Email,
			&ip.Signer,
			&ip.IIN,
			&ip.CompanyCode,
			&ip.AdditionalInformation,
			&ip.UserContract,
			&ip.Status,
			&ip.CreatedAt,
			&ip.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, rows.Err()
}

// For Individual (search by IIN)
func (r *IndividualRepository) GetIndividualsByIIN(ctx context.Context, iin, code string) ([]models.Individual, error) {
	query := `
		SELECT 
			id, 
			full_name, 
			COALESCE(email, '') as email,
			company_code,
			COALESCE(user_contract, '') as user_contract,
			COALESCE(additional_information, '') as additional_information,
			user_contract,
			status,
			created_at,
			updated_at
		FROM Individual
		WHERE iin = ? AND company_code LIKE CONCAT('%', ?, '%') AND status = 2 OR status = 3
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, iin, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var individuals []models.Individual
	for rows.Next() {
		var ind models.Individual
		err = rows.Scan(
			&ind.ID,
			&ind.FullName,
			&ind.IIN,
			&ind.Email,
			&ind.CompanyCode,
			&ind.UserContract,
			&ind.AdditionalInformation,
			&ind.UserContract,
			&ind.Status,
			&ind.CreatedAt,
			&ind.UpdatedAt,
		)
		if err != nil {
			return nil, err
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
