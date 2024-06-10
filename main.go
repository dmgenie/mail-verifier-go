package main

import (
	"fmt"
	
	emailverifier "github.com/AfterShip/email-verifier"
)

var (
    verifier = emailverifier.
        NewVerifier().
        EnableAutoUpdateDisposable()
)

func main() {
    domain := "lapeds.com"
    if verifier.IsDisposable(domain) {
        fmt.Printf("%s is a disposable domain\n", domain)
        return
    }
    fmt.Printf("%s is not a disposable domain\n", domain)
}
