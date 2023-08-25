package requestschemas

type Login struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type Register struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type CreateUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Admin    bool
}
