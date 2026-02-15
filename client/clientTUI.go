package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- 1. Lip Gloss Styling (Defined once, used everywhere) ---
var (
	primaryColor = lipgloss.Color("#FF7CCB") // Pinkish
	peerColor    = lipgloss.Color("#04B575") // Green
	myColor      = lipgloss.Color("#00D2FF") // Blue

	appStyle = lipgloss.NewStyle().Margin(1, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true)

	viewportStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	peerPrefixStyle = lipgloss.NewStyle().Foreground(peerColor).Bold(true)
	myPrefixStyle   = lipgloss.NewStyle().Foreground(myColor).Bold(true)
	systemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Italic(true)
)

type chatMsg string
type errMsg error

type model struct {
	viewport viewport.Model
	textarea textarea.Model
	messages []string
	conn     net.Conn
	reader   *bufio.Reader
	err      error
}

func initialModel(conn net.Conn) model {
	ta := textarea.New()
	ta.Placeholder = "Type a message..."
	ta.Focus()
	ta.Prompt = "💬 "
	// ta.CharLimit = 280   // <-- i am evil. Bufferoverflow
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(0, 0) // Dimensions set dynamically in Update()
	vp.SetContent(systemStyle.Render("Connected to server! Start typing..."))

	return model{
		textarea: ta,
		viewport: vp,
		messages: []string{},
		conn:     conn,
		reader:   bufio.NewReader(conn),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, receiveMessage(m.reader))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// Pass the message to our Bubbles
	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {

	// Let Bubbles handle terminal resizing!
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(headerStyle.Render("TCP Chat Server"))
		footerHeight := lipgloss.Height(m.textarea.View())

		// Calculate available height (Total - Header - Footer - Margins)
		availableHeight := msg.Height - headerHeight - footerHeight - 6

		m.viewport.Width = msg.Width - 6
		m.viewport.Height = availableHeight
		m.textarea.SetWidth(msg.Width - 6)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := strings.TrimSpace(m.textarea.Value())
			if input != "" {
				fmt.Fprintln(m.conn, input) // Send to server

				// Format our own message with Lip Gloss
				formattedMsg := myPrefixStyle.Render("You: ") + input
				m.messages = append(m.messages, formattedMsg)

				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
				m.textarea.Reset()
			}
		}

	case chatMsg:
		// Format peer messages with Lip Gloss
		rawMsg := string(msg)
		var formattedMsg string

		if strings.Contains(rawMsg, "joined") || strings.Contains(rawMsg, "left") {
			formattedMsg = systemStyle.Render(rawMsg)
		} else {
			formattedMsg = peerPrefixStyle.Render("Peer: ") + rawMsg
		}

		m.messages = append(m.messages, formattedMsg)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, receiveMessage(m.reader)

	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nConnection error: %v\n", m.err)
	}

	// 3. Simplified Lip Gloss Layout
	// We stack our components vertically, letting Lip Gloss handle the spacing.
	header := headerStyle.Render("🌐 Go TCP Chat")
	chatLog := viewportStyle.Render(m.viewport.View())
	inputBox := m.textarea.View()
	helpTxt := systemStyle.Render("(ctrl+c / esc to quit)")

	ui := lipgloss.JoinVertical(lipgloss.Left,
		header,
		chatLog,
		inputBox,
		helpTxt,
	)

	return appStyle.Render(ui)
}

func receiveMessage(reader *bufio.Reader) tea.Cmd {
	return func() tea.Msg {
		text, err := reader.ReadString('\n')
		if err != nil {
			return errMsg(err)
		}
		return chatMsg(strings.TrimSpace(text))
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run client.go <server_address:port>")
		return
	}

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := tea.NewProgram(initialModel(conn), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Fatal error starting TUI:", err)
		os.Exit(1)
	}
}
