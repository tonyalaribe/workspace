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
const DeleteSubmissionTriggerEvent TriggerEvent = "DeleteSubmission"

type Trigger struct {
	ID          string
	WorkspaceID string
	FormID      string
	EventType   TriggerEvent
	URL         string
	SecretToken string
}

type Database interface {
	GetInheritance() (string, error)         //GetInheritance Get the current Inheritance tree from database
	SaveInheritance(roles interface{}) error //SaveInheritance persists tuhe current inheritance tree to database

	GetRoles() (string, error)         //GetRoles gets the current roles json string from database
	SaveRoles(roles interface{}) error //SaveRoles persists the current permissions and roles tree

	CreateWorkspace(WorkSpace) error                                         //CreateWorkspace adds a new workspace to the database
	GetWorkspaces() ([]WorkSpace, error)                                     //GetWorkspaces returns all workspaces on the database
	GetWorkspaceUsersAndRoles(workspaceID string) (WorkSpace, []User, error) //GetWorkspaceUsersAndRoles retuns all users of a workspace and their associated roles
	GetWorkspaceBySlug(workspaceID string) (WorkSpace, error)                //GetWorkspaceBySlug Returns a workspace's meta data by slugname

	CreateForm(workspaceID string, form Form) error               //Create Form Creates a new form inthe db
	GetForms(workspaceID string) ([]Form, error)                  //GetForms Get all fforms associated with given workspaceID
	GetFormBySlug(workspaceID, formID string) (Form, error)       //GetFormBySlug gets a form with a given form ID under the given workspaceID
	GetFormJSONBySlug(workspaceID, formID string) ([]byte, error) //GetFormJSONBySlug gets raw json byte content for a form, given its slug.

	NewFormSubmission(workspaceID, formID string, submission SubmissionData) (SubmissionData, error)    //NewFormSubmission persists a new form data submission to the database
	UpdateFormSubmission(workspaceID, formID string, submissionID int, submission SubmissionData) error //UpdateFormSubmission updates a form submission in the db
	GetFormSubmissions(workspaceID, formID string) ([]SubmissionData, error)                            //GetFormSubmissions gets the submissions associated with a given form
	GetFormSubmissionDetails(workspaceID, formID string, submissionID int) (SubmissionData, error)      //GetFormSubmissionDetails gets the details for a given submission ID
	DeleteFormSubmission(workspaceID, formID string, submissionID int) (SubmissionData, error)          //DeleteFormSubmission deletes a form submission from the db given the submissionID

	CreateUser(user User) error //CreateUser Creates a user in the database
	GetUser(username string) (User, error)
	//GetUser returns a user given the username
	GetUserByEmail(email string) (User, error) //GetsUserByEmail returns a user givenits username
	GetAllUsers() ([]User, error)              //GetAllUsers returns all users from the datababse

	AddToSubmissionChangelog(workspaceID, formID string, submissionID int, changelogItem ChangelogItem) error //AddToSubmissionChangelog Stores a submissions changelog to db.
	GetSubmissionChangelog(workspaceID, formID string, submissionID int) ([]ChangelogItem, error)             //GetSubmissionChangelog retrieves a changelog from db based on its associated submission

	UpdateTrigger(trigger Trigger) error //UpdateTrigger updates a trigger in the database
	//Each event would be stored as one trigger
	DeleteTrigger(trigger Trigger) error //DeleteTrigger deletes a trigger from the database

	GetTriggers(workspaceID string, formID string, ID string, event TriggerEvent) (Trigger, error) //GetTriggers gets all trigger from the database of a given event type associated with a form
	GetFormTriggers(WorkSpace string, formID string) ([]Trigger, error)                            //GetFormTriggers gets all triggers associated with a given form
	GetEventTriggers(WorkSpace string, formID string, event TriggerEvent) ([]Trigger, error)       //GetEventTrigggers Gets all triggers associated with a given eventType and form
}
