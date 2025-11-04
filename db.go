package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	dbClient    *mongo.Client
	userColl    *mongo.Collection
	todoColl    *mongo.Collection
	sessionColl *mongo.Collection // <-- New
)

// InitDB initializes the database connection and collections
func InitDB(mongoURI string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
	dbClient = client
	db := client.Database("todoAppDB")
	userColl = db.Collection("users")
	todoColl = db.Collection("todos")
	sessionColl = db.Collection("sessions") // <-- New

	// --- New TTL Index for Sessions ---
	// This tells MongoDB to automatically delete documents from the 'sessions'
	// collection when the 'expiresAt' time is reached.
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},     // Index on the 'expiresAt' field
		Options: options.Index().SetExpireAfterSeconds(0), // Delete immediately after 'expiresAt' time
	}
	_, err = sessionColl.Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		log.Fatalf("Failed to create session TTL index: %v", err)
	}
	// --- End New ---

	// Ensure username is unique
	userColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	})
}

// CreateAdminUser creates the 'admin' user if it doesn't exist
func CreateAdminUser() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if admin user already exists
	count, err := userColl.CountDocuments(ctx, bson.M{"username": "admin"})
	if err != nil {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	if count > 0 {
		log.Println("Admin user already exists.")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin-password"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	adminUser := User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
	}

	_, err = userColl.InsertOne(ctx, adminUser)
	if err != nil {
		log.Printf("Failed to create admin user: %v", err)
	} else {
		log.Println("Admin user created successfully.")
	}
}

// GetUserByUsername finds a user by their username
func GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := userColl.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckPasswordHash compares a plaintext password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateTodo inserts a new todo into the database
func CreateTodo(todo Todo) (*Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := todoColl.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}
	todo.ID = res.InsertedID.(primitive.ObjectID)
	return &todo, nil
}

// GetTodosByUserID fetches all todos for a specific user
func GetTodosByUserID(userID primitive.ObjectID) ([]Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todos []Todo
	cursor, err := todoColl.Find(ctx, bson.M{"userId": userID}, options.Find().SetSort(bson.D{{"_id", -1}})) // Sort by newest first
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

// GetTodoByID finds a single todo by its ID and UserID
func GetTodoByID(todoIDHex string, userID primitive.ObjectID) (*Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(todoIDHex)
	if err != nil {
		return nil, err
	}

	var todo Todo
	err = todoColl.FindOne(ctx, bson.M{"_id": objID, "userId": userID}).Decode(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates a todo's title and content
func UpdateTodo(todoIDHex string, userID primitive.ObjectID, title, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(todoIDHex)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"title":   title,
			"content": content,
		},
	}
	_, err = todoColl.UpdateOne(ctx, bson.M{"_id": objID, "userId": userID}, update)
	return err
}

// DeleteTodo removes a todo from the database
func DeleteTodo(todoIDHex string, userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(todoIDHex)
	if err != nil {
		return err
	}

	// Note: We also need to delete the associated image file from disk.
	// This is handled in the handler.
	_, err = todoColl.DeleteOne(ctx, bson.M{"_id": objID, "userId": userID})
	return err
}

// --- New Session Functions ---

// CreateSession generates a new session for a user
func CreateSession(userID primitive.ObjectID, username string, ttl time.Duration) (*Session, error) {
	// Generate a secure, random session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	sessionToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Create the session object
	session := &Session{
		SessionToken: sessionToken,
		UserID:       userID,
		Username:     username,
		ExpiresAt:    time.Now().Add(ttl),
	}

	// Insert into the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := sessionColl.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by its token
func GetSession(sessionToken string) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session Session
	err := sessionColl.FindOne(ctx, bson.M{"sessionToken": sessionToken}).Decode(&session)
	if err != nil {
		return nil, err // Not found
	}

	// Double-check expiration (although TTL index should handle it)
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	return &session, nil
}

// DeleteSession removes a session from the database (logout)
func DeleteSession(sessionToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sessionColl.DeleteOne(ctx, bson.M{"sessionToken": sessionToken})
	return err
}
