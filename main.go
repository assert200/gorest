package gorest

import "time"

var executeTestCh, testResultCh chan RestTest

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) ResultTallys {
	amountOfTests := len(restTests)

	executeTestCh = make(chan RestTest, 100000)
	testResultCh = make(chan RestTest, 100000)

	for w := 1; w <= workers; w++ {
		go testWorker(executeTestCh, testResultCh)
	}

	for n := 1; n <= workers; n++ {
		go newTestWorker(executeTestCh, testResultCh)
	}

	for restTestIndex := 0; restTestIndex < amountOfTests;  restTestIndex++ {
		// Only set tests off in batches of size of testWorker group
		executeTestCh <- restTests[restTestIndex]
	}

	// REVISIT
	time.Sleep(5 * time.Minute)
	close(executeTestCh)
	close(testResultCh)

	resultTallys := ResultTallys{}
	for resultTally := range testResultCh {
		resultTallys.Add(resultTally)
	}

	return resultTallys
}

func testWorker(executeTestCh chan RestTest, testResultCh chan<- RestTest) {

	for test := range executeTestCh {
		result := ExecuteAndVerify(test)

		testResultCh <- result
	}
}

func newTestWorker(executeTestCh chan RestTest, testResultCh chan<- RestTest) {

	for result := range testResultCh {
		if len(result.RestTestResult.Errors) == 0 {
			if result.Generator != nil {
				newTests := result.Generator(result)

				for _, newTest := range newTests {
					executeTestCh <- newTest
				}
			}
		}
	}
}


