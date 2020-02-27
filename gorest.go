package gorest

import (
	"sync"
	"time"
)

const (
	Timeout = 60 * time.Second
)

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) (ResultTallys, []RestTestResult) {
	amountOfTests := len(restTests)

	testCh := make(chan RestTest, 1000)
	spiderCh := make(chan RestTest, 1000)
	resultCh := make(chan RestTest, 1000)

	var workerWG sync.WaitGroup
	for w := 1; w <= workers; w++ {
		workerWG.Add(1)
		go testWorker(&workerWG, testCh, spiderCh)
	}

	var spiderWG sync.WaitGroup
	spiderWG.Add(1)
	go spiderWorker(&spiderWG, testCh, spiderCh, resultCh)

	var resultsWG sync.WaitGroup
	var allTests []RestTest
	resultsWG.Add(1)
	go resultWorker(&resultsWG, resultCh, &allTests)

	// initial seed of tests to execute
	for restTestIndex := 0; restTestIndex < amountOfTests; restTestIndex++ {
		testCh <- restTests[restTestIndex]
	}

	// Wait for workers to finish - exit condition is no new tests in last 60 seconds, and no items in the spideCh
	workerWG.Wait()

	// Close all the channels to stop all the spiders and result workers
	close(testCh)
	close(spiderCh)
	close(resultCh)

	spiderWG.Wait()
	resultsWG.Wait()

	// Return all the results
	var results []RestTestResult
	tally := ResultTallys{}
	for _, test := range allTests {
		results = append(results, test.RestTestResult)
		tally.Add(test)
	}
	return tally, results
}

func testWorker(wg *sync.WaitGroup, testCh chan RestTest, spiderCh chan RestTest) {
	defer wg.Done()

	for true {
		select {
		case test := <-testCh:
			result := ExecuteAndVerify(test)
			spiderCh <- result

		case <-time.After(Timeout):
			if len(spiderCh) == 0 {
				// we exit when we have had no new tests in the last Timeout
				// AND there are nothing left in the spider (so not going to get any more)
				return
			}
			// otherwise wait for some more tests
		}
	}
}

func spiderWorker(wg *sync.WaitGroup, testCh chan RestTest, spiderCh chan RestTest, resultCh chan RestTest) {
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

func resultWorker(wg *sync.WaitGroup, resultCh chan RestTest, allTests *[]RestTest) {
	defer wg.Done()

	// Get all the results
	for test := range resultCh {
		*allTests = append(*allTests, test)
	}
}
