# Preprocess Logs

The purpose of this document is to provide an AI agent with a framework for doing preprocessing on large
quantities of logs. This framework is needed in order to carefully manage AI context. It allows the agent to
extract useful information without having to load the entire log contents into context. All output files will be
saved to an <analysis_directory>, which should be named "analysis" and placed inside the original log directory.

## Usage Check

First verify that a log directory path was provided. If no path was provided in the command arguments, respond with:

> Error: No log directory path provided.
>
> Usage Example: /preprocess-logs logs_41509396734
>
> This command analyzes log files in the specified directory by:
> 1. Splitting large files into manageable shards
> 2. Searching for various error patterns
> 3. Generating a human-readable report of failures
>
> The analysis and preprocessing artifacts will be saved inside the log directory.

Only proceed with the analysis if a valid directory path was provided.

## Phase 0: Check for Pre-existing Analysis

Before beginning the log preprocessing procedure, check if a previous analysis has already been completed for the target log files.

1. **Check for existing analysis directory**: Look for an `analysis` directory inside the original log directory
(i.e., `<original_log_directory>/analysis/`)

2. **Verify analysis completeness**: If the analysis directory exists, check for the presence of key analysis artifacts:
   - `shards/` directory containing shard files
   - `search_results/` directory containing search result files
   - `<original_log_directory_name>_preprocessing_report.md`

3. **User confirmation for re-analysis**: If a complete analysis is found, ask the user for confirmation before proceeding:

   > Found existing analysis for <original_log_directory>. The analysis includes:
   > - X shard files
   > - Search results
   > - Preprocessing report
   > 
   > Do you want to re-analyze these logs and overwrite the existing analysis?

## Phase 1: Split Large Logs

The first stage is large log files are split into smaller pieces called **shards** to allow context efficient
processing. Each shard contains a fixed number of lines (default 1800) based on the maximum input limits of the
intended analysis tool.

1. Store all shard files in a directory called "shards", inside the <analysis_directory>. The analysis directory
   should be named `analysis` and placed inside the original log directory. Each shard should be named
   `<original_log_name>_shard_<shard_decimal_index>`.

2. **Split Command:** Use the following command to split log files into shards with decimal numbering:

  ```bash
  split -l 1800 -d -a 3 "<original_log_file>" "<original_log_directory>/analysis/shards/<original_log_name>_shard_"
  ```

  **Command explanation:**
  - `-l 1800`: Split every 1800 lines
  - `-d`: Use numeric suffixes instead of alphabetic
  - `-a 3`: Use 3-digit suffixes for better readability and sorting

  Example shard files:
  ```
  log_dump_12/analysis/shards/system_log_shard_001.txt
  log_dump_12/analysis/shards/unit_tests_shard_012.txt
  ```

## Phase 2: Generate Failure Metadata

Find potential errors in log shards using `ripgrep` (`rg`) for pattern matching. Do not read shards into context
at this point: we are simply generating an index of lines that might potentially represent errors.

**Directory Setup:** Create a `search_results/` subdirectory within the analysis directory to organize ripgrep output:
```bash
mkdir -p "<original_log_directory>/analysis/search_results"
```

### Search Profiles

Use targeted search profiles based on the type of failures you're looking for. If the user didn't specify what you are
searching for, you should iteratively search using each profile.

#### Profile 1: Test Failures
For standard test output failures:
```bash
rg --line-number --ignore-case --json -C 5 -- "^[-]{3} FAIL:|\\s+FAIL\$|\\s+FAIL\\t|\\[FAILED\\]|panic: test timed out" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/test_failures_search.jsonl"
```

#### Profile 2: Connection/Network Errors
For network-related issues:
```bash
rg --line-number --ignore-case --json -C 5 "ECONNREFUSED|connection refused|dial.*failed|cannot connect|connection reset" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/connection_errors_search.jsonl"
```

#### Profile 3: Startup/Initialization Errors
For service startup problems:
```bash
rg --line-number --ignore-case --json -C 5 "error starting|failed to start|initialization failed|startup failed|cannot initialize" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/startup_errors_search.jsonl"
```

#### Profile 4: Docker/Container Issues
For container-related problems:
```bash
rg --line-number --ignore-case --json -C 5 "container.*failed|docker.*error|OCI runtime|container.*exit.*[1-9]" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/container_errors_search.jsonl"
```

#### Profile 5: Resource/Timeout Issues
For resource constraints and timeouts:
```bash
rg --line-number --ignore-case --json -C 5 "out of memory|OOM|deadline exceeded|context canceled|timeout waiting" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/resource_errors_search.jsonl"
```

#### Profile 6: Panic/Crash Detection
For application crashes:
```bash
rg --line-number --ignore-case --json -C 5 "panic:|fatal error:|segmentation fault|SIGSEGV|goroutine.*panic" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/panic_errors_search.jsonl"
```

#### Fallback: General Errors
Only use if specific searches yield no results:
```bash
rg --line-number --ignore-case --json -C 5 "ERROR|FAIL|CRITICAL" "<original_log_directory>/analysis/shards/" > "<original_log_directory>/analysis/search_results/general_errors_search.jsonl"
```

### Search Result Management

After running each search profile, split the results into manageable shards:

```bash
# Split search results into 1800-line shards (same size as log splitting)
split -l 1800 -d -a 3 "<original_log_directory>/analysis/search_results/test_failures_search.jsonl" \
  "<original_log_directory>/analysis/search_results/test_failures_shard_"
```

Repeat this for each search profile that generates results.

**Ripgrep JSON Output Structure:**
The ripgrep command outputs JSON lines where each entry has a `type` field:
- `"type":"match"` - Contains the actual match with file path, line number, and matched text
- `"type":"context"` - Contains surrounding context lines with their line numbers
- `"type":"begin"` and `"type":"end"` - File boundaries and summary statistics

## Phase 3: Generate Human Readable Log Preprocessing Report

This phase produces a structured summary for human consumption. Store the report as a **Markdown file** at
`<original_log_directory>/analysis/<original_log_directory_name>_preprocessing_report.md`.

**Formatting Requirements:**
- Target line length: 120 characters
- Lines that would suffer from being split (e.g., URLs, code snippets, file paths) may exceed this limit
- Apply best-effort line wrapping for readability while preserving technical accuracy

### Report Type: Test Output

If the logs represent output from one or more tests, then the report will focus on describing tests that included failures.

- Do not include a given test in the summary unless it failed
- If individual tests in the input logs are sorted into discrete test groups, i.e. CI actions, then this should be
  reflected in the format of the output file.

IMPORTANT: The ripgrep JSON output will help determine which tests failed, but matches in the
ripgrep output alone **do not** indicate a failed test. The search results serve as a *starting
point* for finding failed tests.

The basic format of the `Preprocessing Report` for logs representing tests is as follows:

> # Test Output Preprocessing Report
> 
> ## Search Results Summary
> - Log Type Detected: <test_output|container_logs|system_logs>
> - Total Matches Found:
>   - Test Failures: X matches
>   - Connection Errors: Y matches
>   - [other profiles...]
> 
> ## Test Failures
> 
> <list of failed tests> // see below for details of how test failures should be structured
> 
> ## Failure Clusters
> 
> <list of classes of failures> // see below for details of how failure classes should be structured

For each match entry (`"type":"match"`) in the ripgrep JSON output, perform the following steps:

1. Extract the match details and surrounding context from the JSON output
  - The match entry contains file path, line number, and matched text
  - Context entries provide surrounding lines with their line numbers
  - If the JSON context isn't sufficient, read the entire log shard as a fallback
2. For the entry, determine the following:
  a. if the entry belongs to an actually failed test, or if it's a false positive (e.g., a log in a passing test
     contained one of the search patterns). If you determine that the failure is a false positive, ignore the entry.
  b. **IMPORTANT: Avoid duplicating test suite summaries.** If the failure is a test suite summary that only reports
     the aggregate status of individual tests that have already been identified (e.g., "--- FAIL: TestSuiteName",
     "FAIL TestSuiteName", or summary lines like "2 Failed | 1 Passed"), ignore these entries. Only record 
     individual test failures that provide specific failure details and root causes.
  c. if the failure belongs to a test which failed, determine which specific test it belongs to
  d. if tests are organized into groups, i.e. CI actions, determine which group the test belongs to
  e. the class of failure. Think deeply about the log output, and try to briefly summarize what it conveys.
     e.g. "Root component invalid array access", or "runtime type panic in ServerProcess"
3. Record the test failure in the report:

> ### CI Action: Unit Tests                                     <-- this is the group the test belongs to.
>                                                               <-- if the test group has already been added to the report, add the test failure entry under the existing heading
>
> 1. `TestParallelProcessing`                                   <-- this is the name of the test
>   - failure location: `unit_tests_shard_003` line 62      <-- record where the error can be found in the shard files
>   - failure class: `consistency assertion failed in MainLoop` <-- determined failure class
>   - relevant log lines:                                       <-- try to show a brief selection of log lines that make it easy to understand what happened
>     ```
>     ...
>     ```

Note that a given test should not have multiple entries. If multiple match entries in the ripgrep JSON output correspond
to a single test, try to determine what the "actual" cause of the failure was. If unsure, include all potentially
relevant failures under the test failure entry in the report.

**Example of avoiding duplication:** If you see both:
- `[FAILED] TestSpecificFunction` with detailed error information
- `--- FAIL: TestSuiteName (123.45s)` that contains TestSpecificFunction

Only record the specific test failure (`TestSpecificFunction`), not the suite summary (`TestSuiteName`).

4. In addition to listing failed tests, it can be helpful to group similar failures together. These are called
   "failure clusters". After adding a failed test to the list of failed tests, you should add the test to the
   corresponding failure cluster. For example, if multiple tests are failing due to `invalid configuration: could
   not start system`, then you should add an `Invalid Configuration` failure cluster to the list, and add the test
   name as a sub-bullet

Example failure clusters:

> ## Failure Clusters
>
> 1. Nullptr Access
>   a. `CI Action: Unit Tests::TestNewImpl`
> 2. Invalid Configuration
>   a. `CI Action: Unit Tests::TestProcessing`
>   b. `CI Action: E2E Tests::TestEndToEndInMemory`

### Report Type: Arbitrary Log Output

If the logs represent an arbitrary selection of logs from a running system, then there aren't any "failed tests"
to detail. Instead, you should analyze the entries in the ripgrep JSON output, and generate the discovered set
of failure clusters. To do this, follow the same procedure defined above.

## Context Compaction

Since you will be dealing with large quantities of data, it is likely that you will need to compact context despite
best efforts to limit what's being loaded. 

### Strategies for managing large result sets:

1. **Use samples for initial analysis**: When search results are sharded, start with the `*_sample.jsonl` files
2. **Process shards sequentially**: Load and analyze one shard at a time, maintaining running totals/summaries
3. **Prioritize unique failures**: Focus on distinct error patterns rather than repetitive instances
4. **Discard processed content**: After extracting relevant information from a shard, clear it from context

Discard context related to literal log contents first: retain in context information related to what specific 
tests have failed, and what classes of failure are being observed.
