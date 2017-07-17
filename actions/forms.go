package actions

import (
	slugify "github.com/metal3d/go-slugify"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func CreateForm(workspaceID string, formData database.Form) error {
	formData.ID = slugify.Marshal(formData.Name, true)
	conf := config.Get()
	err := conf.Database.CreateForm(workspaceID, formData)
	return err
}

func GetForms(workspaceID string) ([]database.Form, error) {
	conf := config.Get()
	forms, err := conf.Database.GetForms(workspaceID)
	return forms, err
}

func GetFormBySlug(workspaceID, formID string) (database.Form, error) {
	conf := config.Get()
	form, err := conf.Database.GetFormBySlug(workspaceID, formID)
	return form, err
}
