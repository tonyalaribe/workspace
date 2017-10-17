# boltdb
--
    import "gitlab.com/middlefront/workspace/database/boltdb"


## Usage

#### type BoltDBProvider

```go
type BoltDBProvider struct {
	RootDirectory       string
	AppMetadata         string
	WorkspacesMetadata  string
	WorkspacesContainer string
	UsersBucket         string
	FormsMetadata       string
	ChangelogBucket     string
	Triggers            string
}
```

BoltDBProvider holds information that is required for the boltdb provider to
function. Including a pointer to the database instance.

#### func  New

```go
func New(RootDirectory, AppMetadata, WorkspacesMetadata, WorkspacesContainer, UsersBucket, FormsMetadata string) (*BoltDBProvider, error)
```
New creates a boltdb instance given names of main buckets

#### func (*BoltDBProvider) AddToSubmissionChangelog

```go
func (boltDBProvider *BoltDBProvider) AddToSubmissionChangelog(workspaceID, formID string, submissionID int, changelogItem database.ChangelogItem) error
```
AddToSubmissionChangelog Stores a submissions changelog to db.

#### func (*BoltDBProvider) CreateForm

```go
func (boltDBProvider *BoltDBProvider) CreateForm(workspaceID string, formData database.Form) error
```
Create Form Creates a new form inthe db

#### func (*BoltDBProvider) CreateUser

```go
func (boltDBProvider *BoltDBProvider) CreateUser(user database.User) error
```
CreateUser Creates a user in the database

#### func (*BoltDBProvider) CreateWorkspace

```go
func (boltDBProvider *BoltDBProvider) CreateWorkspace(workspaceData database.WorkSpace) error
```
CreateWorkspace adds a new workspace to the database

#### func (*BoltDBProvider) DeleteFormSubmission

```go
func (boltDBProvider *BoltDBProvider) DeleteFormSubmission(workspaceID, formID string, submissionID int) (database.SubmissionData, error)
```
DeleteFormSubmission deletes a form submission from the db given the
submissionID

#### func (*BoltDBProvider) DeleteTrigger

```go
func (boltDBProvider *BoltDBProvider) DeleteTrigger(trigger database.Trigger) error
```
DeleteTrigger deletes a trigger from the database

#### func (*BoltDBProvider) GetAllUsers

```go
func (boltDBProvider *BoltDBProvider) GetAllUsers() ([]database.User, error)
```
GetAllUsers returns all users from the datababse

#### func (*BoltDBProvider) GetEventTriggers

```go
func (boltDBProvider *BoltDBProvider) GetEventTriggers(WorkspaceID, FormID string, EventType database.TriggerEvent) ([]database.Trigger, error)
```
GetEventTrigggers Gets all triggers associated with a given eventType and form

#### func (*BoltDBProvider) GetFormBySlug

```go
func (boltDBProvider *BoltDBProvider) GetFormBySlug(workspaceID, formID string) (database.Form, error)
```
GetFormBySlug gets a form with a given form ID under the given workspaceID

#### func (*BoltDBProvider) GetFormJSONBySlug

```go
func (boltDBProvider *BoltDBProvider) GetFormJSONBySlug(workspaceID, formID string) ([]byte, error)
```
GetFormJSONBySlug gets raw json byte content for a form, given its slug.

#### func (*BoltDBProvider) GetFormSubmissionDetails

```go
func (boltDBProvider *BoltDBProvider) GetFormSubmissionDetails(workspaceID, formID string, submissionID int) (database.SubmissionData, error)
```
GetFormSubmissionDetails gets the details for a given submission ID

#### func (*BoltDBProvider) GetFormSubmissions

```go
func (boltDBProvider *BoltDBProvider) GetFormSubmissions(workspaceID, formID string) ([]database.SubmissionData, error)
```
GetFormSubmissions gets the submissions associated with a given form

#### func (*BoltDBProvider) GetFormTriggers

```go
func (boltDBProvider *BoltDBProvider) GetFormTriggers(WorkspaceID, FormID string) ([]database.Trigger, error)
```
GetFormTriggers gets all triggers associated with a given form

#### func (*BoltDBProvider) GetForms

```go
func (boltDBProvider *BoltDBProvider) GetForms(workspaceID string) ([]database.Form, error)
```
GetForms Get all fforms associated with given workspaceID

#### func (*BoltDBProvider) GetInheritance

```go
func (boltDBProvider *BoltDBProvider) GetInheritance() (string, error)
```
GetInheritance Get the current Inheritance tree from database

#### func (*BoltDBProvider) GetRoles

```go
func (boltDBProvider *BoltDBProvider) GetRoles() (string, error)
```
GetRoles gets the current roles json string from database

#### func (*BoltDBProvider) GetSubmissionChangelog

```go
func (boltDBProvider *BoltDBProvider) GetSubmissionChangelog(workspaceID, formID string, submissionID int) ([]database.ChangelogItem, error)
```
GetSubmissionChangelog retrieves a changelog from db based on its associated
submission

#### func (*BoltDBProvider) GetTriggers

```go
func (boltDBProvider *BoltDBProvider) GetTriggers(WorkspaceID, FormID, ID string, EventType database.TriggerEvent) (database.Trigger, error)
```
GetTriggers gets all trigger from the database of a given event type associated
with a form

#### func (*BoltDBProvider) GetUser

```go
func (boltDBProvider *BoltDBProvider) GetUser(username string) (database.User, error)
```
GetUser returns a user given the username

#### func (*BoltDBProvider) GetUserByEmail

```go
func (boltDBProvider *BoltDBProvider) GetUserByEmail(email string) (database.User, error)
```
GetsUserByEmail returns a user givenits username

#### func (*BoltDBProvider) GetWorkspaceBySlug

```go
func (boltDBProvider *BoltDBProvider) GetWorkspaceBySlug(workspaceID string) (database.WorkSpace, error)
```
GetWorkspaceBySlug Returns a workspace's meta data by slugname

#### func (*BoltDBProvider) GetWorkspaceUsersAndRoles

```go
func (boltDBProvider *BoltDBProvider) GetWorkspaceUsersAndRoles(workspaceID string) (database.WorkSpace, []database.User, error)
```
GetWorkspaceUsersAndRoles retuns all users of a workspace and their associated
roles

#### func (*BoltDBProvider) GetWorkspaces

```go
func (boltDBProvider *BoltDBProvider) GetWorkspaces() ([]database.WorkSpace, error)
```
GetWorkspaces returns all workspaces on the database

#### func (*BoltDBProvider) NewFormSubmission

```go
func (boltDBProvider *BoltDBProvider) NewFormSubmission(workspaceID, formID string, submission database.SubmissionData) error
```
NewFormSubmission persists a new form data submission to the database

#### func (*BoltDBProvider) SaveInheritance

```go
func (boltDBProvider *BoltDBProvider) SaveInheritance(roles interface{}) error
```
SaveInheritance persists tuhe current inheritance tree to database

#### func (*BoltDBProvider) SaveRoles

```go
func (boltDBProvider *BoltDBProvider) SaveRoles(roles interface{}) error
```
SaveRoles persists the current permissions and roles tree

#### func (*BoltDBProvider) UpdateFormSubmission

```go
func (boltDBProvider *BoltDBProvider) UpdateFormSubmission(workspaceID, formID string, submissionID int, submission database.SubmissionData) error
```
UpdateFormSubmission updates a form submission in the db

#### func (*BoltDBProvider) UpdateTrigger

```go
func (boltDBProvider *BoltDBProvider) UpdateTrigger(trigger database.Trigger) error
```
UpdateTrigger updates a trigger in the database
