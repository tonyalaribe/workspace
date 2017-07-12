package database

type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

type Database interface {
	CreateWorkspace(WorkSpace) error
	GetWorkspaces() ([]WorkSpace, error)
}
