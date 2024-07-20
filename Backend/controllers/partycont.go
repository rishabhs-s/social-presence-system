package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	// "fmt"

	"net/http"

	"github.com/rishabhs-s/db"
	"github.com/rishabhs-s/models"
)

func CreateParty(w http.ResponseWriter, r *http.Request) {
	var partyRequest struct {
		Name     string `json:"name"`
		HostName string `json:"hostname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&partyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the host user by username
	var hostUser models.User
	if err := db.Db.Where("username = ?", partyRequest.HostName).First(&hostUser).Error; err != nil {
		http.Error(w, "Host user not found", http.StatusNotFound)
		return
	}

	// Create a new party
	party := models.Party{
		Name:         partyRequest.Name,
		HostUsername: hostUser.Username,
		Host:         hostUser,
	}

	// Save the party to the database
	if err := db.Db.Create(&party).Error; err != nil {
		http.Error(w, "Failed to create party", http.StatusInternalServerError)
		return
	}

	// Add the party to the host's parties
	if err := db.Db.Model(&hostUser).Association("Parties").Append(&party); err != nil {
		http.Error(w, "Failed to update host's parties", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf("Party %s successfully created by host %s ", partyRequest.Name, partyRequest.HostName))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(party)
}
func InviteToParty(w http.ResponseWriter, r *http.Request) {
	log.Println("InviteToParty")

	var invitation models.PartyInvitation
	if err := json.NewDecoder(r.Body).Decode(&invitation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var friends models.User
	if err := db.Db.Where("username = ?", invitation.FromUserName).First(&friends).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	isFriend := false
	for _, friend := range friends.Friends {
		if friend == invitation.ToUserName {
			isFriend = true
			break
		}
	}
	if !isFriend {
		slog.Error("User is not friends, only friends can invite.")
		http.Error(w, "User is not a friend.", http.StatusForbidden)
		return
	}

	var party models.Party
	if err := db.Db.Where("id = ?", invitation.PartyID).First(&party).Error; err != nil {
		slog.Error("Party not found")
		http.Error(w, "Party not found", http.StatusNotFound)
		return
	}

	for _, invite := range party.Invites {
		if invite == invitation.ToUserName {
			http.Error(w, "User already invited to the party", http.StatusConflict)
			return
		}
	}

	// Add the user to the party's invites
	party.Invites = append(party.Invites, invitation.ToUserName)
	if err := db.Db.Save(&party).Error; err != nil {
		http.Error(w, "Failed to update party invites", http.StatusInternalServerError)
		return
	}

	// Create the invitation
	if err := db.Db.Create(&invitation).Error; err != nil {
		http.Error(w, "Failed to create invitation", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf(" %s successfully invited %s ", invitation.FromUserName, invitation.ToUserName))

	json.NewEncoder(w).Encode(invitation)
}
func LeaveParty(w http.ResponseWriter, r *http.Request) {
	type LeavePartyRequest struct {
		Username string `json:"username"`
		PartyID  uint   `json:"partyid"`
	}

	var leavePartyReq LeavePartyRequest

	if err := json.NewDecoder(r.Body).Decode(&leavePartyReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var party models.Party
	result := db.Db.Preload("Users").Where("id = ? AND host_id != ?", leavePartyReq.PartyID, leavePartyReq.Username).First(&party)

	if result.Error != nil {
		http.Error(w, "Party not found or user is the host", http.StatusNotFound)
		return
	}

	// Retrieve the user
	var user models.User
	if err := db.Db.Where("username = ?", leavePartyReq.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Remove the user from the party's Users relationship
	if err := db.Db.Model(&party).Association("Users").Delete(&user); err != nil {
		http.Error(w, "Failed to leave party", http.StatusInternalServerError)
		return
	}

	// Remove the party from the user's Parties relationship
	if err := db.Db.Model(&user).Association("Parties").Delete(&party); err != nil {
		http.Error(w, "Failed to update user parties", http.StatusInternalServerError)
		return
	}

	slog.Info(fmt.Sprintf(" %s left party %d", leavePartyReq.Username, leavePartyReq.PartyID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Party left successfully"})
}

func AcceptRejectInvitation(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PartyId  uint   `json:"party_id"`
		Status   string `json:"status"` // "accepted" or "rejected"
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var invitation models.PartyInvitation
	if err := db.Db.First(&invitation, request.PartyId).Error; err != nil {
		slog.Error("Invitation not found for user")
		http.Error(w, "Invitation not found", http.StatusNotFound)
		return
	}
	print("-----" + invitation.ToUserName)
	println(request.Username)
	if invitation.ToUserName != request.Username {
		slog.Error("User not invited to this Party : UNAUTHORIZED")

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invitation.Status = request.Status
	if err := db.Db.Save(&invitation).Error; err != nil {
		http.Error(w, "Failed to update invitation status", http.StatusInternalServerError)
		return
	}

	if request.Status == "accepted" {
		var party models.Party
		if err := db.Db.First(&party, invitation.PartyID).Error; err != nil {
			http.Error(w, "Party not found", http.StatusNotFound)
			return
		}
		slog.Info("Party accepted")

		// Add the user to the party's users
		party.Users = append(party.Users, invitation.ToUserName)
		if err := db.Db.Save(&party).Error; err != nil {
			http.Error(w, "Failed to update party users", http.StatusInternalServerError)
			return
		}

		// Add the party to the user's parties
		var user models.User
		if err := db.Db.Where("username = ?", invitation.ToUserName).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err := db.Db.Model(&user).Association("Parties").Append(&party); err != nil {
			http.Error(w, "Failed to update user's parties", http.StatusInternalServerError)
			return
		}
	}
	slog.Info("Invite Rejected.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Invitation updated successfully"})
}

type RemoveUserRequest struct {
	PartyID  uint   `json:"party_id"`
	Username string `json:"username"`
	Host     string `json:"host"`
}

func RemoveUserFromParty(w http.ResponseWriter, r *http.Request) {

	var req RemoveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var party models.Party
	if err := db.Db.Where("id = ?", req.PartyID).First(&party).Error; err != nil {
		slog.Error("No party found for this id")
		http.Error(w, "party not found", http.StatusNotFound)
		return
	}
	if party.HostUsername != req.Host {
		slog.Error("The user is not host, cant remove user")
		http.Error(w, "only the host can remove users", http.StatusUnauthorized)
		return
	}
	var partyInvitation models.PartyInvitation
	if err := db.Db.Where("party_id = ? AND to_user_name = ?", req.PartyID, req.Username).First(&partyInvitation).Error; err != nil {
		slog.Error("User not present in party")
		http.Error(w, "user not present", http.StatusNotFound)
		return
	}
	partyInvitation.Status = "Removed"
	if err := db.Db.Save(&partyInvitation).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Removed user %s from party %d", req.Username, req.PartyID)
	json.NewEncoder(w).Encode(map[string]string{"message": "Host removed User"})

}
