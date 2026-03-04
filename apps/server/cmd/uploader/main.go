package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// signatureResponse mirrors the listing upload-signature payload.
type signatureResponse struct {
	UploadURL string            `json:"upload_url"`
	Params    map[string]string `json:"params"`
}

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "Base server URL (no trailing slash)")
	cookie := flag.String("cookie", "", "Session cookie value, e.g. cc_refresh_token=... (required)")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("usage: uploader -server http://localhost:8080 -cookie cc_refresh_token=... <file1> [file2 ...]")
		os.Exit(1)
	}
	if *cookie == "" {
		fmt.Println("cookie is required: -cookie cc_refresh_token=...")
		os.Exit(1)
	}

	sig, err := fetchSignature(*serverURL, *cookie)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get signature: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: 60 * time.Second}

	for _, f := range files {
		url, err := uploadFile(client, sig, f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "upload failed for %s: %v\n", f, err)
			continue
		}
		fmt.Printf("%s -> %s\n", f, url)
	}
}

func fetchSignature(serverURL, cookie string) (*signatureResponse, error) {
	body := bytes.NewBufferString("{}")
	req, err := http.NewRequest(http.MethodPost, serverURL+"/api/v1/listings/upload-signature", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("signature request failed: %s: %s", resp.Status, string(b))
	}

	var sig signatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&sig); err != nil {
		return nil, err
	}
	if sig.UploadURL == "" || len(sig.Params) == 0 {
		return nil, fmt.Errorf("invalid signature response")
	}
	return &sig, nil
}

func uploadFile(client *http.Client, sig *signatureResponse, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add signature params
	for k, v := range sig.Params {
		if err := writer.WriteField(k, v); err != nil {
			return "", err
		}
	}

	fw, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, file); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, sig.UploadURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s: %s", resp.Status, string(b))
	}

	var out struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.SecureURL != "" {
		return out.SecureURL, nil
	}
	if out.URL != "" {
		return out.URL, nil
	}
	return "", fmt.Errorf("upload succeeded but no URL returned")
}
