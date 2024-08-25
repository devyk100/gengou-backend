package gengou_db

import (
	"bufio"
	"context"
	"fmt"
	"github.com/devyk100/gengou-db/pkg/database"
	"github.com/devyk100/gengou-db/pkg/kafka_internal"
	"github.com/devyk100/gengou-db/pkg/redis_internal"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func database_trial() {

	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	//fmt.Println(dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer conn.Close(context.Background())

	queries := database.New(conn)

	//user, err := queries.InsertUser(ctx, database.InsertUserParams{
	//	Name:            "Yash",
	//	UserID:          "yashk095",
	//	EmailID:         "yash@dev.in",
	//	Phone:           "172873819247",
	//	PastExperiences: pgtype.Text{},
	//	UserType:        database.UserTypeLearner,
	//})
	//if err != nil {
	//	log.Fatal("failed to insert user", err.Error())
	//}
	//
	//fmt.Println(user.Name)
	err = queries.DeleteUser(ctx, "yashk095")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//Users, err := queries.GetUsers(ctx)
	//for _, user := range Users {
	//	fmt.Println(user.Name, user.UserID, user.EmailID)
	//}
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//User, err := queries.GetAUser(context.Background(), 1)
	//fmt.Println(User.UserID, err)
	//fmt.Println("Exiting")
	//var now time.Time
	//err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	//if err != nil {
	//	log.Fatal("failed to execute query", err)
	//}

	//fmt.Println(now)

}

var ctx = context.Background()

func redis_test() {
	err := godotenv.Load("./../../.env")
	dsn := os.Getenv("REDIS_URL")
	db, err := redis_internal.Init(dsn, time.Second*100)
	if err != nil {
		log.Fatal(err.Error())
	}
	now := time.Now()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your input: ")
	scanner.Scan()
	input := scanner.Text()
	if input == "p" {
		for a := 0; a < 100; a++ {
			db.Publish("something", "bla bla bla bhai"+strconv.Itoa(a))
			//scanner.Scan()
			//scanner.Text()
		}
	} else {
		var hell chan struct{}
		db.Subscribe("something", func(payload string) {
			fmt.Println(payload)
		}, &hell)
	}
	fmt.Println("You entered:", input)

	//fmt.Println(db.Get("foo2"))
	//fmt.Println(db.Get("foo"))
	//db.Set("foo2", "bar2")
	db.HSet("somehtingew1", "fewjao2f", "123132", "fnm21awoifawoeg", "2230")
	fmt.Println(db.HGet("somehting1", "fewjaof"))
	fmt.Println(time.Since(now), "Is the time taken")
	//db.Close()
}

func main() {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Println("Error loading .env file", err.Error())
	}
	kafkaUrl := os.Getenv("KAFKA_URL")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	kafkaUserName := os.Getenv("KAFKA_USERNAME")
	//redis_test()
	//kafka_internal.Producer("Whiteboard")
	instance := kafka_internal.KafkaProducer{}
	instance.CreateProducer("Whiteboard", kafkaUrl, kafkaUserName, kafkaPassword)
	instance.Produce("something 1")
	instance.Produce("something 2")
	instance.Produce("something 3")
	instance.CloseProducer()
	//instance := kafka_internal.KafkaConsumer{}
	//instance.CreateConsumer("Whiteboard", "user10", []string{kafkaUrl}, kafkaUserName, kafkaPassword)
	//c := instance.GetMessageChan()
	//instance.StartConsumer()
	//defer instance.CloseConsumer()
	//for {
	//	val := <-*c
	//	fmt.Println(string(val.Value))
	//}
}
