package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/database/entity"
)

// UserController is the controller for users
type UserController struct {
	db *client.DBClient
}

// NewUserController creates a new UserController
func NewUserController(db *client.DBClient) *UserController {
	return &UserController{
		db: db,
	}
}

// GetProfiles returns a list of profiles subscribed to a given topic
func (uc *UserController) GetProfiles(ctx context.Context, topicID string) ([]entity.Profile, error) {
	// Select all users that have subscribed to the topic
	var userTopics []entity.UserTopic
	result := uc.db.GetUsersDB().WithContext(ctx).Preload("User").Where("topic_id = ?", topicID).Find(&userTopics)
	if result.Error != nil {
		fmt.Printf("error retrieving user channels: %v", result.Error)
		return nil, result.Error
	}

	profiles := make([]entity.Profile, 0, len(userTopics))
	for _, userTopic := range userTopics {
		profileString := userTopic.User.Profile

		var profile entity.Profile
		err := json.Unmarshal([]byte(profileString), &profile)
		if err != nil {
			fmt.Printf("error unmarshalling profile: %v", err)
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
