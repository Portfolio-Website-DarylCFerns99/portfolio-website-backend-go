package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"portfolio-website-backend/internal/config"
)

func getAdminURL(endpoint string) string {
	// The server uses localhost:8000/api/v1 by default
	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}
	return fmt.Sprintf("http://localhost:8000%s/admin/%s", apiPrefix, endpoint)
}

func defaultHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("ADMIN_API_KEY")
	if apiKey == "" {
		log.Fatal("ADMIN_API_KEY environment variable is missing (check .env file).")
	}
	req.Header.Set("X-Admin-Api-Key", apiKey)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("User management CLI tool")
		fmt.Println("Usage: manage_user <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  create  Create a new user")
		fmt.Println("  list    List all users")
		fmt.Println("  get     Get a user by ID")
		os.Exit(1)
	}

	// Load configuration to ensure ADMIN_API_KEY is extracted from .env
	config.LoadConfig()

	command := os.Args[1]
	client := &http.Client{}

	switch command {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		email := createCmd.String("email", "", "User's email")
		password := createCmd.String("password", "", "User's password")
		username := createCmd.String("username", "", "User's username (optional, email will be used if not provided)")
		createCmd.Parse(os.Args[2:])

		if *email == "" || *password == "" {
			log.Fatal("--email and --password are required")
		}

		payload := map[string]string{
			"email":    *email,
			"password": *password,
			"username": *username,
		}
		jsonLoad, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, getAdminURL("users"), bytes.NewBuffer(jsonLoad))
		if err != nil {
			log.Fatal(err)
		}
		defaultHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error connecting to server:", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			log.Fatalf("Server returned error %d: %s", resp.StatusCode, string(body))
		}
		log.Printf("User created successfully: %s\n", string(body))

	case "list":
		req, err := http.NewRequest(http.MethodGet, getAdminURL("users"), nil)
		if err != nil {
			log.Fatal(err)
		}
		defaultHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error connecting to server:", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			log.Fatalf("Server returned error %d: %s", resp.StatusCode, string(body))
		}

		log.Printf("List of Users:\n%s\n", string(body))

	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		id := getCmd.String("id", "", "User's ID (UUID)")
		getCmd.Parse(os.Args[2:])

		if *id == "" {
			log.Fatal("--id is required")
		}

		req, err := http.NewRequest(http.MethodGet, getAdminURL("users/"+*id), nil)
		if err != nil {
			log.Fatal(err)
		}
		defaultHeaders(req)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error connecting to server:", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode >= 400 {
			log.Fatalf("Server returned error %d: %s", resp.StatusCode, string(body))
		}

		log.Printf("User Found:\n%s\n", string(body))

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
