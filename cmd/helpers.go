package cmd

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
)

var (
	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3"))

	normalStyle = lipgloss.NewStyle()
	blueStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))
	redStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))
)

func openInEditor(taskFile string) error {
	nvim := exec.Command("nvim", taskFile)
	nvim.Stdin = os.Stdin
	nvim.Stdout = os.Stdout
	nvim.Stderr = os.Stderr

	return nvim.Run()
}
