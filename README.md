# contributors

The `contributors` GitHub Action helps ensuring that contributors get the recognition that they deserve. It will generate code that can be used within markdown files, e.g., `README.md` to show profile picture of people that contributed to a repository and the link to their profiles.

## Example output [overview]

[//]: kontrolplane/generate-contributors-list

<a href="https://github.com/levivannoort"><img src="https://avatars.githubusercontent.com/u/73097785?v=4" title="levivannoort" width="50" height="50"></a>

[//]: kontrolplane/generate-contributors-list

## Example usage [basic - public repository]

`github-action`
```yaml
name: update-contributors

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * 0'

permissions:
  contents: write
  pull-requests: write

jobs:
  update-contributors:
    name: validate-pull-request-title
    runs-on: ubuntu-latest
    steps:
      - name: checkout
        uses: actions/checkout@v4

      - name: update-contributors
        uses: kontrolplane/generate-contributors-list@v1.0.0
        with:
          owner: kontrolplane
          repository: pull-request-title-validator

      - name: check-for-changes
        id: check
        run: |
          if git diff --quiet README.md; then
            echo "No changes to commit"
            echo "changes=false" >> $GITHUB_OUTPUT
            exit 0
          fi
          echo "changes=true" >> $GITHUB_OUTPUT

      - name: open-pull-request
        if: steps.check.outputs.changes == 'true'
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          git config user.name github-actions
          git config user.email github-actions@github.com
          git checkout -b update-contributors
          git add README.md
          git commit -m "chore: update contributors section"
          git push -u origin update-contributors
          gh pr create \
            --title "chore: update contributors" \
            --body "Automatically update contributors section." \
            --base main \
            --head update-contributors
```

`README.md`
```markdown
...

## contributors

[//]: kontrolplane/generate-contributors-list

[//]: kontrolplane/generate-contributors-list
```

## Example usage [basic - private repository]

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
        uses: kontrolplane/generate-contributors-list@latest
        with:
          token: ${secrets.GITHUB_TOKEN}
          owner: kontrolplane
          repository: private-repository

...
```

`README.md`
```markdown
...

## contributors

[//]: kontrolplane/generate-contributors-list

[//]: kontrolplane/generate-contributors-list
```

## Example output [code]

`README.md`

```markdown
[//]: kontrolplane/generate-contributors-list

<a href="https://github.com/levivannoort"><img src="https://avatars.githubusercontent.com/u/73097785?v=4" title="levivannoort" width="50" height="50"></a>

[//]: kontrolplane/generate-contributors-list
```
