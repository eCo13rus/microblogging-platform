package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex;not null"`
	Email         string `gorm:"uniqueIndex;not null"`
	Password      string `gorm:"not null"`
	Posts         []Post
	Comments      []Comment
	Likes         []Like
	Followers     []User `gorm:"many2many:follows;joinForeignKey:followed_id;joinReferences:follower_id"`
	Following     []User `gorm:"many2many:follows;joinForeignKey:follower_id;joinReferences:followed_id"`
	Notifications []Notification
}

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}

type Comment struct {
	gorm.Model
	PostID  uint
	Post    Post
	UserID  uint
	User    User
	Content string `gorm:"not null"`
}

type Like struct {
	gorm.Model
	PostID uint
	Post   Post
	UserID uint
	User   User
}

type Notification struct {
	gorm.Model
	UserID  uint
	User    User
	Type    string `gorm:"not null"`
	Content string `gorm:"not null"`
	IsRead  bool   `gorm:"default:false"`
}

type Analytics struct {
	gorm.Model
	PostID         uint
	Post           *Post
	Views          int
	LikesCount     int
	CommentsCount  int
	SentimentScore float32
}
