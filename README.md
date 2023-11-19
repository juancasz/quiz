## Quiz
My solutions to https://github.com/gophercises/quiz

## Index
- [Usage](#usage)
- [Build](#build)
- [Tests](#tests)

### Usage
```bash
Usage of ./bin/quiz:
  -path string
        path csv quiz file
  -seconds int
        seconds to finish quiz (default 30)
  -shuffle
        shuffle the questions randomly
```

Example:
```bash
./bin/quiz -path=csv/problems.csv -seconds=5 --shuffle
```

### Build
To build new binary after changing code:

```bash
make build
```

### Tests
To run tests:

```bash
make tests
```