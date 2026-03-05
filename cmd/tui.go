package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7FFF00")).
			Bold(true).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Background(lipgloss.Color("#004400"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Bold(true)

	systemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))
)

type model struct {
	index       int
	choices     []string
	width       int
	height      int
	spinner     spinner.Model
	loading     bool
	showMenu    bool
	cpuPercent  float64
	ramPercent  float64
	currentTime string
	err         error
	executing   bool
	statusMsg   string
	lastCommand string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7FFF00"))
	return model{
		index:    0,
		choices:  []string{"⚙ Форматнуть Go", "🧪 Аудит фронта", "🚀 Деплой", "🔧 git-relasens", "📊 Системная инфо", "🌐 Проверка URL", "📁 Сканировать логи", "📜 Просмотр логов", "❌ Выйти"},
		spinner:  s,
		loading:  true,
		showMenu: false,
	}
}

type tickMsg time.Time
type systemInfoMsg struct {
	cpu float64
	ram float64
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
		loadingCmd(),
		updateSystemInfo(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadingCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second) // Simulate loading
		return loadingCompleteMsg{}
	}
}

func updateSystemInfo() tea.Cmd {
	return func() tea.Msg {
		cpuPercent, _ := cpu.Percent(time.Second, false)
		memInfo, _ := mem.VirtualMemory()
		var cpu float64
		if len(cpuPercent) > 0 {
			cpu = cpuPercent[0]
		}
		return systemInfoMsg{
			cpu: cpu,
			ram: memInfo.UsedPercent,
		}
	}
}

type loadingCompleteMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.showMenu && m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.showMenu && m.index < len(m.choices)-1 {
				m.index++
			}
		case "enter":
			if m.showMenu && !m.executing {
				m.executing = true
				m.statusMsg = ""
				m.lastCommand = m.choices[m.index]
				return m, tea.Batch(m.executeChoice(), m.spinner.Tick)
			}
		case "c": // Clear status message
			m.statusMsg = ""
			m.err = nil
		}

	case tickMsg:
		m.currentTime = time.Now().Format("15:04:05")
		return m, tea.Batch(tickCmd(), updateSystemInfo())

	case systemInfoMsg:
		m.cpuPercent = msg.cpu
		m.ramPercent = msg.ram
		return m, updateSystemInfo()

	case loadingCompleteMsg:
		m.loading = false
		m.showMenu = true
		return m, nil

	case errorMsg:
		m.executing = false
		m.statusMsg = fmt.Sprintf("Ошибка: %v", msg.err)
		return m, nil

	case successMsg:
		m.executing = false
		m.statusMsg = string(msg)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.loading {
		return m.loadingView()
	}
	return m.menuView()
}

func (m model) loadingView() string {
	logo := `
██╗    ██╗ █████╗ ██╗  ██╗ █████╗     ██╗      █████╗ ██╗   ██╗ ██████╗██╗██╗
██║    ██║██╔══██╗██║ ██╔╝██╔══██╗    ██║     ██╔══██╗██║   ██║██╔════╝██║██║
██║ █╗ ██║███████║█████╔╝ ███████║    ██║     ███████║██║   ██║██║     ██║██║
██║███╗██║██╔══██║██╔═██╗ ██╔══██║    ██║     ██╔══██║██║   ██║██║     ██║██║
╚███╔███╔╝██║  ██║██║  ██╗██║  ██║    ███████╗██║  ██║╚██████╔╝╚██████╗██║██║
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚══════╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝╚═╝
`

	content := fmt.Sprintf("%s\n\n%s Загрузка ядра...",
		titleStyle.Render(logo),
		m.spinner.View())

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func (m model) menuView() string {
	logo := `
██╗    ██╗ █████╗ ██╗  ██╗ █████╗     ██╗      █████╗ ██╗   ██╗ ██████╗██╗██╗
██║    ██║██╔══██╗██║ ██╔╝██╔══██╗    ██║     ██╔══██╗██║   ██║██╔════╝██║██║
██║ █╗ ██║███████║█████╔╝ ███████║    ██║     ███████║██║   ██║██║     ██║██║
██║███╗██║██╔══██║██╔═██╗ ██╔══██║    ██║     ██╔══██║██║   ██║██║     ██║██║
╚███╔███╔╝██║  ██║██║  ██╗██║  ██║    ███████╗██║  ██║╚██████╔╝╚██████╗██║██║
 ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚══════╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝╚═╝
`

	// System monitor
	systemInfo := fmt.Sprintf("🕒 %s  |  CPU: %.1f%%  |  RAM: %.1f%%",
		m.currentTime, m.cpuPercent, m.ramPercent)

	header := lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render(logo),
		headerStyle.Render("FEMA CLI - Интерактивный режим"),
		systemStyle.Render(systemInfo),
	)

	// Status message with animation
	var statusLine string
	if m.executing {
		statusLine = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render(
			m.spinner.View() + " Выполняется: " + m.lastCommand)
	} else if m.statusMsg != "" {
		statusLine = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render(
			"✅ " + m.statusMsg)
	}

	// Menu items
	var items []string
	for i, choice := range m.choices {
		if i == m.index {
			items = append(items, selectedStyle.Render(" > "+choice))
		} else {
			items = append(items, normalStyle.Render("   "+choice))
		}
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, items...)

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"\n",
		menu,
		"\n",
		normalStyle.Render("Используй ↑↓ для навигации, Enter для выбора, q для выхода, c для очистки"),
	)

	if statusLine != "" {
		content = lipgloss.JoinVertical(lipgloss.Center, content, "\n", statusLine)
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, content)
}

func (m model) executeChoice() tea.Cmd {
	return func() tea.Msg {
		// Start spinner animation during command execution
		cmd := m.spinner.Tick
		go func() {
			for m.executing {
				cmd()
				time.Sleep(100 * time.Millisecond)
			}
		}()
		switch m.index {
		case 0: // ⚙ Форматнуть Go
			return m.runFormatGo()
		case 1: // 🧪 Аудит фронта
			return m.runAuditFront()
		case 2: // 🚀 Деплой
			return m.runDeploy()
		case 3: // 🔧 git-relasens
			return m.runGitRelasens()
		case 4: // 📊 Системная инфо
			return m.runSysinfo()
		case 5: // 🌐 Проверка URL
			return m.runPing()
		case 6: // 📁 Сканировать логи
			return m.runScan()
		case 7: // 📜 Просмотр логов
			return m.runLogs()
		case 8: // ❌ Выйти
			return tea.Quit()
		}
		return nil
	}
}

func (m model) runFormatGo() tea.Msg {
	// Ensure formater.py exists
	ensureFormaterExists()

	cmd := exec.Command("python3", "formater.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка форматирования: %v\n%s", err, output)}
	}
	return successMsg(fmt.Sprintf("✅ Форматирование завершено!\n%s", output))
}

func (m model) runAuditFront() tea.Msg {
	cmd := exec.Command("python", "audit_front.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка аудита фронта: %v\n%s", err, output)}
	}
	return successMsg(fmt.Sprintf("🧪 Аудит фронта завершён:\n%s", output))
}

func (m model) runDeploy() tea.Msg {
	// Placeholder for deploy functionality
	return successMsg("🚀 Деплой выполнен (заглушка)")
}

func (m model) runGitRelasens() tea.Msg {
	cmd := exec.Command("python", "git-relasens.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка git-relasens: %v\n%s", err, output)}
	}
	return successMsg(fmt.Sprintf("🔧 git-relasens выполнен:\n%s", output))
}

func (m model) runSysinfo() tea.Msg {
	cmd := exec.Command("./fema-cli/fema", "sysinfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка системной инфы: %v", err)}
	}
	return successMsg(fmt.Sprintf("📊 Системная информация:\n%s", output))
}

func (m model) runPing() tea.Msg {
	// For now, just test google.com
	cmd := exec.Command("./fema-cli/fema", "ping", "http://localhost:80")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка проверки URL: %v", err)}
	}
	return successMsg(fmt.Sprintf("🌐 Проверка Google:\n%s", output))
}

func (m model) runScan() tea.Msg {
	cmd := exec.Command("./fema-cli/fema", "scan")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка сканирования: %v", err)}
	}
	return successMsg(fmt.Sprintf("📁 Сканирование завершено:\n%s", output))
}

func (m model) runLogs() tea.Msg {
	cmd := exec.Command("./fema-cli/fema", "logs")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errorMsg{err: fmt.Errorf("ошибка просмотра логов: %v", err)}
	}
	return successMsg(fmt.Sprintf("📜 Логи:\n%s", output))
}

type errorMsg struct{ err error }
type successMsg string

func ensureFormaterExists() {
	if _, err := os.Stat("formater.py"); os.IsNotExist(err) {
		originalPath := "./script/formater.py"
		if _, err := os.Stat(originalPath); err == nil {
			input, err := os.ReadFile(originalPath)
			if err == nil {
				os.WriteFile("formater.py", input, 0755)
			}
		}
	}
}

func ensureAuditFrontExists() {
	if _, err := os.Stat("audit_front.py"); os.IsNotExist(err) {
		originalPath := "./script/audit_front.py"
		if _, err := os.Stat(originalPath); err == nil {
			input, err := os.ReadFile(originalPath)
			if err == nil {
				os.WriteFile("audit_front.py", input, 0755)
			}
		}
	}
}

func ensureGitRelasensExists() {
	if _, err := os.Stat("git-relasens.py"); os.IsNotExist(err) {
		originalPath := "./script/git-relasens.py"
		if _, err := os.Stat(originalPath); err == nil {
			input, err := os.ReadFile(originalPath)
			if err == nil {
				os.WriteFile("git-relasens.py", input, 0755)
			}
		}
	}
}

var tuiCmd = &cobra.Command{
	Use:     "tui",
	Short:   "Rich TUI interface with animations and system monitor",
	Long:    "Launch beautiful terminal user interface with ASCII logo, system monitoring, and interactive menu",
	Aliases: []string{"ui", "menu", "interactive"},
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}
