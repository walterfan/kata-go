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
