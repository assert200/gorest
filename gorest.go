package gorest

import (
	"sync"
)

var todoChan, doneChan chan RestTest

// Test away
func Test(restTests []RestTest, workers int) ResultTallys {
	var wg sync.WaitGroup
	amountOfTests := len(restTests)

	todoChan = make(chan RestTest, 100000)
	doneChan = make(chan RestTest, 100000)

	for w := 1; w <= workers; w++ {
		go worker(&wg, w, todoChan, doneChan)
	}

	restTestIndex := 0
	for restTestIndex < amountOfTests {
		// Only set tests off in batches of size of testWorker group
		for i := 0; i < workers; i++ {
			todoChan <- restTests[restTestIndex]
			wg.Add(1)

			restTestIndex++

			if restTestIndex == amountOfTests {
				break
			}
		}

		// Wait utill these tests are complete till kicking off more tests
		wg.Wait()
	}

	close(todoChan)
	close(doneChan)

	resultTallys := ResultTallys{}
	for resultTally := range doneChan {
		resultTallys.Add(resultTally)
	}

	return resultTallys
}

func worker(wg *sync.WaitGroup, id int, todoChan chan RestTest, doneChan chan<- RestTest) {
	for todoTest := range todoChan {
		doneTest := DoAndVerify(todoTest)

		if len(doneTest.RestTestResult.Errors) == 0 {
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
