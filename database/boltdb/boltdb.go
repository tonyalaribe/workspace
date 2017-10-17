package boltdb

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

//BoltDBProvider holds information that is required for the boltdb provider to function. Including a pointer to the database instance.
type BoltDBProvider struct {
	db                  *bolt.DB
	RootDirectory       string
	AppMetadata         string
	WorkspacesMetadata  string
	WorkspacesContainer string
	UsersBucket         string
	FormsMetadata       string
	ChangelogBucket     string
	Triggers            string
}

//New creates a boltdb instance given names of main buckets
func New(RootDirectory, AppMetadata, WorkspacesMetadata, WorkspacesContainer, UsersBucket, FormsMetadata string) (*BoltDBProvider, error) {

	var boltdb BoltDBProvider
	boltdb.AppMetadata = AppMetadata
	boltdb.RootDirectory = RootDirectory
	boltdb.WorkspacesMetadata = WorkspacesMetadata
	boltdb.WorkspacesContainer = WorkspacesContainer
	boltdb.UsersBucket = UsersBucket
	boltdb.FormsMetadata = FormsMetadata
	boltdb.ChangelogBucket = "ChangelogBucket"
	boltdb.Triggers = "Triggers"

	boltFile := filepath.Join(RootDirectory, "workspace.db")

	os.MkdirAll(RootDirectory, os.ModePerm)
	db, err := bolt.Open(boltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(boltdb.AppMetadata))
		tx.CreateBucketIfNotExists([]byte(boltdb.WorkspacesMetadata))
		tx.CreateBucketIfNotExists([]byte(boltdb.WorkspacesContainer))
		tx.CreateBucketIfNotExists([]byte(boltdb.UsersBucket))
		tx.CreateBucketIfNotExists([]byte(boltdb.Triggers))
		tx.CreateBucketIfNotExists([]byte(boltdb.ChangelogBucket))
		return nil
	})

	boltdb.db = db
	return &boltdb, nil
}
