package scripts_test

import (
	"testing"

	"github.com/wuttinanhi/code-judge-system/configs"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeUpdateWithTestcase(t *testing.T) {
	configs.LoadConfig()

	db := databases.NewMySQLDatabase()
	testServiceKit := services.CreateServiceKit(db)
	testServiceKit.KafkaService.OverriddenHost("localhost:9094")

	db.Migrator().DropTable(
		entities.Challenge{},
		entities.ChallengeTestcase{},
		entities.Submission{},
		entities.SubmissionTestcase{},
	)

	databases.StartMigration(db)

	// get a user
	user, err := testServiceKit.UserService.FindUserByID(1)
	if err != nil {
		t.Fatal(err)
	}

	// user, err := testServiceKit.UserService.Register("test@example.com", "testpassword", "testuser")
	// if err != nil {
	// 	t.Error(err)
	// }

	challenge, err := testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
		User:        user,
		Testcases: []*entities.ChallengeTestcase{
			{ID: 0, Input: "1", ExpectedOutput: "1", LimitMemory: 1, LimitTimeMs: 1, ActionFlag: "create"},
			{ID: 0, Input: "2", ExpectedOutput: "2", LimitMemory: 2, LimitTimeMs: 2, ActionFlag: "create"},
		},
	})
	if err != nil {
		t.Error(err)
	}

	challenge.Name = "UPDATED 1"
	challenge.Description = "UPDATED 1"
	challenge.Testcases = append(challenge.Testcases,
		&entities.ChallengeTestcase{
			ID:             0,
			Input:          "ADD",
			ExpectedOutput: "ADD",
			LimitMemory:    1,
			LimitTimeMs:    1,
			ActionFlag:     "create",
		},
		&entities.ChallengeTestcase{
			ID:             0,
			Input:          "ADD",
			ExpectedOutput: "ADD",
			LimitMemory:    1,
			LimitTimeMs:    1,
			ActionFlag:     "create",
		},
		&entities.ChallengeTestcase{
			ID:             0,
			Input:          "ADD",
			ExpectedOutput: "ADD",
			LimitMemory:    1,
			LimitTimeMs:    1,
			ActionFlag:     "create",
		},
	)

	testServiceKit.ChallengeService.UpdateChallengeWithTestcase(challenge)
	if err != nil {
		t.Error(err)
	}

	// deleteByDatabaseID(challenge, 4)
	UpdateByDatabaseID(challenge, 1, "UPDATED", "UPDATED")
	UpdateByDatabaseID(challenge, 5, "UPDATED", "UPDATED")

	testServiceKit.ChallengeService.UpdateChallengeWithTestcase(challenge)
	if err != nil {
		t.Error(err)
	}
}

func DeleteByDatabaseID(challenge *entities.Challenge, id uint) {
	challenge.Testcases[id-1].ActionFlag = "delete"
}

func UpdateByDatabaseID(challenge *entities.Challenge, id uint, input string, expectedOutput string) {
	challenge.Testcases[id-1].ActionFlag = "update"
	challenge.Testcases[id-1].Input = input
	challenge.Testcases[id-1].ExpectedOutput = expectedOutput
}
