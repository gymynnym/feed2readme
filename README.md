# feed2readme

`feed2readme` is a tool that appends the latest posts from RSS feeds to a Markdown file. It supports [LogHub](https://loghub.me/) and [Zenn](https://zenn.dev/) as providers.

## Quick Start

```
$ touch feed2readme.toml README.md
$ echo -e "[feed]\nloghub = [\"gymynnym\"]\nzenn = [\"gymynnym\"]" > feed2readme.toml
$ go install github.com/gymynnym/feed2readme/cmd/feed2readme@latest
$ feed2readme
```

Default option values:

- `-c`: feed2readme.toml
- `-m`: README.md
- `-l`: 5 (max posts per provider)

## Usage on GitHub Actions

> [!IMPORTANT]
> `feed2readme.toml` and `README.md` must be present in the repository root.

```yaml
name: Update latest posts in README

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * *'

jobs:
  update-readme:
    runs-on: ubuntu-latest

    permissions:
      contents: write

    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v6
      with:
        go-version: 'stable'

    - run: go install github.com/gymynnym/feed2readme/cmd/feed2readme@latest
    - run: feed2readme -c feed2readme.toml -m README.md -l 5

    - name: git commit
      run: |
        git config --local user.email "gymynnym@users.noreply.github.com"
        git config --local user.name "gymynnym"
        git add README.md
        git diff --cached --quiet || (git commit -m "docs: update latest posts" && git push origin HEAD:${GITHUB_REF_NAME})
```

## Usage on Local

### Installation

```bash
$ go install github.com/gymynnym/feed2readme/cmd/feed2readme@latest
# or
$ git clone git@github.com/gymynnym/feed2readme.git
$ cd feed2readme
$ go build -o feed2readme ./cmd/feed2readme
```

### Configuration: `feed2readme.toml`

```toml
[feed]
loghub = ["gymynnym"]
zenn = ["gymynnym"]
```

Set only username arrays for each provider.

### Usage

```bash
$ feed2readme -c feed2readme.toml -m README.md -l 5
# config file: feed2readme.toml
# output file: README.md
# limit of items per provider: 5
```

## Markdown Behavior

- If both markers exist, content between them is replaced:
  - `<!-- FEED START -->`
  - `<!-- FEED END -->`
- If markers do not exist, a new feed block is appended at the end of the file.

Example Markdown output:

```markdown
<!-- FEED START -->
## Latest Posts

#### LogHub (한국어)

- [Foo](https://loghub.me/.../articles/...)

#### Zenn (日本語)

- [Bar](https://zenn.dev/.../articles/...)

<!-- FEED END -->
```
