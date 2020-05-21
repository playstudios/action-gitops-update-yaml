# Github Action - Update Yaml

This repository contains a Github action to update values within a yaml file in a remote repository. It can also be used as a library within other go applications.

## Usage as a Github Action

... todo

## Usage as a library

The following yaml

```yaml
path:
    to:
      the:
        - value: hello
        - value: world
      nested: |
        abc: 123
```

and example

```go
repo, err := git.CloneGithub("playstudios/my-repo", "abc123")
file, err := repo.ReadFile("my/file/name.yaml")

document, err := document.Parse(file)
err = document.Set("path.to.the.1.value", "universe")
err = document.Set("path.to.nested.abc", 456)
err = document.Write(file)

err = repo.Commit("updated something")
```

outputs

```yaml
path:
    to:
      the:
        - hello: bob?
        - hello: universe
      nested: |
        abc: 456
```
