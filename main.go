package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

func rgbToHex(R int, G int, B int) string {
	return fmt.Sprintf("#%02x%02x%02x", R, G, B)
}

func lumToHex(l int) string {
	return rgbToHex(l, l, l)
}

func randRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func randRangeS(min int, max int, seed int64) int {
	rand.Seed(seed)
	return rand.Intn(max-min+1) + min
}

type Color struct {
	r int
	g int
	b int
}

type model struct {
	screenWidth  int
	screenHeight int

	imgWidth  int
	imgHeight int
	img       [][]Color

	currentColor Color
	gridOn       bool

	statusBGColor Color
	statusFGColor Color
	statusMessage string

	ev       tea.MouseMsg
	relPos   [2]int
	inCanvas bool
}

func initialModel() model {
	var img [][]Color
	width := 16
	height := 16
	for y := 0; y < height; y++ {
		var row []Color
		for _x := 0; _x < width; _x++ {
			row = append(row, Color{0, 0, 0})
		}
		img = append(img, row)
	}
	return model{
		imgWidth:  width,
		imgHeight: height,
		img:       img,
		gridOn:    false,

		currentColor: Color{255, 0, 0},

		statusBGColor: Color{50, 30, 250},
		statusFGColor: Color{0, 0, 0},
		statusMessage: "Welcome to the Terminal Painter!",
	}
}

func (m model) Init() tea.Cmd {
	m.img[0][0] = Color{255, 255, 255}
	m.img[1][0] = Color{0, 0, 255}
	m.img[0][1] = Color{255, 0, 0}
	m.img[1][1] = Color{255, 0, 255}
	m.img[1][2] = Color{0, 255, 0}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmd = tea.Quit
		case "1":
			m.currentColor = Color{255, 0, 0}
		case "2":
			m.currentColor = Color{0, 0, 0}
		case "ctrl+s":
			save("out.png", m.img)
			m.statusBGColor = Color{20, 250, 50}
			m.statusMessage = "Saved to out.png!"
		}
	case tea.MouseMsg:
		m.ev = msg
		if zone.Get("paintScreen").InBounds(msg) {
			x, y := zone.Get("paintScreen").Pos(msg)
			m.inCanvas = true
			m.relPos = [2]int{x, y}
		} else {
			m.inCanvas = false
		}

		if msg.Type == tea.MouseLeft {
			if zone.Get("paintScreen").InBounds(msg) {
				x, y := zone.Get("paintScreen").Pos(msg)
				if !(x < 0 || x/2 > m.imgWidth || y < 0 || y > m.imgHeight) {
					m.img[y][x/2] = m.currentColor
				}
			} else if zone.Get("palette").InBounds(msg) {

			}
			// fmt.Printf("Clicked %v\n", msg)
			// fmt.Printf("PaintScreen %v %v %v %v\n",
			// 	zone.Get("paintScreen").StartX,
			// 	zone.Get("paintScreen").StartY,
			// 	zone.Get("paintScreen").EndX,
			// 	zone.Get("paintScreen").EndY,
			// )
			// fmt.Printf("Pos: %v %v\n", x, y)
		}
	}

	return m, cmd
}

func (m model) View() string {
	s := ""

	// Status line
	statusLine := lipgloss.NewStyle().
		Background(lipgloss.Color(rgbToHex(
			m.statusBGColor.r,
			m.statusBGColor.g,
			m.statusBGColor.b,
		))).
		Foreground(lipgloss.Color(rgbToHex(
			m.statusFGColor.r,
			m.statusFGColor.g,
			m.statusFGColor.b,
		))).
		Width(m.screenWidth).
		Padding(0, 3).
		Margin(1, 0).
		Render(m.statusMessage)

	// s += "\n"

	paintScreen := ""

	for y := 0; y < m.imgHeight; y++ {
		for x := 0; x < m.imgWidth*2; x++ {
			color1 := m.img[y][x/2]

			str := " "
			if m.gridOn {
				if x%2 == 0 {
					str = "⎣"
				} else {
					str = "⎯"
				}
			}

			paintScreen += lipgloss.NewStyle().
				Background(lipgloss.Color(rgbToHex(color1.r, color1.g, color1.b))).
				Foreground(lipgloss.Color("#ffffff")).
				// Render("█")
				Render(str)
			// Render("#")
		}
		if y != m.imgHeight-1 {
			paintScreen += "\n"
		}
	}

	// Add Paint Screen
	s += lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(zone.Mark("paintScreen", paintScreen))

	// Debug info
	s += fmt.Sprintf("\nMouse event: %v\n", m.ev)
	if m.inCanvas {
		s += fmt.Sprintf("x: %v, y: %v\n", m.relPos[0], m.relPos[1])
	} else {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("196" /* red */)).Render("Out of canvas!")
	}

	// Palette
	s = lipgloss.JoinHorizontal(lipgloss.Left, s, lipgloss.Place(10, m.imgHeight+2, lipgloss.Center, lipgloss.Center, renderPalette()))

	// Status line
	s = lipgloss.JoinVertical(lipgloss.Top, lipgloss.Place(0, 1, 0, 0, statusLine), s)

	return zone.Scan(s)
}

func main() {
	m := initialModel()
	if len(os.Args) == 2 {
		_, err := os.Stat(os.Args[1])
		if err == fs.ErrNotExist {

		} else if err != nil {
			panic(err)
		} else {
			m.img, m.imgWidth, m.imgHeight = load(os.Args[1])
		}
	}

	zone.NewGlobal()
	if err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Start(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
