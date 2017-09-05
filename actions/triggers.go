package actions

import (
	"log"

	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func UpdateTrigger(trigger database.Trigger) error {
	conf := config.Get()
	// Persist workspace
	log.Printf("%#v", trigger)
	err := conf.Database.UpdateTrigger(trigger)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTrigger(trigger database.Trigger) error {
	conf := config.Get()
	// Persist workspace
	log.Printf("%#v", trigger)
	err := conf.Database.DeleteTrigger(trigger)
	if err != nil {
		return err
	}

	return nil
}

func GetFormTriggers(workspaceID string, formID string) ([]database.Trigger, error) {
	conf := config.Get()
	// Persist workspace

	triggers := []database.Trigger{}

	triggers, err := conf.Database.GetFormTriggers(workspaceID, formID)
	if err != nil {
		log.Println(err)
		// return trigger, err
	}

	return triggers, nil
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
