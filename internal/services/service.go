package services

import (
	"context"
	"ozge/internal/models"
	"ozge/internal/repositories"
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

	too.ID = id
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

	ip.ID = id
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

	individual.ID = id
	return individual, nil
}

// In TOOService
func (s *TOOService) SearchTOOsByBIN(ctx context.Context, bin string) ([]models.TOO, error) {
	return s.Repo.GetTOOsByBIN(ctx, bin)
}

// In IPService
func (s *IPService) SearchIPsByIIN(ctx context.Context, iin string) ([]models.IP, error) {
	return s.Repo.GetIPsByIIN(ctx, iin)
}

// In IndividualService
func (s *IndividualService) SearchIndividualsByIIN(ctx context.Context, iin string) ([]models.Individual, error) {
	return s.Repo.GetIndividualsByIIN(ctx, iin)
}
