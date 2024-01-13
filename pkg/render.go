package pkg

import (
	"fmt"
	"os"
	"strings"


	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
	
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type Uper struct {
	By float64
}

type ProgressModel struct {
	// function that ends the main request function instead
	EndReqMaker func() 

	// listens for update to increment 
	Channel  chan Uper

	// More information on the top of prores bar 
	Info     string

	// the default progress bar model
	progress    progress.Model
}


// create a new Progress Model
func NewProgressModel() ProgressModel {
	return ProgressModel{
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

// standard procedure for bubble tea framework, 
// Basicaly it initiate the first listener on `m.Channel``
func (m ProgressModel) Init() tea.Cmd {
	return m.Listen
}

// Start the whole terminal takeover :)
func (m ProgressModel) Run() {

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}


// This keeps the ball running. I.e it listens on the channel for update and send out a message
func (m ProgressModel) Listen()  tea.Msg {
	
	return <-m.Channel
}


// standard procedure for bubble tea framework, 
// it make the update to the model
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Uper:

		cmd := m.progress.IncrPercent(msg.By)		
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}
		return m, tea.Batch(m.Listen, cmd)		

	case tea.KeyMsg:
		m.EndReqMaker()
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

// standard procedure for bubble tea framework, 
// it formats the view
func (m ProgressModel) View() string {
	pad := strings.Repeat(" ", padding)
	return pad + m.Info +
		pad + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}


type PrintModel struct {
	Text string
}


func (p PrintModel) Init() tea.Cmd {
	return func() tea.Msg {return 0}
}

// Start the whole terminal takeover :)
func (p PrintModel) Run() {

	if _, err := tea.NewProgram(p).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (p PrintModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p, tea.Quit
}

func (p PrintModel) View() string {
	l := lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0")).Render
	return l(p.Text)
}


// generate table that displays results
func TableGen(head map[string]string, body map[string]string) string {

	table := ""

	longestLen := 0
	lonegestKey := 0
	lonegestvalue := 0

	for key, value := range body {		
		thelen := len(key) + len(value)

		if len(key) > lonegestKey {
			lonegestKey = len(key)
		}

		if len(value) > lonegestvalue {
			lonegestvalue = len(value)
		}

		if thelen > longestLen {
			longestLen = thelen
		}

	}
	

	hr := "+--" // This sets demacation between fields, +--+---+

	hr += strings.Repeat("-", lonegestKey)
	hr += "+--"
	hr += strings.Repeat("-", lonegestvalue)
	hr += "+\n"

	bodier := ""

	header := ""

	for key, value := range head {
		
		header += "| "
		header +=  key + strings.Repeat(" ", lonegestKey - len(key)) + " | "
		header +=  value +  strings.Repeat(" ", lonegestvalue - len(value)) + " |\n"
	}

	for key, value := range body {
		bodier += "| "  + key + strings.Repeat(" ", lonegestKey - len(key)) + " | " +  value + strings.Repeat(" ", lonegestvalue - len(value))  + " |\n"
	}

	table = "\n\n" + hr + header + hr + bodier + hr

	return table
}

