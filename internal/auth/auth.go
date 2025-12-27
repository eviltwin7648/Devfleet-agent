package auth

import (
	"net/http"
	"fmt"
)

func VerifyAgent(apiKey string) bool {
	resp, err := http.Get("http://localhost:8080/api/v1/agent/verify?apiKey=" + apiKey)
	if err != nil {
		fmt.Println("Error While Verifying Agent",err)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Agent verification failed with status:", resp.Status)
		return false
	}
	fmt.Println("Agent Verified Successfully", resp)
	//should get a jwt token and send with every single request
	defer resp.Body.Close()

	return true
}