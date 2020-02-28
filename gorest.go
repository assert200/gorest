package gorest

import (
	"sync"
)

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) (ResultTallys, []RestTestResult) {
	amountOfTests := len(restTests)

	testCh := make(chan RestTest, 10000)
	resultCh := make(chan RestTest, 100)

	var workWG sync.WaitGroup // All the work to be done
	var resultsWG sync.WaitGroup // All the results are collated

	for w := 1; w <= workers; w++ {
		go worker(&workWG, testCh, resultCh)
	}

	var allTests []RestTest
	resultsWG.Add(1)
	go resultWorker(&resultsWG, resultCh, &allTests)

	// initial seed of tests to execute
	for restTestIndex := 0; restTestIndex < amountOfTests; restTestIndex++ {
		workWG.Add(1)
		testCh <- restTests[restTestIndex]
	}

	// Wait for worker to finish, once it has, it is safe to close the channel
	workWG.Wait()
	close(testCh)

	// Now it is safe to close the result channel and wait for it to finish
	close(resultCh)
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

func worker(workWG *sync.WaitGroup, testCh chan RestTest, resultCh chan RestTest) {
	defer workWG.Done()

	nextChannels := executeTests(testCh, resultCh)
	recurseTests(nextChannels, testCh, resultCh)
}

func recurseTests(nextChannels []chan RestTest, testCh chan RestTest, resultCh chan RestTest) {
	for n := 0; n < len(nextChannels); n++ {
		nextChain := executeTests(nextChannels[n], resultCh)
		recurseTests(nextChain, testCh, resultCh)
	}
}

func executeTests(testCh chan RestTest, resultCh chan RestTest) []chan RestTest{

	var nextChannels []chan RestTest

	for test := range testCh {
		result := ExecuteAndVerify(test)

		if len(result.RestTestResult.Errors) == 0 {
			if result.Generator != nil {
				newTests := result.Generator(result)

				nextCh := make(chan RestTest, len(newTests)+1)
				for _, newTest := range newTests {
					nextCh <- newTest
				}
				nextChannels = append(nextChannels, nextCh)
			}
		}

		resultCh <- result
	}

	return nextChannels
}

func resultWorker(resultsWG *sync.WaitGroup, resultCh chan RestTest, allTests *[]RestTest) {
	defer resultsWG.Done()

	// Get all the results
	for test := range resultCh {
		*allTests = append(*allTests, test)
	}
}
