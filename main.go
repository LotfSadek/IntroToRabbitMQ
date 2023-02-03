package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

func main() {
	// Postgre DB connection
	connStr := "user=postgres password=1234rewQ host=localhost port=5432 dbname=RabbitMQLogger sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// logger, _ := zap.NewDevelopment()
	// manager := usago.NewChannelManager(
	// 	"amqp://guest:guest@localhost:5672/",
	// 	logger,
	// )

	// bldr := usago.NewChannelBuilder().WithQueue(
	// 	"golang-queue1",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// ).WithConfirms(true)

	// chnl, err := manager.NewChannel(*bldr)
	// if err != nil {
	// 	fmt.Printf("failed to create channel")
	// 	return
	// }

	// // consume
	// consumer, err := chnl.RegisterConsumer(
	// 	"golang-queue1",
	// 	"",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// msgs := <-consumer

	// // publish
	// body := "Hello World!"
	// sno, err := chnl.Publish(
	// 	"Notifications",
	// 	"golang-queue1",
	// 	false, // mandatory
	// 	false, // immediate
	// 	amqp091.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte(body),
	// 	},
	// )
	// fmt.Printf("The SNO is %v", sno)
	// New Development Creates a suggared logger -> heavy on I/O ops
	// time		Log Level 	Message
	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}
	logger.Info("app starting up")
	// map[string]interface{}
	// Dict<string,Any>
	logger.Info("Connecting to RabbitMQ")
	// var a error
	// a= new error
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// ALWAYS MAKE SURE TO HANDLE THE ERROR
	if err != nil {
		logger.Error("error Creating Connection", zap.Error(err))
		panic(err)
	}
	// You can assume that the connection is valid
	defer conn.Close() // we can defer here because we wil be blocking

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("error Creating Channel", zap.Error(err)) // log.err + panic
	}
	// // We can assume ch is valid
	defer ch.Close()

	// // TODO: Bind an exchange to the queue, you might have to use a diffrenet name for the queue
	q, err := ch.QueueDeclare(
		"golang-queue1",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal("error Creating Channel", zap.Error(err)) // log.err + panic
	}

	if err := ch.QueueBind(
		q.Name,          // queue name
		"#",             // binding key
		"Notifications", // exchange
		false,           // no-wait
		nil,             // arguments
	); err != nil {
		logger.Fatal("Failed to bind a queue: %v", zap.Error(err))
	}
	// q, err := ch.QueueDeclare(
	// 	"golang-queue", // name
	// 	false,          // durable
	// 	false,          // delete when unused
	// 	false,          // exclusive
	// 	false,          // no-wait
	// 	nil,            // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Fatal("error Creating Channel", zap.Error(err)) // log.err + panic
	}

	// channel diffrent from RabbitMQ Channel
	// go's chan which is an abstraction over os Pipes "<- chan amqp.delivery"
	// <- means that its a recieve only channel
	// chan means its a channel
	// amqp.delivery: message type of the channel
	//=========================================================================
	// go + ANYTHING = means running on a seperate thread
	// func() {} () = anon function
	// youre creating a function that runs on a seperate thread

	go func() {
		for {
			msg := <-msgs // pulling a message from a channel
			// if the channel is empty then execution is blocked forever  :)
			// if not empty it will pull one
			body := msg.Body
			now := time.Now()
			// the "business logic"
			if len(body) > 0 {
				logger.Info("Message", zap.String("body", string(msg.Body)))
				_, err = db.Exec("INSERT INTO logs(date,body) VALUES($1,$2)", now.Format("2006-01-02 15:04:05"), string(msg.Body))
				if err != nil {
					log.Fatalf("Error inserting data: %v", err)
					msg.Nack(false, true)
					continue
				}
				fmt.Println("Data inserted successfully!")
				msg.Ack(true)
				// success ack - not nack
				//
			}
		}
	}()

	select {}
	//  blocks until it recieves a signal
	//  since no signal is provided it will block forever
	// failOnError(err, "Failed to register a consumer")

	// forever := make(chan bool)

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()

	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
}

// commit to github
// install postgresql server using docker - create a db that takes on the log data
// replace the current implementation of rabbitmq to follow whats done on usago
// persist each log to the DB
