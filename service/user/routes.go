package user

import (
	"fmt"
	"go-sample-rest-api/config"
	auth2 "go-sample-rest-api/service/auth"
	"go-sample-rest-api/types"
	"go-sample-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
	auth  auth2.Authenticator
}

func NewHandler(store types.UserStore, auth auth2.Authenticator) *Handler {
	return &Handler{store: store, auth: auth}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

	// admin routes
	router.HandleFunc("/users/{userID}", auth2.WithJWTAuth(h.handleGetUser, h.store)).Methods(http.MethodGet)
}

// handleLogin godoc
// @Summary User login
// @Description Login with email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body types.LoginUserPayload true "Login Credentials"
// @Success 200 {object} map[string]string "token: JWT Token on successful login."
// @Failure 400 {object} types.HTTPError "Bad Request when the payload is invalid."
// @Failure 404 {object} types.HTTPError "Not Found, invalid email or password."
// @Failure 500 {object} types.HTTPError "Internal Server Error"
// @Router /login [post]
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil || u == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !h.auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := h.auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// handleRegister godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body types.RegisterUserPayload true "Register Information"
// @Success 201 {object} nil "Successfully registered and no content returned."
// @Failure 400 {object} types.HTTPError "Bad Request if the payload is invalid or user exists."
// @Failure 500 {object} types.HTTPError "Internal Server Error if database error occurs."
// @Router /register [post]
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	// hash password
	hashedPassword, err := h.auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get detailed information about a user.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} types.User "Successful retrieval of user detail."
// @Failure 400 {object} types.HTTPError "Bad Request if user ID is missing or invalid."
// @Failure 404 {object} types.HTTPError "Not Found if user does not exist."
// @Failure 500 {object} types.HTTPError "Internal Server Error"
// @Router /users/{id} [get]
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["userID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
