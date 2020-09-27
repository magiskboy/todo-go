package main

import (
	amqpbackend "github.com/RichardKnop/machinery/v1/backends/amqp"
	amqpbroker "github.com/RichardKnop/machinery/v1/brokers/amqp"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v2"
	"os"
	"runtime"
	"strconv"
)

func loadConfig() (*config.Config, error) {
	return config.NewFromEnvironment()
	// return &config.Config{
	// Broker:        "redis://localhost:6379/0",
	// DefaultQueue:  "my_queue",
	// ResultBackend: "redis://localhost:6379/0",
	// }, nil
}

func createServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	backend := amqpbackend.New(cnf)
	broker := amqpbroker.New(cnf)
	server := machinery.NewServer(cnf, broker, backend)

	tasks := map[string]interface{}{}

	return server, server.RegisterTasks(tasks)
}

func createConsumerTag() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	pid := os.Getpid()
	return hostname + " - " + strconv.Itoa(pid), nil
}

// StartWorker launch the worker
func StartWorker() error {
	server, err := createServer()
	if err != nil {
		return err
	}
	tag, err := createConsumerTag()
	if err != nil {
		return err
	}
	worker := server.NewWorker(tag, runtime.NumCPU())
	log.INFO.Println("Worker " + tag + " started")
	return worker.Launch()
}
