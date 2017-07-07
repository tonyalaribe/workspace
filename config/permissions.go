package config

import (
	"encoding/json"
	"io/ioutil"
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
	if err := SaveJSON("roles.json", &jsonOutputRoles); err != nil {
		log.Fatal(err)
	}
	// Save inheritance information
	if err := SaveJSON("inheritance.json", &jsonOutputInher); err != nil {
		log.Fatal(err)
	}
}

func GenerateRolesInstance() *gorbac.RBAC {
	rbac := gorbac.New()
	permissions := make(gorbac.Permissions)

	// map[RoleId]PermissionIds
	var jsonRoles map[string][]string
	// map[RoleId]ParentIds
	var jsonInher map[string][]string
	// Load roles information
	if err := LoadJSON("roles.json", &jsonRoles); err != nil {
		log.Fatal(err)
	}
	// Load inheritance information
	if err := LoadJSON("inheritance.json", &jsonInher); err != nil {
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

func SaveJSON(filename string, jsonObject interface{}) error {
	jsonByte, _ := json.Marshal(jsonObject)
	err := ioutil.WriteFile(filename, jsonByte, 0644)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func LoadJSON(filename string, jsonObject interface{}) error {
	fileByte, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(fileByte, &jsonObject)
	if err != nil {
		log.Println(err)
	}
	return nil

}
