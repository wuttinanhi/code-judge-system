package scripts_test

import (
	"fmt"
	"testing"

	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSubmissionProcess(t *testing.T) {
	db := databases.NewMySQLDatabase()
	testServiceKit := services.CreateServiceKit(db)
	testServiceKit.KafkaService.OverriddenHost("localhost:9094")

	submission, err := testServiceKit.SubmissionService.GetSubmissionByID(1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(submission.ID)
	fmt.Println(len(submission.SubmissionTestcases))

	testServiceKit.SubmissionService.ProcessSubmission(submission)
}
