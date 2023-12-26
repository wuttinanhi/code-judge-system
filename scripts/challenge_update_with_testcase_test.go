package scripts_test

import (
	"fmt"
	"testing"

	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeUpdateWithTestcase(t *testing.T) {
	db := databases.NewMySQLDatabase()
	testServiceKit := services.CreateServiceKit(db)
	testServiceKit.KafkaService.OverriddenHost("localhost:9094")

	db.Migrator().DropTable(
		entities.User{},
		entities.Challenge{},
		entities.ChallengeTestcase{},
		entities.Submission{},
		entities.SubmissionTestcase{},
	)

	databases.StartMigration(db)

	// get a user
	user, err := testServiceKit.UserService.Register("test@example.com", "testpassword", "testuser")
	if err != nil {
		t.Error(err)
	}

	challenge, err := testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
		User:        user,
		Testcases: []*entities.ChallengeTestcase{
			{Input: "1", ExpectedOutput: "1", LimitMemory: 100, LimitTimeMs: 1000},
			{Input: "2", ExpectedOutput: "2", LimitMemory: 200, LimitTimeMs: 2000},
		},
	})
	if err != nil {
		t.Error(err)
	}

	challenge.Testcases = []*entities.ChallengeTestcase{
		{Input: "1", ExpectedOutput: "1", LimitMemory: 100, LimitTimeMs: 1000},
		{Input: "3", ExpectedOutput: "3", LimitMemory: 300, LimitTimeMs: 3000},
	}

	testServiceKit.ChallengeService.UpdateChallenge(challenge)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(challenge.Testcases)
}
