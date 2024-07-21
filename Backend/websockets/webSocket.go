
package websockets

import (

	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rishabhs-s/db"
	"github.com/rishabhs-s/models"
)

type Message struct {
	Type    string `json:"type"`
	Partyid int `json:"partyid"`
	User    string `json:"user"`
	Content string `json:"content"`
	Status   string `json:"status"`
}

// type OnlineMessage struct {
// 	Username string `json:"username"`
// 	Status   string `json:"status"`
// }

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan Message)
	// Onlinebroadcast = make(chan OnlineMessage)

	mutex     = &sync.Mutex{}
)

func HandleConnectionsParty(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	username := r.URL.Query().Get("username")
	if username == "" {
		ws.Close()
		return
	}

	mutex.Lock()
	clients[ws] = username
	broadcast <- Message{Type: "status", Content: "online", User: username}
	mutex.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			mutex.Lock()
			delete(clients, ws)
			broadcast <- Message{Type: "status", Content: "offline", User: username,Partyid: msg.Partyid}
			mutex.Unlock()
			break
		}

		// Handle specific message types
		switch msg.Type {
		case "join_party":
			print("join party")
				if msg.User != "" {
					userJoinParty(msg.User, msg.Partyid)
				}
		case "leave_party":
			print("Leave Party")

			    if msg.User != "" {
			        userLeaveParty(msg.User, msg.Partyid)
			    }

			case "online":
				print("Online type")

		default:
			broadcast <- msg
		}
	}
}

func HandleMessagesParty() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
func HandleFriendMessages() {
	for {
		msg := <-broadcast
		friends := getUserFriends(msg.User)

		mutex.Lock()
		for client, username := range clients {
			if contains(friends, username) {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mutex.Unlock()
	}
} 
func getUsersInParty(partyid int) ([]string) {
	var party models.Party

	if err := db.Db.Where("id = ?", partyid).First(&party).Error; err != nil {
		log.Println("Error in getting users")

	}
	users, err := party.Users.ToStringSlice()
    if err != nil {
        log.Println("Error converting users to slice:", err)
        return nil
    }
	return users
}

func userJoinParty(username string, partyID int) {
	print("Inside Join Party func")

	var user models.User
	var party models.Party
	if err := db.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println(err)
		return
	}
	if err := db.Db.First(&party, partyID).Error; err != nil {
		log.Println(err)
		return
	}
	userAlreadyJoined := false
	for _, u := range party.Users {
		if u == user.Username {
			userAlreadyJoined = true
			break
		}
	}

	if userAlreadyJoined {
		log.Printf("%s is already in the party %d", user.Username, partyID)
		return
	}


	party.Users = append(party.Users, user.Username)


	if err := db.Db.Save(&party).Error; err != nil {
		log.Println(err)
		return
	}

	broadcastPartyStatus(partyID, "joined", user.Username)
}

func userLeaveParty(username string, partyID int) {
	print("Inside Leave Party func")
	var user models.User
	var party models.Party
	if err := db.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println(err)
		return
	}
	if err := db.Db.First(&party, partyID).Error; err != nil {
		log.Println(err)
		return
	}

	for i, u := range party.Users {
		if u == user.Username {
			party.Users = append(party.Users[:i], party.Users[i+1:]...)
			break
		}
	}

	if err := db.Db.Save(&party).Error; err != nil {
		log.Println(err)
		return
	}

	broadcastPartyStatus(partyID, "left", user.Username)
}
func broadcastPartyStatus(partyID int, status string, username string) {
	print(partyID)
    users := getUsersInParty(partyID)
    mutex.Lock()
    defer mutex.Unlock()

    for _, user := range users {
        for client, clientUsername := range clients {
            if clientUsername == user {
                err := client.WriteJSON(Message{Type: "party_status", Content: status, User: username, Partyid:partyID })
                if err != nil {
                    log.Println("Error sending message to client:", err)
                }
            }
        }
    }
}


func broadcastOnlineStatus(partyID int, status string, friends[] string) {

}

func userStatusChange(username, status string) {
	msg := Message{User: username, Status: status}
	broadcast <- msg
}

func getUserFriends(username string) []string {
	var user models.User
	if err := db.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println(err)
	}	
	var friends []string
	for _, u := range user.Friends {
		if u == user.Username {
			friends = append(friends,u.(string))
			break
		}
	}
	return friends
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}