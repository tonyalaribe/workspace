package config

import (
	"encoding/json"
	"log"

	"github.com/mikespook/gorbac"
)

func SavePermissions() {
	conf := Get()
	rbac := conf.RolesManager
	// Persist the change
	// map[RoleId]PermissionIds
	jsonOutputRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonOutputInher := make(map[string][]string)
	SaveJsonHandler := func(r gorbac.Role, parents []string) error {
		// WARNING: Don't use gorbac.RBAC instance in the handler,
		// otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.(*gorbac.StdRole).Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID()] = permissions
		jsonOutputInher[r.ID()] = parents
		return nil
	}
	if err := gorbac.Walk(rbac, SaveJsonHandler); err != nil {
		log.Fatalln(err)
	}

	// Save roles information
	if err := conf.Database.SaveRoles(&jsonOutputRoles); err != nil {
		log.Fatal(err)
	}
	// Save inheritance information
	if err := conf.Database.SaveInheritance(&jsonOutputInher); err != nil {
		log.Fatal(err)
	}
}

func GenerateRolesInstance() *gorbac.RBAC {
	rbac := gorbac.New()
	permissions := make(gorbac.Permissions)
	conf := Get()
	// map[RoleId]PermissionIds
	jsonRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonInher := make(map[string][]string)

	// Load roles information
	rolesJSONString, err := conf.Database.GetRoles()
	if err != nil {
		log.Fatal(err)
	}

	if rolesJSONString == "" {
		rolesJSONString = `{"superadmin":["superadmin"]}`
	}
	err = json.Unmarshal([]byte(rolesJSONString), &jsonRoles)
	if err != nil {
		log.Println(err)
	}

	inheritanceJSONString, err := conf.Database.GetInheritance()
	// Load inheritance information
	if err != nil {
		log.Println(err)
	}

	if inheritanceJSONString == "" {
		inheritanceJSONString = `{"superadmin":[]}`
	}
	err = json.Unmarshal([]byte(inheritanceJSONString), &jsonInher)
	if err != nil {
		log.Fatal(err)
	}
	for rid, pids := range jsonRoles {
		role := gorbac.NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = gorbac.NewStdPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac.Add(role)
	}
	for rid, parents := range jsonInher {
		if err := rbac.SetParents(rid, parents); err != nil {
			log.Fatal(err)
		}
	}
	return rbac

}
