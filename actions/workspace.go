package actions

import (
	"log"
	"time"

	"github.com/metal3d/go-slugify"
	"github.com/mikespook/gorbac"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

//CreateWorkspace creates a workspace with given user as owner
func CreateWorkspace(workspaceData database.WorkSpace, user database.User) error {

	workspaceData.Creator = user.Username
	workspaceData.ID = slugify.Marshal(workspaceData.Name, true)
	workspaceData.Created = int(time.Now().UnixNano() / 1000000) //Get the time since epoch in milli seconds (javascript date compatible)
	conf := config.Get()
	// Persist workspace

	err := conf.Database.CreateWorkspace(workspaceData)
	if err != nil {
		return err
	}

	spectator := gorbac.NewStdRole(workspaceData.ID + "-spectator")
	spectator.Assign(gorbac.NewStdPermission("view-" + workspaceData.ID))
	conf.RolesManager.Add(spectator)

	editor := gorbac.NewStdRole(workspaceData.ID + "-editor")
	editor.Assign(gorbac.NewStdPermission("edit-" + workspaceData.ID))
	conf.RolesManager.Add(editor)

	supervisor := gorbac.NewStdRole(workspaceData.ID + "-supervisor")
	supervisor.Assign(gorbac.NewStdPermission("approve-" + workspaceData.ID))
	conf.RolesManager.Add(supervisor)

	admin := gorbac.NewStdRole(workspaceData.ID + "-admin")
	admin.Assign(gorbac.NewStdPermission("admin-" + workspaceData.ID))
	conf.RolesManager.Add(admin)

	conf.RolesManager.SetParent(workspaceData.ID+"-editor", workspaceData.ID+"-spectator")
	conf.RolesManager.SetParent(workspaceData.ID+"-supervisor", workspaceData.ID+"-editor")
	conf.RolesManager.SetParent(workspaceData.ID+"-admin", workspaceData.ID+"-supervisor")

	conf.RolesManager.SetParent("superadmin", workspaceData.ID+"-admin")

	return nil
}

//GetWorkspaces gets workspaces given user can access
func GetWorkspaces(user database.User) ([]database.WorkSpace, error) {
	conf := config.Get()
	finalWorkspaces := []database.WorkSpace{}
	//Get Workspaces
	workspaces, err := conf.Database.GetWorkspaces()
	if err != nil {
		return finalWorkspaces, err
	}

	for _, v := range workspaces {
		workspacePermissionString := "view-" + v.ID
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)
		if gorbac.AnyGranted(conf.RolesManager, user.Roles, workspacePermission, nil) {
			finalWorkspaces = append(finalWorkspaces, v)
		}
	}

	return finalWorkspaces, nil
}

//Gets users associated with workspace and their roles
func GetWorkspaceUsersAndRoles(workspaceID string) ([]database.User, error) {
	conf := config.Get()
	finalUsers := []database.User{}
	workspace, users, err := conf.Database.GetWorkspaceUsersAndRoles(workspaceID)
	if err != nil {
		return finalUsers, err
	}

	for _, u := range users {
		workspacePermissionString := "view-" + workspace.ID
		log.Println(workspacePermissionString)
		workspacePermission := gorbac.NewStdPermission(workspacePermissionString)
		for _, v := range u.Roles {
			if conf.RolesManager.IsGranted(v, workspacePermission, nil) {
				u.CurrentRoleString = v
				finalUsers = append(finalUsers, u)
				continue
			}
		}
	}
	return finalUsers, nil
}

//GetWorkspaceBySlug gets workspace by slug
func GetWorkspaceBySlug(workspaceID string) (database.WorkSpace, error) {
	conf := config.Get()
	workspace, err := conf.Database.GetWorkspaceBySlug(workspaceID)
	if err != nil {
		log.Println(err)
		return workspace, err
	}
	return workspace, nil
}
