package repositories

import (
	"context"
	"database/sql"
	"ozge/internal/models"
)

// TOO Repository
type TOORepository struct {
	Db *sql.DB
}

func (r *TOORepository) CreateTOO(ctx context.Context, too models.TOO) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO TOO (name, bin, registration_file, ceo_name, ceo_order_file, ceo_id_file, 
		representative_poa, representative_id, bank_details, legal_address, actual_address, 
		contact_details, email, egov_file, company_card, company_code) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		too.Name, too.BIN, too.RegistrationFile, too.CEOName, too.CEOOrderFile, too.CEOIDFile,
		too.RepresentativePOA, too.RepresentativeID, too.BankDetails, too.LegalAddress,
		too.ActualAddress, too.ContactDetails, too.Email, too.EgovFile, too.CompanyCard, too.CompanyCode,
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
		SET user_contract = ?
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
		INSERT INTO IP (name, iin, registration_file, representative_poa, representative_id, 
		bank_details, legal_address, actual_address, contact_details, email, company_card, company_code) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ip.Name, ip.IIN, ip.RegistrationFile, ip.RepresentativePOA, ip.RepresentativeID,
		ip.BankDetails, ip.LegalAddress, ip.ActualAddress, ip.ContactDetails, ip.Email, ip.CompanyCard, ip.CompanyCode,
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
		SET user_contract = ?
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
		INSERT INTO Individual (full_name, iin, id_file, bank_details, legal_address, 
		actual_address, contact_details, email, company_code) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		individual.FullName, individual.IIN, individual.IDFile, individual.BankDetails,
		individual.LegalAddress, individual.ActualAddress, individual.ContactDetails, individual.Email, individual.CompanyCode,
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
		SET user_contract = ?
		WHERE id = ?`,
		individual.UserContract, individual.ID,
	)
	if err != nil {
		return err
	}

	return err
}

// For TOO (search by BIN)
func (r *TOORepository) GetTOOsByBIN(ctx context.Context, bin string) ([]models.TOO, error) {
	query := `
		SELECT id, name, bin, ceo_name, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
		FROM TOO
		WHERE bin = ?
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, bin)
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
			&t.CEOName,
			&t.BankDetails,
			&t.LegalAddress,
			&t.ActualAddress,
			&t.ContactDetails,
			&t.Email,
			&t.CompanyCode,
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
func (r *IPRepository) GetIPsByIIN(ctx context.Context, iin string) ([]models.IP, error) {
	query := `
		SELECT id, name, iin, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
		FROM IP
		WHERE iin = ?
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, iin)
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
			&ip.IIN,
			&ip.BankDetails,
			&ip.LegalAddress,
			&ip.ActualAddress,
			&ip.ContactDetails,
			&ip.Email,
			&ip.CompanyCode,
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
func (r *IndividualRepository) GetIndividualsByIIN(ctx context.Context, iin string) ([]models.Individual, error) {
	query := `
		SELECT id, full_name, iin, bank_details, legal_address, actual_address, contact_details, email, company_code, user_contract, created_at, updated_at
		FROM Individual
		WHERE iin = ?
		ORDER BY created_at DESC
	`
	rows, err := r.Db.QueryContext(ctx, query, iin)
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
			&ind.BankDetails,
			&ind.LegalAddress,
			&ind.ActualAddress,
			&ind.ContactDetails,
			&ind.Email,
			&ind.CompanyCode,
			&ind.UserContract,
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
		SELECT id, name, bin, ceo_name, bank_details, legal_address, actual_address, 
		contact_details, email, company_code, user_contract, created_at, updated_at
		FROM TOO WHERE token = ?`, token).
		Scan(&too.ID, &too.Name, &too.BIN, &too.CEOName, &too.BankDetails, &too.LegalAddress,
			&too.ActualAddress, &too.ContactDetails, &too.Email, &too.CompanyCode,
			&too.CreatedAt, &too.UpdatedAt)

	if err != nil {
		return models.TOO{}, err
	}

	return too, nil
}

func (r *IPRepository) FindByToken(ctx context.Context, token string) (models.IP, error) {
	var ip models.IP

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, iin, bank_details, legal_address, actual_address, contact_details, 
		email, company_code, user_contract, created_at, updated_at
		FROM IP WHERE token = ?`, token).
		Scan(&ip.ID, &ip.Name, &ip.IIN, &ip.BankDetails, &ip.LegalAddress, &ip.ActualAddress,
			&ip.ContactDetails, &ip.Email, &ip.CompanyCode, &ip.CreatedAt, &ip.UpdatedAt)

	if err != nil {
		return models.IP{}, err
	}

	return ip, nil
}

func (r *IndividualRepository) FindByToken(ctx context.Context, token string) (models.Individual, error) {
	var individual models.Individual

	err := r.Db.QueryRowContext(ctx, `
		SELECT id, full_name, iin, bank_details, legal_address, actual_address, contact_details, 
		email, company_code, user_contract, created_at, updated_at
		FROM Individual WHERE token = ?`, token).
		Scan(&individual.ID, &individual.FullName, &individual.IIN, &individual.BankDetails,
			&individual.LegalAddress, &individual.ActualAddress, &individual.ContactDetails,
			&individual.Email, &individual.CompanyCode, &individual.UserContract,
			&individual.CreatedAt, &individual.UpdatedAt)

	if err != nil {
		return models.Individual{}, err
	}

	return individual, nil
}
