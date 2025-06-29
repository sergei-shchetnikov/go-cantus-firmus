package main

import (
	"fmt"
	"go-cantus-firmus/internal/cantusgen"
	"go-cantus-firmus/internal/music"
	"go-cantus-firmus/internal/musicxml"
	"math/rand"
	"os"
	"time"
)

// Project: go-cantus-firmus
// Created: 2025-06-21

// selectRandomItems selects 'count' random items from a slice using reservoir sampling algorithm
func selectRandomItems[T any](items []T, count int) []T {
	if count <= 0 || len(items) == 0 {
		return nil
	}
	if count >= len(items) {
		return items
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]T, count)
	copy(result, items[:count])

	for i := count; i < len(items); i++ {
		j := r.Intn(i + 1)
		if j < count {
			result[j] = items[i]
		}
	}

	return result
}

// This is a test/demo program to verify the core functionality:
// 1. Generates cantus firmus melodies
// 2. Converts them to musical notation
// 3. Exports to MusicXML format
func main() {
	// Generate cantus firmus melodies of length 10
	generated := cantusgen.GenerateCantus(10)
	fmt.Printf("Generated %d cantus firmus melodies\n", len(generated))

	if len(generated) == 0 {
		panic("failed to generate any cantus firmus melodies")
	}

	// Select up to 10 random melodies
	selected := selectRandomItems(generated, 10)

	// Convert each melody to A minor realization
	var sequences [][]musicxml.Note
	for _, intervals := range selected {
		var cf music.CantusFirmus
		for _, val := range intervals {
			cf = append(cf, music.Interval(val))
		}

		realization, err := cf.Realize("Minor")
		if err != nil {
			panic(err)
		}

		var notes []musicxml.Note
		for _, note := range realization {
			notes = append(notes, musicxml.Note{
				Step:       note.Step,
				Octave:     note.Octave,
				Alteration: note.Alteration,
			})
		}
		sequences = append(sequences, notes)
	}

	// Create MusicXML file
	xmlStr, err := musicxml.ToMusicXML(sequences)
	if err != nil {
		panic(err)
	}

	// Save to file
	err = os.WriteFile("cantus_firmus.musicxml", []byte(xmlStr), 0644)
	if err != nil {
		panic(err)
	}
}
