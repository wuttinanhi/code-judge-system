package repositories

import "github.com/wuttinanhi/code-judge-system/entities"

// AddTestcase implements ChallengeRepository.
func (r *challengeRepository) AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error) {
	err := r.db.Model(challenge).Association("Testcases").Append(testcase)
	return testcase, err
}

// FindTestcaseByID implements ChallengeRepository.
func (r *challengeRepository) FindTestcaseByID(id uint) (testcase *entities.ChallengeTestcase, err error) {
	result := r.db.First(&testcase, id)
	return testcase, result.Error
}

// AllTestcases implements ChallengeRepository.
func (r *challengeRepository) AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error) {
	err = r.db.Model(challenge).Association("Testcases").Find(&testcases)
	return testcases, err
}

// UpdateTestcase implements ChallengeRepository.
func (r *challengeRepository) UpdateTestcase(testcase *entities.ChallengeTestcase) error {
	result := r.db.Save(testcase)
	return result.Error
}

// DeleteTestcase implements ChallengeRepository.
func (r *challengeRepository) DeleteTestcase(testcase *entities.ChallengeTestcase) error {
	result := r.db.Delete(testcase)
	return result.Error
}
