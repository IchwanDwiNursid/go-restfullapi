package cart

import (
	"fmt"
	"net/http"

	"github.com/IchwanDwiNursid/go_restfullapi/service/auth"
	"github.com/IchwanDwiNursid/go_restfullapi/types"
	"github.com/IchwanDwiNursid/go_restfullapi/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)


type Handler struct{
	store types.OrderStore
	productStore types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store , productStore:  productStore,userStore: userStore }
}

func (h *Handler) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/cart/checkout",auth.WithJWTAuth(h.handleCheckout,h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request){
	userId := auth.GetUserIdFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r,&cart);err != nil {
		utils.WriteError(w,http.StatusBadRequest,err)
		return
	}

	if err := utils.Validate.Struct(cart);err != nil {
		errors := err.(validator.FieldError)
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("invalid payload: %v", errors))
		return
	}

	
	// get products
     productIds , err := getCartItemsIds(cart.Items)
	if err != nil {
		utils.WriteError(w,http.StatusBadRequest,err)
		return 
	}	 

	//check product exist in database
	ps , err := h.productStore.GetProductByIds(productIds)
	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError,err)
		return
	}

	orderID , totalPrice , err := h.createOrder(ps,cart.Items,userId)
	if err != nil {
		utils.WriteError(w,http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w,http.StatusOK,map[string]any{
		"total_price" : totalPrice,
		"order_id" : orderID,
	})

}