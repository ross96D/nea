package nea

import tea "github.com/charmbracelet/bubbletea"

type App struct {
	nav *Navigator
}

func NewApp() App {
	return App{
		nav: NewNavigator(),
	}
}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// resend any tea.Cmd for proper handling
	if msg, ok := msg.(tea.Cmd); ok {
		return app, msg
	}

	cmd := app.nav.Update(msg)
	return app, cmd
}

func (app App) View() string {
	return app.nav.View()
}
