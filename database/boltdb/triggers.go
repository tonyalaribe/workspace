package boltdb

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

//UpdateTrigger updates a trigger in the database
func (boltDBProvider *BoltDBProvider) UpdateTrigger(trigger database.Trigger) error {
	err := boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers))

		trigger.ID = trigger.WorkspaceID + ":" + trigger.FormID + ":" + string(trigger.EventType) + ":" + trigger.URL

		dataByte, err := json.Marshal(trigger)
		if err != nil {
			log.Println(err)
		}
		err = triggersBucket.Put([]byte(trigger.ID), dataByte)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	return err
}

//DeleteTrigger deletes a trigger from the database
func (boltDBProvider *BoltDBProvider) DeleteTrigger(trigger database.Trigger) error {
	err := boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers))

		trigger.ID = trigger.WorkspaceID + ":" + trigger.FormID + ":" + string(trigger.EventType) + ":" + trigger.URL

		err := triggersBucket.Delete([]byte(trigger.ID))
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//GetTriggers gets all trigger from the database of a given event type associated with a form
func (boltDBProvider *BoltDBProvider) GetTriggers(WorkspaceID, FormID, ID string, EventType database.TriggerEvent) (database.Trigger, error) {
	var triggerByte []byte
	var trigger database.Trigger
	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers))

		triggerKey := WorkspaceID + ":" + FormID + ":" + string(EventType) + ":" + ID

		triggerByte = triggersBucket.Get([]byte(triggerKey))

		return nil
	})
	if err != nil {
		return trigger, err
	}

	err = json.Unmarshal(triggerByte, &trigger)
	if err != nil {
		return trigger, err
	}
	return trigger, nil
}

//GetFormTriggers gets all triggers associated with a given form
func (boltDBProvider *BoltDBProvider) GetFormTriggers(WorkspaceID, FormID string) ([]database.Trigger, error) {

	var trigger database.Trigger
	var triggers []database.Trigger
	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers)).Cursor()
		prefix := []byte(WorkspaceID + ":" + FormID)

		var err error
		for k, v := triggersBucket.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = triggersBucket.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			err = json.Unmarshal(v, &trigger)
			if err != nil {
				log.Println(err)
				return err
			}
			triggers = append(triggers, trigger)
			// log.Println(triggers)
		}
		return nil
	})
	if err != nil {
		return triggers, err
	}
	return triggers, nil
}

//GetEventTrigggers Gets all triggers associated with a given eventType and form
func (boltDBProvider *BoltDBProvider) GetEventTriggers(WorkspaceID, FormID string, EventType database.TriggerEvent) ([]database.Trigger, error) {

	var triggers []database.Trigger
	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers)).Cursor()
		prefix := []byte(WorkspaceID + ":" + FormID + ":" + string(EventType))

		var err error
		var trigger database.Trigger
		for k, v := triggersBucket.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = triggersBucket.Next() {
			err = json.Unmarshal(v, &trigger)
			if err != nil {
				return err
			}
			triggers = append(triggers, trigger)
		}
		return nil
	})

	return triggers, err
}
