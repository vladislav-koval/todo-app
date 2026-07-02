package users_postgres_repository

import (
	"github.com/vladislav-koval/todo-app/internal/core/domain"
)

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainFromModel(userModel UserModel) domain.User {
	return domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
}

func userDomainsFromModels(userModels []UserModel) []domain.User {
	userDomains := make([]domain.User, len(userModels))

	for i, userModel := range userModels {
		userDomains[i] = userDomainFromModel(userModel)
	}

	return userDomains
}
