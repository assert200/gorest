package gorest

import (
	"gotest.tools/assert"
	"testing"
)

func googleRequest() RestRequest {
	restRequest := NewRestRequest()
	restRequest.URL.Scheme = "https"
	restRequest.URL.Path = "www.google.com"
	restRequest.Method = "GET"
	return restRequest
}

func firstTest() RestTest {
	restTest1 := RestTest{}
	restTest1.Description = "1 Google"
	restTest1.ExpectedStatusCode = 200
	restTest1.RestRequest = googleRequest()
	restTest1.Generator = secondTestGenerator
	return restTest1
}
func secondTestGenerator(restTestResponse RestTest) []RestTest {
	restTest2a := RestTest{}
	restTest2a.Description = "2a Google"
	restTest2a.ExpectedStatusCode = 200
	restTest2a.RestRequest = googleRequest()
	restTest2a.Generator = thirdTestGenerator

	restTest2b := RestTest{}
	restTest2b.Description = "2b Google"
	restTest2b.ExpectedStatusCode = 200
	restTest2b.RestRequest = googleRequest()
	restTest2b.Generator = thirdTestGenerator

	var restTests []RestTest
	restTests = append(restTests, restTest2a, restTest2b)
	//fmt.Println(restTests)
	return restTests
}

func thirdTestGenerator(restTestResponse RestTest) []RestTest {
	restTest3 := RestTest{}
	restTest3.Description = "3 Google"
	restTest3.ExpectedStatusCode = 200
	restTest3.RestRequest = googleRequest()

	restTests := []RestTest{}
	restTests = append(restTests, restTest3)
	return restTests

}

func TestGorest(t *testing.T) {
	var restTests []RestTest

	restTests = append(restTests, firstTest())
	resultTallys, allResults := RunTest(restTests, 2)

	assert.Assert(t, len(resultTallys) == 4, "resultTallys wrong: %d", len(resultTallys))
	assert.Assert(t, len(allResults) == 5, "allResults Wrong: %d", len(allResults))
}
