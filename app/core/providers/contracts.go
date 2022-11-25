package providers

type Variable struct {
	Key   string
	Value string
}

type Mount struct {
	Source   string
	Target   string
	TypeBind string
}

type TerminalProperties struct {
	Identifier     string
	Isolated       bool
	Name           string
	Width          int
	Heigth         int
	ClientTerminal string
	Token          string
}

type ManagerDetails struct {
	ID      string
	Address string
}
