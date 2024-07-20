package controllers

// import (
// 	"encoding/json"
// 	// "fmt"
// 	"log"
// 	"net/http"
// 	// "reflect"

// 	"github.com/rishabhs-s/db"
// 	"github.com/rishabhs-s/models"
// )

// func CreateParty(w http.ResponseWriter, r *http.Request) {
//     var party models.Party
//     if err := json.NewDecoder(r.Body).Decode(&party); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Save the party to the database
//     db.Db.Create(&party)

//     // Return the created party with its ID
//     json.NewEncoder(w).Encode(party)
// }


// func InviteToParty(w http.ResponseWriter, r *http.Request) {
// 	log.Println("InviteToParty")
// 	// partyID := mux.Vars(r)["id"]
//     var invitation models.PartyInvitation
// 	if err := json.NewDecoder(r.Body).Decode(&invitation); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

 
// 	var friends models.User
//     db.Db.Model(&models.User{}).Where("username = ?", invitation.FromUserName).Find(&friends)

// 	isFriend:=false
// 	for _,friend := range friends.Friends{
// 		if(friend==invitation.ToUserName){
// 			isFriend=true
// 			break
// 		}
// 	}
// 	if !isFriend {
// 		log.Print("")
//         http.Error(w, "User is not a friend.", http.StatusForbidden)
//         return
//     }
//     // Create the invitation
//     db.Db.Create(&invitation)

//     // Return the created invitation
//     json.NewEncoder(w).Encode(invitation)
// }


// func UpdatePartyInviteStatus(w http.ResponseWriter, r *http.Request) {
//     type UpdateInvite struct {
//         ToUserName string `json:"to_user_name"`
//         Status     string `json:"status"`
//     }
// 	var updateInvite UpdateInvite

//     if err := json.NewDecoder(r.Body).Decode(&updateInvite); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}
//     if updateInvite.Status != "accepted" && updateInvite.Status != "rejected" {
// 		http.Error(w, "Invalid status", http.StatusBadRequest)
// 		return
// 	}

// 	var invitation models.PartyInvitation
// 	result := db.Db.Where("to_user_name = ? AND status = ?", updateInvite.ToUserName, "pending").First(&invitation)

// 	if result.Error != nil {
// 		http.Error(w, "Invitation not found or already responded", http.StatusNotFound)
// 		return
// 	}

// 	invitation.Status = updateInvite.Status
// 	db.Db.Save(&invitation)


//     if updateInvite.Status == "accepted" {
// 		// Retrieve the user and update the parties column
// 		var user models.User
// 		if err := db.Db.Where("username = ?", updateInvite.ToUserName).First(&user).Error; err != nil {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		// Assuming the party ID is stored in the invitation
// 		partyID := invitation.PartyID
// 		// Append the new party to the user's parties array
// 		user.Parties = append(user.Parties, partyID)

// 		// Save the updated user record
// 		if err := db.Db.Save(&user).Error; err != nil {
// 			http.Error(w, "Failed to update user parties", http.StatusInternalServerError)
// 			return
// 		}
//     }

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(invitation)
// }



// func LeaveParty(w http.ResponseWriter, r *http.Request) {
//     type LeavePartyRequest struct {
//         Username string `json:"username"`
//         PartyID    uint   `json:"partyid"`
//     }

// 	var leavePartyReq LeavePartyRequest

//     if err := json.NewDecoder(r.Body).Decode(&leavePartyReq); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	var invitation models.PartyInvitation
// 	result := db.Db.Where("to_user_name = ? AND party_id = ? AND status = ?", leavePartyReq.Username, leavePartyReq.PartyID, "accepted").First(&invitation)

// 	if result.Error != nil {
// 		http.Error(w, "Invitation not found or not accepted", http.StatusNotFound)
// 		return
// 	}

// 	// Update invitation status to "left"
// 	invitation.Status = "left"
// 	db.Db.Save(&invitation)

// 	// Retrieve the user and update the parties column
// 	var user models.User
// 	if err := db.Db.Where("username = ?", leavePartyReq.Username).First(&user).Error; err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}
    

// 	// Remove the party from the user's parties array

// 	for i, id := range user.Parties {
//         partyID := float64(leavePartyReq.PartyID)
//         fmt.Println(reflect.TypeOf(id))
// 		if (id) == partyID {
//             print("true")
// 			user.Parties = append(user.Parties[:i], user.Parties[i+1:]...)
// 			break
// 		}

// 	}

// 	// Save the updated user record
// 	if err := db.Db.Save(&user).Error; err != nil {
// 		http.Error(w, "Failed to update user parties", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Party left successfully"})
// }