package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

//AddToSubmissionChangelog Stores a submissions changelog to db.
func (boltDBProvider *BoltDBProvider) AddToSubmissionChangelog(workspaceID, formID string, submissionID int, changelogItem database.ChangelogItem) error {
	dataByte, err := json.Marshal(changelogItem)
	if err != nil {
		log.Println(err)
		return err
	}

	/*Save to boltdb*/
	err = boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		var workspaceBucket *bolt.Bucket
		workspaceBucket, err = tx.Bucket([]byte(boltDBProvider.ChangelogBucket)).CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}

		var formBucket *bolt.Bucket
		formBucket, err = workspaceBucket.CreateBucketIfNotExists([]byte(formID))
		if err != nil {
			log.Println(err)
		}
		var submissionBucket *bolt.Bucket
		submissionBucket, err = formBucket.CreateBucketIfNotExists(itob(submissionID))
		if err != nil {
			log.Println(err)
		}

		var nextSeq uint64
		nextSeq, err = submissionBucket.NextSequence()
		if err != nil {
			log.Println(err)
		}

		err = submissionBucket.Put(itob(int(nextSeq)), dataByte)
		if err != nil {
			log.Println(err)
		}

		return nil
	})

	return err
}

//GetSubmissionChangelog retrieves a changelog from db based on its associated submission
func (boltDBProvider *BoltDBProvider) GetSubmissionChangelog(workspaceID, formID string, submissionID int) ([]database.ChangelogItem, error) {
	/*Save to boltdb*/
	var err error
	changelogItems := []database.ChangelogItem{}
	changelogItem := database.ChangelogItem{}

	err = boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		var workspaceBucket *bolt.Bucket
		workspaceBucket, err = tx.Bucket([]byte(boltDBProvider.ChangelogBucket)).CreateBucketIfNotExists([]byte(workspaceID))
		if err != nil {
			log.Println(err)
		}

		var formBucket *bolt.Bucket
		formBucket, err = workspaceBucket.CreateBucketIfNotExists([]byte(formID))
		if err != nil {
			log.Println(err)
		}

		var submissionBucket *bolt.Bucket
		submissionBucket, err = formBucket.CreateBucketIfNotExists(itob(submissionID))
		if err != nil {
			log.Println(err)
		}

		err = submissionBucket.ForEach(func(k []byte, v []byte) error {
			err = json.Unmarshal(v, &changelogItem)
			if err != nil {
				log.Println(err)
				return err
			}
			changelogItems = append(changelogItems, changelogItem)
			return err
		})
		if err != nil {
			log.Println(err)
		}

		return nil
	})

	reversedChangelogItems := []database.ChangelogItem{}
	for _, v := range changelogItems {
		reversedChangelogItems = append(reversedChangelogItems, v)
	}
	return reversedChangelogItems, err
}
