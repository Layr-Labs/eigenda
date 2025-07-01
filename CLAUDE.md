# CLAUDE.md - EigenDA

> **Purpose** – This file is the onboarding manual for every AI assistant (Claude, Cursor, GPT, etc.) and every
> human who edits this repository. It encodes our coding standards, guard rails, and workflow tricks.

---

## 1. Non-negotiable golden rules

| #:  | AI *must NOT* do                                                                                                                                   |
|-----|----------------------------------------------------------------------------------------------------------------------------------------------------|
| G-0 | Write changes or use tools when you are not sure about something project-specific, or if you don't have context for a particular feature/decision. |
| G-1 | Touch test files, specs, or generated files without EXPLICIT permission. **Humans own tests and specs.**                                           |
| G-2 | Refactor large modules without human guidance. For changes >50 LOC or >3 files, **ask for confirmation**.                                          |

---

## 2. File imports

NOTE: Be aware that whatever you add to this list is automatically loaded into context (due to `@` annotation). It's helpful
   to provide project context, but only within reasonable limits. See the [Anthropic docs](https://docs.anthropic.com/en/docs/claude-code/memory#claude-md-imports)
   for more information.

| File Import                      | Description                          |
|----------------------------------|--------------------------------------|
| @Makefile                        | Build, test, and formatting commands |
| @mise.toml                       | Dependency management                |
| @.golangci.yml                   | Golang linting YAML                  |
| @README.md                       | Main project README file             |
| @docs/spec/src/**/*.md           | All project spec markdown files      |
| @.claude/commands/*              | Project claude commands              |

If there are imports that are relevant only to a particular part of the project, then they should be added to a CLAUDE.md
   file *in the relevant subdirectory*. Then the imports will only be processed when files within that directory are read.

---

## 3. Coding standards

1. **Line Length**: 120 characters. This should also be observed for Markdown files, where manual newlines should be inserted.
   - The line length limit may be ignored in cases where a line break hinders readability, e.g. URLs.
2. **Imports**: Use `goimports` for sorted imports. Group standard library, third-party, and local imports.
3. **Testing**: `*_test.go` files in same package; use `testify` for assertions.

### 3.1. Error handling

1. Return errors explicitly; don't panic except for unrecoverable errors
   - Some exceptions can be made for test code, where returning an error adds more complexity than benefit.
2. Use error wrapping with `fmt.Errorf("context: %w", err)` for additional context
   - Ensure that `%w` is used for error wrapping, *not* `%v`

### 3.2. Code Documentation

1. Write docs for all exported functions/types in production code
2. Write docs for unexported functions/types if they contain non-trivial logic. A good rule of thumb: if you can't understand everything
   there is to know about a function/type by its *name*, you should write a doc.
3. Function/type docs should NOT simply be a rephrasing of the function/type name.
   - E.g. the doc for `computeData` should NOT be "Computes the data".
4. Function docs should consider the following helpful information, if relevant:
   - What are the inputs?
   - Are there any restrictions on what the input values are permitted to be?
   - What is returned in the standard case?
   - What is returned in the error case(s)?
   - What side effects does calling the function have?
   - Are there any performance implications that users should be aware of?
   - Are there any performance optimizations that should/could be undertaken in the future?

### 3.3. Doc Files

1. **Hierarchical organization**: Hierarchical numbering for sections makes referencing easier.
2. **Tabular format for key facts**: Tables are helpful for understanding data at a glance: use them where appropriate.
3. **Use Links**: Links are very helpful to assist a human navigating through the codebase.
   - IMPORTANT: double check that links aren't broken after making changes to doc files. Similarly, if documentation
   contains links directly to code, make sure that code changes are paired with the corresponding doc updates.

---

## 4. Directory-Specific CLAUDE.md Files

1. **Always check for `CLAUDE.md` files in specific directories** before working on code within them. These files contain targeted context.
2. If a directory's `CLAUDE.md` is outdated or incorrect, **update it**.
3. If you make significant changes to a directory's structure, patterns, or critical implementation details, **document these in its `CLAUDE.md`**.
4. If a directory lacks a `CLAUDE.md` but contains complex logic or patterns worth documenting for AI/humans, **suggest creating one**.
5. Use `@` annotation within CLAUDE.md files to automatically load in helpful context, e.g. `@docs/submoduleDocs`. These imports will be automatically
   processed whenever the `CLAUDE.md` file is read.
6. If there is domain-specific terminology relevant to a directory, consider adding a small glossary of terms.

---

## 5. Common pitfalls

1. Forgetting to run `go mod tidy` after adding new dependencies.
2. Not linting before committing code.
3. Wrong working directory when running commands (ensure you're in project root).
4. Large AI refactors in a single commit.
5. Delegating test/spec writing entirely to AI (can lead to false confidence).

---

## 6. Files to NOT modify

These files and directories should generally not be modified without explicit permission:

1. **Generated files**: Any files that are automatically generated during build processes.
   - Smart contract bindings are an important example of autogenerated files that shouldn't be directly modified under any circumstances
2. **Cryptographic resources**: Files in `resources/` (SRS tables, G1/G2 points) are cryptographic parameters.
3. **Dependencies**: `go.mod` and `go.sum` files should only be modified through `go mod` commands.
4. **Documentation**: Security audits and formal specifications should not be modified.
5. **CI/CD configurations**: GitHub workflows and Docker configurations require careful review.
6. **Files that control IDE behavior**:
   - `.gitignore`: Controls version control file exclusions
   - IDE configuration files (if present): `.vscode/`, `.idea/`, etc.

---

## 7. AI Assistant Workflow: Step-by-Step Methodology

When responding to user instructions, the AI assistant (Claude, Cursor, GPT, etc.) should follow this process
   to ensure clarity, correctness, and maintainability:

1. **Consult Relevant Guidance**: When the user gives an instruction, consult the relevant instructions from
   `CLAUDE.md` files (both root and directory-specific) for the request.
2. **Clarify Ambiguities**: Based on what you could gather, see if there's any need for clarifications. If so,
   ask the user targeted questions before proceeding.
3. **Break Down & Plan**: Break down the task at hand and chalk out a rough plan for carrying it out,
   referencing project conventions and best practices.
4. **Trivial Tasks**: If the plan/request is trivial, go ahead and get started immediately.
5. **Non-Trivial Tasks**: Otherwise, present the plan to the user for review and iterate based on their
   feedback.
6. **Track Progress**: Use a to-do list (internally, or optionally in a `TODOS.md` file) to keep track of your
   progress on multi-step or complex tasks.
7. **If Stuck, Re-plan**: If you get stuck or blocked, return to step 3 to re-evaluate and adjust your
   plan.
8. **Check for related updates**: Once the user's request is fulfilled, look for any complementary changes that need to be made:
   - Code documentation / doc files that reference details that have been modified
   - Variable names that need to be updated
   - Error messages that use old terminology
   - Related functions / structures that should be renamed to match new changes
   - Links contained in documentation that were broken by the changes
9. **Lint**: Make sure changes pass linting, and that they adhere to style and coding standards
10. **Test**: Run tests related to the changes that have been made. Short tests should always be run, but ask permission
   before trying to run long tests.
11. **User Review**: After completing the task, ask the user to review what you've done, and repeat the
   process as needed.
12. **Session Boundaries**: If the user's request isn't directly related to the current context and can be
    safely started in a fresh session, suggest starting from scratch to avoid context confusion.

## 8. AI Assistant User Interactions

1. When responding to a prompt with a list of items, number the list for easy reference.
2. Use line numbers and file paths so that the user can easily find elements being referred to.
3. When asked to review something, don't focus on praising what's good about it. Instead, focus on concrete feedback for
   improvement. If nothing can be improved, it's ok to just say so.
