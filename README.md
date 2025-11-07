# Go Tic Tac Toe

## Changelog

**03.11.2025** :

- Upgrade go version to 1.25 and upgrade (ebiten, image) dependencies to latest versions.

To upgrade the dependencies, we used the following commands:

```bash
go mod edit -go=1.25

go get -u ./...

go mod tidy
```

**02.11.2025** :

- Upgrade deprecated ebiten text to v2

The `github.com/hajimehoshi/ebiten/v2/text/v2` package is now used instead of the deprecated `github.com/hajimehoshi/ebiten/v2/text` package. The `text.Draw` signature has changed and now it takes less arguments, so we need to create a DrawOption object that stores translation information.
