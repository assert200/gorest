package gorest

import (
	"sync"
)

var wg sync.WaitGroup
var firstRestTests []RestTest

// Prime Prime
func Prime(newRestTest RestTest) {
	firstRestTests = append(firstRestTests, newRestTest)
}

//Start Start
func Start(workers int) Results {
	todoChan := make(chan RestTest, 100000)
	doneChan := make(chan RestTest, 100000)

	for _, firstRestTest := range firstRestTests {
		todoChan <- firstRestTest
		wg.Add(1)
	}

	firstRestTests = nil

	for w := 1; w <= workers; w++ {
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

		if len(doneTest.Errors) == 0 {
			if todoTest.Generator != nil {
				newTests := todoTest.Generator(doneTest)

				for _, newTest := range newTests {
					todoChan <- newTest
					wg.Add(1)
				}
			}
		}

		doneChan <- doneTest
		wg.Done()
	}
}
