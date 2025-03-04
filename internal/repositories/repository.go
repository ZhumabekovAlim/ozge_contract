package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
	var companyID int
	var storedHashes []struct {
		ID       int
		Password string
	}

	rows, err := r.Db.QueryContext(ctx, `
		SELECT c.id, c.password 
		FROM companies c
		JOIN TOO t ON CAST(SUBSTRING_INDEX(t.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE t.iin = ? OR ? = 'all'
	`, iin, iin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry struct {
			ID       int
			Password string
		}
		if err := rows.Scan(&entry.ID, &entry.Password); err != nil {
			return nil, err
		}
		storedHashes = append(storedHashes, entry)
	}

	passwordValid := false
	for _, entry := range storedHashes {
		if bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(pass)) == nil {
			passwordValid = true
			companyID = entry.ID
			break
		}
	}

	if !passwordValid {
		return nil, fmt.Errorf("❌ Неверный пароль")
	}

	query := `
	SELECT t.id, COALESCE(t.name, ''), COALESCE(t.bin, ''), COALESCE(t.bank_details, ''), 
	       COALESCE(t.email, ''), COALESCE(t.signer, ''), COALESCE(t.iin, ''), 
	       COALESCE(t.company_code, ''), COALESCE(t.additional_information, ''), 
	       COALESCE(t.user_contract, ''), COALESCE(t.status, 0), t.created_at, t.updated_at,
	       COALESCE(d.id, 0), COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
	       COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
	       COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
	       COALESCE(d.token, ''), COALESCE(d.created_at, NOW()), COALESCE(d.updated_at, NOW())
	FROM TOO t
	JOIN companies c ON CAST(SUBSTRING_INDEX(t.company_code, '.', 1) AS UNSIGNED) = c.id
	LEFT JOIN discard d ON t.id = d.contract_id
	WHERE (t.iin = ? OR ? = 'all') AND c.id = ?
	ORDER BY t.created_at DESC
	`

	rows, err = r.Db.QueryContext(ctx, query, iin, iin, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toos []models.TOO
	for rows.Next() {
		var t models.TOO
		var discard models.Discard
		err = rows.Scan(
			&t.ID, &t.Name, &t.BIN, &t.BankDetails, &t.Email, &t.Signer, &t.IIN, &t.CompanyCode,
			&t.AdditionalInformation, &t.UserContract, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber,
			&discard.ContractID, &discard.Reason, &discard.CompanyName, &discard.BIN,
			&discard.Signer, &discard.ContractPath, &discard.Token, &discard.CreatedAt, &discard.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if discard.ID != 0 {
			t.Discard = &discard
		}

		toos = append(toos, t)
	}

	return toos, rows.Err()
}

// For IP (search by IIN)
func (r *IPRepository) GetIPsByIIN(ctx context.Context, iin, pass string) ([]models.IP, error) {
	var companyID int
	var storedHashes []struct {
		ID       int
		Password string
	}

	// 1. Получаем id компании и хеш пароля
	rows, err := r.Db.QueryContext(ctx, `
		SELECT c.id, c.password 
		FROM companies c
		JOIN IP ip ON CAST(SUBSTRING_INDEX(ip.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE ip.iin = ? OR ? = 'all'
	`, iin, iin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry struct {
			ID       int
			Password string
		}
		if err := rows.Scan(&entry.ID, &entry.Password); err != nil {
			return nil, err
		}
		storedHashes = append(storedHashes, entry)
	}

	// 2. Проверяем введенный пароль и сохраняем ID компании
	passwordValid := false
	for _, entry := range storedHashes {
		if bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(pass)) == nil {
			passwordValid = true
			companyID = entry.ID
			break
		}
	}

	if !passwordValid {
		return nil, fmt.Errorf("❌ Неверный пароль")
	}

	// 3. Запрос на получение данных
	query := `
	SELECT ip.id, COALESCE(ip.name, ''), COALESCE(ip.bin, ''), COALESCE(ip.bank_details, ''), 
	       COALESCE(ip.email, ''), COALESCE(ip.signer, ''), COALESCE(ip.iin, ''), 
	       COALESCE(ip.company_code, ''), COALESCE(ip.additional_information, ''), 
	       COALESCE(ip.user_contract, ''), COALESCE(ip.status, 0), ip.created_at, ip.updated_at,
	       COALESCE(d.id, 0), COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
	       COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
	       COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
	       COALESCE(d.token, ''), COALESCE(d.created_at, NOW()), COALESCE(d.updated_at, NOW())
	FROM IP ip
	JOIN companies c ON CAST(SUBSTRING_INDEX(ip.company_code, '.', 1) AS UNSIGNED) = c.id
	LEFT JOIN discard d ON ip.id = d.contract_id
	 WHERE (ip.iin = ? OR ? = 'all') AND c.id = ?
	ORDER BY ip.created_at DESC
	`

	rows, err = r.Db.QueryContext(ctx, query, iin, iin, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []models.IP
	for rows.Next() {
		var ip models.IP
		var discard models.Discard
		err = rows.Scan(
			&ip.ID, &ip.Name, &ip.BIN, &ip.BankDetails, &ip.Email, &ip.Signer, &ip.IIN, &ip.CompanyCode,
			&ip.AdditionalInformation, &ip.UserContract, &ip.Status, &ip.CreatedAt, &ip.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber,
			&discard.ContractID, &discard.Reason, &discard.CompanyName, &discard.BIN,
			&discard.Signer, &discard.ContractPath, &discard.Token, &discard.CreatedAt, &discard.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if discard.ID != 0 {
			ip.Discard = &discard
		}

		ips = append(ips, ip)
	}

	return ips, rows.Err()
}

// For Individual (search by IIN)
func (r *IndividualRepository) GetIndividualsByIIN(ctx context.Context, iin, pass string) ([]models.Individual, error) {
	var companyID int
	var storedHashes []struct {
		ID       int
		Password string
	}

	// Получаем id компании и хеш пароля
	rows, err := r.Db.QueryContext(ctx, `
		SELECT c.id, c.password 
		FROM companies c
		JOIN Individual ind ON CAST(SUBSTRING_INDEX(ind.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE ind.iin = ? OR ? = 'all'
	`, iin, iin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry struct {
			ID       int
			Password string
		}
		if err := rows.Scan(&entry.ID, &entry.Password); err != nil {
			return nil, err
		}
		storedHashes = append(storedHashes, entry)
	}

	// Проверяем введенный пароль и сохраняем ID компании
	passwordValid := false
	for _, entry := range storedHashes {
		if bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(pass)) == nil {
			passwordValid = true
			companyID = entry.ID
			break
		}
	}

	if !passwordValid {
		return nil, fmt.Errorf("❌ Неверный пароль")
	}

	// Запрос на получение данных
	query := `
	SELECT ind.id, COALESCE(ind.full_name, ''), COALESCE(ind.iin, ''), 
	       COALESCE(ind.email, ''), COALESCE(ind.company_code, ''), 
	       COALESCE(ind.user_contract, ''), COALESCE(ind.additional_information, ''), 
	       COALESCE(ind.status, 0), ind.created_at, ind.updated_at,
	       COALESCE(d.id, 0), COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
	       COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
	       COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
	       COALESCE(d.token, ''), COALESCE(d.created_at, NOW()), COALESCE(d.updated_at, NOW())
	FROM Individual ind
	JOIN companies c ON CAST(SUBSTRING_INDEX(ind.company_code, '.', 1) AS UNSIGNED) = c.id
	LEFT JOIN discard d ON ind.id = d.contract_id
	WHERE (ind.iin = ? OR ? = 'all') AND c.id = ?
	ORDER BY ind.created_at DESC
	`

	rows, err = r.Db.QueryContext(ctx, query, iin, iin, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var individuals []models.Individual
	for rows.Next() {
		var ind models.Individual
		var discard models.Discard
		err = rows.Scan(
			&ind.ID, &ind.FullName, &ind.IIN, &ind.Email, &ind.CompanyCode,
			&ind.UserContract, &ind.AdditionalInformation, &ind.Status, &ind.CreatedAt, &ind.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber,
			&discard.ContractID, &discard.Reason, &discard.CompanyName, &discard.BIN,
			&discard.Signer, &discard.ContractPath, &discard.Token, &discard.CreatedAt, &discard.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if discard.ID != 0 {
			ind.Discard = &discard
		}

		individuals = append(individuals, ind)
	}

	return individuals, rows.Err()
}

type CompanyDataRepo struct {
	Db *sql.DB
}

func (r *CompanyDataRepo) GetAllDataByIIN(ctx context.Context, iin, pass string) ([]interface{}, error) {
	var companyID int
	var storedHashes []struct {
		ID       int
		Password string
	}

	// 1. Получаем id компании и хеши паролей
	rows, err := r.Db.QueryContext(ctx, `
		SELECT c.id, c.password 
		FROM companies c
		JOIN TOO t ON CAST(SUBSTRING_INDEX(t.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE t.iin = ? OR ? = 'all'
		UNION
		SELECT c.id, c.password 
		FROM companies c
		JOIN IP ip ON CAST(SUBSTRING_INDEX(ip.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE ip.iin = ? OR ? = 'all'
		UNION
		SELECT c.id, c.password 
		FROM companies c
		JOIN Individual ind ON CAST(SUBSTRING_INDEX(ind.company_code, '.', 1) AS UNSIGNED) = c.id
		WHERE ind.iin = ? OR ? = 'all'
	`, iin, iin, iin, iin, iin, iin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Сохраняем id компании и хеши паролей
	for rows.Next() {
		var entry struct {
			ID       int
			Password string
		}
		if err := rows.Scan(&entry.ID, &entry.Password); err != nil {
			return nil, err
		}
		storedHashes = append(storedHashes, entry)
	}

	// 3. Проверяем введенный пароль и запоминаем ID компании
	passwordValid := false
	for _, entry := range storedHashes {
		if bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(pass)) == nil {
			passwordValid = true
			companyID = entry.ID // Сохраняем ID компании, у которой подошел пароль
			break
		}
	}

	if !passwordValid {
		return nil, fmt.Errorf("❌ Неверный пароль")
	}

	// 4. Получаем данные из всех таблиц (TOO, IP, Individual) с учетом company_id
	query := `
		SELECT * FROM (
			(SELECT 'TOO' as source, t.id, COALESCE(t.name, ''), COALESCE(t.bin, ''), COALESCE(t.bank_details, ''), 
					COALESCE(t.email, ''), COALESCE(t.signer, ''), COALESCE(t.iin, ''), 
					COALESCE(t.company_code, ''), COALESCE(t.additional_information, ''), 
					COALESCE(t.user_contract, ''), COALESCE(t.status, 0), t.created_at, t.updated_at,
					COALESCE(d.id, 0) as discard_id, COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
					COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
					COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
					COALESCE(d.token, ''), COALESCE(d.created_at, NOW()) AS discard_created_at, COALESCE(d.updated_at, NOW()) AS discard_updated_at
			 FROM TOO t
			 JOIN companies c ON CAST(SUBSTRING_INDEX(t.company_code, '.', 1) AS UNSIGNED) = c.id
			 LEFT JOIN discard d ON t.id = d.contract_id
			 WHERE (t.iin = ? OR ? = 'all') AND c.id = ?)
			 
			UNION ALL
			
			(SELECT 'IP' as source, ip.id, COALESCE(ip.name, ''), COALESCE(ip.bin, ''), COALESCE(ip.bank_details, ''), 
					COALESCE(ip.email, ''), COALESCE(ip.signer, ''), COALESCE(ip.iin, ''), 
					COALESCE(ip.company_code, ''), COALESCE(ip.additional_information, ''), 
					COALESCE(ip.user_contract, ''), COALESCE(ip.status, 0), ip.created_at, ip.updated_at,
					COALESCE(d.id, 0)  as discard_id, COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
					COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
					COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
					COALESCE(d.token, ''), COALESCE(d.created_at, NOW()) AS discard_created_at, COALESCE(d.updated_at, NOW()) AS discard_updated_at
			 FROM IP ip
			 JOIN companies c ON CAST(SUBSTRING_INDEX(ip.company_code, '.', 1) AS UNSIGNED) = c.id
			 LEFT JOIN discard d ON ip.id = d.contract_id
			 WHERE (ip.iin = ? OR ? = 'all') AND c.id = ?)
		
			UNION ALL
			
			(SELECT 'Individual' as source, ind.id, COALESCE(ind.full_name, ''), '' AS bin, '' AS bank_details,
					COALESCE(ind.email, ''), '' AS signer, COALESCE(ind.iin, ''), 
					COALESCE(ind.company_code, ''), COALESCE(ind.additional_information, ''), 
					COALESCE(ind.user_contract, ''), COALESCE(ind.status, 0), ind.created_at, ind.updated_at,
					COALESCE(d.id, 0) as discard_id, COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
					COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
					COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
					COALESCE(d.token, ''), COALESCE(d.created_at, NOW()) AS discard_created_at, COALESCE(d.updated_at, NOW()) AS discard_updated_at
			 FROM Individual ind
			 JOIN companies c ON CAST(SUBSTRING_INDEX(ind.company_code, '.', 1) AS UNSIGNED) = c.id
			 LEFT JOIN discard d ON ind.id = d.contract_id
			 WHERE (ind.iin = ? OR ? = 'all') AND c.id = ?)
		) AS combined
		ORDER BY created_at DESC;
	`

	rows, err = r.Db.QueryContext(ctx, query, iin, iin, companyID, iin, iin, companyID, iin, iin, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []interface{}

	for rows.Next() {
		var record struct {
			Source                string
			ID                    int
			Name                  string
			BIN                   string
			BankDetails           string
			Email                 string
			Signer                string
			IIN                   string
			CompanyCode           string
			AdditionalInformation string
			UserContract          string
			Status                int
			CreatedAt             string
			UpdatedAt             string
			Discard               *models.Discard
		}

		var discard models.Discard
		err = rows.Scan(
			&record.Source, &record.ID, &record.Name, &record.BIN, &record.BankDetails,
			&record.Email, &record.Signer, &record.IIN, &record.CompanyCode,
			&record.AdditionalInformation, &record.UserContract, &record.Status,
			&record.CreatedAt, &record.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber,
			&discard.ContractID, &discard.Reason, &discard.CompanyName,
			&discard.BIN, &discard.Signer, &discard.ContractPath,
			&discard.Token, &discard.CreatedAt, &discard.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if discard.ID != 0 {
			record.Discard = &discard
		}

		results = append(results, record)
	}

	return results, rows.Err()
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
	var discard models.Discard

	err := r.Db.QueryRowContext(ctx, `
		SELECT 
		COALESCE(t.id, 0), COALESCE(t.name, ''), COALESCE(t.bin, ''), COALESCE(t.bank_details, ''), 
		COALESCE(t.email, ''), COALESCE(t.signer, ''), COALESCE(t.iin, ''), COALESCE(t.company_code, ''), 
		COALESCE(t.user_contract, '') as user_contract, COALESCE(t.additional_information, '') as additional_information, 
		COALESCE(t.status, 0),     COALESCE(t.contract_name, '') AS ip_contract_name, COALESCE(t.created_at, ''), COALESCE(t.updated_at, ''),
		COALESCE(d.id, 0), COALESCE(d.full_name, ''), COALESCE(d.iin, ''), COALESCE(d.phone_number, ''), 
		COALESCE(d.contract_id, 0), COALESCE(d.reason, ''), COALESCE(d.company_name, ''), 
		COALESCE(d.bin, ''), COALESCE(d.signer, ''), COALESCE(d.contract_path, ''), 
		COALESCE(d.created_at, ''), COALESCE(d.updated_at, '')
	FROM TOO t
	LEFT JOIN discard d ON t.id = d.contract_id
	WHERE t.token = ? OR d.token = ?`, token, token).
		Scan(
			&too.ID, &too.Name, &too.BIN, &too.BankDetails, &too.Email, &too.Signer, &too.IIN,
			&too.CompanyCode, &too.UserContract, &too.AdditionalInformation, &too.Status, &too.ContractName,
			&too.CreatedAt, &too.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber, &discard.ContractID,
			&discard.Reason, &discard.CompanyName, &discard.BIN, &discard.Signer, &discard.ContractPath,
			&discard.CreatedAt, &discard.UpdatedAt,
		)

	if err != nil {
		return models.TOO{}, err
	}

	// Если у discard есть данные, добавляем их в структуру
	if discard.ID != 0 {
		too.Discard = &discard
	}

	return too, nil
}

func (r *IPRepository) FindByToken(ctx context.Context, token string) (models.IP, error) {
	var ip models.IP
	var discard models.Discard

	err := r.Db.QueryRowContext(ctx, `
		SELECT 
    COALESCE(ip.id, 0) AS ip_id, 
    COALESCE(ip.name, '') AS ip_name, 
    COALESCE(ip.bin, '') AS ip_bin, 
    COALESCE(ip.bank_details, '') AS ip_bank_details, 
    COALESCE(ip.email, '') AS ip_email, 
    COALESCE(ip.signer, '') AS ip_signer, 
    COALESCE(ip.iin, '') AS ip_iin, 
    COALESCE(ip.company_code, '') AS ip_company_code, 
    COALESCE(ip.user_contract, '') AS ip_user_contract, 
    COALESCE(ip.additional_information, '') AS ip_additional_information, 
    COALESCE(ip.status, 0) AS ip_status, 
    COALESCE(ip.contract_name, '') AS ip_contract_name,
    COALESCE(ip.created_at, '') AS ip_created_at, 
    COALESCE(ip.updated_at, '') AS ip_updated_at,
    COALESCE(d.id, 0) AS discard_id, 
    COALESCE(d.full_name, '') AS discard_full_name, 
    COALESCE(d.iin, '') AS discard_iin, 
    COALESCE(d.phone_number, '') AS discard_phone_number, 
    COALESCE(d.contract_id, 0) AS discard_contract_id, 
    COALESCE(d.reason, '') AS discard_reason, 
    COALESCE(d.company_name, '') AS discard_company_name, 
    COALESCE(d.bin, '') AS discard_bin, 
    COALESCE(d.signer, '') AS discard_signer, 
    COALESCE(d.contract_path, '') AS discard_contract_path, 
    COALESCE(d.created_at, '') AS discard_created_at, 
    COALESCE(d.updated_at, '') AS discard_updated_at
FROM IP ip
LEFT JOIN discard d ON ip.id = d.contract_id
WHERE ip.token = ? OR d.token = ?`, token, token).
		Scan(
			&ip.ID, &ip.Name, &ip.BIN, &ip.BankDetails, &ip.Email, &ip.Signer, &ip.IIN,
			&ip.CompanyCode, &ip.UserContract, &ip.AdditionalInformation, &ip.Status, &ip.ContractName,
			&ip.CreatedAt, &ip.UpdatedAt,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber, &discard.ContractID,
			&discard.Reason, &discard.CompanyName, &discard.BIN, &discard.Signer, &discard.ContractPath,
			&discard.CreatedAt, &discard.UpdatedAt,
		)

	if err != nil {
		return models.IP{}, err
	}

	// Если у discard есть данные, добавляем их в структуру
	if discard.ID != 0 {
		ip.Discard = &discard
	}

	return ip, nil
}

func (r *IndividualRepository) FindByToken(ctx context.Context, token string) (models.Individual, error) {
	var individual models.Individual
	var discard models.Discard

	err := r.Db.QueryRowContext(ctx, `
		SELECT 
    COALESCE(ind.id, 0) AS ind_id, 
    COALESCE(ind.full_name, '') AS ind_full_name, 
    COALESCE(ind.iin, '') AS ind_iin, 
    COALESCE(ind.email, '') AS ind_email, 
    COALESCE(ind.company_code, '') AS ind_company_code, 
    COALESCE(ind.user_contract, '') AS ind_user_contract, 
    COALESCE(ind.additional_information, '') AS ind_additional_information, 
    COALESCE(ind.status, 0) AS ind_status, 
    COALESCE(ind.created_at, '') AS ind_created_at, 
    COALESCE(ind.updated_at, '') AS ind_updated_at, 
    COALESCE(ind.contract_name, '') AS ind_contract_name,
    COALESCE(d.id, 0) AS discard_id, 
    COALESCE(d.full_name, '') AS discard_full_name, 
    COALESCE(d.iin, '') AS discard_iin, 
    COALESCE(d.phone_number, '') AS discard_phone_number, 
    COALESCE(d.contract_id, 0) AS discard_contract_id, 
    COALESCE(d.reason, '') AS discard_reason, 
    COALESCE(d.company_name, '') AS discard_company_name, 
    COALESCE(d.bin, '') AS discard_bin, 
    COALESCE(d.signer, '') AS discard_signer, 
    COALESCE(d.contract_path, '') AS discard_contract_path, 
    COALESCE(d.created_at, '') AS discard_created_at, 
    COALESCE(d.updated_at, '') AS discard_updated_at
FROM Individual ind
LEFT JOIN discard d ON ind.id = d.contract_id
WHERE ind.token = ? OR d.token = ?`, token, token).
		Scan(
			&individual.ID, &individual.FullName, &individual.IIN, &individual.Email, &individual.CompanyCode,
			&individual.UserContract, &individual.AdditionalInformation, &individual.Status,
			&individual.CreatedAt, &individual.UpdatedAt, &individual.ContractName,
			&discard.ID, &discard.FullName, &discard.IIN, &discard.PhoneNumber, &discard.ContractID,
			&discard.Reason, &discard.CompanyName, &discard.BIN, &discard.Signer, &discard.ContractPath,
			&discard.CreatedAt, &discard.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Individual{}, fmt.Errorf("Ind not found")
	} else if err != nil {
		return models.Individual{}, err
	}

	if err != nil {
		return models.Individual{}, err
	}

	fmt.Println("indiv: ", individual)

	// Если у discard есть данные, добавляем их в структуру
	if discard.ID != 0 {
		individual.Discard = &discard
	}

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
	COALESCE(id, '') AS id, 
    COALESCE(full_name, '') AS full_name, 
    COALESCE(iin, '') AS iin, 
    COALESCE(email, '') AS email, 
    COALESCE(company_code, '') AS company_code, 
   COALESCE(user_contract, '') AS user_contract,
    COALESCE(additional_information, '') AS additional_information, 
    COALESCE(token, '') AS token, 
    COALESCE(status, 0) AS status,
	COALESCE(contract_name, '') AS contract_name,
    COALESCE(created_at, '') AS created_at, 
    COALESCE(updated_at, '') AS updated_at
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
			&individual.ContractName,
			&individual.CreatedAt,
			&individual.UpdatedAt)

	if err != nil {
		return models.Individual{}, err
	}

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

// CreateDiscard создаёт новую запись в Discard
func (r *DiscardRepository) CreateDiscard(ctx context.Context, discard models.Discard) (int, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO discard (full_name, iin, phone_number, contract_id, reason, company_name, bin, signer,contract_path) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?,'')`,
		discard.FullName, discard.IIN, discard.PhoneNumber, discard.ContractID,
		discard.Reason, discard.CompanyName, discard.BIN, discard.Signer,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// UpdateContractPath обновляет путь к контракту в Discard
func (r *DiscardRepository) UpdateContractPath(ctx context.Context, discard models.Discard) error {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE discard SET contract_path = ?, updated_at = NOW() WHERE id = ?`,
		discard.ContractPath, discard.ID,
	)
	return err
}

func (r *CompanyRepository) FindPasswordByID(ctx context.Context, id string) (string, error) {
	var hashedPassword string
	err := r.Db.QueryRowContext(ctx, `SELECT password FROM companies WHERE id = ?`, id).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Компания не найдена")
		}
		return "", err
	}
	return hashedPassword, nil
}

func (r *DiscardRepository) GetCreatedAt(ctx context.Context, id int) (string, error) {
	var created string
	err := r.Db.QueryRowContext(ctx, "SELECT created_at FROM discard WHERE id = ?", id).Scan(&created)
	return created, err
}

func (r *DiscardRepository) UpdateToken(ctx context.Context, id int, token string) error {
	_, err := r.Db.ExecContext(ctx, `UPDATE discard SET token = ? WHERE id = ?`, token, id)
	return err
}
