package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type zabbixAPI struct {
	Key string `json:"result"`
}

var payload = map[string]interface{}{
	"jsonrpc": "2.0",
	"method":  "user.login",
	"params": map[string]interface{}{
		"user":     getEnv("ZABBIX_USER"),
		"password": getEnv("ZABBIX_PASSWORD"),
	},
	"id":   1,
	"auth": nil,
}

var zabbixURL = getEnv("ZABBIX_URL")

func getEnv(key string) (value string) {
	if value = os.Getenv(key); value == "" {
		log.Fatalf("ENV %q not set in environment.", key)
	}
	return value
}

func getZabbixKey(zabbixURL string) {
	// Build the request
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("GET", zabbixURL, body)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-type", "application/json")

	// Send the request via a client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	// Close the body
	defer resp.Body.Close()

	// Get the API key
	var zabbixAPIKey zabbixAPI

	err = json.NewDecoder(resp.Body).Decode(&zabbixAPIKey)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Set ZABBIX_API_KEY in the environment
	err = os.Setenv("ZABBIX_API_KEY", zabbixAPIKey.Key)
	if err != nil {
		log.Printf("Cannot set ZABBIX_API_KEY in environment. Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	getZabbixKey(zabbixURL)
}
