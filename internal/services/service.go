package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"ozge/internal/models"
	"ozge/internal/repositories"
	"time"
)

// TOO Service
type TOOService struct {
	Repo *repositories.TOORepository
}

func (s *TOOService) CreateTOO(ctx context.Context, too models.TOO) (models.TOO, error) {
	id, err := s.Repo.CreateTOO(ctx, too)
	if err != nil {
		return models.TOO{}, err
	}

	createdAt, err := s.Repo.GetCreatedAt(ctx, id)
	if err != nil {
		return models.TOO{}, err
	}

	strCreatedAt, err := time.Parse("2006-01-02 15:04:05", createdAt)

	token := generateToken(id, strCreatedAt)

	err = s.Repo.UpdateToken(ctx, id, token)
	if err != nil {
		return models.TOO{}, err
	}

	too.ID = id
	too.CreatedAt = createdAt
	too.Token = token
	return too, nil
}

func (s *TOOService) UpdateContractTOO(ctx context.Context, too models.TOO) (models.TOO, error) {
	err := s.Repo.UpdateContractTOO(ctx, too)
	if err != nil {
		return models.TOO{}, err
	}

	return too, nil
}

// IP Service
type IPService struct {
	Repo *repositories.IPRepository
}

func (s *IPService) CreateIP(ctx context.Context, ip models.IP) (models.IP, error) {
	id, err := s.Repo.CreateIP(ctx, ip)
	if err != nil {
		return models.IP{}, err
	}

	createdAt, err := s.Repo.GetCreatedAt(ctx, id)
	if err != nil {
		return models.IP{}, err
	}

	strCreatedAt, err := time.Parse("2006-01-02 15:04:05", createdAt)

	token := generateToken(id, strCreatedAt)

	err = s.Repo.UpdateToken(ctx, id, token)
	if err != nil {
		return models.IP{}, err
	}

	ip.ID = id
	ip.CreatedAt = createdAt
	ip.Token = token
	return ip, nil
}

func (s *IPService) UpdateContractIP(ctx context.Context, ip models.IP) (models.IP, error) {
	err := s.Repo.UpdateContractIP(ctx, ip)
	if err != nil {
		return models.IP{}, err
	}

	return ip, nil
}

// Individual Service
type IndividualService struct {
	Repo *repositories.IndividualRepository
}

func (s *IndividualService) CreateIndividual(ctx context.Context, individual models.Individual) (models.Individual, error) {
	id, err := s.Repo.CreateIndividual(ctx, individual)
	if err != nil {
		return models.Individual{}, err
	}

	createdAt, err := s.Repo.GetCreatedAt(ctx, id)
	if err != nil {
		return models.Individual{}, err
	}

	strCreatedAt, err := time.Parse("2006-01-02 15:04:05", createdAt)

	token := generateToken(id, strCreatedAt)

	err = s.Repo.UpdateToken(ctx, id, token)
	if err != nil {
		return models.Individual{}, err
	}

	individual.ID = id
	individual.CreatedAt = createdAt
	individual.Token = token
	return individual, nil
}

func (s *IndividualService) UpdateContractIndividual(ctx context.Context, individual models.Individual) (models.Individual, error) {
	err := s.Repo.UpdateContractIndividual(ctx, individual)
	if err != nil {
		return models.Individual{}, err
	}

	return individual, nil
}

// In TOOService
func (s *TOOService) SearchTOOsByBIN(ctx context.Context, bin, code string) ([]models.TOO, error) {
	return s.Repo.GetTOOsByBIN(ctx, bin, code)
}

// In IPService
func (s *IPService) SearchIPsByIIN(ctx context.Context, iin, code string) ([]models.IP, error) {
	return s.Repo.GetIPsByIIN(ctx, iin, code)
}

// In IndividualService
func (s *IndividualService) SearchIndividualsByIIN(ctx context.Context, iin, code string) ([]models.Individual, error) {
	return s.Repo.GetIndividualsByIIN(ctx, iin, code)
}

func (s *TOOService) SearchTOOByToken(ctx context.Context, token string) (models.TOO, error) {
	return s.Repo.FindByToken(ctx, token)
}

func (s *IPService) SearchIPByToken(ctx context.Context, token string) (models.IP, error) {
	return s.Repo.FindByToken(ctx, token)
}

func (s *IndividualService) SearchIndividualByToken(ctx context.Context, token string) (models.Individual, error) {
	return s.Repo.FindByToken(ctx, token)
}

// generateToken создает токен на основе ID и времени создания
func generateToken(id int, createdAt time.Time) string {
	data := fmt.Sprintf("%d:%d", id, createdAt.Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
