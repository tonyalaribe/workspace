package actions

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/Jeffail/gabs"
	"gitlab.com/middlefront/workspace/config"
	"gitlab.com/middlefront/workspace/database"
	"gitlab.com/middlefront/workspace/storage"
)

//NewFormSubmission creates a new form submission
func NewFormSubmission(workspaceID, formID string, submission database.SubmissionData) error {

	conf := config.Get()
	formInfoByte, err := conf.Database.GetFormJSONBySlug(workspaceID, formID)
	if err != nil {
		return err
	}

	formInfo, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		return err
	}

	schema := formInfo.Path("jsonschema")
	for k, v := range submission.FormData {
		schemaObject := schema.Path("properties").Search(k)
		switch schemaObject.Path("type").Data().(string) {
		case "string":
			itemFormat := ""
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				var file storage.StowFile
				file, err = storage.UploadBase64(filepath.Join(workspaceID, formID, submission.SubmissionName), v.(string))
				if err != nil {
					log.Println(err)
				}
				submission.FormData[k] = file.URL
				break
			default:
				submission.FormData[k] = v.(string)
			}

		case "array":
			switch schemaObject.Path("items.type").Data().(string) {
			case "string":
				switch schemaObject.Path("items.format").Data().(string) {
				case "data-url":
					items := []string{}
					var file storage.StowFile
					for _, item := range v.([]interface{}) {
						file, err = storage.UploadBase64(filepath.Join(workspaceID, formID, submission.SubmissionName), item.(string))
						if err != nil {
							log.Println(err)
						}
						items = append(items, file.URL)
					}
					submission.FormData[k] = items
				}
			}
		case "integer":
			//Using type float64 due to compiler complaints when handling integer types
			submission.FormData[k] = submission.FormData[k].(float64)
		default:
			submission.FormData[k] = submission.FormData[k].(string)
		}
	}

	err = conf.Database.NewFormSubmission(workspaceID, formID, submission)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["workspaceID"] = workspaceID
	data["formID"] = formID
	data["event"] = string(database.NewSubmissionTriggerEvent)
	data["submission"] = submission

	TriggerEvent(workspaceID, formID, database.NewSubmissionTriggerEvent, data)

	return err
}

//UpdateSubmission updates a submission with givven workspaceID
func UpdateSubmission(workspaceID, formID string, submissionIDString string, newSubmission database.SubmissionData) error {
	submissionID, err := strconv.Atoi(submissionIDString)
	if err != nil {
		return err
	}

	//Get the previously updated data
	conf := config.Get()
	oldSubmission, err := conf.Database.GetFormSubmissionDetails(workspaceID, formID, submissionID)
	if err != nil {
		return err
	}

	oldSubmission.Status = newSubmission.Status
	oldSubmission.LastModified = newSubmission.LastModified
	oldSubmission.FormData = newSubmission.FormData

	formInfoByte, err := conf.Database.GetFormJSONBySlug(workspaceID, formID)
	if err != nil {
		return err
	}
	formMetaData, err := gabs.ParseJSON(formInfoByte)
	if err != nil {
		return err
	}

	schema := formMetaData.Path("jsonschema")
	for k, v := range newSubmission.FormData {
		schemaObject := schema.Path("properties").Search(k)
		switch schemaObject.Path("type").Data().(string) {
		case "string":
			itemFormat := ""
			if schemaObject.ExistsP("format") {
				itemFormat = schemaObject.Path("format").Data().(string)
			}
			switch itemFormat {
			case "data-uri", "data-url":
				//file formatting
				file, err := storage.UploadBase64(filepath.Join(workspaceID, formID, newSubmission.SubmissionName), v.(string))
				if err != nil {
					log.Println(err)
				}

				oldSubmission.FormData[k] = file.URL
				break
			default:
				oldSubmission.FormData[k] = v.(string)
			}

		case "array":
			switch schemaObject.Path("items.type").Data().(string) {
			case "string":
				switch schemaObject.Path("items.format").Data().(string) {
				case "data-url":
					items := []string{}
					for _, item := range v.([]interface{}) {
						file, err := storage.UploadBase64(filepath.Join(workspaceID, formID, newSubmission.SubmissionName), item.(string))
						if err != nil {
							log.Println(err)
						}
						items = append(items, file.URL)
					}
					oldSubmission.FormData[k] = items
				}
			}
		case "integer":
			//Using type float64 due to compiler complaints when handling integer types
			oldSubmission.FormData[k] = newSubmission.FormData[k].(float64)
		default:
			oldSubmission.FormData[k] = newSubmission.FormData[k].(string)
		}
	}

	conf.Database.UpdateFormSubmission(workspaceID, formID, submissionID, oldSubmission)
	if err != nil {
		return err
	}

	changelogItem := database.ChangelogItem{
		Created:      oldSubmission.LastModified,
		Notes:        newSubmission.SubmissionNotes,
		WorkspaceID:  workspaceID,
		SubmissionID: submissionID,
		FormID:       formID,
	}
	err = conf.Database.AddToSubmissionChangelog(workspaceID, formID, submissionID, changelogItem)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	data["workspaceID"] = workspaceID
	data["formID"] = formID

	data["old_submission"] = oldSubmission
	data["new_submission"] = newSubmission
	data["changelog"] = changelogItem

	if newSubmission.Status == "published" {
		data["event"] = string(database.ApproveSubmissionTriggerEvent)
		TriggerEvent(workspaceID, formID, database.ApproveSubmissionTriggerEvent, data)
	} else {
		data["event"] = string(database.UpdateSubmissionTriggerEvent)
		TriggerEvent(workspaceID, formID, database.UpdateSubmissionTriggerEvent, data)
	}
	TriggerEvent(workspaceID, formID, database.UpdateSubmissionTriggerEvent, data)

	if err != nil {
		log.Println(err)
	}

	return nil
}

//GetFormSubmissions returns all submissions under a given formID and workspace
func GetFormSubmissions(workspaceID, formID string) ([]database.SubmissionData, error) {
	conf := config.Get()
	submissions, err := conf.Database.GetFormSubmissions(workspaceID, formID)
	return submissions, err
}

//GetFormSubmissionDetails returns the submission details given the submission id
func GetFormSubmissionDetails(workspaceID, formID, submissionIDString string) (database.SubmissionData, error) {
	submissionData := database.SubmissionData{}
	submissionID, err := strconv.Atoi(submissionIDString)
	if err != nil {
		return submissionData, err
	}
	conf := config.Get()
	submissionData, err = conf.Database.GetFormSubmissionDetails(workspaceID, formID, submissionID)
	if err != nil {
		return submissionData, err
	}
	return submissionData, nil
}

//DeleteFormSubmission deletes a form submission given the submission id
func DeleteFormSubmission(workspaceID, formID, submissionIDString string) (database.SubmissionData, error) {
	submissionData := database.SubmissionData{}
	submissionID, err := strconv.Atoi(submissionIDString)
	if err != nil {
		return submissionData, err
	}
	conf := config.Get()
	submissionData, err = conf.Database.DeleteFormSubmission(workspaceID, formID, submissionID)
	if err != nil {
		return submissionData, err
	}

	data := make(map[string]interface{})
	data["workspaceID"] = workspaceID
	data["formID"] = formID
	data["event"] = string(database.DeleteSubmissionTriggerEvent)
	data["submission"] = submissionData

	TriggerEvent(workspaceID, formID, database.DeleteSubmissionTriggerEvent, data)

	if err != nil {
		log.Println(err)
	}

	return submissionData, nil
}
