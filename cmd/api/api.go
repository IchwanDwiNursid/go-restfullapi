package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/IchwanDwiNursid/go_restfullapi/service/cart"
	"github.com/IchwanDwiNursid/go_restfullapi/service/order"
	"github.com/IchwanDwiNursid/go_restfullapi/service/product"
	"github.com/IchwanDwiNursid/go_restfullapi/service/user"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr string
	db *sql.DB
}


func NewApiServer(addr string , db *sql.DB) *ApiServer {
	return  &ApiServer{
		addr: addr,
		db : db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//user-router
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	//product router
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	//order router
	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore,productStore,userStore)
	cartHandler.RegisterRoutes(subrouter)


	log.Println("Server Listen on port" + s.addr)

	return http.ListenAndServe(s.addr,router)
}