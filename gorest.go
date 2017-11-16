package gorest

import (
	"sync"
)

var wg sync.WaitGroup

//Start Start
func Start(startTest RestTest) Results {
	todoChan := make(chan RestTest, 10000)
	doneChan := make(chan RestTest, 10000)

	todoChan <- startTest
	wg.Add(1)

	for w := 1; w <= 20; w++ {
		go worker(&wg, w, todoChan, doneChan)
	}

	wg.Wait()

	close(todoChan)
	close(doneChan)

	results := Results{}
	for testResult := range doneChan {
		results.Add(testResult)
	}

	return results
}

func worker(wg *sync.WaitGroup, id int, todoChan chan RestTest, doneChan chan<- RestTest) {
	for todoTest := range todoChan {
		doneTest := DoAndVerify(todoTest)
		newTests := generateNewTestsFromResponse(doneTest)
		for _, newTest := range newTests {
			todoChan <- newTest
			wg.Add(1)
		}

		doneChan <- doneTest
		wg.Done()
	}
}
