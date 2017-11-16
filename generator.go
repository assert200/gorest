package gorest

// A generator creates new tests from responses from existing tests
type generator func(restTestResponse RestTest) (newTests []RestTest)

var generators map[string]generator

// RegisterGenerator registers new generators that will trigger if description
func RegisterGenerator(key string, gen generator) {
	if generators == nil {
		generators = make(map[string]generator)
	}
	generators[key] = gen
}

// Take a completed test and see if you can generate more (spider time)
func generateNewTestsFromResponse(restTestResponse RestTest) []RestTest {
	if gen, ok := generators[restTestResponse.Description]; ok {
		return gen(restTestResponse)
	}

	return []RestTest{}
}
