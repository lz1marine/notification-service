package entities

import "errors"

type Channels struct {
	ID        string
	Name      string
	IsEnabled bool
}

// TODO: delme
func GetChannels() []Channels {
	return []Channels{
		{
			ID:        "1",
			Name:      "SMS",
			IsEnabled: true,
		},
		{
			ID:        "2",
			Name:      "Email",
			IsEnabled: true,
		},
		{
			ID:        "3",
			Name:      "Slack",
			IsEnabled: true,
		},
	}

}

func GetChannelID(name string) *Channels {
	chans := GetChannels()

	for _, v := range chans {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

type UserChannels struct {
	User      *Users
	Channel   *Channels
	IsEnabled bool
}

// TODO: delme
var userChans = []UserChannels{
	{
		User: GetUser("1"),
		Channel: &Channels{
			ID:        "1",
			Name:      "SMS",
			IsEnabled: true,
		},
		IsEnabled: true,
	},
	{
		User: GetUser("1"),
		Channel: &Channels{
			ID:        "2",
			Name:      "Email",
			IsEnabled: true,
		},
		IsEnabled: true,
	},
}

func GetUserChannel(userID string) []Channels {
	if userID != "1" {
		return []Channels{}
	}

	chans := make([]Channels, 0, len(userChans))
	for _, v := range userChans {
		c := Channels{
			ID:        v.Channel.ID,
			Name:      v.Channel.Name,
			IsEnabled: v.IsEnabled,
		}
		chans = append(chans, c)
	}

	return chans
}

func SetUserChannel(userID string, chans []Channels) error {
	if userID != "1" {
		return errors.New("User not found")
	}

	availChans := GetChannels()

	for i, v := range chans {
		found := false
		indexFound := -1
		for j, c := range userChans {
			if c.Channel.Name == v.Name {
				found = true
				indexFound = j
				break
			}
		}

		if found {
			userChans[indexFound].IsEnabled = chans[i].IsEnabled
		} else {
			for _, c := range availChans {
				if c.Name == v.Name {
					userChans = append(userChans, UserChannels{
						User:      GetUser(userID),
						Channel:   &c,
						IsEnabled: chans[i].IsEnabled,
					})
					break
				}
			}
		}
	}

	return nil
}

type Notifications struct {
	ID         string
	Title      string
	Message    string
	TemplateID string
	Status     int
	ChannelID  string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

// TODO
func AddNotification(notification *Notifications, eventID string) error {
	return nil
}

// TODO
func PatchNotification(notification *Notifications, messageID string) error {
	return nil
}

// TODO: this implementation expects a small amount of templates. If this is not the case, we should use a nosql document store
type Templates struct {
	ID        string
	Template  string
	IsEnabled bool
}

func GetTemplates(templateID string) Templates {
	return Templates{
		ID: "1",
		Template: `<!DOCTYPE html>
<html>
<body>
	<h3>Name:</h3><span>Hello {{.Name}}</span><br/><br/>
	<h3>Email:</h3><span>{{.Email}}</span><br/>
	<h3>Message:</h3><span>{{.Message}}</span><br/>
</body>
</html>`,
		IsEnabled: true,
	}
}
