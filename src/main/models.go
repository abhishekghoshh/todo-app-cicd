package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"passwordHash"`
}

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	ImagePath string             `bson:"imagePath"` // Stores the filename, e.g., "12345.jpg"
	UserID    primitive.ObjectID `bson:"userId"`
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SessionToken string             `bson:"sessionToken"` // The random token
	UserID       primitive.ObjectID `bson:"userId"`
	Username     string             `bson:"username"`  // Store username for convenience
	ExpiresAt    time.Time          `bson:"expiresAt"` // Used for TTL index
}
