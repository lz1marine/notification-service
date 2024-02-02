package entities

type Users struct {
	ID        string
	Username  string
	Password  string
	Profile   string
	IsEnabled bool
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}

func GetUser(id string) *Users {
	return &Users{
		ID:        id,
		Username:  "test",
		Password:  "test",
		Profile:   "test",
		IsEnabled: true,
		CreatedAt: "test",
		UpdatedAt: "test",
		DeletedAt: "test",
	}
}
