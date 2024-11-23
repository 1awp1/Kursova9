package response

import "dim_kurs/internal/domain/model"

type GetUsers struct {
	Page  int
	Limit int
	users []model.User
}
