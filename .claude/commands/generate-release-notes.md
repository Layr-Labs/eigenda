# Release Notes

Your job is to help the user compile release notes for the EigenDA repository. You will assist the user in gathering
and sorting information about new features, bug fixes, improvements, etc., and based on the feedback from the user
you will generate a well-structured release notes document.


# Information you will need to gather

You will need to gather some information from the user to create comprehensive release notes.

1. Optionally, the user may provide a draft of the release notes that you can help polish. The draft might be release
   notes that you have helped work on previously, or they might be notes that they have written themselves. If the user
   doesn't specify, always ask them if they have a draft to use as a starting point. 
      a. If the user is providing a draft, they will often pass a file path when they invoke this command. 
         If you get a file path in this way, it's probably a draft that you should use as a starting point.
      b. The first thing you should do when the user provides a draft is to read it and see if you have a
         "#DRAFT - DO NOT PUBLISH" section at the bottom. This is where you will keep notes to yourself as to
         what steps you have completed, and what steps you still need to complete. If the draft doesn't have this 
         section, you should add it yourself, and assume that no steps have been completed yet.
2. The tag/branch for the prior release (e.g., v1.0.0).
    a. The exact commit for this release. If it's a branch, use the latest 
      commit in that branch. Always use the upstream commit.
    b. Never guess at what this is. Always ask the user. This is important, and it should always be your first
       question. If the user gives you a draft and the draft says what the prior release is, you can use that
       instead of asking the user.
3. The tag/branch for the current release being documented (e.g., v1.1.0).
    a. The exact commit for this release. If it's a branch, use the latest 
       commit in that branch. Always use the upstream commit.
    b. Never guess at what this is. Always ask the user. This is important, and it should always be your next 
       question after you determine the prior release information.  If the user gives you a draft and the draft says 
       what the prior release is, you can use that instead of asking the user.
4. The list of commits between the prior release and the current release.
    a. the general category for each commit. The categories are:
       - Validators
       - Disperser
       - Data API
       - Contracts
       - Integrations
       - Other (for miscellaneous commits that don't fit into the above categories)
    b. The importance of each commit. We use the conventional commits format. The importance levels are:
       - Major: for significant features, changes, or fixes that have a substantial impact.
       - Minor: for smaller improvements, bug fixes, or changes that have a lesser impact.
5. Whether or not this is an optional or a mandatory release for validators. The user will need to be the source
   of this information.
     a. If it's a mandatory release, the reason why it's mandatory.

# How to gather the information

The best way to gather information is to get it from git/github, if it is reasonable to do so. For example, if the user
provides a branch name, you can look it up to get the latest commit. If you have access to the github gui, use it.

Some information must come from the user. Sometimes the user will volunteer this information. Other times, you will
need to prompt them for it. 

When you ask the user for information, only ask for one thing at a time. As a rule of thumb, if it will take the user
multiple sentences to answer, consider breaking it up into multiple questions.

## Sorting and understanding commits

Commit messages can be terse, and you may be lacking context on some of the changes, or on the subject matter in 
general. That's ok, the user should be able to provide context.

For each commit you are unsure about, ask the user for clarification. Be sure to present the user with all information
you have available to you. It's very important as well to give the user a link they can click on to see the commit or 
PR in question.

We use squash merging. So for each commit in the release, there may actually be multiple "inner commits" that got
squashed together. You can go ahead and ignore these inner commits, and only deal with the top level commit.
Each of these top level commits should have a PR in github, which you may optionally look at if you are trying
to gather more information about that commit.

## Verifying information

Once you have sorted commits (i.e. into appropriate categories), it's important to verify the information with the user.

When you initially create the list of commits, include a special "[UNVERIFIED]" tag at the end of each commit line.
As you verify each commit with the user, you will remove the "[UNVERIFIED]" tag.

For each category, do the following:

Tell the user that you'd like to verify the contents of the category.

- Clearly state the category you are working on. (Do not mix categories in the same list.)
- Clearly state whether we are working with major or minor commits. (Do not mix major and minor commits in the same list.)
- Present the user with a list of 8 or fewer commits at a time (i.e. walk through each section in a paginated manner). 
- It's ok if there are fewer than 8 commits that are presented at a time (i.e. if there are only 3 commits in a 
  category, just present those 3). Never mix categories or major/minor importance levels in the same list given 
  to the user.
- Each commit should be in an enumerated list. 
- Tell the user that they should type a list of numbers for commits that are out of place, or if they want to change 
  the importance level. For each commit listed by the user, ask them what category or importance level it should 
  be instead (one at a time). If the user just directly tells you what changes to make, that's ok too.
- Based on the feedback from the user, update the document. If you are confident of the changes, 
  remove the "[UNVERIFIED]" tag.
- If a commit lacks the "[UNVERIFIED]" tag, you can assume it has already been verified by the user, and you don't 
  need to ask about it again.

When you present a list of commits to be verified by the user, use a format something like this:

```
Verifying Validators - Major Commits
  1. feat: LittDB Snapshots in https://github.com/Layr-Labs/eigenda/pull/1657
  2. feat!: validator state cache in https://github.com/Layr-Labs/eigenda/pull/1903

❓ Do any of these need to be moved to a different category or have their importance level changed?
```

THIS IS EXCEPTIONALLY IMPORTANT. VERIFY EACH COMMIT.

At the end, double check that there are no remaining commits with the "[UNVERIFIED]" tag. If there are, 
you need to circle back to the user and verify them.

# Release Notes Template

Below is a rough template for the release notes. Release notes are always markdown files. Sometimes a section
might be empty, and that's okay. If that happens, omit that section from the final output.

Note that sometimes there may be some major features that deserve their own section.

```markdown

# ${CURRENT_RELEASE} - Release Notes

- Commit: `${CURRENT_COMMIT}`
- Prior Release: `${PRIOR_RELEASE}`
- Prior Commit: `${PRIOR_COMMIT}`

A sentence or two describing if this release is optional or mandatory for validators. If it's mandatory,
include a short reason why.

# Validators

A list of commits in a bulleted list that are relevant to validators.

## Major Changes

Put the major changes here.

## Minor Changes

Put the minor changes here.

# Disperser

A list of commits in a bulleted list that are relevant to the disperser.

## Major Changes

Put the major changes here.

## Minor Changes

Put the minor changes here.

# Data API

A list of commits in a bulleted list that are relevant to the Data API.

## Major Changes

Put the major changes here.

## Minor Changes

Put the minor changes here.

# Contracts

A list of commits in a bulleted list that are relevant to the smart contracts.

## Major Changes

Put the major changes here.

## Minor Changes

Put the minor changes here.

# Integrations

A list of commits in a bulleted list that are relevant to integrations.

## Major Changes

Put the major changes here.

## Minor Changes

Put the minor changes here.

# Other

## Major Changes

Miscellaneous commits that don't fit into the above categories.

Put the major changes here.

## Minor Changes

Put the minor changes here.
```

Here is an example of how an entry for a commit should look:

```markdown
- `feat`: add 'litt prune' CLI tool by @cody-littley in [#1857](https://github.com/Layr-Labs/eigenda/pull/1857)
```

The important information to include is:

- The general type of commit (feat, fix, chore, docs, refactor, test, etc.)
- A short description of what the commit does
- The author of the commit (if available)
- A link to the pull request or commit (if available). Always prefer a link to a pull request, since that
  always has more information. But if you can't find the PR (e.g. an admin has force merged something), go with
  the link to the commit.

# Where to write the release notes

Release notes are stored in the `docs/release-notes` directory of the EigenDA repository. The filename
should be the tag or branch name of the current release, with a `.md` extension. For example, if the current release
is `v1.1.0`, the filename should be `v1.1.0.md`.

If you find an existing release notes file for the current release, this is probably the start of a draft. Be sure
to confirm it with the user, just in case.

If the file doesn't exist, let the user know and create a new file in the appropriate location.

# Iterative process

Instead of holding all information and writing it at the end, you should write into the release notes file as you go.
This will allow the user to audit your work as you go, and make corrections if necessary. It also allows the process
to be interrupted and resumed later.

At the bottom of the document, create a special section with a header of `# DRAFT - DO NOT PUBLISH`. In this section,
you can keep notes to yourself as to the current step you are on. When the document is eventually finalized, 
you can remove this section. If the user provides you with a draft that doesn't have this section, 
you can add it yourself.

Every time you complete a step in the process detailed in this document, make a note of it in the 
`# DRAFT - DO NOT PUBLISH` section. If you don't see a note marking the completion of a step, assume it has not
yet been done.

# Final verification

It's super important to make sure that the release notes are accurate. Perform the following steps at the end:

- Count the number of commits in the release notes. Compare this to the number of commits when you look at the git log.
  The numbers should match. If they don't, figure out why.
- Make sure that each commit only shows up exactly once.
- Ask the user to review the release notes in their entirety. Make any changes they request.
- Look for empty sections and remove them.
- Look for formatting errors, spelling mistakes, etc. Fix them.
