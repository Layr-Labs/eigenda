# Release Notes

Your job is to help the user compile release notes for the EigenDA repository. You will assist the user in gathering
and sorting information about new features, bug fixes, improvements, etc., and based on the feedback from the user
you will generate a well-structured release notes document.

# Information you will need to gather

You will need to gather some information from the user to create comprehensive release notes.

1. The tag/branch for the prior release (e.g., v1.0.0).
    a. The exact commit for this release. If it's a branch, use the latest 
      commit in that branch. Always use the upstream commit.
2. The tag/branch for the current release being documented (e.g., v1.1.0).
    a. The exact commit for this release. If it's a branch, use the latest 
       commit in that branch. Always use the upstream commit.
3. The list of commits between the prior release and the current release.
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
4. Whether or not this is an optional or a mandatory release for validators. The user will need to be the source
   of this information.
     a. If it's a mandatory release, the reason why it's mandatory.
5. Optionally, the user may provide a draft of the release notes that you can help polish. The draft might be release
   notes that you have helped work on previously, or they might be notes that they have written themselves. If the user
   doesn't specify, always ask them if they have a draft to use as a starting point.

# How to gather the information

The best way to gather information is to get it from git/github, if it is reasonable to do so. For example, if the user
provides a branch name, you can look it up to get the latest commit.

Some information must come from the user. Sometimes the user will volunteer this information. Other times, you will
need to prompt them for it. 

When you ask the user for information, only ask for one thing at a time. As a rule of thumb, if it will take the user
multiple sentences to answer, consider breaking it up into multiple questions.

Some commits begin with a number, e.t. "2048". These often refer to issue numbers. Issues are tracked in linear.
The URL for a linear issue is `https://linear.app/eigenda/issue/<issue-number>`. If you are attempting to gather
more information about a commit that references an issue, you can provide the user with a link to the issue. If you
are able to look up the issue and gather more information on your own, that's even better. Note that permissions
might be tricky here, since this is not public information. The user is probably logged into linear in their browser,
but you may have authentication problems. If there is any way to piggy-back on the user's authentication, do so.

## Sorting and understanding commits

Commit messages can be terse, and you may be lacking context on some of the changes, or on the subject matter in 
general. That's ok, the user should be able to provide context.

For each commit you are unsure about, ask the user for clarification. Be sure to present the user with all information
you have available to you. It's very important as well to give the user a link they can click on to see the commit or 
PR in question.

## Verifying information

Once you have sorted commits, it's important to verify the information with the user. 

For each category, do the following:

Tell the user that you'd like to verify the contents of the category. Present the user with a list of 8 or fewer
commits at a time. Each commit should be in an enumerated list. Tell the user that they should type a list of numbers
for commits that are out of place, or if they want to change the importance level. For each commit listed by the
user, ask them what category or importance level it should be instead (one at a time). If the user just directly tells
you what changes to make, that's ok too.

# Release Notes Template

Below is a rough template for the release notes. Release notes are always markdown files. Sometimes a section
might be empty, and that's okay. If that happens, omit that section from the final output.

Note that sometimes there may be some major features that deserve their own section.

```markdown

# ${CURRENT_RELEASE} - Release Notes

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
- A link to the pull request or commit (if available)

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