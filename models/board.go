package models

import "time"

type Board struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BoardElement struct {
	ID        int       `json:"id"`
	BoardID   int       `json:"board_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	PositionX int       `json:"position_x"`
	PositionY int       `json:"position_y"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBoardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateBoardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type CreateBoardElementRequest struct {
	Type      string `json:"type" binding:"required"`
	Content   string `json:"content"`
	PositionX int    `json:"position_x"`
	PositionY int    `json:"position_y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type UpdateBoardElementRequest struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	PositionX int    `json:"position_x"`
	PositionY int    `json:"position_y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
