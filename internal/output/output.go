package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/tidwall/pretty"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

type Format string

const (
	FormatHuman Format = "human"
	FormatJSON  Format = "json"
)

var Current Format = FormatHuman

func SetFormat(f string) {
	if f == "json" {
		Current = FormatJSON
	} else {
		Current = FormatHuman
	}
}

func PrintJSON(data []byte) {
	if Current == FormatJSON {
		fmt.Println(string(data))
		return
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println(string(data))
		return
	}
	formatted := pretty.Color(pretty.Pretty(data), nil)
	fmt.Print(string(formatted))
}

func PrintSuccess(msg string) {
	if Current == FormatJSON {
		b, _ := json.Marshal(map[string]string{"status": "ok", "message": msg})
		fmt.Println(string(b))
		return
	}
	fmt.Println(green("✓ ") + msg)
}

func PrintError(msg string) {
	if Current == FormatJSON {
		b, _ := json.Marshal(map[string]string{"status": "error", "message": msg})
		fmt.Println(string(b))
		return
	}
	fmt.Fprintln(os.Stderr, red("✗ ")+msg)
}

func PrintInfo(msg string) {
	if Current == FormatJSON {
		return
	}
	fmt.Println(cyan("→ ") + msg)
}

func PrintHeader(msg string) {
	if Current == FormatJSON {
		return
	}
	fmt.Println(bold(msg))
}

func PrintWarning(msg string) {
	if Current == FormatJSON {
		return
	}
	fmt.Println(yellow("⚠ ") + msg)
}

func PrintAPIError(status int, data []byte) {
	if Current == FormatJSON {
		fmt.Println(string(data))
		return
	}
	fmt.Fprintf(os.Stderr, "%s HTTP %d: %s\n", red("✗"), status, string(data))
}
