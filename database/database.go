package database

type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

type User struct {
	ProviderUserID    string
	Username          string
	Name              string
	Email             string
	Roles             []string
	CurrentRoleString string
}

type Form struct {
	Creator    string                 `json:"creator"`
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	JSONSchema map[string]interface{} `json:"jsonschema"`
	UISchema   map[string]interface{} `json:"uischema"`
}

type SubmissionData struct {
	FormData        map[string]interface{} `json:"formData"`
	Created         int                    `json:"created"`
	LastModified    int                    `json:"lastModified"`
	SubmissionName  string                 `json:"submissionName"`
	Status          string                 `json:"status"`
	ID              int                    `json:"id"`
	SubmissionNotes string                 `json:"submissionNotes"`
}

type ChangelogItem struct {
	Created      int    `json:"created"`
	WorkspaceID  string `json:"workspaceID"`
	FormID       string `json:"formID"`
	SubmissionID int    `json:"submissionID"`

	Notes string `json:"note"`
}

type Files struct {
	Status     string `json:"status"`
	File       string `json:"file"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Path       string `json:"path"`
	CreatedBy  string `json:"createdBy"`
	UploadDate string `json:"uploadDate"`
}

type TriggerEvent string

const NewSubmissionTriggerEvent TriggerEvent = "NewSubmission"
const UpdateSubmissionTriggerEvent TriggerEvent = "UpdateSubmission"
const ApproveSubmissionTriggerEvent TriggerEvent = "ApproveSubmission"

type Trigger struct {
	ID          string
	WorkspaceID string
	FormID      string
	EventType   TriggerEvent
	URL         string
	SecretToken string
}

type Database interface {
	GetInheritance() (string, error)
	SaveInheritance(roles interface{}) error

	GetRoles() (string, error)
	SaveRoles(roles interface{}) error

	CreateWorkspace(WorkSpace) error
	GetWorkspaces() ([]WorkSpace, error)
	GetWorkspaceUsersAndRoles(workspaceID string) (WorkSpace, []User, error)
	GetWorkspaceBySlug(workspaceID string) (WorkSpace, error)

	CreateForm(workspaceID string, form Form) error
	GetForms(workspaceID string) ([]Form, error)
	GetFormBySlug(workspaceID, formID string) (Form, error)
	GetFormJSONBySlug(workspaceID, formID string) ([]byte, error)

	NewFormSubmission(workspaceID, formID string, submission SubmissionData) error
	UpdateFormSubmission(workspaceID, formID string, submissionID int, submission SubmissionData) error
	GetFormSubmissions(workspaceID, formID string) ([]SubmissionData, error)
	GetFormSubmissionDetails(workspaceID, formID string, submissionID int) (SubmissionData, error)

	CreateUser(user User) error
	GetUser(username string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetAllUsers() ([]User, error)

	AddToSubmissionChangelog(workspaceID, formID string, submissionID int, changelogItem ChangelogItem) error
	GetSubmissionChangelog(workspaceID, formID string, submissionID int) ([]ChangelogItem, error)

	UpdateTrigger(trigger Trigger) error //Each event would be stored as one trigger
	//The root triggers (push, etc). Used in constructing the form on settings page
	GetTriggers(workspaceID string, formID string, ID string, event TriggerEvent) (Trigger, error)

	GetEventTriggers(WorkSpace string, formID string, event TriggerEvent) ([]Trigger, error)
}
