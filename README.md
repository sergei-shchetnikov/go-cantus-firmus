# Cantus Firmus Generator

This CLI program generates all possible Cantus Firmi that adhere to the rules of strict style and saves them to a MusicXML file.

## Features

- Generation of Cantus Firmi of a specified length (8 to 16 notes).
- Selection from several musical modes (major, dorian, phrygian, lydian, mixolydian, minor, locrian).
- Specification of the desired number of leaps in the Cantus Firmus.
- Filtering of results based on strict style rules.
- Saving generated Cantus Firmi to a MusicXML file.
- Option to choose how many Cantus Firmi to save (random selection if the number is less than the total).

## Example
Here is an example of a generated Cantus Firmus with the parameters: length 10, major mode, and 3 leaps.
![](./images/cantus_1.PNG)

## Cantus Firmus Rules

- Starting and ending on the tonic.
- Predominantly stepwise motion with a limited number of leaps.
- Leaps greater than a third must be compensated by motion in the opposite direction.
- Absence of excessive repetition of individual notes and note patterns.
- The upper and/or lower climaxes are reached only once.
- Absence of augmented or diminished intervals, including in melodic contours.
- For minor mode, the 6th and 7th degrees are raised when necessary.

## Project Structure

```
.
├── cmd
│   └── main.go               # Main executable file of the program
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

### How to Install and Run:

1.  **Download the archive** for your operating system from the "Assets" section below.

2.  **Unpack the archive** to a convenient location.

      * **Windows:**
          * Download `cantusgen_windows_amd64.zip`.
          * Unzip it and run `cantusgen_windows_amd64.exe` from your Command Prompt or PowerShell.
      * **Linux/macOS:**
          * Download `cantusgen_darwin_amd64.tar.gz` (for macOS) or `cantusgen_linux_amd64.tar.gz` (for Linux).
          * Unpack the archive using: `tar -xzf cantusgen_darwin_amd64.tar.gz`
          * Open your terminal, navigate to the folder with the unpacked file.
          * Make the file executable: `chmod +x cantusgen_darwin_amd64`
          * Run the program: `./cantusgen_darwin_amd64`

### Running

After unpacking the downloaded archive, you can run the program. The method depends on your operating system:

### Windows

1.  Open **Command Prompt** or **PowerShell**.
2.  Navigate to the directory where you unpacked the `cantusgen_windows_amd64.zip` archive.
3.  Run the executable:
    ```bash
    cantusgen_windows_amd64.exe
    ```

### macOS

1.  Open **Terminal**.
2.  Navigate to the directory where you unpacked the `cantusgen_darwin_amd64.tar.gz` archive.
3.  Ensure the file is executable (if you haven't already):
    ```bash
    chmod +x cantusgen_darwin_amd64
    ```
4.  Run the program:
    ```bash
    ./cantusgen_darwin_amd64
    ```

### Linux

1.  Open **Terminal**.
2.  Navigate to the directory where you unpacked the `cantusgen_linux_amd64.tar.gz` archive.
3.  Ensure the file is executable (if you haven't already):
    ```bash
    chmod +x cantusgen_linux_amd64
    ```
4.  Run the program:
    ```bash
    ./cantusgen_linux_amd64
    ```

The program will ask you questions in the console:

1. Desired length of the Cantus Firmus (from 8 to 16 notes).
2. Mode (major, dorian, phrygian, lydian, mixolydian, minor, locrian).
3. Desired number of leaps.

After entering the data, the program will generate Cantus Firmi and ask how many of them to save. The MusicXML file will be saved in the current directory with a name including generation parameters and a timestamp, for example: `cantus_length10_major_leaps1_20250621_150405.musicxml`.

## License

MIT

## Author

Sergei Shchetnikov