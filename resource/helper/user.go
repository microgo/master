package helper

import (
	"master/model/postgres"
)

func (r *Helper) GetUserByID(ID int) (*postgres.User, error) {
	user := postgres.User{}
	err := r.PostgreSql.
		Where("id = ?", ID).
		Limit(1).
		Find(&user).Error
	return &user, err
}
