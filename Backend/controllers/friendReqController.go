package controllers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rishabhs-s/db"
	"github.com/rishabhs-s/models"
)

func AddFriend(w http.ResponseWriter, r *http.Request) {
	var req models.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Error in decoding body of user ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var fromUser models.User
	if err := db.Db.Where("username = ?", req.FromUserName).First(&fromUser).Error; err != nil {
		slog.Error("From user not found: ", err)
		http.Error(w, "From user not found", http.StatusNotFound)
		return
	}

	// Check if the ToUserName exists
	var toUser models.User
	if err := db.Db.Where("username = ?", req.ToUserName).First(&toUser).Error; err != nil {
		slog.Error("To user not found: ", err)
		http.Error(w, "To user not found", http.StatusNotFound)
		return
	}
	var existingRequest models.FriendRequest
	if err := db.Db.Where("(from_user_name = ? AND to_user_name = ?) OR (from_user_name = ? AND to_user_name = ?)",
		req.FromUserName, req.ToUserName, req.ToUserName, req.FromUserName).First(&existingRequest).Error; err == nil {
		if existingRequest.Status == "pending" {
			http.Error(w, "Friend request already sent", http.StatusConflict)
			return
		} else if existingRequest.Status == "accepted" {
			http.Error(w, "Users are already friends", http.StatusConflict)
			return
		}
	}

	req.Status = "pending"
	if err := db.Db.Create(&req).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Friend request from User ID %s to User ID %s", req.FromUserName, req.ToUserName)

	json.NewEncoder(w).Encode(req)
}

func AcceptRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	var req models.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Error in decoding body of user ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fr models.FriendRequest
	if err := db.Db.First(&fr, req.ID).Error; err != nil {
		slog.Error("No Request with this id found", err)
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	fr.Status = req.Status
	if err := db.Db.Save(&fr).Error; err != nil {
		slog.Error("Not able to save friend req %s: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fr.Status == "accepted" {
		var fromUser, toUser models.User
		if err := db.Db.Where("username = ?", fr.FromUserName).First(&fromUser).Error; err != nil {
			log.Printf("Error fetching fromUser: %v", err)
			http.Error(w, "From user not found", http.StatusNotFound)
			return
		}

		if err := db.Db.Where("username = ?", fr.ToUserName).First(&toUser).Error; err != nil {
			log.Printf("Error fetching toUser: %v", err)
			http.Error(w, "To user not found", http.StatusNotFound)
			return
		}

		// Update Friends array for both users
		fromUser.Friends = append(fromUser.Friends, toUser.Username)
		toUser.Friends = append(toUser.Friends, fromUser.Username)

		// Save updated users back to the database
		if err := db.Db.Save(&fromUser).Error; err != nil {
			log.Printf("Error saving fromUser: %v", err)
			http.Error(w, "Failed to update fromUser friends", http.StatusInternalServerError)
			return
		}
		if err := db.Db.Save(&toUser).Error; err != nil {
			log.Printf("Error saving toUser: %v", err)
			http.Error(w, "Failed to update toUser friends", http.StatusInternalServerError)
			return
		}

		log.Printf("Users %s and %s added to each other's friends list", fromUser.Username, toUser.Username)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fr)
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUsername   string `json:"from_user_name"`
		RemoveUsername string `json:"remove_user_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fromUser, toUser models.User

	if err := db.Db.Where("username = ?", req.FromUsername).First(&fromUser).Error; err != nil {
		slog.Error("User(%s) not found", req.FromUsername)
		http.Error(w, "from user not found", http.StatusNotFound)
		return
	}
	if err := db.Db.Where("username = ?", req.RemoveUsername).First(&toUser).Error; err != nil {
		slog.Error("User(%s) not found", req.RemoveUsername)
		http.Error(w, "to user not found", http.StatusNotFound)
		return
	}
	// Remove toUser's ID from fromUser friends array
	for i, frienduserName := range fromUser.Friends {
		log.Print(frienduserName)
		log.Printf("Removing from %s", fromUser.Username)
		if frienduserName == toUser.Username {
			fromUser.Friends = append(fromUser.Friends[:i], fromUser.Friends[i+1:]...)
			break
		}
	}

	// Remove fromUser ID from toUser friends array
	for i, frienduserName := range toUser.Friends {
		log.Printf("Removing from %s", toUser.Username)
		log.Print(frienduserName)

		if frienduserName == fromUser.Username {
			slog.Debug("list2", toUser.Friends)
			toUser.Friends = append(toUser.Friends[:i], toUser.Friends[i+1:]...)
			slog.Debug("list2", toUser.Friends)
			break
		}
	}
	slog.Info("Both removed from each others friend list")

	// Save updated users back to the database
	if err := db.Db.Save(&fromUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Db.Save(&toUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("friend list Update Success")

	json.NewEncoder(w).Encode(map[string]string{"message": "friend removed successfully"})

}

func ViewFriendList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var user models.User
	if err := db.Db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("User no t found %s", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	slog.Info("Friend list fetch successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Friends)
}
