package main

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

var palette = []string{
	"124", "160", "196", // Red
	"130", "172", "214", // Orange
	"142", "184", "226", // Yellow
	"34", "76", "118", //   Green
	"39", "81", "123", //   Cyan
	"19", "20", "21", //    Blue
	"55", "56", "57", //    Purple
	"93", "135", "177", //  Magenta
	"", "", "",
}

func renderPalette() string {
	s := ""

	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(255, 0, 0))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(255, 255, 0))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(lumToHex(255))).Render("  ")
	// s += "\n"
	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(0, 255, 0))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(0, 255, 255))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(lumToHex(128))).Render("  ")
	// s += "\n"
	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(0, 0, 255))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(rgbToHex(255, 0, 255))).Render("  ")
	// s += lipgloss.NewStyle().Background(lipgloss.Color(lumToHex(0))).Render("  ")

	for i, color := range palette {
		s += lipgloss.NewStyle().Background(lipgloss.Color(color)).Render("  ")
		if (i+1)%3 == 0 {
			s += "\n"
		}
	}

	return s
}
