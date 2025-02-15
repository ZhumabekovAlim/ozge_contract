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
		contact_details, email, egov_file, company_card, company_code, user_contract) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		too.Name, too.BIN, too.RegistrationFile, too.CEOName, too.CEOOrderFile, too.CEOIDFile,
		too.RepresentativePOA, too.RepresentativeID, too.BankDetails, too.LegalAddress,
		too.ActualAddress, too.ContactDetails, too.Email, too.EgovFile, too.CompanyCard, too.CompanyCode, too.UserContract,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// IP Repository
type IPRepository struct {
	Db *sql.DB
}

func (r *IPRepository) CreateIP(ctx context.Context, ip models.IP) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO IP (name, iin, registration_file, representative_poa, representative_id, 
		bank_details, legal_address, actual_address, contact_details, email, company_card, company_code, user_contract) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ip.Name, ip.IIN, ip.RegistrationFile, ip.RepresentativePOA, ip.RepresentativeID,
		ip.BankDetails, ip.LegalAddress, ip.ActualAddress, ip.ContactDetails, ip.Email, ip.CompanyCard, ip.CompanyCode, ip.UserContract,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// Individual Repository
type IndividualRepository struct {
	Db *sql.DB
}

func (r *IndividualRepository) CreateIndividual(ctx context.Context, individual models.Individual) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO Individual (full_name, iin, id_file, bank_details, legal_address, 
		actual_address, contact_details, email, company_code, user_contract) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		individual.FullName, individual.IIN, individual.IDFile, individual.BankDetails,
		individual.LegalAddress, individual.ActualAddress, individual.ContactDetails, individual.Email, individual.CompanyCode, individual.UserContract,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
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
		SELECT id, full_name, iin, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
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
