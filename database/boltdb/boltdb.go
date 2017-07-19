package boltdb

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

type BoltDBProvider struct {
	db                  *bolt.DB
	RootDirectory       string
	AppMetadata         string
	WorkspacesMetadata  string
	WorkspacesContainer string
	UsersBucket         string
	FormsMetadata       string
}

func New(RootDirectory, AppMetadata, WorkspacesMetadata, WorkspacesContainer, UsersBucket, FormsMetadata string) (*BoltDBProvider, error) {

	var boltdb BoltDBProvider
	boltdb.AppMetadata = AppMetadata
	boltdb.RootDirectory = RootDirectory
	boltdb.WorkspacesMetadata = WorkspacesMetadata
	boltdb.WorkspacesContainer = WorkspacesContainer
	boltdb.UsersBucket = UsersBucket
	boltdb.FormsMetadata = FormsMetadata
	boltFile := filepath.Join(RootDirectory, "workspace.db")

	os.MkdirAll(RootDirectory, os.ModePerm)
	db, err := bolt.Open(boltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Println(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(AppMetadata))
		tx.CreateBucketIfNotExists([]byte(WorkspacesMetadata))
		tx.CreateBucketIfNotExists([]byte(WorkspacesContainer))
		tx.CreateBucketIfNotExists([]byte(UsersBucket))
		return nil
	})

	boltdb.db = db
	return &boltdb, nil
}

//
// func (BoltDB) Initialize() {
// 	conf := config.Get()
// 	conf.BoltFile = filepath.Join(conf.RootDirectory, "workspace.db")
//
// 	os.MkdirAll(conf.RootDirectory, os.ModePerm)
// 	db, err := bolt.Open(conf.BoltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	db.Update(func(tx *bolt.Tx) error {
// 		tx.CreateBucketIfNotExists([]byte(conf.WorkspacesMetadata))
// 		tx.CreateBucketIfNotExists([]byte(conf.WorkspacesContainer))
// 		tx.CreateBucketIfNotExists([]byte(conf.UsersBucket))
// 		return nil
// 	})
//
// 	log.Println(db.GoString())
// 	conf.DB = db
//
// }
//
//
// func (BoltDB) Initialize(){
//   conf := config.Get()
//   conf.BoltFile = filepath.Join(conf.RootDirectory, "workspace.db")
//   conf.SubmissionsBucket = []byte("submissions")
//   os.MkdirAll(conf.RootDirectory, os.ModePerm)
//   db, err := bolt.Open(conf.BoltFile, 0600, &bolt.Options{Timeout: 3 * time.Second})
//   if err != nil {
//     log.Println(err)
//   }
//   db.Update(func(tx *bolt.Tx) error {
//     tx.CreateBucketIfNotExists([]byte(conf.WorkspacesMetadata))
//     tx.CreateBucketIfNotExists([]byte(conf.WorkspacesContainer))
//     tx.CreateBucketIfNotExists([]byte(conf.UsersBucket))
//     return nil
//   })
//
//   log.Println(db.GoString())
//   conf.DB = db
//
// }
