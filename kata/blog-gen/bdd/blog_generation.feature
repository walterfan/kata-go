Feature: Blog Generation
  As a technical blogger
  I want to generate a daily blog with weather information
  So that I can publish consistent, high-quality content quickly

  Scenario: Successful blog generation with mock weather
    Given a valid environment with API key "sk-test-key"
    And a mock LLM API server is running
    And weather API is not configured
    When I run "go run main.go --idea \"WebRTC recorder with Pion\" --location \"Hefei\""
    Then a blog markdown file for today should be created
    And the blog content should include the idea "WebRTC recorder with Pion"
    And the blog content should include "Weather:"
    And the blog title should be "my blog at <today>"

  Scenario: Blog generation fails without API key
    Given an environment without API key
    And a mock LLM API server is running
    When I run "go run main.go --idea \"Missing API key test\" --location \"Hefei\""
    Then the process should exit with a non-zero status
    And the error output should contain "OpenAI API key is not set"
