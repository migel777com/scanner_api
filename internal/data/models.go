package data

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	User    UserModel
	Token   TokenModel
	Product ProductModel
}

/*func NewModels(db *sql.DB) Models {
	return Models{
		UserModel{DB: db},
		TokenModel{DB: db},
	}
}*/

func NewModels() Models {
	return Models{
		Product: ProductModel{},
	}
}
