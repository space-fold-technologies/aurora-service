package teams

type TeamRepository interface {
	Create(team *Team) error
	List() ([]*Team, error)
	Remove(name string) error
}
