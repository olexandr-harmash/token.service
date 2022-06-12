package user

type (
	User struct {
		Login string
		Password string
		Email string
		Confirmed bool
	}
)