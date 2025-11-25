// service/user_service.go
package service

import (
	"errors"

	"github.com/Kars07/rest-grpc-demo/models"
	"github.com/Kars07/rest-grpc-demo/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	existing, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id int64, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
