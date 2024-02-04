package entities

type Users struct {
	ID        string
	Username  string
	Password  string
	Email     string
	Phone     string
	Address   string
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
		IsEnabled: true,
		CreatedAt: "test",
		UpdatedAt: "test",
		DeletedAt: "test",
	}
}

func GetEmails(topicID string) []string {
	return []string{
		"A@B@C@gmail.com",
		"testmailgmail.com",
		"afwfaw3fghaw3rh45tyhaw43t@gmail.com",
	}
}

func GetTopic(messageID string) string {
	return "1"
}
