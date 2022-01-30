package source

import (
	"fmt"
	"log"
	"strconv"
)

//AskData запрашивает данные из консоли,
//поддерживает команду exit для выхода из программы
func AskData(exit chan bool) <-chan int {
	data := make(chan int)
	fmt.Println("Введите данные: \n(команда exit для выхода)")
	go func() {
		defer close(data)
		var val string
		for {
			_, err := fmt.Scanln(&val)
			if err != nil {
				log.Println("incorrect input")
				fmt.Println("Некорректный ввод")
				continue
			}
			if val == "exit" {
				log.Println("exit")
				exit <- true
				break
			}
			num, err := strconv.Atoi(val)
			if err != nil {
				log.Println("incorrect input integer")
				fmt.Println("Введено некорректное число")
				continue
			}
			log.Println("entered number ", num)
			data <- num
		}
	}()
	return data
}
