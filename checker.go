package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GitHubReleaseInfo struct to unmarshal the relevant parts of the GitHub API JSON response.
type GitHubReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`        // For more descriptive release name
	Body    string `json:"body"`        // For release notes/changelog
	HTMLURL string `json:"html_url"`    // Link to the release page
}

// getLatestVersionGitHub fetches the latest release tag name for a given appIdentifier (owner/repo).
// For testability, apiBaseURL can be provided to point to a mock server.
// If apiBaseURL is empty, it defaults to "https://api.github.com".
// This is the internal implementation.
func getLatestVersionGitHubImpl(appIdentifier string, apiBaseURL string) (string, error) {
	if !strings.Contains(appIdentifier, "/") {
		return "", fmt.Errorf("invalid application identifier: expected 'owner/repo', got '%s'", appIdentifier)
	}

	baseURL := "https://api.github.com"
	if apiBaseURL != "" {
		baseURL = apiBaseURL // Use mock server URL for testing
	}

	url := fmt.Sprintf("%s/repos/%s/releases/latest", baseURL, appIdentifier)

	client := &http.Client{} // Consider setting a timeout: client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("internal error creating request for %s: %w", appIdentifier, err)
	}
	req.Header.Set("User-Agent", "ShouldUpdateApp/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error fetching release info for %s from %s: %w", appIdentifier, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorMsg strings.Builder
		errorMsg.WriteString(fmt.Sprintf("GitHub API error for %s (status %d)", appIdentifier, resp.StatusCode))
		// Attempt to read GitHub's error message
		var ghError struct {
			Message          string `json:"message"`
			DocumentationURL string `json:"documentation_url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&ghError); err == nil && ghError.Message != "" {
			errorMsg.WriteString(fmt.Sprintf(": %s", ghError.Message))
			if ghError.DocumentationURL != "" {
				errorMsg.WriteString(fmt.Sprintf(" (see %s)", ghError.DocumentationURL))
			}
		} else {
			// Fallback if parsing GitHub's specific error fails
			errorMsg.WriteString(fmt.Sprintf(" (URL: %s)", url))
		}
		return "", errors.New(errorMsg.String())
	}

	var releaseInfo GitHubReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return "", fmt.Errorf("error decoding JSON response for %s from %s: %w", appIdentifier, url, err)
	}

	if releaseInfo.TagName == "" {
		return "", fmt.Errorf("no version tag (tag_name) found in the latest release for %s (URL: %s)", appIdentifier, url)
	}

	// Clean "v" prefix, if any
	tagName := strings.TrimPrefix(releaseInfo.TagName, "v")
	return tagName, nil
}

// getLatestVersion is a package-level variable that points to the actual implementation.
// Tests can override this variable to mock the GitHub API interaction.
var getLatestVersion = getLatestVersionGitHubImpl
