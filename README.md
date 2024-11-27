# contributors

The `contributors` GitHub Action helps ensuring that contributors get the recognition that they deserve. It will generate code that can be used within markdown files, e.g., `README.md` to show profile picture of people that contributed to a repository and the link to their profiles.

## Example usage

`github-action`

```yaml
name: update-contributors

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * 0'

jobs:
  update-contributors:
    name: validate-pull-request-title
    runs-on: ubuntu-latest
    steps:
      - name: checkout
        uses: actions/checkout@v4

      - name: update-contributors
        uses: kontrolplane/contributors@latest

      - name: open-pull-request
        uses: 
```

`README.md`

```markdown
...

## contributors

[//]: kontrolplane/contributors

[//]: kontrolplane/contributors
```

## Example output [code]

`README.md`

```markdown
[//]: kontrolplane/contributors

<a href="https://github.com/levivannoort"><img src="https://avatars.githubusercontent.com/u/73097785?v=4" title="levivannoort" width="50" height="50"></a>

[//]: kontrolplane/contributors
```

## Example output [overview]

[//]: kontrolplane/contributors

<a href="https://github.com/levivannoort"><img src="https://avatars.githubusercontent.com/u/73097785?v=4" title="levivannoort" width="50" height="50"></a>

[//]: kontrolplane/contributors