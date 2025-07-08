# go-cantus-firmus

This project is a Cantus Firmus generator written in Go. It creates all possible Cantus Firmi in whole notes that satisfy the rules of strict style and saves them to a MusicXML file.

## Table of Contents
- [Description](#description)
- [Features](#features)
- [Project Structure](#project-structure)
- [Usage](#usage)
  - [Building](#building)
  - [Running](#running)
- [Cantus Firmus Rules](#cantus-firmus-rules)
- [License](#license)
- [Author](#author)

## Description

The `go-cantus-firmus` program is designed for musicians, composers, and students studying counterpoint. It automates the process of generating Cantus Firmi, ensuring adherence to the basic rules of strict style, such as:

- Starting and ending on the tonic.
- Predominantly stepwise motion with a limited number of leaps.
- Absence of augmented or diminished intervals, including in melodic contours.
- For minor mode, the 6th and 7th degrees are raised when necessary.

The results are saved in MusicXML format, allowing them to be easily opened in most music notation software.

## Features

- Generation of Cantus Firmi of a specified length (8 to 16 notes).
- Selection from several musical modes (major, dorian, phrygian, lydian, mixolydian, minor, locrian).
- Specification of the desired number of leaps in the Cantus Firmus.
- Filtering of results based on strict style rules.
- Saving generated Cantus Firmi to a MusicXML file.
- Option to choose how many Cantus Firmi to save (random selection if the number is less than the total).

## Project Structure

```
.
├── cmd
│   └── main.go               # Main executable file of the program
├── go.mod                    # Go module file
└── internal
    ├── cantusgen             # Logic for generating interval sequences
    │   ├── cantusgen.go
    │   └── cantusgen_test.go
    ├── music                 # Core musical structures (notes, intervals, Cantus Firmus)
    │   ├── cantus.go
    │   ├── cantus_test.go
    │   ├── interval.go
    │   ├── interval_test.go
    │   ├── note.go
    │   └── note_test.go
    ├── musicxml              # Functions for converting and saving to MusicXML format
    │   ├── musicxml_generator.go
    │   └── musicxml_generator_test.go
    ├── rules                 # Implementation of strict style and modal rules
    │   ├── moderules.go
    │   ├── moderules_test.go
    │   ├── rules.go
    │   └── rules_test.go
    └── utils                 # Helper utilities (e.g., selecting random items)
        ├── utils.go
        └── utils_test.go
```

## Usage

### Building

Before you begin, ensure you have Go installed.

To build the executable, navigate to the project's root directory and run:

```bash
go build -o cantusgen ./cmd
```

This will create an executable file named `cantusgen` (or `cantusgen.exe` on Windows) in the project's root directory.

### Running

After building, you can run the program:

```bash
./cantusgen
```

The program will ask you questions in the console:

1. Desired length of the Cantus Firmus (from 8 to 16 notes).
2. Mode (major, dorian, phrygian, lydian, mixolydian, minor, locrian).
3. Desired number of leaps (from 0 to length-4).

After entering the data, the program will generate Cantus Firmi and ask how many of them to save. The MusicXML file will be saved in the current directory with a name including generation parameters and a timestamp, for example: `cantus_length10_major_leaps1_20250621_150405.musicxml`.

## Cantus Firmus Rules

This generator applies the following main rules of strict style:

- Starts and ends on the tonic of the mode.
- Predominantly stepwise motion.
- Limited number of leaps (intervals greater than a second).
- Leaps greater than a third must be compensated by motion in the opposite direction.
- Absence of augmented and diminished intervals, including in melodic contours.
- Absence of excessive repetition of individual notes and note patterns.
- The upper and/or lower climaxes are reached only once.

## License

MIT

## Author

Sergei Shchetnikov