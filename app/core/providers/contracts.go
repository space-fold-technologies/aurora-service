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
	Name           string
	Width          int
	Heigth         int
	ClientTerminal string
}
