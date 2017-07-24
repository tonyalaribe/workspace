package storage

type FileManager interface {
	Save(workspaceID string, formID string, submissionName string, b64Data string) (string, error)
}
