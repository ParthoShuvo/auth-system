package permission

type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Store interface {
	ReadUserPermissions(string) ([]*Permission, error)
}

type Table struct {
	store Store
}

func NewTable(s Store) *Table {
	return &Table{s}
}

func (t *Table) ReadUserPermissions(login string) ([]*Permission, error) {
	return t.store.ReadUserPermissions(login)
}
