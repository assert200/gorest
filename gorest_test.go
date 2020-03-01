package gorest

import (
	"fmt"
	"testing"
)

func googleRequst() RestRequest {
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
	restTest1.RestRequest = googleRequst()
	restTest1.Generator = secondTestGenerator
	return restTest1
}
func secondTestGenerator(restTestResponse RestTest) []RestTest {
	restTest2a := RestTest{}
	restTest2a.Description = "2a Google"
	restTest2a.ExpectedStatusCode = 200
	restTest2a.RestRequest = googleRequst()
	restTest2a.Generator = thirdTestGenerator

	restTest2b := RestTest{}
	restTest2b.Description = "2b Google"
	restTest2b.ExpectedStatusCode = 200
	restTest2b.RestRequest = googleRequst()
	restTest2b.Generator = thirdTestGenerator

	restTests := []RestTest{}
	restTests = append(restTests, restTest2a, restTest2b)
	fmt.Println(restTests)
	return restTests
}

func thirdTestGenerator(restTestResponse RestTest) []RestTest {
	restTest3 := RestTest{}
	restTest3.Description = "3 Google"
	restTest3.ExpectedStatusCode = 200
	restTest3.RestRequest = googleRequst()

	restTests := []RestTest{}
	restTests = append(restTests, restTest3)
	return restTests

}

func TestGorest(t *testing.T) {
	restTests := []RestTest{}

	restTests = append(restTests, firstTest())
	resultTallys, allResults := RunTest(restTests, 2)

	if len(resultTallys) != 4 {
		t.Errorf("resultTallys wrong")
	}

	fmt.Println(len(allResults))
	if len(allResults) != 5 {
		t.Errorf("allResults Wrong")
	}
}
