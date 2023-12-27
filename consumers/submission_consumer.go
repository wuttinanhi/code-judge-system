package consumers

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
	"github.com/wuttinanhi/code-judge-system/services"
)

func StartSubmissionConsumer(serviceKit *services.ServiceKit) {
	kafkaHost := viper.GetString("KAFKA_HOST")
	if kafkaHost == "" {
		log.Fatal("KAFKA_HOST is not set")
	}

	topicName := viper.GetString("KAFKA_SUBMISSION_TOPIC")
	if topicName == "" {
		log.Fatal("SUBMISSION_TOPIC is not set")
	}

	groupID := viper.GetString("KAFKA_SUBMISSION_GROUP")
	if groupID == "" {
		log.Fatal("SUBMISSION_GROUP is not set")
	}

	if serviceKit.KafkaService == nil {
		log.Fatal("Kafka service is not initialized")
	}

	if !serviceKit.KafkaService.IsTopicExist(topicName) {
		log.Fatal("Topic does not exist")
	}

	messageC, errorC := serviceKit.KafkaService.Consume(topicName, groupID)

	log.Println("Start consuming submission topic...")

	for {
		select {
		case message := <-messageC:
			log.Println("Receiving submission ID:", message)

			// parse message as submissionID
			submissionID, err := strconv.ParseUint(message, 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			// get submission
			submission, err := serviceKit.SubmissionService.GetSubmissionByID(uint(submissionID))
			if err != nil {
				log.Println(err)
				continue
			}

			// process submission
			submission, err = serviceKit.SubmissionService.ProcessSubmission(submission)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("Submission processed:", submission.ID)
		case err := <-errorC:
			log.Println(err)
		}
	}
}
