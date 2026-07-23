package browser

import (
	"os/exec"
	"runtime"
)

// Open launches the user's default browser at url. Non-blocking.
func Open(url string) error {
	var name string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		name, args = "open", []string{url}
	case "windows":
		name, args = "cmd", []string{"/c", "start", url}
	default:
		name, args = "xdg-open", []string{url}
	}
	return exec.Command(name, args...).Start()
}
