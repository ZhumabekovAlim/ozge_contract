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

// For TOO (search by BIN)
func (r *TOORepository) GetTOOByBIN(ctx context.Context, bin string) (models.TOO, error) {
	var too models.TOO
	// Only select columns that are not file paths
	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, bin, ceo_name, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
		FROM TOO
		WHERE bin = ?`, bin).Scan(
		&too.ID, &too.Name, &too.BIN, &too.CEOName, &too.BankDetails,
		&too.LegalAddress, &too.ActualAddress, &too.ContactDetails, &too.Email, &too.CompanyCode, &too.CreatedAt, &too.UpdatedAt,
	)
	return too, err
}

// For IP (search by IIN)
func (r *IPRepository) GetIPByIIN(ctx context.Context, iin string) (models.IP, error) {
	var ip models.IP
	err := r.Db.QueryRowContext(ctx, `
		SELECT id, name, iin, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
		FROM IP
		WHERE iin = ?`, iin).Scan(
		&ip.ID, &ip.Name, &ip.IIN, &ip.BankDetails,
		&ip.LegalAddress, &ip.ActualAddress, &ip.ContactDetails, &ip.Email, &ip.CompanyCode, &ip.CreatedAt, &ip.UpdatedAt,
	)
	return ip, err
}

// For Individual (search by IIN)
func (r *IndividualRepository) GetIndividualByIIN(ctx context.Context, iin string) (models.Individual, error) {
	var individual models.Individual
	err := r.Db.QueryRowContext(ctx, `
		SELECT id, full_name, iin, bank_details, legal_address, actual_address, contact_details, email, company_code, created_at, updated_at
		FROM Individual
		WHERE iin = ?`, iin).Scan(
		&individual.ID, &individual.FullName, &individual.IIN, &individual.BankDetails,
		&individual.LegalAddress, &individual.ActualAddress, &individual.ContactDetails, &individual.Email, &individual.CompanyCode, &individual.CreatedAt, &individual.UpdatedAt,
	)
	return individual, err
}
