package main

import (
	"encoding/csv"
	redisbackend "github.com/RichardKnop/machinery/v1/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v1/brokers/redis"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v2"
	"io"
	"os"
	"runtime"
	"strconv"
)

func loadConfig() (*config.Config, error) {
	return &config.Config{
		Broker:        "redis://localhost:6379/0",
		DefaultQueue:  "my_queue",
		ResultBackend: "redis://localhost:6379/0",
	}, nil
}

func processCSVFile(filename string) error {
	file, err := os.Open("upload/" + filename)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		log.DEBUG.Println(record)
	}
	return nil
}

func createServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	backend := redisbackend.New(cnf, "localhost:6379", "", "", 0)
	broker := redisbroker.New(cnf, "localhost:6379", "", "", 0)
	server := machinery.NewServer(cnf, broker, backend)

	tasks := map[string]interface{}{
		"process_csv": processCSVFile,
	}

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
