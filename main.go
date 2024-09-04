package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	icon        = flag.String("icon", "", "Location of the icon file")
	comment     = flag.String("comment", "", "Comment/tooltip for the application")
	categories  = flag.String("categories", "Application", "Categories (comma-separated)")
	terminal    = flag.Bool("terminal", false, "Whether the app requires a terminal")
	systemWide  = flag.Bool("system", false, "Install system-wide (requires sudo)")
	genericName = flag.String("generic", "", "Generic name of the application")
	help        = flag.Bool("help", false, "Display help information")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <name> <exec>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(icon, "i", "", "Location of the icon file")
	flag.StringVar(comment, "c", "", "Comment/tooltip for the application")
	flag.StringVar(categories, "C", "Application", "Categories (comma-separated)")
	flag.BoolVar(terminal, "t", false, "Whether the app requires a terminal")
	flag.BoolVar(systemWide, "s", false, "Install system-wide (requires sudo)")
	flag.StringVar(genericName, "g", "", "Generic name of the application")
	flag.BoolVar(help, "h", false, "Display help information")
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Check for positional arguments
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Error: Name and Exec are required as positional arguments")
		flag.Usage()
		os.Exit(1)
	}

	name := args[0]
	exec := args[1]

	if *genericName == "" {
		*genericName = generateGenericName(name)
	}

	desktopEntry := fmt.Sprintf(`[Desktop Entry]
Encoding=UTF-8
Version=1.0
Name=%s
GenericName=%s
Exec=%s
Terminal=%t
Icon=%s
Type=Application
Categories=%s
Comment=%s
`, name, *genericName, exec, *terminal, *icon, *categories, *comment)

	installPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")
	if *systemWide {
		installPath = "/usr/share/applications"
	}

	filename := strings.ToLower(name) + ".desktop"
	fullPath := filepath.Join(installPath, filename)

	err := os.WriteFile(fullPath, []byte(desktopEntry), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Desktop entry created: %s\n", fullPath)
}

func generateGenericName(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	words := strings.Fields(name)
	caser := cases.Title(language.Und)
	for i, word := range words {
		words[i] = caser.String(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}
