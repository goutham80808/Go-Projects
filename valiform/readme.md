# Valiform (Validation + format)


**Valiform** (Validate + Format) is a powerful and flexible command-line tool for validating structured data files against a user-defined set of rules. After validation, it intelligently separates valid and invalid data, writing the clean data to a new file and providing a detailed report of all errors.

It's the perfect helper for data cleaning, pre-processing pipelines, and ensuring data quality before it enters your systems.

## Features

-   **CLI-Driven:** Easy to use and integrate into automated scripts.
-   **Configurable Rules:** Define complex validation logic in a simple, human-readable `YAML` file. No need to recompile to change rules!
-   **Multi-Format Support:** Currently supports CSV for input and JSON for output. (Extensible for more formats).
-   **Rich Validation:** Supports a variety of rules:
    -   `required` fields
    -   `min`/`max` for numerical values
    -   `regex` for complex string patterns
    -   (Coming soon: `enum`, `minLength`/`maxLength`, and more)
-   **Detailed Error Reporting:** For every failed validation, Valiform reports the row number, the field name, and a clear message explaining the issue.

---

## Installation

### Prerequisites

-   Go (version 1.18 or higher) installed on your system.
-   Git

### Building from Source

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/goutham80808/valiform.git
    cd valiform
    ```

2.  **Build the binary:**
    ```bash
    go build -o valiform ./cmd/valiform
    ```
    This will create an executable file named `valiform` in the current directory.

---

## Quick Start

Let's validate a sample CSV file of student grades.

**1. Create a data file (`data.csv`):**

```csv
StudentID,FirstName,Grade,Status
S12345,John,85,Enrolled
S67890,Jane,92,Enrolled
SABCDE,Peter,-10,Enrolled
S54321,Mary,105,Graduated
F11111,InvalidID,75,Unknown
```

**2. Create a rules file (`rules.yaml`):**

This file defines what constitutes "valid" data.

```yaml
file_type: "csv"
has_header: true
fields:
  - name: "StudentID"
    type: "string"
    rules:
      required: true
      regex: "^S[0-9]{5}$" 
  - name: "FirstName"
    type: "string"
    rules:
      required: true
  - name: "Grade"
    type: "integer"
    rules:
      required: true
      min: 0
      max: 100
  - name: "Status"
    type: "string"
    rules:
      # This rule is not yet implemented, but shows future capability
      # enum: ["Enrolled", "Withdrawn", "Graduated"]
      required: false
```

**3. Run Valiform:**

Execute the tool, pointing it to your configuration and data files.

```bash
  go run ./cmd/valiform -c rules.yaml -i data.csv -o valid_grades.json
```

**4. Check the Results:**

*   **Terminal Output:** Valiform will print a summary of its findings and a detailed list of all errors to your terminal.

    ```text
    --- Valiform Starting ---
    ...
    --- Validation Complete ---
    Successfully wrote 2 valid records to valid_data.json

    Found 3 validation errors:
    - row 4, field 'StudentID': value 'SABCDE' does not match regex pattern '^S[0-9]{5}$'
    - row 4, field 'Grade': value -10 is less than min 0
    - row 5, field 'Grade': value 105 is greater than max 100
    - row 6, field 'StudentID': value 'F11111' does not match regex pattern '^S[0-9]{5}$'
    ```

*   **Output File (`valid_data.json`):** A new file will be created containing only the records that passed all validation checks, formatted as clean JSON.

    ```json
    [
      {
        "FirstName": "John",
        "Grade": "85",
        "Status": "Enrolled",
        "StudentID": "S12345"
      },
      {
        "FirstName": "Jane",
        "Grade": "92",
        "Status": "Enrolled",
        "StudentID": "S67890"
      }
    ]
    ```

---

## Usage

```
valiform -c <config_file> -i <input_file> [flags]
```

### Flags

| Flag           | Shorthand | Description                                           | Required | Default             |
| -------------- | --------- | ----------------------------------------------------- | -------- | ------------------- |
| `--config`     | `-c`      | Path to the validation rules YAML file.               | **Yes**  |                     |
| `--input`      | `-i`      | Path to the input data file.                          | **Yes**  |                     |
| `--output`     | `-o`      | Path for the formatted output file.                   | No       | `valid_output.json` |
| `--help`       | `-h`      | Show the help message.                                | No       |                     |

---

## Future Development

This project is a great starting point. Future improvements could include:

-   [ ] **Improved Type Conversion:** Store validated data in its proper type (e.g., numbers in JSON output as `85`, not `"85"`).
-   [ ] **More Validators:** Implement `enum`, `minLength`, `maxLength`, date formats, etc.
-   [ ] **More I/O Formats:** Add support for reading JSON and writing to CSV or formatted text.
-   [ ] **Data Aggregation:** Add an option to calculate basic statistics on valid data (e.g., counts, averages).
-   [ ] **Unit Tests:** Increase test coverage for robustness and reliability.
