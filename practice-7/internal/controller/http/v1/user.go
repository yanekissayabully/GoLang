package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	uc usecase.UserInterface
}

func NewUserRoutes(handler *gin.RouterGroup, uc usecase.UserInterface) {
	r := &userRoutes{uc}
	
	h := handler.Group("/users")
	{
		h.POST("/register", r.Register)
		h.POST("/login", r.Login)
		
		protected := h.Group("/")
		protected.Use(utils.JWTAuth())
		{
			protected.GET("/me", r.GetMe)
		}
		
		h.PATCH("/promote/:id", utils.JWTAuth(), utils.RoleCheck("admin"), r.Promote)
	}
}

func (r *userRoutes) Register(c *gin.Context) {
	var req entity.CreateUserDTO
	c.ShouldBindJSON(&req)
	
	hashed, _ := utils.HashPassword(req.Password)
	
	role := req.Role
	if role != "admin" {
		role = "user"
	}
	
	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		Role:     role,
	}
	
	created, _, _ := r.uc.RegisterUser(&user)
	c.JSON(http.StatusCreated, gin.H{"user": created})
}

func (r *userRoutes) Login(c *gin.Context) {
	var req entity.LoginUserDTO
	c.ShouldBindJSON(&req)
	
	token, _ := r.uc.LoginUser(&req)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, _ := r.uc.GetUserByID(userID)
	c.JSON(http.StatusOK, gin.H{
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (r *userRoutes) Promote(c *gin.Context) {
	userID := c.Param("id")
	r.uc.PromoteUser(userID)
	c.JSON(http.StatusOK, gin.H{"message": "promoted to admin"})
}