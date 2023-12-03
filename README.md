## Quiz
My solutions to https://github.com/gophercises/quiz

## Index
- [Usage](#usage)
- [Build](#build)
  - [Local](#local)
  - [Github Actions](#github-actions)
      - [Supported Platforms](#supported-platforms)
- [Tests](#tests)

### Usage
```bash
Usage of ./quiz:
  -path string
        path csv quiz file
  -seconds int
        seconds to finish quiz (default 30)
  -shuffle
        shuffle the questions randomly
```

Example:
```bash
./quiz -path=csv/problems.csv -seconds=5 --shuffle
```

### Build

#### Local
To build new binary after changing code:

```bash
make build
```

This will build a binary for architecture `GOOS=linux GOARCH=amd64` in uncommitted folder `bin`

#### Github Actions
Code pushed to master branch will automatically generate a release with binaries available for different architectures.

##### Supported Platforms
- `linux/amd64`
- `darwin/amd64`
- `windows/amd64`

### Tests
To run tests:

```bash
make tests
```

Test are also run in configured github actions workflow each time new commits are pushed to `master`