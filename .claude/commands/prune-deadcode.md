# Prune Dead Code

Systematically identify and remove dead code from a directory or module.

## Usage Check

First verify that a target directory was provided. If no target was provided in the command arguments, respond with:

> Error: No target directory provided.
>
> Usage: `/prune-deadcode [directory]`
>
> Example: `/prune-deadcode core/encoding`
>
> This command analyzes code to find and remove:
> - Symbols (functions, types, constants, variables) that are never used
> - Entire files or modules that are unused
> - Dead code chains (symbols only used by other dead symbols)

Only proceed if a valid target was provided.

## 1. Scope Assessment

Before searching, understand the scope:

1. Identify the language(s) in the target directory
2. Count source files (excluding tests and generated files)
3. For large scopes, use parallel exploration agents to search different subdirectories concurrently
4. Skip generated files and test files during symbol extraction

## 2. Dead Code Detection

For each symbol in the target, determine if it's used:

**Exported/public symbols:**
1. Search the entire repository for usage outside the symbol's own module
2. Exclude test files from "production usage" determination

**Private/non-exported symbols:**
1. Search only within the same file or module for usage
2. Simpler analysis since scope is contained

**Both:**
- Account for transitive dependencies: a symbol used only by dead code is also dead

### Classification (for exported symbols)

| Category | Criteria |
|----------|----------|
| **Actively Used** | Found in production code outside target module |
| **Test-Only** | Only found in test files outside target module |
| **Self-Test Only** | Only found in target module's own test files |
| **Dead** | Not used externally, and not transitively required by any used symbol |

For private symbols, classification is simpler: either used within their scope, or dead.

## 3. What to Target

Focus on (in priority order):

1. **Entire dead modules**: Directories where nothing is imported externally
2. **Entire dead files**: Files where all symbols are unused
3. **Standalone dead functions**: Top-level functions never called
4. **Dead types**: Structs/classes/interfaces where the type itself is never referenced

## 4. What NOT to Target

**Do NOT suggest removing individual methods from utilities that are actively used.**

If a utility (type/class) is in production use, its methods are presumed to have future value even if not currently
called. Only target methods when:

1. The entire type/utility is dead, OR
2. The method is clearly vestigial (deprecated, commented as unused)

Edge cases to handle carefully:

- **Mocks**: Dead if only used by dead test code
- **Interfaces**: Check if any implementation is used
- **Entry points**: Functions in main/CLI modules may be intentionally uncalled

## 5. Report Format

Present findings organized by impact:

> ## Dead Code Report: `<target>`
>
> ### High Impact
> - `<module_path>/` - Entire module unused
> - `<file_path>` - Entire file unused (N lines)
>
> ### Individual Symbols
>
> #### 1. `<file_path>:<line>` - `<SymbolName>` (function/type/const/var)
> **Evidence**: No production usage found outside module
> **Dependencies**: Removing this also removes `<other_symbol>`

## 6. Interactive Walkthrough

After presenting the report:

1. Start with high-impact items (entire modules/files) before individual symbols
2. Present one item at a time
3. Show the code snippet
4. Ask: "Delete this dead code? (yes/no/skip)"
5. If yes: Delete the code, then run verification
6. **Do not advance** until user explicitly responds (next/done/skip)
7. Continue until all items processed

## 7. Post-Deletion Verification

After each deletion:

1. Run the project's lint/build command to verify compilation
2. If it fails, revert and report the issue

After all deletions:

1. Run dependency cleanup if applicable
2. Summarize: "Removed N symbols, M lines of code"
