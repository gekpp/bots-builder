package main

import (
	"fmt"

	"github.com/gekpp/bots-builder/questionnaire"
	"github.com/gekpp/env"
	"github.com/google/uuid"
)

var (
	argQuestionnaireID = env.MustString("QNR_ID")
	qnrService         = questionnaire.NewDummy(uuid.MustParse(argQuestionnaireID))
)

func main() {
	fmt.Printf("Bot for questionnaire id=%s started\n", argQuestionnaireID)

	// ctx := context.Background()

	// qnrService.Start(ctx,)
	// qnrService.Answer()
}
