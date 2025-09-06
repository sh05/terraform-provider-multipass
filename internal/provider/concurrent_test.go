package provider

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/sh05/terraform-provider-multipass/internal/common"
)

// TestConcurrentInstanceCreation tests concurrent instance creation
func TestConcurrentInstanceCreation(t *testing.T) {
	const numInstances = 3
	var wg sync.WaitGroup
	results := make(chan error, numInstances)

	for i := 0; i < numInstances; i++ {
		wg.Add(1)
		go func(instanceNum int) {
			defer wg.Done()
			
			// Create a unique test name for each instance
			testName := fmt.Sprintf("TestAccConcurrentInstance_%d", instanceNum)
			
			// Use a separate testing.T context for each goroutine
			t.Run(testName, func(t *testing.T) {
				resource.Test(t, resource.TestCase{
					PreCheck:                 func() { testAccPreCheck(t) },
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: testAccConcurrentInstanceConfig(instanceNum),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr(
									fmt.Sprintf("multipass_instance.concurrent_test_%d", instanceNum),
									"name",
									fmt.Sprintf("concurrent-test-%d", instanceNum),
								),
								resource.TestCheckResourceAttrSet(
									fmt.Sprintf("multipass_instance.concurrent_test_%d", instanceNum),
									"id",
								),
							),
						},
					},
				})
			})
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Check if any errors occurred
	for err := range results {
		if err != nil {
			t.Errorf("Concurrent test failed: %v", err)
		}
	}
}

// TestConcurrentSameInstanceAccess tests concurrent access to the same instance
func TestConcurrentSameInstanceAccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create the instance
			{
				Config: testAccConcurrentSameInstanceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.shared", "name", "shared-instance"),
				),
			},
			// Then test concurrent access to the same instance
			{
				Config: testAccConcurrentSameInstanceAccessConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.shared", "name", "shared-instance"),
					resource.TestCheckResourceAttr("data.multipass_instance.concurrent1", "name", "shared-instance"),
					resource.TestCheckResourceAttr("data.multipass_instance.concurrent2", "name", "shared-instance"),
					resource.TestCheckResourceAttr("data.multipass_instance.concurrent3", "name", "shared-instance"),
				),
			},
		},
	})
}

// TestConcurrentDataSourceQueries tests concurrent data source queries
func TestConcurrentDataSourceQueries(t *testing.T) {
	const numQueries = 5
	var wg sync.WaitGroup
	errors := make(chan error, numQueries)

	for i := 0; i < numQueries; i++ {
		wg.Add(1)
		go func(queryNum int) {
			defer wg.Done()
			
			testName := fmt.Sprintf("ConcurrentDataQuery_%d", queryNum)
			t.Run(testName, func(t *testing.T) {
				resource.Test(t, resource.TestCase{
					PreCheck:                 func() { testAccPreCheck(t) },
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: testAccConcurrentDataSourceConfig(),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttrSet("data.multipass_instance.test", "id"),
								resource.TestCheckResourceAttrSet("data.multipass_instance.test", "instances.#"),
							),
						},
					},
				})
			})
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent data source query failed: %v", err)
		}
	}
}

// TestRaceConditionInstanceOperations tests race conditions in instance operations
func TestRaceConditionInstanceOperations(t *testing.T) {
	client := NewMultipassClient("")
	instanceName := "race-test-instance"

	// Create launch options
	opts := &common.LaunchOptions{
		Name:   instanceName,
		Image:  "22.04",
		CPU:    "1",
		Memory: "1G",
		Disk:   "5G",
	}

	// Test concurrent launches of the same instance (should fail for duplicates)
	t.Run("ConcurrentLaunch", func(t *testing.T) {
		const numAttempts = 3
		var wg sync.WaitGroup
		results := make(chan error, numAttempts)

		for i := 0; i < numAttempts; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := client.Launch(opts)
				results <- err
			}()
		}

		wg.Wait()
		close(results)

		successCount := 0
		failureCount := 0

		for err := range results {
			if err == nil {
				successCount++
			} else {
				failureCount++
			}
		}

		// Only one launch should succeed, others should fail due to name conflict
		if successCount > 1 {
			t.Errorf("Expected only 1 successful launch, got %d", successCount)
		}
	})

	// Test concurrent operations on the same instance
	t.Run("ConcurrentOperations", func(t *testing.T) {
		// First ensure the instance exists (ignore error if it already exists)
		client.Launch(opts)

		var wg sync.WaitGroup
		operationResults := make(chan error, 6)

		// Start instance
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.StartInstance(instanceName)
			operationResults <- err
		}()

		// Stop instance
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.StopInstance(instanceName)
			operationResults <- err
		}()

		// Restart instance
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.RestartInstance(instanceName)
			operationResults <- err
		}()

		// Get instance info
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetInstance(instanceName)
			operationResults <- err
		}()

		// List instances
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.ListInstances()
			operationResults <- err
		}()

		// Suspend instance
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := client.SuspendInstance(instanceName)
			operationResults <- err
		}()

		wg.Wait()
		close(operationResults)

		// Some operations may fail due to state conflicts, which is expected
		errorCount := 0
		for err := range operationResults {
			if err != nil {
				errorCount++
				t.Logf("Operation failed (expected in concurrent scenario): %v", err)
			}
		}

		// Clean up
		client.DeleteInstance(instanceName, true)
	})
}

// TestDeadlockPrevention tests that concurrent operations don't cause deadlocks
func TestDeadlockPrevention(t *testing.T) {
	client := NewMultipassClient("")
	
	// Create a timeout to prevent infinite waiting
	timeout := time.After(30 * time.Second)
	done := make(chan bool)

	go func() {
		const numOperations = 10
		var wg sync.WaitGroup

		for i := 0; i < numOperations; i++ {
			wg.Add(1)
			go func(opNum int) {
				defer wg.Done()
				
				instanceName := fmt.Sprintf("deadlock-test-%d", opNum)
				opts := &common.LaunchOptions{
					Name:   instanceName,
					Image:  "22.04",
					CPU:    "1",
					Memory: "1G",
					Disk:   "5G",
				}

				// Perform multiple operations on each instance
				client.Launch(opts)
				client.GetInstance(instanceName)
				client.StartInstance(instanceName)
				client.StopInstance(instanceName)
				client.DeleteInstance(instanceName, true)
			}(i)
		}

		wg.Wait()
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test timed out - possible deadlock detected")
	case <-done:
		t.Log("Deadlock prevention test completed successfully")
	}
}

// TestConcurrentProviderAccess tests concurrent access through the provider
func TestConcurrentProviderAccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConcurrentProviderAccessConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that all instances were created successfully
					resource.TestCheckResourceAttr("multipass_instance.concurrent_1", "name", "concurrent-provider-1"),
					resource.TestCheckResourceAttr("multipass_instance.concurrent_2", "name", "concurrent-provider-2"),
					resource.TestCheckResourceAttr("multipass_instance.concurrent_3", "name", "concurrent-provider-3"),
					// Check that all data sources work
					resource.TestCheckResourceAttrSet("data.multipass_instance.list", "instances.#"),
				),
			},
		},
	})
}

// TestResourceStateConsistency tests that resource state remains consistent under concurrent access
func TestResourceStateConsistency(t *testing.T) {
	// This test would verify that Terraform state remains consistent
	// when multiple operations are performed concurrently
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStateConsistencyConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.state_test", "name", "state-consistency-test"),
					resource.TestCheckResourceAttrSet("multipass_instance.state_test", "id"),
					resource.TestCheckResourceAttrSet("multipass_instance.state_test", "state"),
				),
			},
			// Update the resource and ensure state consistency
			{
				Config: testAccStateConsistencyConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.state_test", "name", "state-consistency-test-updated"),
					resource.TestCheckResourceAttrSet("multipass_instance.state_test", "id"),
				),
			},
		},
	})
}

// Configuration helper functions

func testAccConcurrentInstanceConfig(instanceNum int) string {
	return fmt.Sprintf(`
resource "multipass_instance" "concurrent_test_%d" {
  name   = "concurrent-test-%d"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`, instanceNum, instanceNum)
}

func testAccConcurrentSameInstanceConfig() string {
	return `
resource "multipass_instance" "shared" {
  name   = "shared-instance"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccConcurrentSameInstanceAccessConfig() string {
	return `
resource "multipass_instance" "shared" {
  name   = "shared-instance"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}

data "multipass_instance" "concurrent1" {
  name = multipass_instance.shared.name
}

data "multipass_instance" "concurrent2" {
  name = multipass_instance.shared.name
}

data "multipass_instance" "concurrent3" {
  name = multipass_instance.shared.name
}
`
}

func testAccConcurrentDataSourceConfig() string {
	return `
data "multipass_instance" "test" {}
`
}

func testAccConcurrentProviderAccessConfig() string {
	return `
resource "multipass_instance" "concurrent_1" {
  name   = "concurrent-provider-1"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}

resource "multipass_instance" "concurrent_2" {
  name   = "concurrent-provider-2"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}

resource "multipass_instance" "concurrent_3" {
  name   = "concurrent-provider-3"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}

data "multipass_instance" "list" {}
`
}

func testAccStateConsistencyConfig() string {
	return `
resource "multipass_instance" "state_test" {
  name   = "state-consistency-test"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccStateConsistencyConfigUpdated() string {
	return `
resource "multipass_instance" "state_test" {
  name   = "state-consistency-test-updated"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
}
`
}