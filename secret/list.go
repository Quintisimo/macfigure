package secret

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quintisimo/macfigure/gen/secret"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := filepath.Base(fmt.Sprintf("%s", i))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type listModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return "\n" + m.list.View()
}

func List(secret []secret.Secret) (string, error) {
	const defaultWidth = 20

	secretPath := ""
	secretLen := len(secret)

	if secretLen == 0 {
		return "", nil
	}

	if secretLen == 1 {
		return secret[0].Source, nil
	}

	items := make([]list.Item, secretLen)
	for i, secretItem := range secret {
		items[i] = item(secretItem.Source)
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select secret to edit"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := listModel{list: l}
	runModel, runErr := tea.NewProgram(m).Run()
	if runErr != nil {
		return "", runErr
	}
	secretPath = runModel.(listModel).choice

	return secretPath, nil
}
