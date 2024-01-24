package testrunner

// TestRunner is the interface that wraps the basic test runner functionality
type TestRunner interface {
	DiscoverTests() ([]string, error)
	ExecuteTests(tests []string) error
	ReportResults() error
}

// NewTestRunner returns a new TestRunner
func NewTestRunner() TestRunner {
	return &testRunner{}
}

// testRunner is the implementation of the TestRunner interface
type testRunner struct {
	testDiscoverer TestDiscoverer
	testRunner     TestRunner
	testIsolator   TestIsolator
	testFixturer   TestFixturer
	testCLI        TestCLI
}

// DiscoverTests discovers the tests to be executed
func (tr *testRunner) DiscoverTests() ([]string, error) {
	return tr.testDiscoverer.DiscoverTests()
}

// ExecuteTests executes the tests
func (tr *testRunner) ExecuteTests(tests []string) error {
	return tr.testRunner.ExecuteTests(tests)
}

// ReportResults reports the test results
func (tr *testRunner) ReportResults() error {
	return tr.testRunner.ReportResults()
}

// TestDiscoverer is the interface that wraps the test discovery functionality
type TestDiscoverer interface {
	DiscoverTests() ([]string, error)
}

// NewTestDiscoverer returns a new TestDiscoverer
func NewTestDiscoverer() TestDiscoverer {
	return &testDiscoverer{}
}

// testDiscoverer is the implementation of the TestDiscoverer interface
type testDiscoverer struct {
}

// DiscoverTests discovers the tests to be executed
func (td *testDiscoverer) DiscoverTests() ([]string, error) {
	return []string{}, nil
}
