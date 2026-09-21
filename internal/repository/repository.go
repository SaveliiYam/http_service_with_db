package repository

import "github.com/SaveliiYam/http_service_with_db/internal/models"

type UserRepository struct {
	users  []models.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: []models.User{
			{ID: 1, Name: "Alex"},
			{ID: 2, Name: "Maria"},
		},
		nextID: 3,
	}
}

func (r *UserRepository) GetByID(id int) (models.User, bool) {
	for _, user := range r.users {
		if user.ID == id {
			return user, true
		}
	}

	return models.User{}, false
}

func (r *UserRepository) Create(name string) models.User {
	user := models.User{
		ID:   r.nextID,
		Name: name,
	}
	r.users = append(r.users, user)
	r.nextID++

	return user
}
