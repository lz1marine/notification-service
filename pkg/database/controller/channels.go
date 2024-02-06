package controller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lz1marine/notification-service/pkg/client"
	"github.com/lz1marine/notification-service/pkg/database/entity"
)

// ChannelController is the controller for channels
type ChannelController struct {
	db *client.DBClient
}

// NewChannelController creates a new ChannelController
func NewChannelController(db *client.DBClient) *ChannelController {
	return &ChannelController{
		db: db,
	}
}

// GetAllChannels returns all channels, both enabled and disabled
func (cc *ChannelController) GetAllChannels(ctx context.Context) []entity.Channel {
	var chans []entity.Channel
	result := cc.db.GetNotificationsDB().WithContext(ctx).Find(&chans)
	if result.Error != nil {
		fmt.Printf("error retrieving channels: %v", result.Error)
	}

	return chans
}

// GetChannelID returns the channel ID given a channel name
func (cc *ChannelController) GetChannelID(ctx context.Context, name string) (string, error) {
	var channel entity.Channel
	result := cc.db.GetNotificationsDB().WithContext(ctx).Where("name = ?", name).Find(&channel)
	if result.Error != nil {
		fmt.Printf("error retrieving channels: %v", result.Error)
		return "", result.Error
	}

	return channel.ID, nil
}

// GetUserChannels returns the user channels
func (cc *ChannelController) GetUserChannels(ctx context.Context, userID string) ([]entity.UserChannel, error) {
	var userChannels []entity.UserChannel

	result := cc.db.GetNotificationsDB().WithContext(ctx).Preload("Channel").Where("user_id = ?", userID).Find(&userChannels)
	if result.Error != nil {
		fmt.Printf("error retrieving user channels: %v", result.Error)
		return nil, result.Error
	}

	return userChannels, nil
}

// SetUserChannels adds or updates the user channel
func (cc *ChannelController) SetUserChannels(ctx context.Context, userID string, chans []entity.Channel) error {
	userChans, err := cc.GetUserChannels(ctx, userID)
	if err != nil {
		return err
	}

	// Loop through the channels
	for i, v := range chans {
		cur := v
		found := false
		indexFound := -1

		// Check if the channel already exists for the user
		for j, c := range userChans {
			if c.Channel.Name == cur.Name {
				found = true
				indexFound = j
				break
			}
		}

		if found {
			userChans[indexFound].IsEnabled = chans[i].IsEnabled

			// Update
			result := cc.db.GetNotificationsDB().WithContext(ctx).Save(&userChans[indexFound])
			if result.Error != nil {
				fmt.Printf("error updating user channel: %v", result.Error)
				return result.Error
			}
		} else {
			channelID, err := cc.GetChannelID(ctx, cur.Name)
			if err != nil {
				return err
			}

			newEntity := entity.UserChannel{
				ID:        uuid.New().String(),
				UserID:    userID,
				ChannelID: channelID,
				IsEnabled: chans[i].IsEnabled,
			}

			// Insert
			result := cc.db.GetNotificationsDB().WithContext(ctx).Create(&newEntity)
			if result.Error != nil {
				fmt.Printf("error creating user channel: %v", result.Error)
				return result.Error
			}
		}
	}

	return nil
}
