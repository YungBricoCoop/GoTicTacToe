# Go Tic Tac Toe

Description : 3D Tic-Tac-Toe game using raycasting with the Ebiten library.
References :

- https://lodev.org/cgtutor/raycasting.html

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

**21.11.2025** :

- Add golangci-lint configuration file

Creation of `.golangci.yml` file, the file contains a known config from : https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322. This will help improve code quality over the whole project.

**28.11.2025** :

- Add GitHub Actions workflow for golangci-lint

Creation of `.github/workflows/golangci-lint.yml` file, the action used it the official one from golangci-lint : https://github.com/golangci/golangci-lint-action
