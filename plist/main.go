package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	Label             string
	ProgramArguments  string
	StandardOutPath   string
	StandardErrorPath string
	Debug             bool
}

func main() {
	defaultConfig := Config{
		Label:            "clipd.plist",
		ProgramArguments: "/Users/normanchen/Documents/clipd/daemon/daemon",
		// "~/go/bin/daemon"
		StandardOutPath: "/Users/normanchen/Desktop/output2.log",
		// "/var/log/clipd/output.log"
		StandardErrorPath: "/Users/normanchen/Desktop/error2.log",
		// "/var/log/clipd/error.log"
		Debug: true,
	}

	config := promptConfig(defaultConfig)

	tmpl := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Label}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.ProgramArguments}}</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.StandardOutPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.StandardErrorPath}}</string>
    <key>Debug</key>
    <{{if .Debug}}true{{else}}false{{end}}/>
</dict>
</plist>`

	t := template.Must(template.New("plist").Parse(tmpl))

	file, err := os.Create("clipd.plist")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Execute the template with the config values and write to file
	if err := t.Execute(file, config); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Plist file generated successfully.")
}

func promptConfig(defaults Config) Config {
	label := prompt("Label (default: "+defaults.Label+"): ", defaults.Label)

	programArgs := prompt("Program Arguments (default: "+defaults.ProgramArguments+"): ", defaults.ProgramArguments)

	stdOutPath := prompt("Standard Out Path (default: "+defaults.StandardOutPath+"): ", defaults.StandardOutPath)

	stdErrPath := prompt("Standard Error Path (default: "+defaults.StandardErrorPath+"): ", defaults.StandardErrorPath)

	debugStr := prompt("Debug (true/false) (default: "+fmt.Sprintf("%t", defaults.Debug)+"): ", fmt.Sprintf("%t", defaults.Debug))
	debug := strings.ToLower(debugStr) == "true"

	return Config{
		Label:             label,
		ProgramArguments:  programArgs,
		StandardOutPath:   stdOutPath,
		StandardErrorPath: stdErrPath,
		Debug:             debug,
	}
}

func prompt(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}
