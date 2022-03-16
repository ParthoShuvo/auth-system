package role

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Store interface {
	ReadUserRoles(string) ([]*Role, error)
}

type Table struct {
	store Store
}

func NewTable(s Store) *Table {
	return &Table{s}
}

func (t *Table) ReadUserRoles(login string) ([]*Role, error) {
	return t.store.ReadUserRoles(login)
}
