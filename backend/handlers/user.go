package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"guandan-world/backend/auth"
	"guandan-world/backend/storage"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo auth.UserRepository
	storage  storage.Storage
}

func NewUserHandler(userRepo auth.UserRepository, storage storage.Storage) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		storage:  storage,
	}
}

const (
	maxAvatarSize = 5 << 20 // 5MB
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not found in context",
		})
		return
	}

	userIDStr := userID.(string)
	ctx := context.Background()

	nickname := c.PostForm("nickname")
	var nicknamePtr *string
	if nickname != "" {
		nickname = strings.TrimSpace(nickname)
		if len(nickname) > 50 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_nickname",
				Message: "Nickname must be 50 characters or less",
			})
			return
		}
		if nickname != "" {
			nicknamePtr = &nickname
		}
	}

	var avatarKeyPtr *string
	file, header, err := c.Request.FormFile("avatar")
	if err == nil {
		defer file.Close()

		if header.Size > maxAvatarSize {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "file_too_large",
				Message: "Avatar file size must be less than 5MB",
			})
			return
		}

		contentType := header.Header.Get("Content-Type")
		ext, ok := allowedImageTypes[contentType]
		if !ok {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_file_type",
				Message: "Only JPEG, PNG, GIF, and WebP images are allowed",
			})
			return
		}

		timestamp := time.Now().UnixNano()
		avatarKey := fmt.Sprintf("avatars/%s/%d%s", userIDStr, timestamp, ext)

		if err := h.storage.Save(ctx, avatarKey, file, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "upload_failed",
				Message: "Failed to save avatar",
			})
			return
		}

		avatarKeyPtr = &avatarKey
	}

	var oldAvatarKey string
	if avatarKeyPtr != nil {
		user, _ := h.userRepo.FindByID(ctx, userIDStr)
		if user != nil && user.AvatarKey.Valid && user.AvatarKey.String != "" {
			oldAvatarKey = user.AvatarKey.String
		}
	}

	if err := h.userRepo.UpdateProfile(ctx, userIDStr, nicknamePtr, avatarKeyPtr); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: "Failed to update profile",
		})
		return
	}

	if oldAvatarKey != "" && avatarKeyPtr != nil && oldAvatarKey != *avatarKeyPtr {
		_ = h.storage.Delete(ctx, oldAvatarKey)
	}

	userEntity, err := h.userRepo.FindByID(ctx, userIDStr)
	if err != nil || userEntity == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to fetch updated user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userEntity.ToUser(),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid user data in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
