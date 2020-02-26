package gorest

import (
	"fmt"
	"sync"
)

var executeTestCh, testResultCh chan RestTest

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) ResultTallys {
	amountOfTests := len(restTests)

	executeTestCh = make(chan RestTest, 100000)
	testResultCh = make(chan RestTest, 100000)
	defer close(executeTestCh)
	defer close(testResultCh)

	var testWorkerWG sync.WaitGroup
	for w := 1; w <= workers; w++ {
		testWorkerWG.Add(1)
		go testWorker(testWorkerWG, executeTestCh, testResultCh)
	}

	var newTestWG sync.WaitGroup
	for n := 1; n <= workers; n++ {
		newTestWG.Add(1)
		go newTestWorker(testWorkerWG, executeTestCh, testResultCh)
	}

	for restTestIndex := 0; restTestIndex < amountOfTests;  restTestIndex++ {
		// Only set tests off in batches of size of testWorker group
		executeTestCh <- restTests[restTestIndex]
	}

	// Wait for workers to finish
	testWorkerWG.Wait()
	newTestWG.Wait()

	fmt.Println("Main waking up....")

	// Get all the results
	resultTallys := ResultTallys{}
	for resultTally := range testResultCh {
		resultTallys.Add(resultTally)
	}

	return resultTallys
}

func testWorker(wg sync.WaitGroup, executeTestCh chan RestTest, testResultCh chan<- RestTest) {
	defer wg.Done()

	fmt.Println("Starting testWorker")
	for test := range executeTestCh {
		result := ExecuteAndVerify(test)

		fmt.Println("testWorker: Adding result to channel")
		testResultCh <- result
	}

	fmt.Println("Finished testWorker")
}

func newTestWorker(wg sync.WaitGroup, executeTestCh chan<- RestTest, testResultCh chan RestTest) {
	defer wg.Done()

	fmt.Println("Starting newTestWorker")
	for result := range testResultCh {
		if len(result.RestTestResult.Errors) == 0 {
			if result.Generator != nil {
				newTests := result.Generator(result)

				for _, newTest := range newTests {
					fmt.Println("newTestWorker: Adding newTest to channel")
					executeTestCh <- newTest
				}
			}
		}
	}

	fmt.Println("Finished newTestWorker")
}


