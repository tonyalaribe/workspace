package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

func (boltDBProvider *BoltDBProvider) AddToSubmissionChangelog(workspaceID, formID string, submissionID int, changelogItem database.ChangelogItem) error {
	dataByte, err := json.Marshal(changelogItem)
	if err != nil {
		log.Println(err)
		return err
	}

	/*Save to boltdb*/
	err = boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		log.Println("Xxx")
		workspaceBucket, err := tx.Bucket([]byte(boltDBProvider.ChangelogBucket)).CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}
		formBucket, err := workspaceBucket.CreateBucketIfNotExists([]byte(formID))
		if err != nil {
			log.Println(err)
		}
		submissionBucket, err := formBucket.CreateBucketIfNotExists(itob(submissionID))
		if err != nil {
			log.Println(err)
		}

		log.Println("1")
		nextSeq, err := submissionBucket.NextSequence()
		if err != nil {
			log.Println(err)
		}
		log.Println("2")
		err = submissionBucket.Put(itob(int(nextSeq)), dataByte)
		if err != nil {
			log.Println(err)
		}
		log.Println(3)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (boltDBProvider *BoltDBProvider) GetSubmissionChangelog(workspaceID, formID string, submissionID int) ([]database.ChangelogItem, error) {
	/*Save to boltdb*/
	var err error
	changelogItems := []database.ChangelogItem{}
	changelogItem := database.ChangelogItem{}
	err = boltDBProvider.db.Update(func(tx *bolt.Tx) error {

		workspaceBucket, err := tx.Bucket([]byte(boltDBProvider.ChangelogBucket)).CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}
		formBucket, err := workspaceBucket.CreateBucketIfNotExists([]byte(formID))
		if err != nil {
			log.Println(err)
		}
		submissionBucket, err := formBucket.CreateBucketIfNotExists(itob(submissionID))
		if err != nil {
			log.Println(err)
		}

		err = submissionBucket.ForEach(func(k []byte, v []byte) error {
			err := json.Unmarshal(v, &changelogItem)
			if err != nil {
				log.Println(err)
			}
			changelogItems = append(changelogItems, changelogItem)

			return nil
		})
		if err != nil {
			log.Println(err)
		}

		return nil
	})
	if err != nil {
		return changelogItems, err
	}

	return changelogItems, nil
}
