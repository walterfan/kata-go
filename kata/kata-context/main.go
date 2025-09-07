package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Example 1: Parent context cancellation affects child contexts
func demoParentCancellationAffectsChild() {
	fmt.Println("=== Demo 1: Parent Context Cancellation Affects Child ===")

	// Create parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Create child context from parent
	childCtx, childCancel := context.WithCancel(parentCtx)

	var wg sync.WaitGroup

	// Start child goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer childCancel() // Clean up child context

		fmt.Println("Child: Starting work...")
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Child: Work completed normally")
		case <-childCtx.Done():
			fmt.Printf("Child: Cancelled due to: %v\n", childCtx.Err())
		}
	}()

	// Start parent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer parentCancel() // Clean up parent context

		fmt.Println("Parent: Starting work...")
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Parent: Work completed, cancelling parent context")
			parentCancel()
		case <-parentCtx.Done():
			fmt.Printf("Parent: Cancelled due to: %v\n", parentCtx.Err())
		}
	}()

	wg.Wait()
	fmt.Println("Demo 1 completed")
}

// Example 2: Child context cancellation does NOT affect parent context
func demoChildCancellationDoesNotAffectParent() {
	fmt.Println("=== Demo 2: Child Context Cancellation Does NOT Affect Parent ===")

	// Create parent context with cancellation
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	// Create child context from parent
	childCtx, childCancel := context.WithCancel(parentCtx)

	var wg sync.WaitGroup

	// Start parent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("Parent: Starting work...")
		select {
		case <-time.After(4 * time.Second):
			fmt.Println("Parent: Work completed normally (child cancellation didn't affect parent)")
		case <-parentCtx.Done():
			fmt.Printf("Parent: Cancelled due to: %v\n", parentCtx.Err())
		}
	}()

	// Start child goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer childCancel()

		fmt.Println("Child: Starting work...")
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Child: Work completed, cancelling child context")
			childCancel()
		case <-childCtx.Done():
			fmt.Printf("Child: Cancelled due to: %v\n", childCtx.Err())
		}
	}()

	wg.Wait()
	fmt.Println("Demo 2 completed")
}

// Example 3: Multiple levels of context hierarchy
func demoMultiLevelContextHierarchy() {
	fmt.Println("=== Demo 3: Multi-level Context Hierarchy ===")

	// Create context hierarchy: grandparent -> parent -> child
	grandparentCtx, grandparentCancel := context.WithCancel(context.Background())
	parentCtx, parentCancel := context.WithCancel(grandparentCtx)
	childCtx, childCancel := context.WithCancel(parentCtx)

	var wg sync.WaitGroup

	// Grandparent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer grandparentCancel()

		fmt.Println("Grandparent: Starting work...")
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("Grandparent: Cancelling - this will cascade to all descendants")
			grandparentCancel()
		case <-grandparentCtx.Done():
			fmt.Printf("Grandparent: Cancelled due to: %v\n", grandparentCtx.Err())
		}
	}()

	// Parent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer parentCancel()

		fmt.Println("Parent: Starting work...")
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Parent: Work completed normally")
		case <-parentCtx.Done():
			fmt.Printf("Parent: Cancelled due to: %v\n", parentCtx.Err())
		}
	}()

	// Child goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer childCancel()

		fmt.Println("Child: Starting work...")
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Child: Work completed normally")
		case <-childCtx.Done():
			fmt.Printf("Child: Cancelled due to: %v\n", childCtx.Err())
		}
	}()

	wg.Wait()
	fmt.Println("Demo 3 completed")
}

// Example 4: Timeout-based cancellation
func demoTimeoutCancellation() {
	fmt.Println("=== Demo 4: Timeout-based Cancellation ===")

	// Create parent context with timeout
	parentCtx, parentCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer parentCancel()

	// Create child context from parent (inherits timeout)
	childCtx, childCancel := context.WithCancel(parentCtx)
	defer childCancel()

	var wg sync.WaitGroup

	// Parent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("Parent: Starting work with 3-second timeout...")
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Parent: Work completed normally")
		case <-parentCtx.Done():
			fmt.Printf("Parent: Cancelled due to: %v\n", parentCtx.Err())
		}
	}()

	// Child goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("Child: Starting work (inherits parent's timeout)...")
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Child: Work completed normally")
		case <-childCtx.Done():
			fmt.Printf("Child: Cancelled due to: %v\n", childCtx.Err())
		}
	}()

	wg.Wait()
	fmt.Println("Demo 4 completed")
}

// Example 5: Deadline-based cancellation with different deadlines
func demoDeadlineCancellation() {
	fmt.Println("=== Demo 5: Deadline-based Cancellation ===")

	// Create parent context with deadline
	parentDeadline := time.Now().Add(4 * time.Second)
	parentCtx, parentCancel := context.WithDeadline(context.Background(), parentDeadline)
	defer parentCancel()

	// Create child context with earlier deadline
	childDeadline := time.Now().Add(2 * time.Second)
	childCtx, childCancel := context.WithDeadline(parentCtx, childDeadline)
	defer childCancel()

	var wg sync.WaitGroup

	// Parent goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("Parent: Starting work with deadline at %v\n", parentDeadline.Format("15:04:05"))
		select {
		case <-time.After(6 * time.Second):
			fmt.Println("Parent: Work completed normally")
		case <-parentCtx.Done():
			fmt.Printf("Parent: Cancelled due to: %v\n", parentCtx.Err())
		}
	}()

	// Child goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Printf("Child: Starting work with earlier deadline at %v\n", childDeadline.Format("15:04:05"))
		select {
		case <-time.After(6 * time.Second):
			fmt.Println("Child: Work completed normally")
		case <-childCtx.Done():
			fmt.Printf("Child: Cancelled due to: %v\n", childCtx.Err())
		}
	}()

	wg.Wait()
	fmt.Println("Demo 5 completed")
}

// Helper function to simulate work with context cancellation check
func simulateWork(ctx context.Context, name string, duration time.Duration) error {
	fmt.Printf("%s: Starting work for %v\n", name, duration)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	deadline := time.Now().Add(duration)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: Work cancelled: %v\n", name, ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				fmt.Printf("%s: Work completed successfully\n", name)
				return nil
			}
			fmt.Printf("%s: Still working...\n", name)
		}
	}
}

// Example 6: Proper context cancellation handling in realistic scenario
func demoRealisticScenario() {
	fmt.Println("=== Demo 6: Realistic Scenario with Proper Context Handling ===")

	// Main context for the entire operation
	mainCtx, mainCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mainCancel()

	var wg sync.WaitGroup

	// Database operation
	wg.Add(1)
	go func() {
		defer wg.Done()

		dbCtx, dbCancel := context.WithTimeout(mainCtx, 3*time.Second)
		defer dbCancel()

		if err := simulateWork(dbCtx, "Database", 4*time.Second); err != nil {
			fmt.Println("Database operation failed, but main operation continues")
		}
	}()

	// API call operation
	wg.Add(1)
	go func() {
		defer wg.Done()

		apiCtx, apiCancel := context.WithTimeout(mainCtx, 5*time.Second)
		defer apiCancel()

		if err := simulateWork(apiCtx, "API Call", 2*time.Second); err != nil {
			fmt.Println("API call failed")
		}
	}()

	// File processing operation
	wg.Add(1)
	go func() {
		defer wg.Done()

		fileCtx, fileCancel := context.WithCancel(mainCtx)
		defer fileCancel()

		if err := simulateWork(fileCtx, "File Processing", 6*time.Second); err != nil {
			fmt.Println("File processing was cancelled")
		}
	}()

	wg.Wait()
	fmt.Println("Demo 6 completed")
}

func main() {
	fmt.Println("Go Context Cancellation Demo")
	fmt.Println("=============================")

	// Run all demos
	demoParentCancellationAffectsChild()
	time.Sleep(1 * time.Second) // Brief pause between demos

	demoChildCancellationDoesNotAffectParent()
	time.Sleep(1 * time.Second)

	demoMultiLevelContextHierarchy()
	time.Sleep(1 * time.Second)

	demoTimeoutCancellation()
	time.Sleep(1 * time.Second)

	demoDeadlineCancellation()
	time.Sleep(1 * time.Second)

	demoRealisticScenario()

	fmt.Println("All demos completed!")
}
