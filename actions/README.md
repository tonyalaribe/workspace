# actions
--
    import "gitlab.com/middlefront/workspace/actions"


## Usage

#### func  ChangeUserWorkspacePermission

```go
func ChangeUserWorkspacePermission(workspaceID string, permissions map[string]interface{}) error
```
ChangeUserWorkspacePermission change permissions for a given workspace

#### func  CreateForm

```go
func CreateForm(workspaceID string, formData database.Form) error
```
CreateForm creates a form under the given workspaceID with the given form
details

#### func  CreateWorkspace

```go
func CreateWorkspace(workspaceData database.WorkSpace, user database.User) error
```
CreateWorkspace creates a workspace with given user as owner

#### func  DeleteFormSubmission

```go
func DeleteFormSubmission(workspaceID, formID, submissionIDString string) (database.SubmissionData, error)
```
DeleteFormSubmission deletes a form submission given the submission id

#### func  DeleteTrigger

```go
func DeleteTrigger(trigger database.Trigger) error
```
DeleteTrigger deletes a trigger given a struct with the trigger data

#### func  GetEventTriggers

```go
func GetEventTriggers(workspaceID string, formID string, event database.TriggerEvent) ([]database.Trigger, error)
```
GetEventTriggers returns triggers attached to given event

#### func  GetFormBySlug

```go
func GetFormBySlug(workspaceID, formID string) (database.Form, error)
```
GetFormBySlug returns a form from db with given formID and workspace

#### func  GetFormSubmissionDetails

```go
func GetFormSubmissionDetails(workspaceID, formID, submissionIDString string) (database.SubmissionData, error)
```
GetFormSubmissionDetails returns the submission details given the submission id

#### func  GetFormSubmissions

```go
func GetFormSubmissions(workspaceID, formID string) ([]database.SubmissionData, error)
```
GetFormSubmissions returns all submissions under a given formID and workspace

#### func  GetFormTriggers

```go
func GetFormTriggers(workspaceID string, formID string) ([]database.Trigger, error)
```
?GetFormTriggers returns triggers associated with given form

#### func  GetForms

```go
func GetForms(workspaceID string) ([]database.Form, error)
```
GetForms returns all forms under the given workspace

#### func  GetSubmissionChangelog

```go
func GetSubmissionChangelog(workspaceID, formID, submissionIDString string) ([]database.ChangelogItem, error)
```
GetSubmissionChangelog returns the Changelogs for a given submission

#### func  GetUsersAndWorkspaceRoles

```go
func GetUsersAndWorkspaceRoles() ([]database.User, error)
```
GetUsersAndWorkspaceRoles retuns all users of a workspace and their associated
roles

#### func  GetWorkspaceBySlug

```go
func GetWorkspaceBySlug(workspaceID string) (database.WorkSpace, error)
```
GetWorkspaceBySlug gets workspace by slug

#### func  GetWorkspaceUsersAndRoles

```go
func GetWorkspaceUsersAndRoles(workspaceID string) ([]database.User, error)
```
Gets users associated with workspace and their roles

#### func  GetWorkspaces

```go
func GetWorkspaces(user database.User) ([]database.WorkSpace, error)
```
GetWorkspaces gets workspaces given user can access

#### func  NewFormSubmission

```go
func NewFormSubmission(workspaceID, formID string, submission database.SubmissionData) error
```
NewFormSubmission creates a new form submission

#### func  PostToURL

```go
func PostToURL(url string, secretToken string, body interface{})
```
PostToURL sends event info to registered webhook

#### func  SetupSuperAdmin

```go
func SetupSuperAdmin(adminUsername string) error
```
Setup a default super admin

#### func  TriggerEvent

```go
func TriggerEvent(workspaceID, formID string, event database.TriggerEvent, data map[string]interface{})
```
TriggerEvent triggers relevant webhooks for given event

#### func  UpdateSubmission

```go
func UpdateSubmission(workspaceID, formID string, submissionIDString string, newSubmission database.SubmissionData) error
```
UpdateSubmission updates a submission with givven workspaceID

#### func  UpdateTrigger

```go
func UpdateTrigger(trigger database.Trigger) error
```
UpdateTrigger updates a trigger given a struct with trigger data
