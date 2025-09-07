# Go Context Cancellation Demo

This program demonstrates how context cancellation works in Go, specifically showing the relationship between parent and child contexts.

## Key Concepts

### 1. Context Hierarchy
- Contexts form a tree structure where child contexts inherit from parent contexts
- When a parent context is cancelled, **all its child contexts are automatically cancelled**
- When a child context is cancelled, **the parent context is NOT affected**

### 2. Context Cancellation Types
- **Manual Cancellation**: Using `context.WithCancel()`
- **Timeout Cancellation**: Using `context.WithTimeout()`
- **Deadline Cancellation**: Using `context.WithDeadline()`

## Demos Included

### Demo 1: Parent Cancellation Affects Child
Shows how cancelling a parent context automatically cancels all child contexts derived from it.

### Demo 2: Child Cancellation Does NOT Affect Parent
Demonstrates that cancelling a child context leaves the parent context unaffected.

### Demo 3: Multi-level Context Hierarchy
Shows a three-level context hierarchy (grandparent → parent → child) and how cancellation cascades down.

### Demo 4: Timeout-based Cancellation
Demonstrates automatic cancellation using timeouts and how child contexts inherit parent timeouts.

### Demo 5: Deadline-based Cancellation
Shows how different deadlines work in parent and child contexts.

### Demo 6: Realistic Scenario
A practical example showing how to use contexts in concurrent operations like database calls, API requests, and file processing.

## How to Run

```bash
go run main.go
```

## Expected Output

The program will run through all 6 demos sequentially, showing:
- How parent cancellation cascades to children
- How child cancellation doesn't affect parents
- Timeout and deadline behavior
- Proper context handling patterns

## Best Practices Demonstrated

1. **Always call cancel functions**: Use `defer cancel()` to ensure cleanup
2. **Check context cancellation**: Use `select` with `<-ctx.Done()` to handle cancellation
3. **Proper error handling**: Check `ctx.Err()` to understand why cancellation occurred
4. **Timeout inheritance**: Child contexts automatically inherit parent timeouts
5. **Resource cleanup**: Use defer statements for proper resource management

## Context Error Types

- `context.Canceled`: Context was explicitly cancelled
- `context.DeadlineExceeded`: Context deadline was reached
