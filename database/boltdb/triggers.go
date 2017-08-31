package boltdb

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

func (boltDBProvider *BoltDBProvider) UpdateTrigger(trigger database.Trigger) error {
	err := boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers))

		// prefix := []byte(trigger.WorkspaceID + ":" + trigger.FormID + ":" + string(trigger.EventType))

		// triggersCursor := triggersBucket.Cursor()
		// lastKey := 0

		// for k, _ := triggersCursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = triggersCursor.Next() {
		// 	// fmt.Printf("key=%s, value=%s\n", k, v)
		// 	numbersStr := strings.Split(string(k), string(prefix))[0]
		// 	log.Println(numbersStr)
		// 	var err error
		// 	number, err := strconv.Atoi(numbersStr)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	if number > lastKey {
		// 		lastKey = number
		// 	}
		// }
		//
		// trigger.ID = trigger.WorkspaceID + ":" + trigger.FormID + ":" + string(trigger.EventType) + ":" + strconv.Itoa(lastKey+1)

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
	if err != nil {
		return err
	}
	return nil
}

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
			log.Println(triggers)
		}
		return nil
	})
	if err != nil {
		return triggers, err
	}
	return triggers, nil
}

func (boltDBProvider *BoltDBProvider) GetEventTriggers(WorkspaceID, FormID string, EventType database.TriggerEvent) ([]database.Trigger, error) {

	var trigger database.Trigger
	var triggers []database.Trigger
	err := boltDBProvider.db.View(func(tx *bolt.Tx) error {
		triggersBucket := tx.Bucket([]byte(boltDBProvider.Triggers)).Cursor()
		prefix := []byte(WorkspaceID + ":" + FormID + ":" + string(EventType))

		var err error
		for k, v := triggersBucket.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = triggersBucket.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			err = json.Unmarshal(v, &trigger)
			if err != nil {
				return err
			}
			triggers = append(triggers)
		}
		return nil
	})

	if err != nil {
		return triggers, err
	}

	return triggers, nil
}
