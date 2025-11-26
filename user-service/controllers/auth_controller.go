package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/config"
	"user-service/middleware"
	"user-service/models"
	"user-service/utils"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	type RequestDTO struct {
		FirstName string `json:"first_name" validate:"required,max=100"`
		LastName  string `json:"last_name" validate:"required,max=100"`
		Address   string `json:"address" validate:"required,max=255"`
		Email     string `json:"email" validate:"required,email"`
		Phone     string `json:"phone" validate:"required"`
		Password  string `json:"password" validate:"required,min=8"`
	}

	var requestDTO RequestDTO
	json.NewDecoder(r.Body).Decode(&requestDTO)
	if err := validate.Struct(requestDTO); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeValidationError(w, err)
		return
	}

	var existingUser models.User
	result := config.DB.Where("email = ?", requestDTO.Email).First(&existingUser)

	if result.Error == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email is already taken"})
		return
	}

	result = config.DB.Where("phone = ?", requestDTO.Phone).First(&existingUser)

	if result.Error == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Phone number is already taken"})
		return
	}

	user = models.User{
		FirstName: requestDTO.FirstName,
		LastName:  requestDTO.LastName,
		Address:   requestDTO.Address,
		Email:     requestDTO.Email,
		Phone:     requestDTO.Phone,
		Password:  requestDTO.Password,
	}

	hashedPassword, err := utils.HashPassword(requestDTO.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	user_results := config.DB.Create(&user)
	if user_results.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": user_results.Error.Error()})
		return
	}

	notification := map[string]string{
		"user_id":   user.ID,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
		"Phone":     user.Phone,
		"type":      "welcome_email",
	}

	notificationBody, _ := json.Marshal(notification)
	err = utils.Publish("notification_queue", notificationBody)
	if err != nil {
		http.Error(w, "Failed to publish notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!", "user_id": user.ID})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	type LoginDTO struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var requestDTO LoginDTO
	json.NewDecoder(r.Body).Decode(&requestDTO)

	if err := validate.Struct(requestDTO); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		writeValidationError(w, err)
		return
	}

	var user models.User
	results := config.DB.Model(&models.User{}).Where("email = ?", requestDTO.Email).First(&user)

	if results.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
		return
	}

	if user.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
		return
	}

	if err := utils.CheckPasswordHash(requestDTO.Password, user.Password); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
		return
	}

	if user.IsActive == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account is inactive. Please contact support."})
		return
	}

	if user.IsVerified == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account is not verified. Please verify your email."})
		return
	}

	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID missing in context", http.StatusUnauthorized)
		return
	}
	userID := id

	var user models.User
	result := config.DB.Model(&models.User{}).Where("id = ?", userID).First(&user)

	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to fetch user"})
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	sqlDB, err := config.DB.DB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "User service is healthy"})
}

func writeValidationError(w http.ResponseWriter, err error) {
	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())
		errors[field] = fmt.Sprintf("failed on '%s' rule", fieldErr.Tag())
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Validation failed",
		"errors":  errors,
	})
}
