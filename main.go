package main

import (
	"fmt"

	"task_20.2_pipeline/pipeline"
	"task_20.2_pipeline/source"
)

func main() {
	exit := make(chan bool)
	defer close(exit)
	data := pipeline.Buffering(exit, pipeline.FilterMultipleNum(exit, pipeline.FilterNegativeNum(exit, source.AskData(exit))))

	for val := range data {
		fmt.Printf("Получены данные: %v\n", val)
	}
}
