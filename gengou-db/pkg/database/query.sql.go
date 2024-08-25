// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const copyFlashcardDeck = `-- name: CopyFlashcardDeck :one
WITH NewDeck AS (
INSERT INTO "FlashcardDeck" (
    title,
    max_review_limit_per_day,
    graduating_interval,
    learning_steps,
    new_cards_limit_per_day,
    easy_interval
)
SELECT
    title,
    max_review_limit_per_day,
    graduating_interval,
    learning_steps,
    new_cards_limit_per_day,
    easy_interval
FROM "FlashcardDeck" AS old_deck
WHERE old_deck.id = $1
    RETURNING id AS new_deck_id
)
SELECT new_deck_id FROM NewDeck
`

func (q *Queries) CopyFlashcardDeck(ctx context.Context, id int32) (int32, error) {
	row := q.db.QueryRow(ctx, copyFlashcardDeck, id)
	var new_deck_id int32
	err := row.Scan(&new_deck_id)
	return new_deck_id, err
}

const copyFlashcardsForDeck = `-- name: CopyFlashcardsForDeck :many
WITH CopiedFlashcards AS (
INSERT INTO "Flashcard" (
    front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
    review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
)
SELECT
    front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
    review_factor, review_interval, due_date, unreviewed_priority_num, $2  -- new deck ID
FROM "Flashcard" AS old_flashcard
WHERE old_flashcard.deck_id = $1  -- old deck ID
    RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
              review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
)
SELECT id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id FROM CopiedFlashcards
`

type CopyFlashcardsForDeckParams struct {
	DeckID   int32
	DeckID_2 int32
}

type CopyFlashcardsForDeckRow struct {
	ID                    int32
	FrontSide             string
	RearSide              string
	FrontAudio            pgtype.Text
	RearAudio             pgtype.Text
	FrontImage            pgtype.Text
	RearImage             pgtype.Text
	ReviewFactor          int32
	ReviewInterval        int32
	DueDate               pgtype.Timestamp
	UnreviewedPriorityNum int32
	DeckID                int32
}

func (q *Queries) CopyFlashcardsForDeck(ctx context.Context, arg CopyFlashcardsForDeckParams) ([]CopyFlashcardsForDeckRow, error) {
	rows, err := q.db.Query(ctx, copyFlashcardsForDeck, arg.DeckID, arg.DeckID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CopyFlashcardsForDeckRow
	for rows.Next() {
		var i CopyFlashcardsForDeckRow
		if err := rows.Scan(
			&i.ID,
			&i.FrontSide,
			&i.RearSide,
			&i.FrontAudio,
			&i.RearAudio,
			&i.FrontImage,
			&i.RearImage,
			&i.ReviewFactor,
			&i.ReviewInterval,
			&i.DueDate,
			&i.UnreviewedPriorityNum,
			&i.DeckID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createCopyFlashcardDecKMapping = `-- name: CreateCopyFlashcardDecKMapping :one
INSERT INTO "FlashcardDeckToCopiers" (deck_id, user_id, copied_deck_id)
VALUES ($1, $3, $2) RETURNING id, deck_id, user_id, copied_deck_id
`

type CreateCopyFlashcardDecKMappingParams struct {
	DeckID       int32
	CopiedDeckID int32
	UserID       string
}

func (q *Queries) CreateCopyFlashcardDecKMapping(ctx context.Context, arg CreateCopyFlashcardDecKMappingParams) (FlashcardDeckToCopier, error) {
	row := q.db.QueryRow(ctx, createCopyFlashcardDecKMapping, arg.DeckID, arg.CopiedDeckID, arg.UserID)
	var i FlashcardDeckToCopier
	err := row.Scan(
		&i.ID,
		&i.DeckID,
		&i.UserID,
		&i.CopiedDeckID,
	)
	return i, err
}

const createFlashcard = `-- name: CreateFlashcard :one
INSERT INTO "Flashcard" (
                         front_side, rear_side, deck_id, review_factor, review_interval, due_date, unreviewed_priority_num
) VALUES (
          $1,$2, $3, $4, $5, $6, $7
         ) RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type CreateFlashcardParams struct {
	FrontSide             string
	RearSide              string
	DeckID                int32
	ReviewFactor          int32
	ReviewInterval        int32
	DueDate               pgtype.Timestamp
	UnreviewedPriorityNum int32
}

func (q *Queries) CreateFlashcard(ctx context.Context, arg CreateFlashcardParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, createFlashcard,
		arg.FrontSide,
		arg.RearSide,
		arg.DeckID,
		arg.ReviewFactor,
		arg.ReviewInterval,
		arg.DueDate,
		arg.UnreviewedPriorityNum,
	)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const createFlashcardDeck = `-- name: CreateFlashcardDeck :one
WITH new_deck AS (
INSERT INTO "FlashcardDeck" (title)
VALUES ($1)
    RETURNING id
    )
INSERT INTO "FlashcardDeckToEditors" (deck_id, user_id)
SELECT id, $2
FROM new_deck RETURNING id, deck_id, user_id
`

type CreateFlashcardDeckParams struct {
	Title  string
	UserID int32
}

func (q *Queries) CreateFlashcardDeck(ctx context.Context, arg CreateFlashcardDeckParams) (FlashcardDeckToEditor, error) {
	row := q.db.QueryRow(ctx, createFlashcardDeck, arg.Title, arg.UserID)
	var i FlashcardDeckToEditor
	err := row.Scan(&i.ID, &i.DeckID, &i.UserID)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "User"
WHERE user_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userID string) error {
	_, err := q.db.Exec(ctx, deleteUser, userID)
	return err
}

const flashcardReview = `-- name: FlashcardReview :one
UPDATE "Flashcard"
SET review_factor = $2, review_interval = $3, due_date = $4, unreviewed_priority_num = $5
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type FlashcardReviewParams struct {
	ID                    int32
	ReviewFactor          int32
	ReviewInterval        int32
	DueDate               pgtype.Timestamp
	UnreviewedPriorityNum int32
}

func (q *Queries) FlashcardReview(ctx context.Context, arg FlashcardReviewParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, flashcardReview,
		arg.ID,
		arg.ReviewFactor,
		arg.ReviewInterval,
		arg.DueDate,
		arg.UnreviewedPriorityNum,
	)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const getAFlashcard = `-- name: GetAFlashcard :many
SELECT id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
FROM "Flashcard"
WHERE deck_id = $1
LIMIT $2
OFFSET $3
`

type GetAFlashcardParams struct {
	DeckID int32
	Limit  int32
	Offset int32
}

func (q *Queries) GetAFlashcard(ctx context.Context, arg GetAFlashcardParams) ([]Flashcard, error) {
	rows, err := q.db.Query(ctx, getAFlashcard, arg.DeckID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Flashcard
	for rows.Next() {
		var i Flashcard
		if err := rows.Scan(
			&i.ID,
			&i.FrontSide,
			&i.RearSide,
			&i.FrontAudio,
			&i.RearAudio,
			&i.FrontImage,
			&i.RearImage,
			&i.ReviewFactor,
			&i.ReviewInterval,
			&i.DueDate,
			&i.UnreviewedPriorityNum,
			&i.DeckID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAUseWithId = `-- name: GetAUseWithId :one
SELECT id, user_id, name, user_type, email_id, phone, profile_picture, password, is_password_set, register_method, is_verified, is_user_id_set, is_phone_set, past_experiences FROM "User"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAUseWithId(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getAUseWithId, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.UserType,
		&i.EmailID,
		&i.Phone,
		&i.ProfilePicture,
		&i.Password,
		&i.IsPasswordSet,
		&i.RegisterMethod,
		&i.IsVerified,
		&i.IsUserIDSet,
		&i.IsPhoneSet,
		&i.PastExperiences,
	)
	return i, err
}

const getAUserWithUserId = `-- name: GetAUserWithUserId :one
SELECT id, user_id, name, user_type, email_id, phone, profile_picture, password, is_password_set, register_method, is_verified, is_user_id_set, is_phone_set, past_experiences FROM "User"
WHERE user_id = $1
`

func (q *Queries) GetAUserWithUserId(ctx context.Context, userID string) (User, error) {
	row := q.db.QueryRow(ctx, getAUserWithUserId, userID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.UserType,
		&i.EmailID,
		&i.Phone,
		&i.ProfilePicture,
		&i.Password,
		&i.IsPasswordSet,
		&i.RegisterMethod,
		&i.IsVerified,
		&i.IsUserIDSet,
		&i.IsPhoneSet,
		&i.PastExperiences,
	)
	return i, err
}

const getAllFlashcards = `-- name: GetAllFlashcards :many
SELECT id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
FROM "Flashcard"
WHERE deck_id = $1
LIMIT $2
OFFSET $3
`

type GetAllFlashcardsParams struct {
	DeckID int32
	Limit  int32
	Offset int32
}

func (q *Queries) GetAllFlashcards(ctx context.Context, arg GetAllFlashcardsParams) ([]Flashcard, error) {
	rows, err := q.db.Query(ctx, getAllFlashcards, arg.DeckID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Flashcard
	for rows.Next() {
		var i Flashcard
		if err := rows.Scan(
			&i.ID,
			&i.FrontSide,
			&i.RearSide,
			&i.FrontAudio,
			&i.RearAudio,
			&i.FrontImage,
			&i.RearImage,
			&i.ReviewFactor,
			&i.ReviewInterval,
			&i.DueDate,
			&i.UnreviewedPriorityNum,
			&i.DeckID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFlashcardDecks = `-- name: GetFlashcardDecks :many
SELECT id, title, max_review_limit_per_day, graduating_interval, learning_steps, new_cards_limit_per_day, easy_interval
FROM "FlashcardDeck"
WHERE id IN (
    SELECT deck_id
    FROM "FlashcardDeckToEditors"
    WHERE user_id = $1
)
ORDER BY id
    LIMIT $2
OFFSET $3
`

type GetFlashcardDecksParams struct {
	UserID int32
	Limit  int32
	Offset int32
}

func (q *Queries) GetFlashcardDecks(ctx context.Context, arg GetFlashcardDecksParams) ([]FlashcardDeck, error) {
	rows, err := q.db.Query(ctx, getFlashcardDecks, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FlashcardDeck
	for rows.Next() {
		var i FlashcardDeck
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.MaxReviewLimitPerDay,
			&i.GraduatingInterval,
			&i.LearningSteps,
			&i.NewCardsLimitPerDay,
			&i.EasyInterval,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsers = `-- name: GetUsers :many
SELECT id, user_id, name, user_type, email_id, phone, profile_picture, password, is_password_set, register_method, is_verified, is_user_id_set, is_phone_set, past_experiences FROM "User"
LIMIT 5
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.UserType,
			&i.EmailID,
			&i.Phone,
			&i.ProfilePicture,
			&i.Password,
			&i.IsPasswordSet,
			&i.RegisterMethod,
			&i.IsVerified,
			&i.IsUserIDSet,
			&i.IsPhoneSet,
			&i.PastExperiences,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertInstructor = `-- name: InsertInstructor :one
INSERT INTO "User" (
                    name, user_id, email_id, phone, past_experiences, user_type
) VALUES (
          $1, $2, $3, $4, $5, 'Instructor'
         ) RETURNING id, user_id, name, user_type, email_id, phone, profile_picture, password, is_password_set, register_method, is_verified, is_user_id_set, is_phone_set, past_experiences
`

type InsertInstructorParams struct {
	Name            string
	UserID          string
	EmailID         string
	Phone           string
	PastExperiences pgtype.Text
}

func (q *Queries) InsertInstructor(ctx context.Context, arg InsertInstructorParams) (User, error) {
	row := q.db.QueryRow(ctx, insertInstructor,
		arg.Name,
		arg.UserID,
		arg.EmailID,
		arg.Phone,
		arg.PastExperiences,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.UserType,
		&i.EmailID,
		&i.Phone,
		&i.ProfilePicture,
		&i.Password,
		&i.IsPasswordSet,
		&i.RegisterMethod,
		&i.IsVerified,
		&i.IsUserIDSet,
		&i.IsPhoneSet,
		&i.PastExperiences,
	)
	return i, err
}

const insertLearner = `-- name: InsertLearner :one
INSERT INTO "User" (
    name, user_id, email_id, phone, user_type
) VALUES (
             $1, $2, $3, $4, 'Learner'
         ) RETURNING id, user_id, name, user_type, email_id, phone, profile_picture, password, is_password_set, register_method, is_verified, is_user_id_set, is_phone_set, past_experiences
`

type InsertLearnerParams struct {
	Name    string
	UserID  string
	EmailID string
	Phone   string
}

func (q *Queries) InsertLearner(ctx context.Context, arg InsertLearnerParams) (User, error) {
	row := q.db.QueryRow(ctx, insertLearner,
		arg.Name,
		arg.UserID,
		arg.EmailID,
		arg.Phone,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.UserType,
		&i.EmailID,
		&i.Phone,
		&i.ProfilePicture,
		&i.Password,
		&i.IsPasswordSet,
		&i.RegisterMethod,
		&i.IsVerified,
		&i.IsUserIDSet,
		&i.IsPhoneSet,
		&i.PastExperiences,
	)
	return i, err
}

const updateFlashcardFrontAudio = `-- name: UpdateFlashcardFrontAudio :one
UPDATE "Flashcard"
SET front_audio = $2
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardFrontAudioParams struct {
	ID         int32
	FrontAudio pgtype.Text
}

func (q *Queries) UpdateFlashcardFrontAudio(ctx context.Context, arg UpdateFlashcardFrontAudioParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardFrontAudio, arg.ID, arg.FrontAudio)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const updateFlashcardFrontImage = `-- name: UpdateFlashcardFrontImage :one
UPDATE "Flashcard"
SET front_image = $2
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardFrontImageParams struct {
	ID         int32
	FrontImage pgtype.Text
}

func (q *Queries) UpdateFlashcardFrontImage(ctx context.Context, arg UpdateFlashcardFrontImageParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardFrontImage, arg.ID, arg.FrontImage)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const updateFlashcardFrontSide = `-- name: UpdateFlashcardFrontSide :one
UPDATE "Flashcard"
SET front_side = $2
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardFrontSideParams struct {
	ID        int32
	FrontSide string
}

func (q *Queries) UpdateFlashcardFrontSide(ctx context.Context, arg UpdateFlashcardFrontSideParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardFrontSide, arg.ID, arg.FrontSide)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const updateFlashcardRearAudio = `-- name: UpdateFlashcardRearAudio :one
UPDATE "Flashcard"
SET rear_audio = $2
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardRearAudioParams struct {
	ID        int32
	RearAudio pgtype.Text
}

func (q *Queries) UpdateFlashcardRearAudio(ctx context.Context, arg UpdateFlashcardRearAudioParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardRearAudio, arg.ID, arg.RearAudio)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const updateFlashcardRearImage = `-- name: UpdateFlashcardRearImage :one
UPDATE "Flashcard"
SET rear_image = $2
WHERE id = $1
RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardRearImageParams struct {
	ID        int32
	RearImage pgtype.Text
}

func (q *Queries) UpdateFlashcardRearImage(ctx context.Context, arg UpdateFlashcardRearImageParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardRearImage, arg.ID, arg.RearImage)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}

const updateFlashcardRearSide = `-- name: UpdateFlashcardRearSide :one
UPDATE "Flashcard"
SET rear_side = $2
WHERE id = $1 RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image, review_factor, review_interval, due_date, unreviewed_priority_num, deck_id
`

type UpdateFlashcardRearSideParams struct {
	ID       int32
	RearSide string
}

func (q *Queries) UpdateFlashcardRearSide(ctx context.Context, arg UpdateFlashcardRearSideParams) (Flashcard, error) {
	row := q.db.QueryRow(ctx, updateFlashcardRearSide, arg.ID, arg.RearSide)
	var i Flashcard
	err := row.Scan(
		&i.ID,
		&i.FrontSide,
		&i.RearSide,
		&i.FrontAudio,
		&i.RearAudio,
		&i.FrontImage,
		&i.RearImage,
		&i.ReviewFactor,
		&i.ReviewInterval,
		&i.DueDate,
		&i.UnreviewedPriorityNum,
		&i.DeckID,
	)
	return i, err
}
