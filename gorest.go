package gorest

import (
	"sync"
)

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) (ResultTallys, []RestTestResult) {
	amountOfTests := len(restTests)

	testCh := make(chan RestTest, 10000)
	resultCh := make(chan RestTest, 10000)

	var workWG sync.WaitGroup // All the work to be done
	var spiderWG sync.WaitGroup
	var resultsWG sync.WaitGroup

	for w := 1; w <= workers; w++ {
		go testWorker(&workWG, testCh, resultCh)
	}

	var allTests []RestTest
	resultsWG.Add(1)
	go resultWorker(&resultsWG, resultCh, &allTests)

	// initial seed of tests to execute
	for restTestIndex := 0; restTestIndex < amountOfTests; restTestIndex++ {
		workWG.Add(1)
		testCh <- restTests[restTestIndex]
	}

	// Wait for work to finish
	workWG.Wait()

	// Clost chanels to tell workers they are done
	close(testCh)

	// Close all the channels to stop all the spiders and result workers

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

func testWorker(workWG *sync.WaitGroup, testCh chan RestTest, resultCh chan RestTest) {
	for test := range testCh {
		result := ExecuteAndVerify(test)

		if len(result.RestTestResult.Errors) == 0 {
			if result.Generator != nil {
				newTests := result.Generator(result)

				for _, newTest := range newTests {
					workWG.Add(1)
					testCh <- newTest
				}
			}
		}

		resultCh <- result

		workWG.Done()
	}
}

func resultWorker(resultsWG *sync.WaitGroup, resultCh chan RestTest, allTests *[]RestTest) {
	defer resultsWG.Done()

	// Get all the results
	for test := range resultCh {
		*allTests = append(*allTests, test)
	}
}
