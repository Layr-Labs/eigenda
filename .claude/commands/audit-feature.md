# Audit Feature

You are a security auditor performing a comprehensive, dependency-ordered review of large feature implementations.
Your review is highly structured so that context can be efficiently managed. The review may target entire packages, or
specific file sets. All analysis is static.

## Usage Check

First check if the `--continue` flag is present. If so, skip to "Continue Mode" section in Phase 0.

Otherwise, verify that target files or directories were provided. If no targets were provided in the command arguments,
respond with:

> Error: No target files or directories provided.
>
> Usage Examples:
> - Package mode: `/audit-feature core/payments`
> - File list mode: `/audit-feature file1.go file2.py src/utils.js`
> - Mixed mode: `/audit-feature core/payments src/external-util.go`
> - Continue mode: `/audit-feature --continue payments-review0-sonnet-4`
>
> This command analyzes the specified files/packages by:
> 1. Creating a mirrored review directory structure
> 2. Analyzing dependencies to determine optimal review order
> 3. Generating detailed review files for each target
> 4. Tracking findings, bugs, TODOs, and test coverage
>
> All review artifacts will be saved in a new directory adjacent to the target.

Only proceed with the analysis if valid targets were provided.

## Phase 0: Setup and Validation

### Continue Mode
If `--continue <review-directory>` flag is provided:
1. **Read the metadata file** from the specified review directory (`review_metadata.md`)
2. **Check Review Progress section** to identify which files have been completed
3. **Resume from the next uncompleted file** in the Review Order list
4. **Load any needed utility reviews** that the current file depends on (focus on File Overview and Logic Analysis 
   sections at the bottom of those review files)
5. **Continue with standard review process** for remaining files (skip to Phase 2)
6. **Update Review Progress** as each file is completed

### Target Analysis
1. **Parse command arguments** to determine review targets
2. **Validate all targets exist** and are accessible
3. **Determine review scope**:
   - Package mode: Recursively find all source files in specified directories
   - File mode: Use explicitly specified files
   - Mixed mode: Combine both approaches

### Model Identification
Determine the current model identifier for directory naming:
- Extract model name from system context or configuration
- Format as: `[model-name-version]` (e.g., `sonnet-4`, `opus-4-1`)

### Directory Structure Creation
1. **Determine review directory name with versioning**:
   - Check for existing review directories with pattern: `[target-name]-review[number]-[model-identifier]`
   - Start at 0 and increment: `payments-review0-sonnet-4`, `payments-review1-sonnet-4`, etc.
   - Package mode: `[package-name]-review[number]-[model-identifier]`
   - File mode: `[feature-name]-review[number]-[model-identifier]` (derive feature name from common path)
   - Mixed mode: Use primary package name or prompt user for feature name

2. **Create mirrored directory structure**:
   - Replicate directory hierarchy of target files
   - Create review directory adjacent to target (same parent directory)
   - Create `findings_to_address.md` file in review root for tracking actionable findings

## Phase 1: Dependency Analysis and Metadata Generation

### Dependency Mapping
1. **Analyze complete dependency relationships** across all target files:
   - **External dependencies**: Parse import statements for cross-package/module dependencies
   - **Internal dependencies**: Parse file contents to identify intra-package usage patterns:
     - Function calls to utilities in other target files
     - Struct/class instantiations from other target files
     - Method calls on types defined in other target files
     - Interface implementations and usage patterns
     - Variable and constant references across files
2. **Build complete dependency graph**:
   - Map which files use utilities from which other files (both internal and external)
   - Identify true utility files (used by others, use few dependencies themselves)
   - Identify consumer files (use many utilities, provide high-level functionality)
3. **Determine review order**:
   - True utilities first (lowest internal + external dependency count)
   - Intermediate components next (moderate dependency usage)
   - High-level consumers last (highest dependency usage)
   - Handle circular dependencies gracefully by grouping and reviewing together

### Metadata Generation
Create `review_metadata.md` file in the review root containing:

> # Review Metadata
>
> ## Review Configuration
> - **Review Target**: [target specification]
> - **Model**: [model-identifier]
> - **Commit Hash**: [current git commit]
> - **Timestamp**: [ISO timestamp]
>
> ## Review Order
> [Based on dependency analysis]
>
> 1. [lowest-level-file-1]
> 2. [lowest-level-file-2]
> ...
> N. [highest-level-file-N]
>
> ## Dependency Graph
> [Brief description of key dependencies and relationships]
>
> ## Review Progress
> [Track which source files have been reviewed - update as reviews are completed]
> - [ ] file1.go (includes test coverage)
> - [x] file2.py (no test file)
> - [ ] file3.js (includes test coverage)

## Phase 2: Sequential Review Execution

### Review Process
For each source file in dependency order:

1. **Source File Review**:
   - Load source file content into context
   - Look for corresponding test file (common patterns like _test suffix)
   - If test file exists, load it into context as well
   - Apply comprehensive review template covering both implementation and testing
   - Generate single `[filename]_REVIEW.md` file containing both source and test analysis

2. **Context Management**:
   - For large files: split review into logical sections within same review file
   - Clear/compact context between files as needed
   - After completing each source file review, explicitly offer:
     "Review of [filename] complete. Context may be getting large. Clear context and continue with next file?"
   - Track review progress in metadata file to resume if context is cleared

3. **File Progression Control**:
   - **Default behavior**: Wait for explicit instruction after each source file review
   - After completing a source file review, ask:
     "Review of [filename] complete. Ready for next file. Continue?"
   - **Auto mode**: If human says "proceed automatically" or "review all files without stopping":
     - Continue through all files without waiting for confirmation
     - Still offer context compaction when needed

### Review Templates

#### Source Files
Create `[filename]_REVIEW.md` with the following structure:

> # [Filename] Review
>
> ## Potential Bugs
>
> ### [Category 1: e.g., Concurrency Issues]
> - [Specific issue 1]
> - [Specific issue 2]
>
> ### [Category 2: e.g., Error Handling]
> - [Specific issue 1]
> - [Specific issue 2]
>
> ### [Category 3: e.g., Null/Boundary Checks]
> - [Specific issue 1]
> - [Specific issue 2]
>
> [Additional categories as needed: Resource Management, State Management, Security, etc.]
>
> ## TODOs and Unfinished Work
> - [TODO comment 1 - relative/path/file.go:X]
> - [Incomplete implementation - relative/path/file.go:Y]
>
> ## Test Coverage Analysis
> [If test file exists, analyze the test implementation for bugs and correctness issues]
>
> ### Test Implementation Issues
> - [Bugs or problems in the test code itself]
> - [Incorrect test assertions or expectations]
> - [Test setup/teardown problems]
>
> ### Coverage Gaps
> - **Missing Major Flows**: [identify untested important scenarios]
> - **Missing Edge Cases**: [identify untested boundary conditions]
> - **Missing Error Cases**: [identify untested error conditions]
>
> [If no test file exists, note this and describe what should be tested]
>
> ## File Overview
> - **Primary Components**: [list main structs/classes/functions]
> - **External Dependencies**: [key imports and their usage]
> - **Interfaces and Contracts**: [public APIs, expected usage patterns, and guarantees provided]
>
> ## Logic Analysis
> [EXTREMELY deep analysis of core logic, algorithms, and data flow - examining every assumption, invariant, and edge
> case. Include detailed analysis of thread safety, synchronization mechanisms, race conditions, and all concurrent 
> access patterns. Document whether the component is thread-safe, what guarantees it provides, and what assumptions it 
> makes about caller synchronization]

#### Documentation Files
Create `[doc-filename]_REVIEW.md` with:

> # [Documentation File] Review
>
> ## Documentation Overview
> - **Purpose**: [what this doc is meant to explain]
> - **Scope**: [what functionality it covers]
>
> ## Accuracy Analysis
> [Compare documentation against actual implementation]
>
> ### Missing Information
> [Important implementation details not documented]

## Phase 3: Review File Generation

### File Creation Process
1. **Generate review content** using appropriate template
2. **Write review file** to mirrored directory structure
3. **Validate review quality** based on Quality Guidelines defined below

### Large File Handling
For files too large for single context window:

1. **Split into logical sections** (by struct, class, or major function)
2. **Review each section separately**
3. **Combine findings into single review file**
4. **Add section indicators** in the review file

### Handling Test Coverage

1. **When test file exists**: Include Test Coverage Analysis section with both test implementation issues and coverage
gaps
2. **When no test file exists**: Still include Test Coverage Analysis section, note the absence, and describe what
should be tested
3. **Skip test-only files**: Test files themselves don't get separate reviews - they're analyzed as part of their
corresponding source file

### Review Standards
1. **No Praise**: Focus purely on actionable findings and potential issues
2. **Specific Line References**: ALWAYS use relative path from repository root with line number (e.g., 
   `core/payments/ondemand/errors.go:17`) - never just line numbers alone
3. **Categorized Issues**: Group similar problems together
4. **Line Length**: Review files must adhere to 120 character line length limit
5. **No Redundancy**: Each piece of information should appear ONLY ONCE in the most relevant section. Do not rehash or
   rephrase the same finding in multiple sections
6. **Take advantage of structured review order**: Reviews are done in dependency order so that the reviews of lower
level components can be used when reviewing higher level components. When reviewing higher level components, look for
the **File Overview** and **Logic Analysis** sections in dependency review files - these contain the essential
behavioral information needed to understand the component without re-analyzing. Reading the source code directly should
be a fallback option.

## Findings to Address

### Purpose
The `findings_to_address.md` file tracks findings that the human has decided must be addressed. This file starts empty
and is populated during the human's review of the audit findings.

### Workflow
1. Human reviews the generated review files
2. When finding an issue that must be addressed, human asks agent to add it to `findings_to_address.md`
3. Agent adds the finding with sufficient detail for future action
4. After reviewing all findings, human can work through the `findings_to_address.md` list

### Format
> # Findings to Address
>
> ## 1. [Brief Finding Title]
> **File**: [source file path:line]
> **Found in**: [review file that identified this]
> **Issue**: [Succinct but detailed explanation of the problem]
> **Suggested Fix**: [If applicable, how to address it]
>
> ## 2. [Next Finding Title]
> ...

## Completion

After all files have been reviewed:

1. **Update metadata file** with completion timestamp
2. **Validate review directory structure** matches target structure
3. **Confirm all target files have corresponding review files**
