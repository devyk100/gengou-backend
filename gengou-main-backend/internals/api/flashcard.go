package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gengou-main-backend/internals/database"
	flashcard_generate "gengou-main-backend/internals/flashcard-generate"
	"gengou-main-backend/internals/redis"
	db_connector "github.com/devyk100/gengou-db/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CreateFlashcardRequest struct {
	Title               string `json:"title"`
	GraduatingInterval  int32  `json:"graduatingInterval"`
	LearningSteps       string `json:"learningSteps"`
	NewCardsLimitPerDay int32  `json:"newCardsLimitPerDay"`
	EasyInterval        int32  `json:"easyInterval"`
}
type GetFlashcardDecksResponsePayload struct {
	Title              string `json:"title"`
	Id                 int32  `json:"id"`
	NewCards           int    `json:"newCards"`
	ReviewCards        int    `json:"reviewCards"`
	NewCardsLimit      int    `json:"newCardsLimit"`
	ReviewCardsLimit   int    `json:"reviewCardsLimit"`
	LearningSteps      string `json:"learningSteps"`
	GraduatingInterval int    `json:"graduatingInterval"`
}

type FlashcardReponse struct {
	Success bool `json:"success"`
}

type CreateFlashcardRequestPayload struct {
	FrontAudioUrl  string `json:"frontAudioUrl"`
	RearAudioUrl   string `json:"rearAudioUrl"`
	FrontImageUrl  string `json:"frontImageUrl"`
	RearImageUrl   string `json:"rearImageUrl"`
	FrontContent   string `json:"frontContent"`
	RearContent    string `json:"rearContent"`
	DeckId         int32  `json:"deckId"`
	ReviewFactor   int32  `json:"reviewFactor"`
	ReviewInterval int32  `json:"reviewInterval"`
}

type UpdateFlashcardRequestPayload struct {
	FrontAudioUrl  string `json:"frontAudioUrl"`
	RearAudioUrl   string `json:"rearAudioUrl"`
	FrontImageUrl  string `json:"frontImageUrl"`
	RearImageUrl   string `json:"rearImageUrl"`
	FrontContent   string `json:"frontContent"`
	RearContent    string `json:"rearContent"`
	DeckId         int32  `json:"deckId"`
	ReviewFactor   int32  `json:"reviewFactor"`
	ReviewInterval int32  `json:"reviewInterval"`
	LearningStepNo int32  `json:"learningStepNo"`
}

type GetReviewableFlashcardRequestPayload struct {
	Timestamp int64 `json:"timestamp"`
}

func FlashcardApiRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/get-decks", func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value("userIdString")
		userId := r.Context().Value("userIdString").(string)
		fmt.Println(val, userId)
		decks, err := database.Queries.GetFlashcardDecks(context.Background(), db_connector.GetFlashcardDecksParams{
			UserID: userId,
			Limit:  100,
			Offset: 0,
		})
		fmt.Println(decks, "Are the decks")

		if err != nil {
			panic(err.Error())
			return
		}
		var flashcardDecks []GetFlashcardDecksResponsePayload
		for _, deck := range decks {
			fmt.Println(deck.Title)
			flashcardDecks = append(flashcardDecks, GetFlashcardDecksResponsePayload{
				Title:              deck.Title,
				Id:                 deck.ID,
				NewCards:           2,
				ReviewCards:        3,
				NewCardsLimit:      int(deck.NewCardsLimitPerDay),
				ReviewCardsLimit:   int(deck.MaxReviewLimitPerDay),
				LearningSteps:      deck.LearningSteps,
				GraduatingInterval: int(deck.GraduatingInterval),
			})
		}
		responseBody, _ := json.Marshal(flashcardDecks)
		log.Println(val, "Is the user accessed on the context")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseBody)
		if err != nil {
			return
		}
	})

	router.Post("/create-deck", func(w http.ResponseWriter, r *http.Request) {
		//var body AuthenticationBody
		//err := json.NewDecoder(r.Body).Decode(&body)
		//if err != nil {
		//	fmt.Println(err.Error())
		//}
		//fmt.Println(body.Token)
		//fmt.Println(r.Header)
		//fmt.Println(r.Host)
		//fmt.Println(r.URL)
		//fmt.Println(r.Body)
		var Body CreateFlashcardRequest
		err := json.NewDecoder(r.Body).Decode(&Body)
		if err != nil {
			panic(err.Error())
		}
		userId := r.Context().Value("userIdString").(string)
		flashcardDeck, err := database.Queries.CreateFlashcardDeck(context.Background(), db_connector.CreateFlashcardDeckParams{
			Title:               Body.Title,
			UserID:              userId,
			GraduatingInterval:  Body.GraduatingInterval,
			LearningSteps:       Body.LearningSteps,
			NewCardsLimitPerDay: Body.NewCardsLimitPerDay,
			EasyInterval:        Body.EasyInterval,
		})
		if err != nil {
			panic(err.Error())
		}
		log.Println(userId, "Is the user accessed on the context", flashcardDeck)
		w.WriteHeader(http.StatusCreated)
		reponseBody, _ := json.Marshal(FlashcardReponse{Success: true})
		fmt.Println("The request to flashcard deck creation was successful")
		_, err = w.Write(reponseBody)
		if err != nil {
			panic(err.Error())
			return
		}
	})

	router.Post("/create-card", func(w http.ResponseWriter, r *http.Request) {
		var Body CreateFlashcardRequestPayload
		err := json.NewDecoder(r.Body).Decode(&Body)
		if err != nil {
			panic(err.Error())
		}
		//userId := r.Context().Value("userIdString").(string)
		flashcard, err := database.Queries.CreateFlashcard(context.Background(), db_connector.CreateFlashcardParams{
			FrontSide:      Body.FrontContent,
			RearSide:       Body.RearContent,
			DeckID:         Body.DeckId,
			ReviewFactor:   Body.ReviewFactor,
			ReviewInterval: Body.ReviewFactor,
			DueDate:        0,
			IsNew:          true,
			FrontAudio: pgtype.Text{
				String: Body.FrontAudioUrl,
				Valid:  true,
			},
			RearAudio: pgtype.Text{
				String: Body.RearAudioUrl,
				Valid:  true,
			},
			FrontImage: pgtype.Text{
				String: Body.FrontImageUrl,
				Valid:  true,
			},
			RearImage: pgtype.Text{
				String: Body.RearImageUrl,
				Valid:  true,
			},
		})
		if err != nil {
			panic(err.Error())
			return
		}
		log.Println(flashcard.FrontSide)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Received"))
		if err != nil {
			panic(err.Error())
			return
		}
	})
	router.Post("/get-review-cards/{deckId}-{limit}-{offset}", func(w http.ResponseWriter, r *http.Request) {
		var Body GetReviewableFlashcardRequestPayload
		err := json.NewDecoder(r.Body).Decode(&Body)
		deckId := chi.URLParam(r, "deckId")
		limit := chi.URLParam(r, "limit")
		offset := chi.URLParam(r, "offset")

		key := "deck-" + deckId
		deckIdInt, err := (strconv.Atoi(deckId))
		if err != nil {
			panic(err.Error())
		}
		limitInt, err := (strconv.Atoi(limit))
		if err != nil {
			panic(err.Error())
		}
		offsetInt, err := (strconv.Atoi(offset))
		if err != nil {
			panic(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		if val, err := redis.Instance.Get(key); val == "" || err != nil {
			fmt.Println("NOT FOUND")
			Deck, err2 := database.Queries.GetAFlashcardDeck(context.Background(), int32(deckIdInt))
			if err2 != nil {
				return
			}
			maxNewCardLimit := Deck.NewCardsLimitPerDay
			maxReviewCardLimit := Deck.MaxReviewLimitPerDay
			if err != nil {
				fmt.Println("Ereror", err.Error())
			}
			newFlashcard, err := database.Queries.GetNewFlashcard(context.Background(), db_connector.GetNewFlashcardParams{
				DeckID: int32(deckIdInt),
				Limit:  int32(limitInt),
				Offset: int32(offsetInt),
			})
			var newFlashcards []flashcard_generate.FlashcardType
			fmt.Println("max", maxNewCardLimit)
			for i := 0; i < len(newFlashcard) && i < int(maxNewCardLimit); i++ {
				newFlashcards = append(newFlashcards, flashcard_generate.FlashcardType{
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
			if err != nil {
				panic(err.Error())
			}
			reviewFlashcard, err := database.Queries.GetReviewFlashcard(context.Background(), db_connector.GetReviewFlashcardParams{
				DeckID:  int32(deckIdInt),
				Limit:   int32(limitInt),
				Offset:  int32(offsetInt),
				DueDate: flashcard_generate.EndOfDay(Body.Timestamp),
			})
			if err != nil {
				panic(err.Error())
			}
			var reviewFlashcards []flashcard_generate.FlashcardType
			for i := 0; i < len(reviewFlashcard) && i < int(maxReviewCardLimit); i++ {
				reviewFlashcards = append(reviewFlashcards, flashcard_generate.FlashcardType{
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

			graduateFlashcard, err := database.Queries.GetGraduateFlashcard(context.Background(), db_connector.GetGraduateFlashcardParams{
				DeckID: int32(deckIdInt),
				Limit:  int32(limitInt),
				Offset: int32(offsetInt),
			})
			if err != nil {
				panic(err.Error())
			}
			var graduateFlashcards []flashcard_generate.FlashcardType
			for i := 0; i < len(graduateFlashcard); i++ {
				graduateFlashcards = append(graduateFlashcards, flashcard_generate.FlashcardType{
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
			newFlashcardsJson, err := json.Marshal(newFlashcards)
			if err != nil {
				panic(err.Error())
			}
			reviewFlashcardsJson, err := json.Marshal(reviewFlashcards)
			if err != nil {
				panic(err.Error())
			}
			graduateFlashcardsJson, err := json.Marshal(graduateFlashcards)
			if err != nil {
				panic(err.Error())
			}
			key = "new-" + deckId
			err = redis.Instance.Set(key, string(newFlashcardsJson), time.Millisecond*time.Duration(flashcard_generate.EndOfDay(Body.Timestamp)-Body.Timestamp))
			if err != nil {
				panic(err.Error())
				return
			}
			key = "review-" + deckId
			err = redis.Instance.Set(key, string(reviewFlashcardsJson), time.Millisecond*time.Duration(flashcard_generate.EndOfDay(Body.Timestamp)-Body.Timestamp))
			if err != nil {
				panic(err.Error())
				return
			}
			key = "graduate-" + deckId
			err = redis.Instance.Set(key, string(graduateFlashcardsJson), time.Millisecond*time.Duration(flashcard_generate.EndOfDay(Body.Timestamp)-Body.Timestamp))
			if err != nil {
				panic(err.Error())
				return
			}
			response, err := json.Marshal(GetReviewableFlashcardResponsePayload{
				New:      newFlashcards,
				Review:   reviewFlashcards,
				Graduate: graduateFlashcards,
			})
			key = "deck-" + deckId
			err = redis.Instance.Set(key, "valid", time.Millisecond*time.Duration(flashcard_generate.EndOfDay(Body.Timestamp)-Body.Timestamp))
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(response)
			if err != nil {
				panic(err.Error())
				return
			}
		} else {
			key = "new-" + deckId
			var newFlashcard []flashcard_generate.FlashcardType
			var reviewFlashcard []flashcard_generate.FlashcardType
			var graduateFlashcard []flashcard_generate.FlashcardType
			newFlashcardJson, err := redis.Instance.Get(key)
			if err != nil {
				panic(err.Error())
				return
			}
			err = json.Unmarshal([]byte(newFlashcardJson), &newFlashcard)
			key = "review-" + deckId
			reviewFlashcardJson, err := redis.Instance.Get(key)
			if err != nil {
				panic(err.Error())
				return
			}
			err = json.Unmarshal([]byte(reviewFlashcardJson), &reviewFlashcard)
			key = "graduate-" + deckId
			graduateFlashcardJson, err := redis.Instance.Get(key)
			if err != nil {
				panic(err.Error())
				return
			}
			err = json.Unmarshal([]byte(graduateFlashcardJson), &graduateFlashcard)
			payload := GetReviewableFlashcardResponsePayload{
				New:      newFlashcard,
				Review:   reviewFlashcard,
				Graduate: graduateFlashcard,
			}

			// Marshal the payload to JSON
			response, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
				return
			}

			// Write response header and body
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(response)
			if err != nil {
				http.Error(w, "Error writing response", http.StatusInternalServerError)
				return
			}

		}

		//flashcard, err := database.Queries.GetReviewFlashcard(context.Background(), db_connector.GetReviewFlashcardParams{
		//	DeckID:  int32(deckIdInt),
		//	IsNew:   true,
		//	DueDate: 0,
		//	Limit:   int32(limitInt),
		//	Offset:  int32(offsetInt),
		//})
		//fmt.Println(len(flashcard), "Is the total number of flashcard", "The timestamp got is", Body.Timestamp)
		//redis.Instance.HSet("", "")
		//deckIdInt, err := (strconv.Atoi(deckId))
		//if err != nil {
		//	panic(err.Error())
		//}
		//limitInt, err := (strconv.Atoi(limit))
		//if err != nil {
		//	panic(err.Error())
		//}
		//offsetInt, err := (strconv.Atoi(offset))
		//if err != nil {
		//	panic(err.Error())
		//}
		//flashcard, err := database.Queries.GetAllFlashcards(context.Background(), db_connector.GetAllFlashcardsParams{
		//	DeckID: int32(deckIdInt),
		//	Limit:  int32(limitInt),
		//	Offset: int32(offsetInt),
		//})
		////fmt.Println(flashcard[0].FrontSide)
		//json, err := json.Marshal(flashcard)
		//if err != nil {
		//	panic(err.Error())
		//	return
		//}
		//w.WriteHeader(http.StatusOK)
		//_, err = w.Write([]byte("json"))
		//if err != nil {
		//	panic(err.Error())
		//	return
		//}
	})
	router.Get("/get-cards/{deckId}-{limit}-{offset}", func(w http.ResponseWriter, r *http.Request) {
		deckId := chi.URLParam(r, "deckId")
		limit := chi.URLParam(r, "limit")
		offset := chi.URLParam(r, "offset")
		deckIdInt, err := (strconv.Atoi(deckId))
		if err != nil {
			panic(err.Error())
		}
		limitInt, err := (strconv.Atoi(limit))
		if err != nil {
			panic(err.Error())
		}
		offsetInt, err := (strconv.Atoi(offset))
		if err != nil {
			panic(err.Error())
		}
		flashcard, err := database.Queries.GetAllFlashcards(context.Background(), db_connector.GetAllFlashcardsParams{
			DeckID: int32(deckIdInt),
			Limit:  int32(limitInt),
			Offset: int32(offsetInt),
		})
		//fmt.Println(flashcard[0].FrontSide)
		json, err := json.Marshal(flashcard)
		if err != nil {
			panic(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(json)
		if err != nil {
			panic(err.Error())
			return
		}
	})
	router.Post("/update-card", func(w http.ResponseWriter, r *http.Request) {
		var Body UpdateFlashcardRequestPayload
		err := json.NewDecoder(r.Body).Decode(&Body)
		if err != nil {
			panic(err.Error())
		}
		//userId := r.Context().Value("userIdString").(string)
		flashcard, err := database.Queries.UpdateFlashcard(context.Background(), db_connector.UpdateFlashcardParams{
			FrontSide:      Body.FrontContent,
			RearSide:       Body.RearContent,
			ReviewFactor:   Body.ReviewFactor,
			ReviewInterval: Body.ReviewFactor,
			DueDate:        0,
			IsNew:          true,
			FrontAudio: pgtype.Text{
				String: Body.FrontAudioUrl,
				Valid:  true,
			},
			RearAudio: pgtype.Text{
				String: Body.RearAudioUrl,
				Valid:  true,
			},
			FrontImage: pgtype.Text{
				String: Body.FrontImageUrl,
				Valid:  true,
			},
			RearImage: pgtype.Text{
				String: Body.RearImageUrl,
				Valid:  true,
			},
			LearningStepNo: Body.LearningStepNo,
		})
		if err != nil {
			panic(err.Error())
			return
		}
		log.Println(flashcard.FrontSide)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Received"))
		if err != nil {
			panic(err.Error())
			return
		}
	})
	return router
}
