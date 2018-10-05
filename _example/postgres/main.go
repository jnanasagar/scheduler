package main

import (
	"io"
	"log"
	"time"

	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
)

var count = 0

var s scheduler.Scheduler

type tmp struct {
	scheduler            *scheduler.Scheduler
}

func (t *tmp ) TaskWithoutArgs() {
	count = count +1
	if count == 5{
		return
	}
	log.Println("TaskWithoutArgs is executed - ",count)
	if taskId, err := t.scheduler.RunAfter(10*time.Second, t.TaskWithoutArgs); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("With in the task, submitted task: %s", taskId)
	}
}

func TaskWithArgs(message string) {
	log.Println("TaskWithArgs is executed. message:", message)
}

func main() {
	store, err := storage.NewPostgresStorage(
		storage.PostgresDBConfig{
			DbURL: "postgresql://jnana:password@localhost:5432/scheduler?sslmode=disable",
		},
	)
	if err != nil {
		log.Fatalf("Couldn't create scheduler storage : %v", err)
	}

	// memStore := storage.NewMemoryStorage()

	test := &tmp{}

	s = scheduler.New(store)
	test.scheduler = &s

	go func(s scheduler.Scheduler, store io.Closer) {
		time.Sleep(time.Second * 50)
		log.Println("timeup, stopping scheduler")
		s.Stop()
	}(s, store)
	// Start a task without arguments
	if taskId, err := s.RunAfter(5*time.Second, test.TaskWithoutArgs); err != nil {
		log.Fatal(err)
	}else {
		log.Printf("submitted task: %s", taskId)
	}
	// Start a task with arguments
	// if _, err := s.RunEvery(5*time.Second, TaskWithArgs, "Hello from recurring task 1"); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // Start the same task as above with a different argument
	// if _, err := s.RunEvery(10*time.Second, TaskWithArgs, "Hello from recurring task 2"); err != nil {
	// 	log.Fatal(err)
	// }
	// defer s.Stop()
	// time.Sleep(10*time.Second)
	s.Start()
	s.Wait()
}
