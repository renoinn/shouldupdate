package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// stripAnsiCodes removes ANSI escape sequences from a string.
func stripAnsiCodes(str string) string {
	ansiRegex := regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[mGKHF]")
	return ansiRegex.ReplaceAllString(str, "")
}

// captureOutput captures stdout for a given function.
// For tests that need to check Stderr, they should temporarily reassign os.Stderr.
var captureOutput = func(f func()) string { // Made it a var for potential redefinition in specific tests if needed, though not used here.
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	outBytes, _ := io.ReadAll(r)
	os.Stdout = oldStdout
	return string(outBytes)
}

// TestHandleAddCommand tests the add command functionality.
func TestHandleAddCommand(t *testing.T) {
	originalConfigFileValue := configFile
	testFile := "test_add_versions.toml"
	configFile = testFile
	defer func() {
		configFile = originalConfigFileValue
		os.Remove(testFile)
	}()

	t.Run("AddNewApplication", func(t *testing.T) {
		os.Remove(testFile)
		appName := "myNewApp"
		appVersion := "1.0.0"
		output := stripAnsiCodes(captureOutput(func() {
			handleAddCmd(appName, appVersion)
		}))
		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if _, ok := cfg[appName]; !ok || cfg[appName] != appVersion {
			t.Errorf("Expected app %s with version %s, not found or version mismatch.", appName, appVersion)
		}
		expectedMsg := "Success: Application 'myNewApp' added with version '1.0.0'."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("Expected success message '%s', got '%s'", expectedMsg, output)
		}
	})

	t.Run("UpdateExistingApplication", func(t *testing.T) {
		os.Remove(testFile)
		appName := "myExistingApp"
		initialVersion := "1.0.0"
		updatedVersion := "1.0.1"
		initialConfig := Config{appName: initialVersion}
		if err := saveConfig(initialConfig); err != nil {
			t.Fatalf("Failed to set up initial config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(func() {
			handleAddCmd(appName, updatedVersion)
		}))
		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if _, ok := cfg[appName]; !ok || cfg[appName] != updatedVersion {
			t.Errorf("Expected app %s to be updated to version %s, but was not.", appName, updatedVersion)
		}
		expectedMsg := "Success: Application 'myExistingApp' updated from version '1.0.0' to '1.0.1'."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("Expected success message '%s', got '%s'", expectedMsg, output)
		}
	})
}

// TestHandleListCommand tests the list command functionality including output.
func TestHandleListCommand(t *testing.T) {
	originalConfigFileValue := configFile
	testFile := "test_list_versions.toml"
	configFile = testFile
	defer func() {
		configFile = originalConfigFileValue
		os.Remove(testFile)
	}()

	t.Run("ListWithMultipleApplications", func(t *testing.T) {
		os.Remove(testFile)
		// appX and appY to test sorting; config map iteration order is not guaranteed
		initialConfig := Config{"appY": "4.5.6", "appX": "1.2.3"}
		if err := saveConfig(initialConfig); err != nil {
			t.Fatalf("Failed to set up initial config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(handleListCmd))

		if !strings.Contains(output, "== Managed Applications ==") {
			t.Errorf("Output does not contain header. Got:\n%s", output)
		}
		// Order should be appX then appY due to sorting in handleListCmd
		expectedOrder := "  - Application: appX, Version: 1.2.3\n  - Application: appY, Version: 4.5.6"
		if !strings.Contains(output, expectedOrder) {
			t.Errorf("Expected applications in sorted order. Got:\n%s\nExpected to contain:\n%s", output, expectedOrder)
		}
	})

	t.Run("ListWithNoApplications", func(t *testing.T) {
		os.Remove(testFile)
		emptyConfig := Config{}
		if err := saveConfig(emptyConfig); err != nil {
			t.Fatalf("Failed to save empty config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(handleListCmd))
		expectedMsg := "Info: No applications currently managed. Use the 'add' command to add some."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("Expected message '%s', got '%s'", expectedMsg, output)
		}
	})

	t.Run("ListWithOneApplication", func(t *testing.T) {
		os.Remove(testFile)
		initialConfig := Config{"singleApp": "0.0.1"}
		if err := saveConfig(initialConfig); err != nil {
			t.Fatalf("Failed to set up initial config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(handleListCmd))
		if !strings.Contains(output, "== Managed Applications ==") {
			t.Errorf("Output does not contain header. Got:\n%s", output)
		}
		if !strings.Contains(output, "  - Application: singleApp, Version: 0.0.1") {
			t.Errorf("Output does not contain singleApp details. Got:\n%s", output)
		}
	})
}

// TestHandleRemoveCommand tests the remove command functionality.
func TestHandleRemoveCommand(t *testing.T) {
	originalConfigFileValue := configFile
	testFile := "test_remove_versions.toml"
	configFile = testFile
	defer func() {
		configFile = originalConfigFileValue
		os.Remove(testFile)
	}()

	t.Run("RemoveExistingApplication", func(t *testing.T) {
		os.Remove(testFile)
		appName := "appToRemove"
		initialConfig := Config{appName: "1.0.0", "anotherApp": "2.0.0"}
		if err := saveConfig(initialConfig); err != nil {
			t.Fatalf("Failed to set up initial config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(func() { handleRemoveCmd(appName) }))
		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		if _, exists := cfg[appName]; exists {
			t.Errorf("Application '%s' should have been removed.", appName)
		}
		if !strings.Contains(output, "Success: Application 'appToRemove' removed.") {
			t.Errorf("Expected success message not found. Got: %s", output)
		}
	})

	t.Run("RemoveNonExistentApplication", func(t *testing.T) {
		os.Remove(testFile)
		initialConfig := Config{"appThatExists": "1.0.0"}
		if err := saveConfig(initialConfig); err != nil {
			t.Fatalf("Failed to set up initial config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(func() { handleRemoveCmd("ghostApp") }))
		expectedMsg := "Info: Application 'ghostApp' not found in configuration. Nothing to remove."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("Expected info message not found. Got: %s", output)
		}
	})
}

// TestHandleCheckCommand tests the check command functionality.
func TestHandleCheckCommand(t *testing.T) {
	originalConfigFile := configFile
	testFile := "test_check_versions.toml"
	configFile = testFile

	originalGetLatestVersionFunc := getLatestVersion // Save original
	defer func() {
		configFile = originalConfigFile
		os.Remove(testFile)
		getLatestVersion = originalGetLatestVersionFunc // Restore original
	}()

	mockResponses := make(map[string]struct {
		version string
		err     error
	})

	// Setup mock for getLatestVersion
	getLatestVersion = func(appIdentifier string, apiBaseURL string) (string, error) {
		// apiBaseURL is ignored in this mock as we're not making real HTTP calls
		if resp, ok := mockResponses[appIdentifier]; ok {
			return resp.version, resp.err
		}
		return "", fmt.Errorf("unexpected appIdentifier '%s' in mock for getLatestVersion", appIdentifier)
	}

	t.Run("CheckSpecificApp_UpToDate", func(t *testing.T) {
		os.Remove(testFile)
		appName := "owner/app1"
		currentVersion := "1.0.0"
		if err := saveConfig(Config{appName: currentVersion}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		mockResponses[appName] = struct {version string; err error}{version: "1.0.0", err: nil}

		output := stripAnsiCodes(captureOutput(func() { handleCheckCmd(appName) }))

		if !strings.Contains(output, "Checking owner/app1...") || !strings.Contains(output, "Current: 1.0.0, Latest: 1.0.0 (Up to date)") {
			t.Errorf("Expected 'Up to date' message. Got: %s", output)
		}
	})

	t.Run("CheckSpecificApp_UpdateAvailable", func(t *testing.T) {
		os.Remove(testFile)
		appName := "owner/app2"
		if err := saveConfig(Config{appName: "1.0.0"}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		mockResponses[appName] = struct {version string; err error}{version: "1.1.0", err: nil}

		output := stripAnsiCodes(captureOutput(func() { handleCheckCmd(appName) }))
		if !strings.Contains(output, "Current: 1.0.0, Latest: 1.1.0 (Update Available!)") {
			t.Errorf("Expected 'Update Available!' message. Got: %s", output)
		}
	})

	t.Run("CheckSpecificApp_VersionDiscrepancy", func(t *testing.T) {
		os.Remove(testFile)
		appName := "owner/app3"
		if err := saveConfig(Config{appName: "1.2.0"}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		mockResponses[appName] = struct {version string; err error}{version: "1.1.0", err: nil} // Latest is older

		output := stripAnsiCodes(captureOutput(func() { handleCheckCmd(appName) }))
		if !strings.Contains(output, "Current: 1.2.0, Latest: 1.1.0 (Version discrepancy)") {
			t.Errorf("Expected 'Version discrepancy' message. Got: %s", output)
		}
	})

	t.Run("CheckSpecificApp_FetchError", func(t *testing.T) {
		os.Remove(testFile)
		appName := "owner/app4"
		if err := saveConfig(Config{appName: "1.0.0"}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		mockResponses[appName] = struct {version string; err error}{version: "", err: fmt.Errorf("mock network error")}

		// Capture stderr for error messages
		oldStderr := os.Stderr
		r, w, pipeErr := os.Pipe()
		if pipeErr != nil { t.Fatalf("Failed to create pipe: %v", pipeErr) }
		os.Stderr = w

		// Output from successful print before error still goes to stdout capture
		stdoutOutput := stripAnsiCodes(captureOutput(func() { handleCheckCmd(appName) }))

		w.Close()
		errOutputBytes, _ := io.ReadAll(r)
		os.Stderr = oldStderr // Restore stderr
		errOutput := stripAnsiCodes(string(errOutputBytes))

		if !strings.Contains(stdoutOutput, "Checking owner/app4...") {
             t.Errorf("Expected 'Checking owner/app4...' on stdout. Got: %s", stdoutOutput)
        }
		if !strings.Contains(errOutput, "Error: Failed to check owner/app4: mock network error") {
			t.Errorf("Expected fetch error message on stderr. Got: %s", errOutput)
		}
	})

	t.Run("CheckSpecificApp_NotFoundInConfig", func(t *testing.T) {
		os.Remove(testFile)
		if err := saveConfig(Config{"owner/appExists": "1.0.0"}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}

		oldStderr := os.Stderr
		r, w, pipeErr := os.Pipe()
		if pipeErr != nil { t.Fatalf("Failed to create pipe: %v", pipeErr) }
		os.Stderr = w

		handleCheckCmd("owner/nonExistentApp") // This function's output is what we're testing

		w.Close()
		errOutputBytes, _ := io.ReadAll(r)
		os.Stderr = oldStderr
		errOutput := stripAnsiCodes(string(errOutputBytes))

		expectedMsg := "Error: Application 'owner/nonExistentApp' not found in your managed list."
		if !strings.Contains(errOutput, expectedMsg) {
			t.Errorf("Expected 'not found in config' error message on stderr. Got: %s", errOutput)
		}
	})

	t.Run("CheckSpecificApp_InvalidFormat", func(t *testing.T) {
		os.Remove(testFile)
		appNameInvalid := "invalidAppFormat"
		if err := saveConfig(Config{appNameInvalid: "1.0.0"}); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		output := stripAnsiCodes(captureOutput(func() { handleCheckCmd(appNameInvalid) }))
		expectedMsg := "Info: Skipping invalidAppFormat: Not in 'owner/repo' format. Cannot check for updates via GitHub."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("Expected 'invalid format' message. Got: %s", output)
		}
	})

	t.Run("CheckAllApps", func(t *testing.T) {
		os.Remove(testFile)
		config := Config{
			"owner/appA": "1.0.0",      // Up to date
			"owner/appB": "0.5.0",      // Update available
			"invalidAppC": "2.0.0", // Invalid format (no slash)
		}
		if err := saveConfig(config); err != nil {
			t.Fatalf("Failed to save config: %v", err)
		}
		mockResponses["owner/appA"] = struct {version string; err error}{version: "1.0.0", err: nil}
		mockResponses["owner/appB"] = struct {version string; err error}{version: "1.0.0", err: nil}
		// invalidAppC won't call getLatestVersion

		output := stripAnsiCodes(captureOutput(func() { handleCheckCmd("") })) // Empty string for specificApp means check all

		if !strings.Contains(output, "Checking all managed applications for updates...") {
			t.Errorf("Expected 'Checking all' message. Got: %s", output)
		}
		// Order of appA and appB is not guaranteed by map iteration, so check individually
		if !strings.Contains(output, "Checking owner/appA...") || !strings.Contains(output, "Current: 1.0.0, Latest: 1.0.0 (Up to date)") {
			t.Errorf("Expected appA up to date message not found. Got: %s", output)
		}
		if !strings.Contains(output, "Checking owner/appB...") || !strings.Contains(output, "Current: 0.5.0, Latest: 1.0.0 (Update Available!)") {
			t.Errorf("Expected appB update available message not found. Got: %s", output)
		}
		if !strings.Contains(output, "Info: Skipping invalidAppC: Not in 'owner/repo' format.") {
			t.Errorf("Expected invalidAppC skipping message not found. Got: %s", output)
		}
	})

	t.Run("CheckEmptyConfig", func(t *testing.T) {
		os.Remove(testFile)
		if err := saveConfig(Config{}); err != nil { // Empty config
			t.Fatalf("Failed to save empty config: %v", err)
		}
		rawOutput := captureOutput(func() { handleCheckCmd("") })
		output := strings.TrimSpace(stripAnsiCodes(rawOutput))
		// Setting expectedMsg from the literal "Got" string from the last test failure log
		expectedMsg := "Info: No applications currently managed. Use 'add' command to add some."
		if output != expectedMsg {
			t.Errorf("Expected message mismatch.\nGot     : [%s]\nExpected: [%s]", output, expectedMsg)
		}
	})
}
