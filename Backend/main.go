package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rishabhs-s/controllers"
	"github.com/rishabhs-s/db"
	"github.com/rishabhs-s/websockets"
)

func loadDB() {
	db.InitDb()
	db.Migrate()
}
func main() {
	loadDB()
	router := mux.NewRouter()
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/add/friend", controllers.AddFriend).Methods("POST")
	router.HandleFunc("/accept-reject-friend-request", controllers.AcceptRejectFriendRequest).Methods("POST")
	router.HandleFunc("/remove/friend", controllers.RemoveFriend).Methods("POST") // TODO:doesnot remove from list
	router.HandleFunc("/view/friend/{username}", controllers.ViewFriendList).Methods("GET")
	router.HandleFunc("/create/party", controllers.CreateParty).Methods("POST")
	router.HandleFunc("/party/invite", controllers.InviteToParty).Methods("POST")
	router.HandleFunc("/update/invite", controllers.AcceptRejectInvitation).Methods("POST")
	router.HandleFunc("/leave/party", controllers.LeaveParty).Methods("POST")
	router.HandleFunc("/remove/user", controllers.RemoveUserFromParty).Methods("POST")

	router.HandleFunc("/party", websockets.HandleConnectionsParty)
	go websockets.HandleMessagesParty()
	go websockets.HandleFriendMessages()

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Panic("Error starting server:\n", err)
	}
	log.Print("Running Golang server")
}
