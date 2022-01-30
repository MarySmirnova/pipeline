package pipeline

import (
	"log"
	"sync"
	"time"
)

//BufferSize определяет размер буфера
const BufferSize int = 5

//BufferClearTime определяет время очистки буфера
const BufferClearTime time.Duration = (10 * time.Second)

type Buffer struct {
	values   []int
	size     int
	start    int
	end      int
	full     bool
	bufMutex sync.Mutex
}

//NewBuffer конструктор нового буфера
func NewBuffer(size int) *Buffer {
	return &Buffer{make([]int, size), size, 0, 0, false, sync.Mutex{}}
}

//Push добавляет в буфер значение
func (b *Buffer) Push(val int) {
	b.bufMutex.Lock()
	defer b.bufMutex.Unlock()

	b.values[b.end] = val
	if b.end == b.size-1 {
		b.end = 0
	} else {
		b.end++
	}

	if b.end == b.start {
		b.full = true
	}
}

//Get возвращает из буфера значение, первое в очереди
func (b *Buffer) Get() int {
	b.bufMutex.Lock()
	defer b.bufMutex.Unlock()

	val := b.values[b.start]
	if b.start == b.size-1 {
		b.start = 0
	} else {
		b.start++
	}

	if b.start == b.end {
		b.full = false
	}
	return val
}

//Clean очищает буфер, возвращает все содержимое
func (b *Buffer) Clean() []int {
	bufValues := make([]int, 0, b.size)

	for i := 0; i < b.size; i++ {
		if b.start == b.end && !b.full {
			break
		}
		val := b.Get()
		bufValues = append(bufValues, val)
	}
	return bufValues
}

//Buffering функция буферизации данных
func Buffering(exit <-chan bool, input <-chan int) <-chan int {
	buf := NewBuffer(BufferSize)
	output := make(chan int)

	//отправляет содержимое буфера в канал
	checkBuf := func() {
		for _, bufVal := range buf.Clean() {
			output <- bufVal
		}
	}

	//добавляет входящие значения в буфер,
	//при переполнении буфера очищает его
	//и отправляет значения в канал
	go func() {
		defer close(output)
		for {
			select {
			case <-exit:
				return
			case _, ok := <-input:
				if !ok {
					return
				}
			default:
			}

			select {
			case val := <-input:
				log.Println("number", val, "is added to the buffer")
				buf.Push(val)
				if buf.full {
					log.Println("buffer full, clearing")
					checkBuf()
				}
			case <-exit:
				return
			}
		}
	}()

	//очищает буфер и выводит содержимое в канал
	//раз в промежуток времени, определенный BufferClearTime
	go func() {
		defer close(output)
		for {
			select {
			case <-exit:
				return
			case _, ok := <-input:
				if !ok {
					return
				}
			default:
			}

			select {
			case <-time.After(BufferClearTime):
				log.Println("it's time to clear the buffer, clearing")
				checkBuf()
			case <-exit:
				return
			}
		}
	}()

	return output
}
