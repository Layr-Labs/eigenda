## Style Guide

This style guide contains coding style guidelines for the EigenDA project. This guide is not exhaustive, but rather
builds on top of the guidelines expressed in [Effective Go](https://go.dev/doc/effective_go). It is intended as a guide
for human engineers, and to provide AI agents with a checklist for code review.

### 1. Style Enforcement Guidelines

1. Style guidelines should be enforced for all new code and documentation.
2. The decision of whether to modify pre-existing code to adhere to the guidelines must be made on a case-by-case basis:
   - If a line is being modified, it's probably reasonable to fix any style issues that exist on that line.
   - If style issues exist in close proximity to changes being made, it *may* make sense to fix the issues.
   - Style fixes shouldn't be allowed to overshadow the main point of a PR.
   - If a large quantity of style fixes are necessary, it's best to split them into a separate PR. E.g. don't turn a
   5 line PR into a 50 line PR just for the sake of style fixes!
3. Recognize that everyone has unique preferences, and be respectful of alternate viewpoints:
   - Pursuing personal style *opinions* on code you are changing is perfectly acceptable: by touching the code, your
   preferences supersede the preferences of the previous engineer.
   - Changes may be made in surrounding code for the sake of readability, but there's a fine line between
   "improving readability", and "aggressively imposing personal preference".
   - If there is a disagreement between engineers about style, the team should come to consensus and enshrine the
   result as an entry in this style guide.

### 2. Error Handling

1. Return errors explicitly; don't panic except for unrecoverable errors, where returning an error is not plausible.
   - Exceptions may be made for test code, where returning an error adds more complexity than benefit.
2. Use error wrapping with `fmt.Errorf("context: %w", err)` for additional context.
   - Ensure that `%w` is used for error wrapping, *not* `%v`.
   - Note that this rule only applies to `fmt.Errorf` specifically! It does NOT apply to `fmt.Sprintf`.

### 3. Code Documentation

1. Document all exported functions, structs, constants, and interfaces in production code.
2. Document unexported functions/types that contain non-trivial logic.
   - A good rule of thumb: if you can't understand everything there is to know about a function/type by its *name*,
   you should write a doc.
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
   - Documented function example:
   ```go
   // This preceding comment describes the function in detail, and isn't simply a rephrasing of the function name
   //
   // It contains the sort of information listed in `3.4`.
   //
   // It describes what is returned.
   func FunctionName(
      // common parameters like context, testing, and logger don't require documentation,
      // unless they're being used in an unusual way
      ctx context.Context,
      // similarly, documentation *may* be omitted for parameters with blatantly obvious purpose
      enabled bool,
      // parameters without blatantly obvious purpose should contain helpful documentation which isn't just a
      // rephrasing of the parameter name
      param1 int,
      ) error {
         // ...
   }
   ```
5. TODO comments should be added to denote future work.
   - TODO comments should clearly describe the future work, with enough detail that an engineer lacking context
   can understand.
   - TODO comments that must be addressed *prior* to merging a PR should clearly be marked,
   e.g. `// TODO: MUST BE ADDRESSED PRIOR TO MERGE`
   - TODO comments that are intended to be merged into `master` should be attributed to the engineer adding the TODO,
   e.g. `// TODO(litt3): we should consider optimizing this algorithm`

### 4. Spelling and Grammar

Proper spelling and grammar are important, because they help keep code and documentation unambiguous, easy to read, 
and professional. They should be checked and carefully maintained.

1. Overly strict adherence to arbitrary grammar and spelling "rules" that don't impact readability is not beneficial.
   Some examples of "rules" that shouldn't be enforced or commented on are:
   - "Don't end a sentence with a preposition"
   - "Never split the infinitive"
   - "Don't use passive voice"
   - "Always spell out numbers"
   - "Don't begin a sentence with 'And', 'But', or 'Because'"
   - Perfect canonical comma usage
   - Use "okay" instead of "ok"
2. Some things are technically correct grammatically, yet hinder readability. Despite being "grammatically correct",
   the following things should not be tolerated:
   - Sentences with ambiguous interpretations
   - Run-on sentences
3. Spelling should be checked, with some caveats:
   - If there are multiple correct spellings for a word, no one "correct" spelling should be asserted over another
   - Neologisms are permitted
4. Colloquial language that is appropriate in a professional setting is acceptable: don't be the "fun police".

### 5. Naming

Good code has good names. Bad names yield bad code.

1. Using names that are too succinct hinders readability:
   - `i` -> `nodeIndex`
   - `req` -> `dispersalRequest`
   - `status` -> `operatorStatus`
   - An exception is made for golang receiver names, are permitted to be a *single character* by convention
2. Consistency is key. A single concept should have a single term, ideally across the entire codebase.
   - The exception here is with local scoping. E.g. if you have an `OperatorId` throughout the codebase, it would be
   reasonable to refer to it as an `id` inside the `Operator` struct.
3. Do not overload terms.
4. Avoid attributing special technical meaning to common generic terms.
   - E.g., you shouldn't try to usurp the word `Component` to mean a specific part of the system, since it's already
   used in many generic contexts.

### 6. Code Structure

1. Keep functions short and readable.
   - A good rule of thumb is to keep functions <50 lines, but this isn't a strict limit.
   - Just because a function is <50 lines doesn't mean it shouldn't be split!
   - Some good candidates for logic to split out of complex functions are:
      - The logic inside a `for` loop or `if` block
      - Input validation
      - Complex calculations
2. Keep nesting as shallow as possible. Ideally, you'd never have > 1 block deep of nesting. Practically, some amount of
   multi-level nesting is unavoidable, but efforts should be made to keep it to a minimum:
   - Split out helper functions
   - Consider using "early-out" logic, to decrease nesting by 1 level:

        Before:
        ```go
        if success {
            for _, item := range items {
                processItem(item)           // <-- nesting here is 2 blocks deep
            }
            return nil
        }
        return error
        ```

        After:
        ```go
        if !success {
            // early-out
            return error
        }
        for _, item := range items {
            processItem(item)               // <-- now it's only 1 block deep
        }
        return nil
        ```
3. Place the most important functions at the top of the file.
4. Public static functions that lack a tight coupling to a specific struct (e.g. a constructor) should be placed in
files with a `_utils` suffix.
5. Don't export things that don't need to be exported
   - Member variables should almost always be unexported
   - Structs, interfaces, and constants should only be exported if necessary

### 7. Defensive Coding

1. Prefer using constructors over raw struct instantiation.
   - Raw struct instantiation is bug-prone: fields can be removed by mistake, or newly added fields may not be
   universally added to all usages.
   - Constructors are a convenient place to validate new struct instantiations.
2. If it is even remotely possible that something could be `nil`, *check it*.
   - Even if it doesn't seem likely that something could be `nil`, it's easy to miss edge cases, and future changes can
   invalidate original assumptions.
   - At minimum, any situation where a `nil` check is skipped must be explicitly commented, stating the reason that
   it's safe.

### 8. TODO(litt3): Missing Guidelines

The following topics are good candidates for future additions to this style guide. Anyone with a strong opinion
should consider creating a PR to add a new section.

1. Package organization and naming
2. Interface/struct design and naming
3. Solidity style
