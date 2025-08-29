// Package integration demonstrates integration testing patterns with GoWise.
// This example shows how to test database operations, HTTP services, and
// other integration scenarios using GoWise assertions.
package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"gowise/pkg/assertions"
)

// User represents a user entity for integration testing.
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

// MockDB simulates a database for integration testing.
type MockDB struct {
	users   map[int]User
	nextID  int
	mutex   sync.RWMutex
	latency time.Duration // Simulated database latency
}

func NewMockDB() *MockDB {
	return &MockDB{
		users:   make(map[int]User),
		nextID:  1,
		latency: 10 * time.Millisecond,
	}
}

func (db *MockDB) CreateUser(username, email string) (User, error) {
	// Simulate database latency
	time.Sleep(db.latency)

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if username == "" {
		return User{}, fmt.Errorf("username cannot be empty")
	}

	// Check for duplicate username
	for _, user := range db.users {
		if user.Username == username {
			return User{}, fmt.Errorf("username already exists: %s", username)
		}
	}

	user := User{
		ID:       db.nextID,
		Username: username,
		Email:    email,
		Active:   true,
		Created:  time.Now(),
	}

	db.users[db.nextID] = user
	db.nextID++

	return user, nil
}

func (db *MockDB) GetUser(id int) (User, error) {
	time.Sleep(db.latency)

	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, exists := db.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found: %d", id)
	}

	return user, nil
}

func (db *MockDB) UpdateUser(id int, updates map[string]interface{}) error {
	time.Sleep(db.latency)

	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, exists := db.users[id]
	if !exists {
		return fmt.Errorf("user not found: %d", id)
	}

	// Apply updates
	if username, ok := updates["username"].(string); ok {
		user.Username = username
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if active, ok := updates["active"].(bool); ok {
		user.Active = active
	}

	db.users[id] = user
	return nil
}

func (db *MockDB) DeleteUser(id int) error {
	time.Sleep(db.latency)

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.users[id]; !exists {
		return fmt.Errorf("user not found: %d", id)
	}

	delete(db.users, id)
	return nil
}

func (db *MockDB) ListUsers() ([]User, error) {
	time.Sleep(db.latency)

	db.mutex.RLock()
	defer db.mutex.RUnlock()

	users := make([]User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}

	return users, nil
}

// UserService provides business logic for user operations.
type UserService struct {
	db *MockDB
}

func NewUserService(db *MockDB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) RegisterUser(username, email string) (User, error) {
	// Validate input
	if username == "" {
		return User{}, fmt.Errorf("username is required")
	}
	if email == "" {
		return User{}, fmt.Errorf("email is required")
	}
	if !strings.Contains(email, "@") {
		return User{}, fmt.Errorf("invalid email format")
	}

	return s.db.CreateUser(username, email)
}

func (s *UserService) GetActiveUsers() ([]User, error) {
	allUsers, err := s.db.ListUsers()
	if err != nil {
		return nil, err
	}

	var activeUsers []User
	for _, user := range allUsers {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, nil
}

// HTTP Handler for REST API testing
func (s *UserService) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := s.RegisterUser(req.Username, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// TestDatabaseIntegration demonstrates database integration testing patterns.
func TestDatabaseIntegration(t *testing.T) {
	assert := assertions.New(t)
	db := NewMockDB()
	service := NewUserService(db)

	t.Run("UserRegistration", func(t *testing.T) {
		// Test successful user creation
		user, err := service.RegisterUser("alice", "alice@example.com")

		assert.NoError(err)
		assert.Equal(user.Username, "alice")
		assert.Equal(user.Email, "alice@example.com")
		assert.True(user.Active)
		assert.True(user.ID > 0)
		assert.True(time.Since(user.Created) < time.Minute)
	})

	t.Run("DuplicateUserValidation", func(t *testing.T) {
		// Create first user
		_, err := service.RegisterUser("bob", "bob@example.com")
		assert.NoError(err)

		// Attempt to create duplicate
		_, err = service.RegisterUser("bob", "bob2@example.com")
		assert.HasError(err)
		assert.ErrorContains(err, "already exists")
		assert.ErrorContains(err, "bob")
	})

	t.Run("InputValidation", func(t *testing.T) {
		// Test empty username
		_, err := service.RegisterUser("", "test@example.com")
		assert.HasError(err)
		assert.ErrorContains(err, "username is required")

		// Test empty email
		_, err = service.RegisterUser("testuser", "")
		assert.HasError(err)
		assert.ErrorContains(err, "email is required")

		// Test invalid email format
		_, err = service.RegisterUser("testuser", "invalid-email")
		assert.HasError(err)
		assert.ErrorContains(err, "invalid email format")
	})

	t.Run("ActiveUsersFiltering", func(t *testing.T) {
		// Create multiple users
		_, _ = service.RegisterUser("active1", "active1@example.com")
		_, _ = service.RegisterUser("active2", "active2@example.com")
		user3, _ := service.RegisterUser("inactive", "inactive@example.com")

		// Deactivate one user
		err := db.UpdateUser(user3.ID, map[string]interface{}{"active": false})
		assert.NoError(err)

		// Get active users
		activeUsers, err := service.GetActiveUsers()
		assert.NoError(err)
		assert.Len(activeUsers, 4) // alice, bob, active1, active2

		// Verify all returned users are active
		for _, user := range activeUsers {
			assert.True(user.Active)
		}
	})
}

// TestHTTPIntegration demonstrates HTTP service integration testing.
func TestHTTPIntegration(t *testing.T) {
	assert := assertions.New(t)
	db := NewMockDB()
	service := NewUserService(db)

	// Create HTTP test server
	mux := http.NewServeMux()
	mux.HandleFunc("/users", service.CreateUserHandler)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("SuccessfulUserCreation", func(t *testing.T) {
		payload := `{"username": "httpuser", "email": "httpuser@example.com"}`
		resp, err := http.Post(server.URL+"/users", "application/json",
			strings.NewReader(payload))

		assert.NoError(err)
		assert.Equal(resp.StatusCode, http.StatusCreated)
		assert.Equal(resp.Header.Get("Content-Type"), "application/json")

		defer resp.Body.Close()

		var user User
		err = json.NewDecoder(resp.Body).Decode(&user)
		assert.NoError(err)
		assert.Equal(user.Username, "httpuser")
		assert.Equal(user.Email, "httpuser@example.com")
		assert.True(user.Active)
	})

	t.Run("InvalidJSONRequest", func(t *testing.T) {
		payload := `{"username": "test", invalid json`
		resp, err := http.Post(server.URL+"/users", "application/json",
			strings.NewReader(payload))

		assert.NoError(err)
		assert.Equal(resp.StatusCode, http.StatusBadRequest)
		resp.Body.Close()
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Empty username
		payload := `{"username": "", "email": "test@example.com"}`
		resp, err := http.Post(server.URL+"/users", "application/json",
			strings.NewReader(payload))

		assert.NoError(err)
		assert.Equal(resp.StatusCode, http.StatusBadRequest)
		resp.Body.Close()
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/users")

		assert.NoError(err)
		assert.Equal(resp.StatusCode, http.StatusMethodNotAllowed)
		resp.Body.Close()
	})
}

// TestConcurrentAccess demonstrates testing concurrent operations.
func TestConcurrentAccess(t *testing.T) {
	assert := assertions.New(t)
	db := NewMockDB()
	service := NewUserService(db)

	t.Run("ConcurrentUserCreation", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		usersChan := make(chan User, numGoroutines)
		errorsChan := make(chan error, numGoroutines)

		// Create users concurrently
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				username := fmt.Sprintf("user%d", id)
				email := fmt.Sprintf("user%d@example.com", id)

				user, err := service.RegisterUser(username, email)
				if err != nil {
					errorsChan <- err
				} else {
					usersChan <- user
				}
			}(i)
		}

		wg.Wait()
		close(usersChan)
		close(errorsChan)

		// Collect results
		var users []User
		var errors []error

		for user := range usersChan {
			users = append(users, user)
		}

		for err := range errorsChan {
			errors = append(errors, err)
		}

		// All operations should succeed (no username conflicts)
		assert.Len(errors, 0)
		assert.Len(users, numGoroutines)

		// Verify all users have unique IDs and usernames
		userIDs := make(map[int]bool)
		usernames := make(map[string]bool)

		for _, user := range users {
			assert.True(user.ID > 0)
			assert.False(userIDs[user.ID])         // No duplicate IDs
			assert.False(usernames[user.Username]) // No duplicate usernames

			userIDs[user.ID] = true
			usernames[user.Username] = true
		}
	})
}

// TestAsyncOperations demonstrates testing asynchronous operations.
func TestAsyncOperations(t *testing.T) {
	assert := assertions.New(t)
	db := NewMockDB()

	// Increase latency to make async nature visible
	db.latency = 100 * time.Millisecond

	t.Run("EventualConsistency", func(t *testing.T) {
		var createdUser User
		var err error
		var mu sync.Mutex
		var completed bool

		// Start async user creation
		go func() {
			user, userErr := db.CreateUser("asyncuser", "async@example.com")
			mu.Lock()
			createdUser = user
			err = userErr
			completed = true
			mu.Unlock()
		}()

		// Use Eventually to wait for async operation completion
		assert.Eventually(func() bool {
			mu.Lock()
			defer mu.Unlock()
			return completed
		}, 1*time.Second, 50*time.Millisecond)

		// Verify successful creation (no need to lock here as Eventually ensures completion)
		mu.Lock()
		finalUser := createdUser
		finalErr := err
		mu.Unlock()
		
		assert.NoError(finalErr)
		assert.Equal(finalUser.Username, "asyncuser")
		assert.True(finalUser.ID > 0)
	})

	t.Run("NeverViolatesInvariant", func(t *testing.T) {
		// Create a user that should remain active
		user, err := db.CreateUser("stableuser", "stable@example.com")
		assert.NoError(err)

		// Start background process that shouldn't modify this user
		go func() {
			// Simulate background operations that shouldn't affect our user
			for i := 0; i < 5; i++ {
				time.Sleep(20 * time.Millisecond)
				db.CreateUser(fmt.Sprintf("bg%d", i), fmt.Sprintf("bg%d@example.com", i))
			}
		}()

		// Assert that our user never becomes inactive during background operations
		assert.Never(func() bool {
			fetchedUser, err := db.GetUser(user.ID)
			return err != nil || !fetchedUser.Active
		}, 200*time.Millisecond, 25*time.Millisecond)
	})

	t.Run("WithinTimeout", func(t *testing.T) {
		// Test that database operations complete within reasonable time
		assert.WithinTimeout(func() {
			_, err := db.CreateUser("timeouttest", "timeout@example.com")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}, 200*time.Millisecond) // Should complete within 200ms
	})
}
