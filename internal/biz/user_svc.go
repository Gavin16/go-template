package biz

import (
	models "minsky/go-template/internal/models"
)

type UserSvc interface {
	GetUserById(id int64)
}

type UserSvcImpl struct {
}

func (userSvc *UserSvcImpl) GetUserById(id int64) models.User {
	user := models.User{}
	user.Id = id
	user.Name = "Mostly"
	user.Age = 31
	user.Title = "marketing"
	user.Email = "Mostly@gmail.com"
	user.Nation = "USA"
	return user
}
