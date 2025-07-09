Purpose of this devcontainer is to be a claude code sandbox to allow running `claude --dangerously-skip-permissions`.
Copied from https://github.com/anthropics/claude-code/tree/main/.devcontainer

You can use this installing https://github.com/devcontainers/cli and running from the root of the repository:
```
devcontainer up --workspace-folder . --config .devcontainer/claude-code/devcontainer.json
```
Then you just have to `docker exec -it CONTAINER_ID bash` to get a shell in the container, from which you can run `claude --dangerously-skip-permissions` to start claude code.