package mocks

// MockTestDiscoverer is a mock implementation of the TestDiscoverer interface
type MockTestDiscoverer struct {
}

// DiscoverTests mocks the test discovery functionality
func (m *MockTestDiscoverer) DiscoverTests() ([]string, error) {
	return []string{}, nil
}

// MockTestRunner is a mock implementation of the TestRunner interface
type MockTestRunner struct {
}

// ExecuteTests mocks the test execution functionality
func (m *MockTestRunner) ExecuteTests(tests []string) error {
	return nil
}

// ReportResults mocks the test result reporting functionality
func (m *MockTestRunner) ReportResults() error {
	return nil
}

// MockTestIsolator is a mock implementation of the TestIsolator interface
type MockTestIsolator struct {
}

// IsolateTests mocks the test isolation functionality
func (m *MockTestIsolator) IsolateTests(tests []string) error {
	return nil
}

// MockTestFixturer is a mock implementation of the TestFixturer interface
type MockTestFixturer struct {
}

// SetupFixtures mocks the fixture setup functionality
func (m *MockTestFixturer) SetupFixtures() error {
	return nil
}

// TeardownFixtures mocks the fixture teardown functionality
func (m *MockTestFixturer) TeardownFixtures() error {
	return nil
}

// MockTestCLI is a mock implementation of the TestCLI interface
type MockTestCLI struct {
}

// ParseCLI mocks the command-line interface functionality
func (m *MockTestCLI) ParseCLI() error {
	return nil
}

// MockTestParallelizer is a mock implementation of the TestParallelizer interface
type MockTestParallelizer struct {
}

// ParallelizeTests mocks the parallel execution functionality
func (m *MockTestParallelizer) ParallelizeTests(tests []string) error {
	return nil
}

// MockTestCoverageAnalyzer is a mock implementation of the TestCoverageAnalyzer interface
type MockTestCoverageAnalyzer struct {
}
