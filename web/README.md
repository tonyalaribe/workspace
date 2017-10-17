# web
--
    import "gitlab.com/middlefront/workspace/web"


## Usage

#### func  App

```go
func App()
```

#### func  ChangeUserWorkspacePermission

```go
func ChangeUserWorkspacePermission(w http.ResponseWriter, r *http.Request)
```
ChangeUserWorkspacePermission chamges workspace permissions, but only if the
user making the action has enough (admin or similar) permissions

#### func  CreateFormHandler

```go
func CreateFormHandler(w http.ResponseWriter, r *http.Request)
```
CreateFormHandler Creates a form in db

#### func  CreateWorkspaceHandler

```go
func CreateWorkspaceHandler(w http.ResponseWriter, r *http.Request)
```
CreateWorkspaceHandler create workspace with database.WorkSpace{} as body.

#### func  DeleteSubmissionHandler

```go
func DeleteSubmissionHandler(w http.ResponseWriter, r *http.Request)
```
DeleteSubmissionHandler deletes a submission given submissionID and formID

#### func  DeleteTriggerHandler

```go
func DeleteTriggerHandler(w http.ResponseWriter, r *http.Request)
```
DeleteTriggerHandler Deletes a trigger (removes all trigger event types)

#### func  GetFormBySlugHandler

```go
func GetFormBySlugHandler(w http.ResponseWriter, r *http.Request)
```
GetFormBySlugHandler gets the form by slug

#### func  GetFormTriggersHandler

```go
func GetFormTriggersHandler(w http.ResponseWriter, r *http.Request)
```
GetFormTriggersHandler Get all triggers associated with a form

#### func  GetFormsHandler

```go
func GetFormsHandler(w http.ResponseWriter, r *http.Request)
```
GetFormsHandler Gets the form with given workspaceID

#### func  GetSubmissionChangelogHandler

```go
func GetSubmissionChangelogHandler(w http.ResponseWriter, r *http.Request)
```
GetSubmissionChangelogHandler returns changelogs associated with submission

#### func  GetSubmissionInfoHandler

```go
func GetSubmissionInfoHandler(w http.ResponseWriter, r *http.Request)
```
GetSubmissionInfoHandler returns information about given submission

#### func  GetSubmissionsHandler

```go
func GetSubmissionsHandler(w http.ResponseWriter, r *http.Request)
```
GetSubmissionsHandler returns the submission based on formID

#### func  GetUploadedFile

```go
func GetUploadedFile(w http.ResponseWriter, r *http.Request, p httprouter.Params)
```
GetUploadedFile streams the uploaded file to the user irrespective of data store

#### func  GetUserInfoFromToken

```go
func GetUserInfoFromToken(next http.Handler) http.Handler
```
GetUserInfoFromToken rerturns json encoding of user info

#### func  GetWorkspaceBySlugHandler

```go
func GetWorkspaceBySlugHandler(w http.ResponseWriter, r *http.Request)
```
GetWorkspaceBySlugHandler Get a workspaces details given the workspaceID(slug)
as param.

#### func  GetWorkspaceUsersAndRolesHandler

```go
func GetWorkspaceUsersAndRolesHandler(w http.ResponseWriter, r *http.Request)
```
GetWorkspaceUsersAndRolesHandler get users and their roles attached to a
workspace

#### func  GetWorkspacesHandler

```go
func GetWorkspacesHandler(w http.ResponseWriter, r *http.Request)
```
GetWorkspacesHandler Get workspaces a user has access to

#### func  HomePageHandler

```go
func HomePageHandler(w http.ResponseWriter, r *http.Request)
```

#### func  LoggingHandler

```go
func LoggingHandler(next http.Handler) http.Handler
```
LoggingHandler Logs request time, method and duration of handler/request
execution

#### func  NewFormSubmissionHandler

```go
func NewFormSubmissionHandler(w http.ResponseWriter, r *http.Request)
```
NewFormSubmissionHandler is triggered when a new submission is made. This
handler is saving the actual variable submission info

#### func  RecoverHandler

```go
func RecoverHandler(next http.Handler) http.Handler
```
RecoverHandler catches all panics, so the serverdoesnt go down ocmpletely, jsust
because of a panic, that could be in one handler request by one user, affecting
every other user.

#### func  TestTriggerHandler

```go
func TestTriggerHandler(w http.ResponseWriter, r *http.Request)
```
TestTriggerHandler Send a test trigger to registered urls tied to trigger
actions

#### func  UpdateSubmissionHandler

```go
func UpdateSubmissionHandler(w http.ResponseWriter, r *http.Request)
```
UodateSubmissionHandler update submission

#### func  UpdateTriggerHandler

```go
func UpdateTriggerHandler(w http.ResponseWriter, r *http.Request)
```
UpdateTriggerHandler Update a trigger in case of changes to the trigger

#### func  UsersAndWorkspaceRoles

```go
func UsersAndWorkspaceRoles(w http.ResponseWriter, r *http.Request)
```
UsersAndWorkspaceRoles Get users and their roles in the workspace

#### type Router

```go
type Router struct {
	*httprouter.Router
}
```

Router struct would carry the httprouter instance, so its methods could be
verwritten and replaced with methds with wraphandler

#### func  NewRouter

```go
func NewRouter() *Router
```
NewRouter is a wrapper that makes the httprouter struct a child of the router
struct

#### func (*Router) Delete

```go
func (r *Router) Delete(path string, handler http.Handler)
```
Delete is an endpoint to only accept requests of method DELETE

#### func (*Router) Get

```go
func (r *Router) Get(path string, handler http.Handler)
```
Get is an endpoint to only accept requests of method GET

#### func (*Router) Post

```go
func (r *Router) Post(path string, handler http.Handler)
```
Post is an endpoint to only accept requests of method POST

#### func (*Router) Put

```go
func (r *Router) Put(path string, handler http.Handler)
```
Put is an endpoint to only accept requests of method PUT

#### type TriggerJSON

```go
type TriggerJSON struct {
	ID          string
	URL         string
	Endpoint    string
	SecretToken string

	NewSubmission     bool
	UpdateSubmission  bool
	ApproveSubmission bool
	DeleteSubmission  bool
}
```

TriggerJSON encodes information about triggers and what events are active

#### type WorkSpace

```go
type WorkSpace struct {
	Creator string `json:"creator"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}
```
