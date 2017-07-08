package storage

type FileManager interface {
	Save(name string, workspace string, b64Data string) (string, error)
}
