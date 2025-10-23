package repository

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"web3_task_blog/internal/repository/entity"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := GetTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&entity.User{})
	return db
}

func TestUserRepository_CreateUser(t *testing.T) {
	//
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	now := time.Now()
	user := &entity.User{
		Username:  "testuser" + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000)),
		Password:  "password123",
		Email:     "testuser@example.com",
		Age:       20,
		Birthday:  &now,
		PostCount: 0,
	}

	err := repo.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("User created successfully!")

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	now := time.Now()
	user := &entity.User{
		Username:  "testuser",
		Password:  "password123",
		Email:     "testuser@example.com",
		Age:       20,
		Birthday:  &now,
		PostCount: 0,
	}

	err := repo.CreateUser(user)
	assert.NoError(t, err)

	foundUser, err := repo.GetUserByID(uint32(user.ID))
	assert.NoError(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
}