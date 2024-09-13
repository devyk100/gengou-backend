-- name: GetAUseWithId :one
SELECT * FROM "User"
WHERE id = $1 LIMIT 1;

-- name: GetAUserWithUserId :one
SELECT * FROM "User"
WHERE user_id = $1;

-- name: GetUsers :many
SELECT * FROM "User"
LIMIT 5;

-- name: InsertInstructor :one
INSERT INTO "User" (
                    name, user_id, email_id, phone, past_experiences, user_type
) VALUES (
          $1, $2, $3, $4, $5, 'Instructor'
         ) RETURNING *;

-- name: InsertLearner :one
INSERT INTO "User" (
    name, user_id, email_id, phone, user_type
) VALUES (
             $1, $2, $3, $4, 'Learner'
         ) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "User"
WHERE user_id = $1;

-- name: CreateFlashcardDeck :one
WITH new_deck AS (
INSERT INTO "FlashcardDeck" (title,max_review_limit_per_day,graduating_interval,learning_steps,new_cards_limit_per_day,easy_interval)
VALUES ($1,$3,$4,$5,$6,$7)
    RETURNING id
    )
INSERT INTO "FlashcardDeckToEditors" (deck_id, user_id)
SELECT id, $2
FROM new_deck RETURNING *;

-- name: GetAFlashcardDeck :one
    SELECT * FROM "FlashcardDeck" WHERE id = $1;


-- name: CreateFlashcard :one
INSERT INTO "Flashcard" (
                             front_side, rear_side, deck_id, review_factor, review_interval, due_date, is_new, front_audio, rear_audio, front_image, rear_image, learning_step_no
) VALUES (
          $1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 0
         ) RETURNING *;

-- name: UpdateFlashcard :one
UPDATE "Flashcard"
SET front_side = $2, rear_side = $3, front_audio = $4, rear_audio = $5,  front_image = $6, rear_image = $7, review_factor = $8, review_interval = $9, is_new = $10, due_date = $11, learning_step_no = $12
WHERE id = $1 RETURNING *;

-- name: GetAllFlashcards :many
SELECT *
FROM "Flashcard"
WHERE deck_id = $1
LIMIT $2
OFFSET $3;

-- name: GetAFlashcard :many
SELECT *
FROM "Flashcard"
WHERE deck_id = $1
LIMIT $2
OFFSET $3;

-- name: GetFlashcardDecks :many
SELECT *
FROM "FlashcardDeck"
WHERE id IN (
    SELECT deck_id
    FROM "FlashcardDeckToEditors"
    WHERE user_id = $1
)
ORDER BY id
    LIMIT $2
OFFSET $3;


-- name: CopyFlashcardDeck :one
WITH NewDeck AS (
INSERT INTO "FlashcardDeck" (
    title,
    max_review_limit_per_day,
    graduating_interval,
    learning_steps,
    new_cards_limit_per_day,
    easy_interval,
    learning_step_no
)
SELECT
    title,
    max_review_limit_per_day,
    graduating_interval,
    learning_steps,
    new_cards_limit_per_day,
    easy_interval,
    learning_step_no
FROM "FlashcardDeck" AS old_deck
WHERE old_deck.id = $1
    RETURNING id AS new_deck_id
)
SELECT new_deck_id FROM NewDeck;

-- name: CreateCopyFlashcardDecKMapping :one
INSERT INTO "FlashcardDeckToCopiers" (deck_id, user_id, copied_deck_id)
VALUES ($1, $3, $2) RETURNING *;

-- name: CopyFlashcardsForDeck :many
WITH CopiedFlashcards AS (
INSERT INTO "Flashcard" (
    front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
    review_factor, review_interval, due_date, is_new, deck_id, learning_step_no
)
SELECT
    front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
    review_factor, review_interval, due_date, is_new, learning_step_no, $2  -- new deck ID
FROM "Flashcard" AS old_flashcard
WHERE old_flashcard.deck_id = $1  -- old deck ID
    RETURNING id, front_side, rear_side, front_audio, rear_audio, front_image, rear_image,
              review_factor, review_interval, due_date, is_new, deck_id
)
SELECT * FROM CopiedFlashcards;




-- name: FlashcardReview :one
UPDATE "Flashcard"
SET review_factor = $2, review_interval = $3, due_date = $4, is_new = $5
WHERE id = $1 RETURNING *;


-- name: GetReviewFlashcard :many
SELECT *
FROM "Flashcard"
WHERE deck_id = $1
  AND is_new = false
  AND due_date <= $2
    AND learning_step_no = -1
    LIMIT $3
OFFSET $4;

-- name: GetNewFlashcard :many
SELECT *
FROM "Flashcard"
WHERE deck_id = $1
AND is_new = true
    AND learning_step_no = 0
LIMIT $2
OFFSET $3;

-- name: GetGraduateFlashcard :many
SELECT *
FROM "Flashcard"
WHERE deck_id = $1
  AND is_new = true
    AND learning_step_no > 0
    LIMIT $2
OFFSET $3;


-- name: CreateReviewGenerated :one
INSERT INTO "ReviewGenerated" (
    date
) VALUES (
             $1
         ) RETURNING *;


-- name: CreateFlashcardReview :one
INSERT INTO "DailyCardReviews" (  "card_id" ,
                                  "deck_id",
                                  "old_learning_step_no",
                                  "old_is_new",
                                  "old_due_date",
                                  "is_review_complete", "review_id")
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;