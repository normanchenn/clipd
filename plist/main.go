package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
	"text/template"
)

type Config struct {
	Label              string
	ExecutablePath     string
	StandardOutputPath string
	StandardErrorPath  string
	Debug              bool
}

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user: ", err)
	}
	baseDir := user.HomeDir

	defaultConfig := Config{
		Label:              "clipd.plist",
		ExecutablePath:     baseDir + "/go/bin/daemon",
		StandardOutputPath: baseDir + "/clipd/logs/clipd-output.log",
		StandardErrorPath:  baseDir + "/clipd/logs/clipd-error.log",
		Debug:              true,
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
        <string>{{.ExecutablePath}}</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutputPath</key>
    <string>{{.StandardOutputPath}}</string>
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

	if err := t.Execute(file, config); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Plist file generated successfully.")
}

func promptConfig(defaults Config) Config {
	label := prompt("Label (default: "+defaults.Label+"): ", defaults.Label)

	programArgs := prompt("Executable Path (default: "+defaults.ExecutablePath+"): ", defaults.ExecutablePath)

	stdOutPath := prompt("Standard Output Path (default: "+defaults.StandardOutputPath+"): ", defaults.StandardOutputPath)

	stdErrPath := prompt("Standard Error Path (default: "+defaults.StandardErrorPath+"): ", defaults.StandardErrorPath)

	debugStr := prompt("Debug (true/false) (default: "+fmt.Sprintf("%t", defaults.Debug)+"): ", fmt.Sprintf("%t", defaults.Debug))
	debug := strings.ToLower(debugStr) == "true"

	return Config{
		Label:              label,
		ExecutablePath:     programArgs,
		StandardOutputPath: stdOutPath,
		StandardErrorPath:  stdErrPath,
		Debug:              debug,
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
