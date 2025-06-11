# GitHub Release Management Process

## Table of Contents

1. [Feature Freeze & Release Branch Creation](#1-feature-freeze--release-branch-creation)
2. [Changes to a Release Branch](#2-changes-to-a-release-branch)
   - [Change Policy](#change-policy)
   - [Change Process](#change-process)
3. [Tagging a Release](#3-tagging-a-release)
4. [Updates After Release Tag Has Been Cut](#4-updates-after-release-tag-has-been-cut)
5. [Official Release & Notes](#5-official-release--notes)
   - [Creating the Release](#creating-the-release)
   - [Release Notes (TODO)](#release-notes-todo)
6. [Post-Release Updates](#6-post-release-updates)

---

### 1. **Feature Freeze & Release Branch Creation**

Enacting a feature freeze helps to ensure that the code we publish to production environments is well tested and
mature. As a general rule of thumb, we should put in place a feature freeze two weeks prior to a release on testnet.

#### Feature Freeze

- **Trigger**: Either a date scheduled in advance, or tied to the completion of a key feature
- **Action**: Announce feature freeze date and time to the team
- **Result**: Development against `master` continues uninterrupted while release is prepared

#### Create Release Branch

- From latest `master` commit:
  - `git checkout master && git pull`
  - `git checkout -b release-<MAJOR>.<MINOR>`
    - **Note:** there is no patch number in the branch name. The same branch is used across multiple patch versions.
  - Example: `release-1.1`
- Push and protect the branch:
  - `git push origin release-1.1`
  - From the GitHub UI, set branch protections to prevent direct pushes and branch deletion

---

### 2. **Changes to a Release Branch**

#### Change Policy

- **High bar for inclusion**: Only critical bugfixes or business-critical features
    - Even bugfixes should not be reflexively included: only high-severity issues
- **Team consensus required**: Single engineer cannot make the decision
- **Public visibility**: Must have team discussion (e.g., Slack thread) before proceeding

#### Change Process

- **If change is also needed on `master`:**
  1. Submit PR and merge into `master` first
  2. Cherry-pick the squashed commit into the release branch
- **If change is release-only:**
  - Submit PR directly against the release branch
- **⚠️ NEVER push directly to the release branch**

---

### 3. **Tagging a Release**

- **When ready**, tag from HEAD of release branch:
  - Tag format: `v<MAJOR>.<MINOR>.<PATCH>`
  - Example: `v1.1.0`
  - `git checkout release-1.1`
  - `git tag v1.1.0`
  - `git push origin v1.1.0`
- **⚠️ Tags are immutable:**
  - NEVER force-push a tag to a different commit
  - If a mistake is made, create a new tag with incremented version

---

### 4. **Updates After Release Tag Has Been Cut**

- **Additional fixes** after initial tag may be required
- **Follow same change policy** as described in [Section 2](#2-changes-to-a-release-branch)
- **Do not tag reflexively** after every merge:
  - Accumulate changes until a meaningful patch set is ready
  - Create new release tag with incremented patch version (e.g., `v1.1.1`, `v1.1.2`)
- **Continue iteratively** until all critical issues are resolved

---

### 5. **Official Release & Notes**

#### Creating the Release

- **When ready to make the release public:**
  - If necessary, tag final patch version from release branch HEAD
  - Create GitHub release via UI, targeting the most recent tag
  - **Note**: Release will likely have non-zero patch version

#### Release Notes (TODO)

- Define a consistent format and content structure
- Include tooling guidance (likely an LLM prompt)
- Notes should include:
  - Feature summary
  - Fixes included
  - Contributor acknowledgments

---

### 6. **Post-Release Updates**

- **Additional patch updates** after release shipping follow the same process:
  - Merge new code into the release branch following [Section 2](#2-changes-to-a-release-branch) policy
  - Cut new tags as needed following [Section 3](#3-tagging-a-release) process
  - Create new GitHub releases targeting updated tags following [Section 5](#5-official-release--notes) process
