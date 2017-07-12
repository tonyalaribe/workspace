package boltdb

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"gitlab.com/middlefront/workspace/database"
)

func (boltDBProvider *BoltDBProvider) NewFormSubmission(workspaceID, formID string, submission database.SubmissionData) error {

	/*Save to boltdb*/
	boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		nextID, err := formBucket.NextSequence()
		if err != nil {
			log.Println(err)
		}

		submission.ID = int(nextID)
		dataByte, err := json.Marshal(submission)
		if err != nil {
			log.Println(err)
		}

		err = formBucket.Put(itob(int(nextID)), dataByte)
		if err != nil {
			log.Println(err)
		}

		return nil
	})

	return nil
}

func (boltDBProvider *BoltDBProvider) UpdateFormSubmission(workspaceID, formID string, submissionID int, submission database.SubmissionData) error {

	dataByte, err := json.Marshal(submission)
	if err != nil {
		log.Println(err)
	}

	/*Save to boltdb*/
	boltDBProvider.db.Update(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		err = formBucket.Put(itob(submissionID), dataByte)
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	return nil
}

func (boltDBProvider *BoltDBProvider) GetFormSubmissions(workspaceID, formID string) ([]database.SubmissionData, error) {

	submissions := []database.SubmissionData{}
	var err error
	err = boltDBProvider.db.View(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		c := formBucket.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			f := database.SubmissionData{}

			err = json.Unmarshal(v, &f)
			if err != nil {
				log.Println(err)
			}
			submissions = append(submissions, f)
		}
		return nil
	})
	if err != nil {
		return submissions, err
	}

	return submissions, nil
}

func (boltDBProvider *BoltDBProvider) GetFormSubmissionDetails(workspaceID, formID string, submissionID int) (database.SubmissionData, error) {

	submission := database.SubmissionData{}
	var err error
	err = boltDBProvider.db.View(func(tx *bolt.Tx) error {
		formBucket := tx.Bucket([]byte(boltDBProvider.WorkspacesContainer)).Bucket([]byte(workspaceID)).Bucket([]byte(formID))

		err = json.Unmarshal(formBucket.Get(itob(submissionID)), &submission)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		return submission, err
	}

	return submission, nil
}