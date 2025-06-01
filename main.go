package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Old color constants are removed from here, will use ui.go

func main() {
	// Define common flag sets for subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)

	// Custom usage for subcommands to ensure they are displayed correctly
	addCmd.Usage = func() {
		PrintUsageMessage("Usage: %s add <application_name> <version>", os.Args[0])
		PrintUsageMessage("Example: %s add myapp 1.0.2", Colorize(os.Args[0], colorCyanFg))
	}
	removeCmd.Usage = func() {
		PrintUsageMessage("Usage: %s remove <application_name>", os.Args[0])
		PrintUsageMessage("Example: %s remove myapp", Colorize(os.Args[0], colorCyanFg))
	}
	listCmd.Usage = func() {
		PrintUsageMessage("Usage: %s list", os.Args[0])
	}
	checkCmd.Usage = func() {
		PrintUsageMessage("Usage: %s check [<application_name>]", os.Args[0])
		PrintUsageMessage("Example: %s check myapp", Colorize(os.Args[0], colorCyanFg))
		PrintUsageMessage("Example: %s check", Colorize(os.Args[0], colorCyanFg))
	}

	if len(os.Args) < 2 {
		printOverallUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if len(addCmd.Args()) < 2 {
			PrintError("Missing application name and/or version for 'add' command.")
			addCmd.Usage()
			os.Exit(1)
		}
		appName := addCmd.Args()[0]
		appVersion := addCmd.Args()[1]
		handleAddCmd(appName, appVersion)
	case "remove":
		removeCmd.Parse(os.Args[2:])
		if len(removeCmd.Args()) < 1 {
			PrintError("Missing application name for 'remove' command.")
			removeCmd.Usage()
			os.Exit(1)
		}
		appName := removeCmd.Args()[0]
		handleRemoveCmd(appName)
	case "list":
		listCmd.Parse(os.Args[2:])
		if len(listCmd.Args()) > 0 {
			PrintError("'list' command does not take any arguments.")
			listCmd.Usage()
			os.Exit(1)
		}
		handleListCmd()
	case "check":
		checkCmd.Parse(os.Args[2:])
		specificApp := ""
		if len(checkCmd.Args()) > 1 { // check can have 0 or 1 arg
			PrintError("'check' command accepts at most one application name.")
			checkCmd.Usage()
			os.Exit(1)
		}
		if len(checkCmd.Args()) == 1 {
			specificApp = checkCmd.Args()[0]
		}
		handleCheckCmd(specificApp)
	default:
		PrintError("Unknown command '%s'.", os.Args[1])
		printOverallUsage()
		os.Exit(1)
	}
}

func printOverallUsage() {
	PrintUsageMessage("Usage: %s <command> [arguments]", os.Args[0])
	PrintUsageMessage("Commands:")
	// Using PrintMessage for command descriptions to allow custom coloring within them
	// Colorize parts of the string for more detailed control if needed.
	PrintMessage("  %s %s\tAdd a new application to monitor", Colorize("add", colorGreenFg), Colorize("<name> <version>", colorFgDefault))
	PrintMessage("  %s %s\tRemove an application from monitoring", Colorize("remove", colorRedFg), Colorize("<name>", colorFgDefault))
	PrintMessage("  %s\t\t\tList all monitored applications", Colorize("list", colorBlueFg))
	PrintMessage("  %s %s\tCheck for updates, optionally for a specific app", Colorize("check", colorYellowFg), Colorize("[<name>]", colorFgDefault))
	PrintUsageMessage("\nUse \"%s <command> --help\" for more information about a command (not yet implemented).", os.Args[0])
}

func handleAddCmd(appName string, appVersion string) {
	config, err := loadConfig()
	if err != nil {
		PrintError("Could not load configuration: %v", err)
		return
	}

	oldVersion, exists := config[appName]
	config[appName] = appVersion

	err = saveConfig(config)
	if err != nil {
		PrintError("Could not save configuration for '%s': %v", appName, err)
		return
	}

	if exists {
		PrintSuccess("Application '%s' updated from version '%s' to '%s'.",
			Colorize(appName, colorYellowFg),
			Colorize(oldVersion, colorMagentaFg),
			Colorize(appVersion, colorCyanFg))
	} else {
		PrintSuccess("Application '%s' added with version '%s'.",
			Colorize(appName, colorYellowFg),
			Colorize(appVersion, colorCyanFg))
	}
}

func handleRemoveCmd(appName string) {
	config, err := loadConfig()
	if err != nil {
		PrintError("Could not load configuration: %v", err)
		return
	}

	if _, exists := config[appName]; !exists {
		PrintInfo("Application '%s' not found in configuration. Nothing to remove.", Colorize(appName, colorMagentaFg))
		return
	}

	delete(config, appName)

	err = saveConfig(config)
	if err != nil {
		PrintError("Could not save configuration after removing '%s': %v", appName, err)
		return
	}
	PrintSuccess("Application '%s' removed.", Colorize(appName, colorYellowFg))
}

func handleListCmd() {
	config, err := loadConfig()
	if err != nil {
		PrintError("Could not load configuration: %v", err)
		return
	}

	if len(config) == 0 {
		PrintInfo("No applications currently managed. Use the 'add' command to add some.")
		return
	}

	PrintHeader("Managed Applications")

	// Sort keys for consistent output order
	var sortedAppNames []string
	for appName := range config {
		sortedAppNames = append(sortedAppNames, appName)
	}
	sort.Strings(sortedAppNames)

	for _, appName := range sortedAppNames {
		appVersion := config[appName]
		PrintMessage("  - Application: %s, Version: %s",
			Colorize(appName, colorYellowFg),
			Colorize(appVersion, colorCyanFg))
	}
}

func handleCheckCmd(specificApp string) {
	config, err := loadConfig()
	if err != nil {
		PrintError("Could not load configuration: %v", err)
		return
	}

	if len(config) == 0 {
		PrintInfo("No applications currently managed. Use 'add' command to add some.")
		return
	}

	if specificApp != "" {
		currentVersion, exists := config[specificApp]
		if !exists {
			PrintError("Application '%s' not found in your managed list.", Colorize(specificApp, colorYellowFg))
			return
		}
		checkAppVersion(specificApp, currentVersion)
	} else {
		PrintMessage("%sChecking all managed applications for updates...%s", colorBlueFg, colorReset) // Using PrintMessage for specific coloring
		for appName, currentVersion := range config {
			checkAppVersion(appName, currentVersion)
		}
	}
}

func checkAppVersion(appName, currentVersion string) {
	if !strings.Contains(appName, "/") {
		PrintInfo("Skipping %s: Not in 'owner/repo' format. Cannot check for updates via GitHub.", Colorize(appName, colorMagentaFg))
		return
	}

	// Using PrintMessage directly for more control over the line ending and formatting
	fmt.Printf("%sChecking %s... %s", colorFgDefault, Colorize(appName, colorYellowFg), colorReset)
	latestVersion, err := getLatestVersion(appName, "") // Call the func variable
	if err != nil {
		// PrintError already adds a newline.
		// Need to ensure the "Checking..." line gets a newline if an error occurs here.
		fmt.Println() // Add newline after "Checking..." before printing error
		PrintError("Failed to check %s: %v", Colorize(appName, colorMagentaFg), err)
		return
	}

	if latestVersion == currentVersion {
		fmt.Printf("%s Current: %s, Latest: %s (%s)%s\n",
			colorFgDefault,
			Colorize(currentVersion, colorCyanFg),
			Colorize(latestVersion, colorGreenFg),
			Colorize("Up to date", colorGreenFg),
			colorReset)
	} else {
		isUpgrade := latestVersion > currentVersion // Lexicographical comparison
		if isUpgrade {
			fmt.Printf("%s Current: %s, Latest: %s (%s)%s\n",
				colorFgDefault,
				Colorize(currentVersion, colorCyanFg),
				Colorize(latestVersion, colorRedFg),
				Colorize("Update Available!", colorRedFg),
				colorReset)
		} else {
			fmt.Printf("%s Current: %s, Latest: %s (%s)%s\n",
				colorFgDefault,
				Colorize(currentVersion, colorCyanFg),
				Colorize(latestVersion, colorYellowFg),
				Colorize("Version discrepancy", colorYellowFg),
				colorReset)
		}
	}
}
