package gorest

import (
	"log"
	"sync"
)

// RunTest main entry point
func RunTest(restTests []RestTest, workers int) (ResultTallys, []RestTestResult) {
	amountOfTests := len(restTests)

	testCh := make(chan RestTest, 10000)
	resultCh := make(chan RestTest, 100)

	var workWG sync.WaitGroup // All the work to be done
	var resultsWG sync.WaitGroup // All the results are collated

	// initial seed of tests to execute
	for restTestIndex := 0; restTestIndex < amountOfTests; restTestIndex++ {
		testCh <- restTests[restTestIndex]
	}

	// set workers off to process
	for w := 1; w <= workers; w++ {
		workWG.Add(1)
		go worker(w, &workWG, testCh, resultCh)
	}

	// capture the results
	var allTests []RestTest
	resultsWG.Add(1)
	go resultWorker(&resultsWG, resultCh, &allTests)

	// Wait for workers to finish, once it has, it is safe to close the channel
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

func worker(wid int, workWG *sync.WaitGroup, testCh chan RestTest, resultCh chan RestTest) {
	defer workWG.Done()

	nextChannels := executeTests(wid, 1, testCh, resultCh)
	recurseTests(wid, 1, nextChannels, testCh, resultCh)
}

func recurseTests(wid int, depth int, nextChannels []chan RestTest, testCh chan RestTest, resultCh chan RestTest) {
	depth = depth+1
	for n := 0; n < len(nextChannels); n++ {
		nextChain := executeTests(wid, depth, nextChannels[n], resultCh)
		recurseTests(wid, depth, nextChain, testCh, resultCh)
	}
}

func executeTests(wid int, depth int, testCh chan RestTest, resultCh chan RestTest) ([]chan RestTest) {
	//log.Printf("wid=%d depth=%d: executing test...", wid, depth)

	var nextChannels []chan RestTest

	readData := true
	for readData {
		select {
		case test := <- testCh:
			log.Printf("wid=%d depth=%d: executing test '%s'", wid, depth, test.Description)
			result := ExecuteAndVerify(test)

			if len(result.RestTestResult.Errors) == 0 {
				if result.Generator != nil {
					newTests := result.Generator(result)

					nextCh := make(chan RestTest, len(newTests)+1)
					for _, newTest := range newTests {
						log.Printf("wid=%d depth=%d: queuing new child test '%s'", wid, depth, newTest.Description)
						nextCh <- newTest
					}
					nextChannels = append(nextChannels, nextCh)
				}
			}

			resultCh <- result
		default:
			readData = false
		}
	}

	//log.Printf("wid=%d depth=%d: finished executing all tests...", wid, depth)
	return nextChannels
}

func resultWorker(resultsWG *sync.WaitGroup, resultCh chan RestTest, allTests *[]RestTest) {
	defer resultsWG.Done()

	// Get all the results
	for test := range resultCh {
		*allTests = append(*allTests, test)
	}
}
