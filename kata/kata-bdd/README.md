# Encoder BDD Demo (Go + Godog)

This project demonstrates Behavior-Driven Development (BDD) in Go using `godog`.
We implement a simple encoder/decoder that supports Base64, Hex, and URL encoding.

## What is BDD?

BDD is a collaborative development practice that defines software behavior through examples written in a business-readable language. These examples become automated acceptance tests that guide development.

Key points:
- Specifications are written as executable examples (Scenarios)
- Common structure: `Given` (context), `When` (action), `Then` (outcome)
- Shared understanding across developers, QA, and stakeholders

## What is Godog?

Godog is the Cucumber implementation for Go. It lets you write behavior in Gherkin (`.feature` files) and bind each step to Go functions. Godog executes the scenarios and reports results, serving as executable specifications and acceptance tests.

Two common ways to use Godog:
- As a CLI tool (runs features directly)
- Integrated with `go test` (runs inside your test suite)

### Using Godog in this project

- Features live in `bdd/*.feature`
- Step definitions and test runner are in `bdd/encoder_steps_test.go`
- We integrate via `go test` by constructing a `godog.TestSuite` and pointing it at the feature paths

Optional: install the Godog CLI for direct execution

```
go install github.com/cucumber/godog/cmd/godog@latest
# from project root, run features under bdd/
godog bdd
```

## Project Structure

```
.
├── bdd
│   ├── encoder.feature          # Feature and scenarios (Gherkin)
│   └── encoder_steps_test.go    # Step definitions + test runner
├── cmd
│   └── encoder
│       └── main.go              # CLI for manual usage
├── pkg
│   └── encoder
│       └── encoder.go           # Encode/Decode implementation
└── go.mod
```

## The Feature (Gherkin)

The behavior is described in `bdd/encoder.feature`:

```
Feature: Encoding and decoding text
  As a user of the encoder tool
  I want to encode and decode strings using base64, hex, and url schemes
  So that I can transform text reliably

  Scenario Outline: Encode text
    Given I have the text "<plain>"
    When I encode using "<type>"
    Then the result should be "<encoded>"

    Examples:
      | type   | plain        | encoded                         |
      | base64 | hello world  | aGVsbG8gd29ybGQ=                |
      | hex    | hello        | 68656c6c6f                      |
      | url    | hello world! | hello+world%21                  |

  Scenario Outline: Decode text
    Given I have the text "<encoded>"
    When I decode using "<type>"
    Then the result should be "<plain>"

    Examples:
      | type   | encoded                         | plain        |
      | base64 | aGVsbG8gd29ybGQ=                | hello world  |
      | hex    | 68656c6c6f                      | hello        |
      | url    | hello+world%21                  | hello world! |
```

## Running the BDD Tests

Prerequisites:
- Go 1.20+

Install dependencies and run tests:

```
go test ./...
```

You should see the Godog suite execute and pass.

If you want to run only the BDD package:

```
go test ./bdd -run TestGodog
```

## Manual CLI Usage

Build and run the encoder CLI:

```
go run ./cmd/encoder --mode encode --type base64 --text "hello world"
# outputs: aGVsbG8gd29ybGQ=

go run ./cmd/encoder --mode decode --type url --text "hello+world%21"
# outputs: hello world!
```

## How BDD Guides the Code

1. Write the feature file (executable specification)
2. Implement step definitions to bind Gherkin steps to Go code
3. Write the minimal implementation to make scenarios pass
4. Refactor with confidence while keeping scenarios green

This demo shows a tight feedback loop: Features explain desired behavior; steps assert expectations; code evolves until behavior is satisfied.

## References
- Godog: https://github.com/cucumber/godog
- Cucumber BDD: https://cucumber.io/docs/bdd/
