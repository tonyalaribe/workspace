package actions

import (
	"log"
	"strconv"

	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
)

func GetSubmissionChangelog(workspaceID, formID, submissionIDString string) ([]database.ChangelogItem, error) {
	submissionID, err := strconv.Atoi(submissionIDString)
	if err != nil {
		log.Println(err)
	}
	conf := config.Get()
	changelogItems, err := conf.Database.GetSubmissionChangelog(workspaceID, formID, submissionID)
	if err != nil {
		return changelogItems, err
	}
	return changelogItems, nil
}
