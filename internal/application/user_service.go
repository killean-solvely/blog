package application

import (
	"errors"
	"log"
	"slices"

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
	domainUserID := domain.NewUserID(userID)

	// Ensure the user exists
	if exists, err := s.userRepo.Exists(domainUserID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("user doesn't exist")
	}

	// Convert parameters to domain variables
	domainUserRoles := []domain.UserRole{}
	for _, role := range userRoles {
		domainUserRoles = append(domainUserRoles, domain.UserRole(role))
	}

	// Get the user
	user, err := s.userRepo.FindByID(domainUserID)
	if err != nil {
		return err
	}

	// Remove the roles that don't exist in the new roles, add the roles that do
	existingRoles := user.UserRoles()
	for _, role := range existingRoles {
		if !slices.Contains(domainUserRoles, role) {
			user.RemoveRole(role)
		}
	}

	existingRoles = user.UserRoles()
	for _, role := range domainUserRoles {
		if !slices.Contains(existingRoles, role) {
			user.AddRole(role)
		}
	}

	// Persist
	if err := s.userRepo.UpdateRoles(domainUserID, user.UserRoles()); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateDescription(userID, description string) error {
	domainUserID := domain.NewUserID(userID)

	// Ensure the user exists
	if exists, err := s.userRepo.Exists(domainUserID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("user doesn't exist")
	}

	// Get the user and update the description
	user, err := s.userRepo.FindByID(domainUserID)
	if err != nil {
		return err
	}

	if err := user.UpdateDescription(description); err != nil {
		return err
	}

	// Persist
	if err := s.userRepo.UpdateDescription(domainUserID, description); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdatePassword(userID, password string) error {
	domainUserID := domain.NewUserID(userID)

	// Ensure the user exists
	if exists, err := s.userRepo.Exists(domainUserID); !exists || err != nil {
		if err != nil {
			return err
		}
		return errors.New("user doesn't exist")
	}

	// Hash the password before passing it into the domain
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	// Get the user and update the description
	user, err := s.userRepo.FindByID(domainUserID)
	if err != nil {
		return err
	}

	if err := user.UpdatePassword(string(passwordHash)); err != nil {
		return err
	}

	// Persist
	if err := s.userRepo.UpdatePasswordHash(domainUserID, string(passwordHash)); err != nil {
		return err
	}

	// Dispatch the events
	if err := s.dispatchAggregateEvents(user); err != nil {
		return err
	}

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
