package main

import (
	"bufio"
	"fmt"
	"go-cantus-firmus/internal/cantusgen"
	"go-cantus-firmus/internal/music"
	"go-cantus-firmus/internal/musicxml"
	"go-cantus-firmus/internal/rules"
	"go-cantus-firmus/internal/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Project: go-cantus-firmus
// Created: 2025-06-21

// This is a test/demo program to verify the core functionality
func main() {
	fmt.Println("=== Cantus Firmus Generator ===")
	fmt.Println("This program generates all possible cantus firmi in whole notes")
	fmt.Println("that satisfy the rules of strict style and saves them to a MusicXML file.")
	fmt.Println()

	// Get user input
	length := getIntegerInput("Enter desired length (8-16 notes): ", 8, 16)
	mode := getModeInput()
	leaps := getIntegerInput(fmt.Sprintf("Enter desired number of leaps in the cantus firmus (0-%d): ", length-4), 0, length-4)

	fmt.Println("\nGenerating... Please wait...")
	startTime := time.Now()

	// Generate interval sequences with length-1 and leaps as part of allowed intervals
	intervalSequences := cantusgen.GenerateCantus(length-1, []int{leaps})
	if len(intervalSequences) == 0 {
		fmt.Println("Generation failed: no sequences could be generated.")
		return
	}

	var validRealizations []music.Realization

	// Process each sequence
	for _, seq := range intervalSequences {
		// Convert []int to []music.Interval
		intervals := make(music.CantusFirmus, len(seq))
		for i, val := range seq {
			intervals[i] = music.Interval(val)
		}

		// Realize the sequence in the chosen mode (with capitalized mode name)
		realization, err := intervals.Realize(strings.Title(mode))
		if err != nil {
			continue // Skip sequences with realization errors
		}

		// Check for augmented/diminished intervals
		if rules.IsFreeOfAugmentedDiminished(realization) {
			validRealizations = append(validRealizations, realization)
		}
	}

	generationTime := time.Since(startTime).Round(time.Millisecond)
	fmt.Printf("\nGeneration completed in %s\n", generationTime)
	fmt.Printf("Found %d valid cantus firmi\n", len(validRealizations))

	if len(validRealizations) == 0 {
		fmt.Println("No valid cantus firmi were generated.")
		return
	}

	// Ask how many to save
	maxToSave := len(validRealizations)
	saveCount := getIntegerInput(
		fmt.Sprintf("How many cantus firmi to save? (1-%d, selection will be random if less than total): ", maxToSave),
		1, maxToSave*2) // Allow numbers larger than max

	var toSave []music.Realization
	if saveCount >= maxToSave {
		toSave = validRealizations
		fmt.Printf("Saving all %d cantus firmi...\n", maxToSave)
	} else {
		toSave = utils.SelectRandomItems(validRealizations, saveCount)
		fmt.Printf("Randomly selecting %d out of %d cantus firmi to save...\n", saveCount, maxToSave)
	}

	// Generate filename with parameters
	filename := fmt.Sprintf("cantus_length%d_%s_leaps%d_%s.musicxml",
		length, strings.ToLower(mode), leaps, time.Now().Format("20060102_150405"))

	// Convert to MusicXML format
	xmlSequences := musicxml.ConvertRealizationsToXMLNotes(toSave)

	// Save to file
	err := musicxml.GenerateAndSaveMusicXML(xmlSequences, filename)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	fmt.Printf("\nSuccessfully saved %d cantus firmi to %s\n", len(toSave), filename)
}

func getIntegerInput(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		value, err := strconv.Atoi(input)
		if err != nil || value < min || value > max {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		return value
	}
}

func getModeInput() string {
	modes := []string{"major", "dorian", "phrygian", "lydian", "mixolydian", "minor", "locrian"}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter mode (major, dorian, phrygian, lydian, mixolydian, minor, locrian): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		for _, mode := range modes {
			if input == mode {
				return mode
			}
		}

		fmt.Println("Invalid mode. Please choose from the available options.")
	}
}
