package auth

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/BenjaminArmijo3/bank/internal/db/sqlc"
	"github.com/BenjaminArmijo3/bank/internal/pkg/utils/crypt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type AuthParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *AuthApp) register(c *gin.Context) {
	var user AuthParams

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := crypt.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := sqlc.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}

	newUser, err := a.server.Store.CreateUser(context.Background(), arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{}.ToUserResponse(&newUser))
}

func (a *AuthApp) login(c *gin.Context) {
	user := new(AuthParams)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbUser, err := a.server.Store.Queries.GetUserByEmail(context.Background(), user.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Email"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println(dbUser.ID)

	if err := crypt.VerifyPassword(user.Password, dbUser.HashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect Email or Password"})
		return
	}

	jwt, err := a.server.Token.CreateToken(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwt})

}

func (u RegisterResponse) ToUserResponse(user *sqlc.User) *RegisterResponse {
	return &RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
