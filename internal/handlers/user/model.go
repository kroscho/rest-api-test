package user

type User struct {
	ID           string `json:"id" bson:"_id, omitempty"` // bson для mongo: _id - уникальный идентификатор, omitempty - поле не может быть пустым
	Email        string `json:"email" bson:"email"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}

type CreateUserDTO struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}
