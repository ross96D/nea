package nea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stack struct {
	list []tea.Model
}

func (s *stack) Last() tea.Model {
	return s.list[len(s.list)-1]
}

func (s *stack) SetLast(m tea.Model) {
	s.list[len(s.list)-1] = m
}

func (s *stack) Push(v tea.Model) {
	s.list = append(s.list, v)
}

func (s *stack) Pop() tea.Model {
	if len(s.list) == 0 {
		return nil
	}

	l := len(s.list)
	if l == 1 {
		return s.list[l-1]
	}
	s.list = s.list[:l-1]

	return s.list[l-2]
}

type navigatorPush tea.Model
type navigatorPop struct{}
type escHandlerMsg bool

func NavigatorPush(m tea.Model) tea.Cmd {
	return func() tea.Msg {
		return navigatorPush(m)
	}
}
func NavigatorPop() tea.Msg {
	return navigatorPop{}
}
func EscHandler(h bool) tea.Cmd {
	return func() tea.Msg {
		return escHandlerMsg(h)
	}
}

type NavModel interface {
	tea.Model
	Enter() tea.Cmd
	Out() tea.Cmd
}

type Navigator struct {
	s         stack
	handleEsc bool
}

func NewNavigator() *Navigator {
	return &Navigator{handleEsc: true, s: stack{list: []tea.Model{}}}
}

func (nav *Navigator) Push(m tea.Model) (tea.Model, tea.Cmd) {
	nav.s.Push(m)
	if m, ok := m.(NavModel); ok {
		// TODO define if Init or Enter should be called.. probably both should be called? or not
		// is a bit confusing having both here
		return m, tea.Sequence(m.Init(), m.Enter())
	}
	return m, m.Init()
}

func (nav *Navigator) Pop() tea.Model {
	return nav.s.Pop()
}

func (nav *Navigator) View() string {
	return nav.s.Last().View()
}

func (nav *Navigator) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case navigatorPop:
		var cmd [2]tea.Cmd = [2]tea.Cmd{nil, nil}
		popped := nav.Pop()
		if popped == nil {
			return nil
		}
		if m, ok := popped.(NavModel); ok {
			cmd[0] = m.Out()
		}
		if m, ok := nav.s.Last().(NavModel); ok {
			cmd[1] = m.Enter()
		}
		return tea.Batch(cmd[0], cmd[1])

	case navigatorPush:
		_, cmd := nav.Push(msg)
		return cmd

	case escHandlerMsg:
		nav.handleEsc = bool(msg)
		return nil

	case tea.KeyMsg:
		// navigation go back
		if msg.Type == tea.KeyEscape && nav.handleEsc {
			return NavigatorPop
		}
	}

	m, cmd := nav.s.Last().Update(msg)
	nav.s.SetLast(m)
	return cmd
}
