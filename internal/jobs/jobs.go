package jobs

import(	
"fmt"
"net/http"

)

func StartPolling(apiKey string, agentId string) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		poll(apiKey, agentId)
		<-ticker.C
	}
	fmt.Println("Polling for jobs...")
}


func poll(apiKey string, agentId string){
resp, err := http.get("http://localhost:8080/api/v1/agent/poll?apiKey=" + apiKey)
if err != nil {
	fmt.Println("Error while polling for jobs", err)
	return
}

if resp.StatusCode != http.StatusOK {
	fmt.Println("Polling for jobs failed with status:", resp.Status)
	return
}
fmt.Println("Polled for jobs successfully", resp)
}

