package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IchwanDwiNursid/go_restfullapi/config"
	"github.com/IchwanDwiNursid/go_restfullapi/service/auth"
	"github.com/IchwanDwiNursid/go_restfullapi/types"
	"github.com/IchwanDwiNursid/go_restfullapi/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
    store types.UserStore
}

func NewHandler(store types.UserStore) *Handler{
	return &Handler{
		store : store,
	}
}


func (h *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister ).Methods("POST")
	router.HandleFunc("/users/{id}",auth.WithJWTAuth(h.handleGetUserById,h.store)).Methods("GET")
}


// -----------------------Handle Login------------------------------------

func (h *Handler) handleLogin(w http.ResponseWriter,r *http.Request){
	var payload types.LoginUserPayload
	
	if err := utils.ParseJSON(r,&payload) ; err != nil {
		utils.WriteError(w,http.StatusBadRequest,err)
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("invalid payload %v", errors))
		return 
	}


	// check user by Email

	u , err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("not found , invalid email or password"))
		return 
	}

	if !auth.ComparePasswords(u.Password,[]byte(payload.Password)){
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("not found , invalid email or password"))
		return 
	}

	secret := []byte(config.Envs.JWTSecret)
	token , err := auth.CreateJwt(secret,u.ID)
	if err != nil {
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("not found , invalid email or password"))
		return 
	}


	utils.WriteJSON(w,http.StatusOK, map[string]string{"token":token})


}

// ---------------------Handle Register------------------------------------

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request){

	var payload types.RegisterUserPayload
	
	if err := utils.ParseJSON(r,&payload) ; err != nil {
		utils.WriteError(w,http.StatusBadRequest,err)
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("invalid payload %v", errors))
		return 
	}


	// check if user exist
	 _,err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("user with email %s alredy exists", payload.Email))
		return
	}

	hashedPassword , err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError,err)
		return
	}

	// if it doesnt we create the new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError,err)
		return
	}

	utils.WriteJSON(w,http.StatusCreated, nil)


}

func(h *Handler) handleGetUserById(w http.ResponseWriter,r *http.Request) {
	param := mux.Vars(r)

	idStr := param["id"]

	id , err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError,err)
		return
	}

	user , err := h.store.GetUserById(id)
	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError,err)
		return
	}


	utils.WriteJSON(w,http.StatusOK, user)
}