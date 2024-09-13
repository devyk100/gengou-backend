package flashcard_generate

import (
	"context"
	"encoding/json"
	"fmt"
	"gengou-main-backend/internals/database"
	"gengou-main-backend/internals/redis"
	db_connector "github.com/devyk100/gengou-db/pkg/database"
	"github.com/jackc/pgx/v5"
	"strconv"
	"time"
)

type FlashcardType struct {
	ID             int32  `json:"id"`
	FrontSide      string `json:"frontSide"`
	RearSide       string `json:"rearSide"`
	FrontAudio     string `json:"frontAudio"`
	RearAudio      string `json:"rearAudio"`
	FrontImage     string `json:"frontImage"`
	RearImage      string `json:"rearImage"`
	ReviewFactor   int32  `json:"reviewFactor"`
	ReviewInterval int32  `json:"reviewInterval"`
	DueDate        int64  `json:"dueDate"`
	IsNew          bool   `json:"isNew"`
	DeckID         int32  `json:"deckId"`
	LearningStepNo int32  `json:"learningStepNo"`
}

type GetReviewableFlashcardResponsePayload struct {
	New      []FlashcardType `json:"new"`
	Review   []FlashcardType `json:"review"`
	Graduate []FlashcardType `json:"graduate"`
}

func EndOfDay(unixTimestampMs int64) int64 {
	// Convert the millisecond UNIX timestamp to a time.Time object
	t := time.Unix(0, unixTimestampMs*int64(time.Millisecond))

	// Set the time to 11:59:59.999 PM
	endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999000000, t.Location())

	// Convert back to UNIX timestamp in milliseconds
	return endOfDay.UnixNano() / int64(time.Millisecond)
}

func GenerateFlashcardsForReviews(Timestamp int64, deckId int) (error, []byte) {
	key := "deck-" + strconv.Itoa(deckId)
	_, _ = redis.Instance.Get(key)
	// Handling the initial case of the review cards for today generation.
	tx, err := database.Conn.Begin(context.Background())
	if err != nil {
		return err, []byte("")
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			panic(err.Error())
		}
	}(tx, context.Background())

	qtx := database.Queries.WithTx(tx)

	// Get the info of the deck
	Deck, err := qtx.GetAFlashcardDeck(context.Background(), int32(deckId))
	if err != nil {
		panic(err.Error())
	}
	maxNewCardLimit := Deck.NewCardsLimitPerDay
	maxReviewCardLimit := Deck.MaxReviewLimitPerDay

	// Create the unique review entry for today
	generated, err := qtx.CreateReviewGenerated(context.Background(), EndOfDay(Timestamp))
	if err != nil {
		return err, nil
	}

	// Get all the new flashcards possible defined by the limit which I have set using the deck parameters
	newFlashcard, err := qtx.GetNewFlashcard(context.Background(), db_connector.GetNewFlashcardParams{
		DeckID: int32(deckId),
		Limit:  maxNewCardLimit,
		Offset: 0,
	})
	var newFlashcards []FlashcardType
	fmt.Println("max", maxNewCardLimit)

	// For the required format without pgx types, i have to unwrap and send them again
	// In the same for loop we also create all the reviews that are there for today
	for i := 0; i < len(newFlashcard); i++ {
		_, err2 := qtx.CreateFlashcardReview(context.Background(), db_connector.CreateFlashcardReviewParams{
			CardID:            newFlashcard[i].ID,
			DeckID:            newFlashcard[i].DeckID,
			OldLearningStepNo: newFlashcard[i].LearningStepNo,
			OldIsNew:          newFlashcard[i].IsNew,
			OldDueDate:        newFlashcard[i].DueDate,
			IsReviewComplete:  false,
			ReviewID:          generated.ID,
		})
		if err2 != nil {
			return err2, nil
		}
		newFlashcards = append(newFlashcards, FlashcardType{
			ID:             newFlashcard[i].ID,
			FrontSide:      newFlashcard[i].FrontSide,
			RearSide:       newFlashcard[i].RearSide,
			FrontAudio:     newFlashcard[i].FrontAudio.String,
			RearAudio:      newFlashcard[i].RearAudio.String,
			FrontImage:     newFlashcard[i].FrontImage.String,
			RearImage:      newFlashcard[i].RearImage.String,
			ReviewFactor:   newFlashcard[i].ReviewFactor,
			ReviewInterval: newFlashcard[i].ReviewInterval,
			DueDate:        newFlashcard[i].DueDate,
			IsNew:          newFlashcard[i].IsNew,
			DeckID:         newFlashcard[i].DeckID,
			LearningStepNo: newFlashcard[i].LearningStepNo,
		})
	}

	// similar for review flashcards
	reviewFlashcard, err := qtx.GetReviewFlashcard(context.Background(), db_connector.GetReviewFlashcardParams{
		DeckID:  int32(deckId),
		Limit:   maxReviewCardLimit,
		Offset:  0,
		DueDate: EndOfDay(Timestamp),
	})
	if err != nil {
		panic(err.Error())
	}
	var reviewFlashcards []FlashcardType
	for i := 0; i < len(reviewFlashcard); i++ {
		_, err2 := qtx.CreateFlashcardReview(context.Background(), db_connector.CreateFlashcardReviewParams{
			CardID:            reviewFlashcard[i].ID,
			DeckID:            reviewFlashcard[i].DeckID,
			OldLearningStepNo: reviewFlashcard[i].LearningStepNo,
			OldIsNew:          reviewFlashcard[i].IsNew,
			OldDueDate:        reviewFlashcard[i].DueDate,
			IsReviewComplete:  false,
			ReviewID:          generated.ID,
		})
		if err2 != nil {
			return err2, nil
		}
		reviewFlashcards = append(reviewFlashcards, FlashcardType{
			ID:             reviewFlashcard[i].ID,
			FrontSide:      reviewFlashcard[i].FrontSide,
			RearSide:       reviewFlashcard[i].RearSide,
			FrontAudio:     reviewFlashcard[i].FrontAudio.String,
			RearAudio:      reviewFlashcard[i].RearAudio.String,
			FrontImage:     reviewFlashcard[i].FrontImage.String,
			RearImage:      reviewFlashcard[i].RearImage.String,
			ReviewFactor:   reviewFlashcard[i].ReviewFactor,
			ReviewInterval: reviewFlashcard[i].ReviewInterval,
			DueDate:        reviewFlashcard[i].DueDate,
			IsNew:          reviewFlashcard[i].IsNew,
			DeckID:         reviewFlashcard[i].DeckID,
			LearningStepNo: reviewFlashcard[i].LearningStepNo,
		})
	}

	graduateFlashcard, err := qtx.GetGraduateFlashcard(context.Background(), db_connector.GetGraduateFlashcardParams{
		DeckID: int32(deckId),
		Limit:  100,
		Offset: 0,
	})

	if err != nil {
		panic(err.Error())
	}
	var graduateFlashcards []FlashcardType
	for i := 0; i < len(graduateFlashcard); i++ {
		_, err2 := qtx.CreateFlashcardReview(context.Background(), db_connector.CreateFlashcardReviewParams{
			CardID:            graduateFlashcard[i].ID,
			DeckID:            graduateFlashcard[i].DeckID,
			OldLearningStepNo: graduateFlashcard[i].LearningStepNo,
			OldIsNew:          graduateFlashcard[i].IsNew,
			OldDueDate:        graduateFlashcard[i].DueDate,
			IsReviewComplete:  false,
			ReviewID:          generated.ID,
		})
		if err2 != nil {
			return err2, nil
		}
		graduateFlashcards = append(graduateFlashcards, FlashcardType{
			ID:             graduateFlashcard[i].ID,
			FrontSide:      graduateFlashcard[i].FrontSide,
			RearSide:       graduateFlashcard[i].RearSide,
			FrontAudio:     graduateFlashcard[i].FrontAudio.String,
			RearAudio:      graduateFlashcard[i].RearAudio.String,
			FrontImage:     graduateFlashcard[i].FrontImage.String,
			RearImage:      graduateFlashcard[i].RearImage.String,
			ReviewFactor:   graduateFlashcard[i].ReviewFactor,
			ReviewInterval: graduateFlashcard[i].ReviewInterval,
			DueDate:        graduateFlashcard[i].DueDate,
			IsNew:          graduateFlashcard[i].IsNew,
			DeckID:         graduateFlashcard[i].DeckID,
			LearningStepNo: graduateFlashcard[i].LearningStepNo,
		})
	}

	fmt.Println("graduate", len(graduateFlashcards), "review", len(reviewFlashcards), "new", len(newFlashcards))
	//newFlashcardsJson, err := json.Marshal(newFlashcards)
	//if err != nil {
	//	panic(err.Error())
	//}
	//reviewFlashcardsJson, err := json.Marshal(reviewFlashcards)
	//if err != nil {
	//	panic(err.Error())
	//}
	//graduateFlashcardsJson, err := json.Marshal(graduateFlashcards)
	//if err != nil {
	//	panic(err.Error())
	//}
	response, err := json.Marshal(GetReviewableFlashcardResponsePayload{
		New:      newFlashcards,
		Review:   reviewFlashcards,
		Graduate: graduateFlashcards,
	})
	return tx.Commit(context.Background()), response
}
