package application

import (
	"errors"
	"log"

	"blog/internal/domain"
	"blog/pkg/ddd"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo        domain.UserRepository
	eventDispatcher ddd.EventDispatcher
}

func NewUserService(
	userRepo domain.UserRepository,
	eventDispatcher ddd.EventDispatcher,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		eventDispatcher: eventDispatcher,
	}
}

func (s *UserService) CreateUser(
	email, password, username string,
	userRoles []string,
) (*UserDTO, error) {
	// First things first, make sure the email / username aren't already being used
	if exists, err := s.userRepo.EmailExists(email); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("email already in use")
	}

	if exists, err := s.userRepo.UsernameExists(username); !exists || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("username already in use")
	}

	// Hash the password before passing it into the domain
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	domainUserRoles := []domain.UserRole{}
	for _, role := range userRoles {
		domainUserRoles = append(domainUserRoles, domain.UserRole(role))
	}

	// Create the user
	user, err := domain.NewUser(email, username, string(passwordHash), "", domainUserRoles)
	if err != nil {
		return nil, err
	}

	// Persist
	if _, err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(user); err != nil {
		return nil, err
	}

	userDTO := UserDTO{}
	userDTO.FromDomain(user)

	return &userDTO, nil
}

func (s *UserService) SetUserRoles(userID string, userRoles []string) error {
	return nil
}

func (s *UserService) UpdateDescription(userID, description string) error {
	return nil
}

func (s *UserService) UpdatePassword(userID, password string) error {
	return nil
}

// Helper method to dispatch events for any aggregate with AggregateBase
func (s *UserService) dispatchAggregateEvents(aggregate ddd.EventAggregate) error {
	events := aggregate.GetUncommittedEvents()
	for _, event := range events {
		if err := s.eventDispatcher.Dispatch(event); err != nil {
			log.Printf("Failed to dispatch event: %v", err)
		}
	}
	aggregate.MarkEventsAsCommitted()
	return nil
}
