package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type WarTicket struct {
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	Qty           int    `json:"qty"`
	TicketID      int    `json:"ticket_id"`
	Status        string `json:"status"`
}

const (
	fileName   = "war_ticket.json"
	baseURL    = "http://localhost:8001/v1/transaction"
	username   = "user"
	password   = "pass"
	concurrent = 20
)

func main() {
	// baca file json
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// parsing json ke struct
	var tickets []WarTicket

	err = json.Unmarshal(file, &tickets)
	if err != nil {
		panic(err)
	}

	fmt.Printf("total data: %d\n", len(tickets))

	// worker limiter
	sem := make(chan struct{}, concurrent)

	var wg sync.WaitGroup

	start := time.Now()

	for i, ticket := range tickets {
		wg.Add(1)

		go func(index int, data WarTicket) {
			defer wg.Done()

			// acquire semaphore
			sem <- struct{}{}

			defer func() {
				// release semaphore
				<-sem
			}()

			sendRequest(index, data)

		}(i, ticket)
	}

	wg.Wait()

	fmt.Printf("\nDone in %s\n", time.Since(start))
}

func sendRequest(index int, ticket WarTicket) {
	body, err := json.Marshal(ticket)
	if err != nil {
		fmt.Printf("[%d] marshal error: %v\n", index, err)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		fmt.Printf("[%d] request error: %v\n", index, err)
		return
	}

	// basic auth
	auth := base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[%d] http error: %v\n", index, err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Printf(
		"[%d] status=%d response=%s\n",
		index,
		resp.StatusCode,
		string(respBody),
	)
}
