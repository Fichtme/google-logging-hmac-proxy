package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/logs", logsHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println("Error closing request body:", err)
		}
	}()

	secret := getSecret()

	if !validateSignature(r, body, secret) {
		http.Error(w, "Invalid Signature", http.StatusForbidden)
		return
	}

	respBody, err := forwardRequest(body)
	if err != nil {
		http.Error(w, "Error forwarding request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respBody)
	if err != nil {
		log.Println("Error writing response body:", err)
	}
}

func logRequest(r *http.Request) {
	log.Printf("Received %s request for %s from %s\n", r.Method, r.RequestURI, r.RemoteAddr)
}

func getSecret() string {
	secret := os.Getenv("HMAC_SECRET")
	if secret == "" {
		secret = "default_secret"
	}
	return secret
}

func validateSignature(r *http.Request, body []byte, secret string) bool {
	receivedSig := r.Header.Get("X-Signature")
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(body)
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return receivedSig == expectedSig
}

func forwardRequest(body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://logging.googleapis.com/v2/entries:write", io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}

	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
	if token == "" {
		token = "your_access_token"
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
