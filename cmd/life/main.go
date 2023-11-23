package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/k6mil6/game/internal/application"
	"github.com/k6mil6/game/pkg/config"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// Exit завершает программу с заданным кодом
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	cfg := config.Config{
		Width:  10,
		Height: 10,
	}
	app := application.New(cfg)
	// Запускаем приложение
	if err := app.Run(ctx); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Processing cancelled.")
		default:
			log.Println("Application run error", err)
		}
		// Возвращаем значение, не равное нулю, чтобы обозначить ошибку
		return 1
	}
	// Выход без ошибок
	return 0
}

func run() error {
	if len(os.Args) != 4 {
		return fmt.Errorf("usage: %s <rows> <columns> <percentage>", os.Args[0])
	}

	rows := os.Args[1]
	convertedRows, err := strconv.Atoi(rows)
	if err != nil || convertedRows <= 0 {
		return fmt.Errorf("invalid rows value: %v", err)
	}

	columns := os.Args[2]
	convertedColumns, err := strconv.Atoi(columns)
	if err != nil || convertedColumns <= 0 {
		return fmt.Errorf("invalid columns value: %v", err)
	}

	percentage := os.Args[3]
	convertedPercentage, err := strconv.Atoi(percentage)
	if err != nil || convertedPercentage <= 0 || convertedPercentage > 100 {
		return fmt.Errorf("invalid percentage value: %v", err)
	}

	cfg := fmt.Sprintf("%sх%s %s%%", rows, columns, percentage)

	file, err := os.Create("config.txt")
	if err != nil {
		return fmt.Errorf("error creating config.txt: %v", err)
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}(file)

	_, err = file.WriteString(cfg)
	if err != nil {
		return fmt.Errorf("error writing to config.txt: %v", err)
	}

	return nil
}
