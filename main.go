package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier().EnableAutoUpdateDisposable()

type DomainCheckRequest struct {
	Domain string `json:"domain"`
}

type DomainCheckResponse struct {
	Disposable bool `json:"disposable"`
}

type EmailCheckRequest struct {
    Email string `json:"email"`
}

type EmailCheckResponse struct {
    Verified string `json:"verified"`
    Error    string `json:"error,omitempty"` // Optional field for error messages
}

func checkDomain(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request DomainCheckRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request: %v", err)
		return
	}

	disposable := verifier.IsDisposable(request.Domain)
	response := DomainCheckResponse{Disposable: disposable}
	json.NewEncoder(w).Encode(response)
}

// func verifyMail(w http.ResponseWriter, r *http.Request) {
//     decoder := json.NewDecoder(r.Body)
// 	var request DomainCheckRequest
// 	err := decoder.Decode(&request)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, "Error decoding request: %v", err)
// 		return
// 	}

//     verified, err := verifier.Verify(request.email)
// 	response := DomainCheckResponse{Verified: verified}
// 	json.NewEncoder(w).Encode(response)

// }

func resToString(res *emailverifier.Result) string {
    ans := ""
    ans += "Email: " + res.Email
    ans += " |Reschable: " + res.Reachable 
    if res.Syntax.Valid {
        ans += " |syntax.valid: " + "True"
    } else {
        ans += " |syntax.valid: " + "False"
    }
    
    if res.Disposable{
        ans += " |Disposable: " + "True"
    }else{
        ans += " |Disposable: " +"Flase"
    }
    if res.RoleAccount{
        ans += " |RoleAccount: " + " True"
    } else {
        ans += " |RoleAccount: " +"False"
    }
    if res.Free {
        ans += " |FreeProvide: " + "True"
    }else {
        ans += " |FreeProvide: " + "False"
    }
    if res.HasMxRecords{
        ans += " |MXRecords: " + "True" 
    } else {
        ans += " |MXRecords: " + "False" 
    }
    return ans
}

func verifyMail(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var request EmailCheckRequest
    err := decoder.Decode(&request)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, "Error decoding request: %v", err)
      return
    }
    fmt.Println(request.Email)
    result, err := verifier.Verify(request.Email)
    if err != nil {
      // Handle verification error (log, return specific error code with message)
      fmt.Println("Error during email verification:", err)
      w.WriteHeader(http.StatusInternalServerError)
      response := EmailCheckResponse{Error: "Internal server error during verification"}
      json.NewEncoder(w).Encode(response)
      return
    }
  
    verified := resToString(result) // Check if the email is deliverable
  
    response := EmailCheckResponse{Verified: verified}
    json.NewEncoder(w).Encode(response)
  }


func main() {
	http.HandleFunc("/check", checkDomain)
	http.HandleFunc("/verify", verifyMail)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
