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
    Email        string `json:"email"`
    Reachable    string `json:"reachable"`
    SyntaxValid  bool   `json:"syntax_valid"`
    Disposable   bool   `json:"disposable"`
    RoleAccount  bool   `json:"role_account"`
    FreeProvider bool   `json:"free_provider"`
    HasMxRecords bool   `json:"has_mx_records"`
    Error        string `json:"error,omitempty"` // Optional field for error messages
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

    response := EmailCheckResponse{
        Email:        result.Email,
        Reachable:    string(result.Reachable),
        SyntaxValid:  result.Syntax.Valid,
        Disposable:   result.Disposable,
        RoleAccount:  result.RoleAccount,
        FreeProvider: result.Free,
        HasMxRecords: result.HasMxRecords,
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/check", checkDomain)
    http.HandleFunc("/verify", verifyMail)
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}
