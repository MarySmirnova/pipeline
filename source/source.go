package source

import (
	"fmt"
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
				fmt.Println("Некорректный ввод")
				continue
			}
			if val == "exit" {
				exit <- true
				break
			}
			num, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("Введено некорректное число")
				continue
			}
			data <- num
		}
	}()
	return data
}
