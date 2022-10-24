package environments

type EnvironmentRepository interface {
	Add(Scope, Target string, Entries []*Entry) error
	List(Scope, Target string) ([]*Entry, error)
	Remove(Keys []string) error
}
