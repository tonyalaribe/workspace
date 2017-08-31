package actions

import (
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func UpdateTrigger(trigger database.Trigger) error {
	conf := config.Get()
	// Persist workspace
	err := conf.Database.UpdateTrigger(trigger)
	if err != nil {
		return err
	}

	return nil
}

func GetTriggers(workspaceID string, formID string, ID string) (database.Trigger, error) {
	conf := config.Get()
	// Persist workspace

	trigger, err := conf.Database.GetTriggers(workspaceID, formID, ID, database.NewSubmissionTriggerEvent)
	if err != nil {
		return trigger, err
	}

	return trigger, nil
}

func GetEventTriggers(workspaceID string, formID string, event database.TriggerEvent) ([]database.Trigger, error) {
	conf := config.Get()
	// Persist workspace
	triggers, err := conf.Database.GetEventTriggers(workspaceID, formID, event)
	if err != nil {
		return triggers, err
	}

	return triggers, nil
}
