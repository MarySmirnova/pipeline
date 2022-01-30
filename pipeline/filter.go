package pipeline

import "log"

//FilterNegativeNum отбрасывает отрицательные значения
func FilterNegativeNum(exit <-chan bool, input <-chan int) <-chan int {
	withoutNeg := make(chan int)
	go func() {
		defer close(withoutNeg)
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
				if val >= 0 {
					log.Println("number ", val, " passed the filter of negative numbers")
					withoutNeg <- val
				}
			case <-exit:
				return
			}
		}
	}()
	return withoutNeg
}

//FilterMultipleNum отбрасывает значения, не кратные 3, включая 0
func FilterMultipleNum(exit <-chan bool, input <-chan int) <-chan int {
	withoutMult := make(chan int)
	go func() {
		defer close(withoutMult)
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
				if val%3 == 0 && val != 0 {
					log.Println("number ", val, " passed the filter of multiples of three")
					withoutMult <- val
				}
			case <-exit:
				return
			}
		}
	}()
	return withoutMult
}
