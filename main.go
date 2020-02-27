package gorest

import (
	"fmt"
	"sync"
	"time"
)


const (
	SLEEP_FOR = 5 * time.Second
	MAX_SLEEPS = 20
)


// RunTest main entry point
func RunTest(restTests []RestTest, workers int) ResultTallys {
	amountOfTests := len(restTests)

	testCh := make(chan RestTest, 100)
	spiderCh := make(chan RestTest, 100)
	resultCh := make(chan RestTest, 100)

	var workerWG sync.WaitGroup
	for w := 1; w <= workers; w++ {
		workerWG.Add(1)
		go testWorker(workerWG, testCh, spiderCh)
	}

	var spiderWG sync.WaitGroup
	spiderWG.Add(1)
	go spiderWorker(spiderWG, testCh, spiderCh, resultCh)

	var resultsWG sync.WaitGroup
	resultsWG.Add(1)
	tally := &ResultTallys{}
	go resultWorker(resultsWG, resultCh, tally)

	// initial seed of tests to execute
	for restTestIndex := 0; restTestIndex < amountOfTests;  restTestIndex++ {
		testCh <- restTests[restTestIndex]
	}

	// Wait for workers to finish
	workerWG.Wait()
	close(testCh)
	close(spiderCh)
	close(resultCh)

	spiderWG.Wait()
	resultsWG.Wait()

	// Return all the results
	return *tally
}

func testWorker(wg sync.WaitGroup, testCh chan RestTest, spiderCh chan<- RestTest) {
	defer wg.Done()

	/*
	for test := range testCh {
		result := ExecuteAndVerify(test)
		spiderCh <- result
	}
	 */

	sleeps := 0
	for sleeps < MAX_SLEEPS {
		select {
		case test := <- testCh:
			sleeps = 0
			result := ExecuteAndVerify(test)
			spiderCh <- result
		default:
			fmt.Println("worker: no tests waiting...")
			time.Sleep(SLEEP_FOR)
			sleeps++
		}
	}

}

func spiderWorker(wg sync.WaitGroup, testCh chan<- RestTest, spiderCh chan RestTest, resultCh chan<- RestTest) {
	defer wg.Done()

	for result := range spiderCh {
		if len(result.RestTestResult.Errors) == 0 {
			if result.Generator != nil {
				newTests := result.Generator(result)

				for _, newTest := range newTests {
					testCh <- newTest
				}
			}
		}

		resultCh <- result
	}
}

func resultWorker(wg sync.WaitGroup, resultCh chan RestTest, resultTallys *ResultTallys) {
	defer wg.Done()

	// Get all the results
	for resultTally := range resultCh {
		resultTallys.Add(resultTally)
	}

}
