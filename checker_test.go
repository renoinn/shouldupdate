package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLatestVersionGitHubImpl(t *testing.T) { // Renamed test function to match target
	// Test case 1: Valid response
	t.Run("ValidResponse", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/owner/repo/releases/latest" {
				t.Errorf("Expected to request '/repos/owner/repo/releases/latest', got: %s", r.URL.Path)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			if r.Header.Get("User-Agent") != "ShouldUpdateApp/1.0" {
				t.Errorf("Expected User-Agent 'ShouldUpdateApp/1.0', got: '%s'", r.Header.Get("User-Agent"))
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"tag_name": "v1.2.3", "name": "Release 1.2.3"}`)
		}))
		defer server.Close()

		version, err := getLatestVersionGitHubImpl("owner/repo", server.URL) // Corrected function call
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if version != "1.2.3" {
			t.Errorf("Expected version '1.2.3', got: '%s'", version)
		}
	})

	// Test case 2: Valid response with non-v prefix tag
	t.Run("ValidResponseNoVPrefix", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"tag_name": "2.0.0", "name": "Release 2.0.0"}`)
		}))
		defer server.Close()

		version, err := getLatestVersionGitHubImpl("owner/repo", server.URL) // Corrected function call
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if version != "2.0.0" {
			t.Errorf("Expected version '2.0.0', got: '%s'", version)
		}
	})

	// Test case 3: GitHub API error (e.g., 404 Not Found)
	t.Run("GitHubAPIError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Found", http.StatusNotFound)
		}))
		defer server.Close()

		_, err := getLatestVersionGitHubImpl("owner/nonexistent_repo", server.URL) // Corrected function call
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "GitHub API error for owner/nonexistent_repo (status 404)") {
			t.Errorf("Expected error to contain 'GitHub API error for owner/nonexistent_repo (status 404)', got: %v", err)
		}
	})

	// Test case 4: Malformed JSON response
	t.Run("MalformedJSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"tag_name": "v1.0"`) // Missing closing brace
		}))
		defer server.Close()

		_, err := getLatestVersionGitHubImpl("owner/repo", server.URL) // Corrected function call
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "error decoding JSON response") {
			t.Errorf("Expected error to contain 'error decoding JSON response', got: %v", err)
		}
	})

	// Test case 5: Empty tag_name in response
	t.Run("EmptyTagName", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"name": "Release without tag"}`)
		}))
		defer server.Close()

		_, err := getLatestVersionGitHubImpl("owner/repo", server.URL) // Corrected function call
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "no version tag (tag_name) found") {
			t.Errorf("Expected error to contain 'no version tag (tag_name) found', got: %v", err)
		}
	})

	// Test case 6: Invalid appIdentifier format
	t.Run("InvalidAppIdentifier", func(t *testing.T) {
		// No server needed as this should be caught before HTTP request
		_, err := getLatestVersionGitHubImpl("ownerrepo", "") // Corrected function call; No slash
		if err == nil {
			t.Fatal("Expected an error for invalid appIdentifier, got nil")
		}
		if !strings.Contains(err.Error(), "invalid application identifier") {
			t.Errorf("Expected error about 'invalid application identifier', got: %v", err)
		}
	})

	// Test case 7: Network error (simulated by closing server immediately)
	t.Run("NetworkErrorSimulation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handler won't be reached
		}))
		serverURL := server.URL
		server.Close()

		_, err := getLatestVersionGitHubImpl("owner/repo", serverURL) // Corrected function call
		if err == nil {
			t.Fatal("Expected a network error, got nil")
		}
		if !strings.Contains(err.Error(), "network error fetching release info") {
			t.Errorf("Expected error to contain 'network error fetching release info', got: %v", err)
		}
	})
}
