package musicxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"go-cantus-firmus/internal/music"
	"os"
)

// ScorePartwise represents the root element of a MusicXML score.
type ScorePartwise struct {
	XMLName  xml.Name `xml:"score-partwise"`
	PartList PartList `xml:"part-list"`
	Part     Part     `xml:"part"`
}

// PartList contains the score-parts.
type PartList struct {
	XMLName   xml.Name  `xml:"part-list"`
	ScorePart ScorePart `xml:"score-part"`
}

// ScorePart represents a single part in the score.
type ScorePart struct {
	XMLName  xml.Name `xml:"score-part"`
	ID       string   `xml:"id,attr"`
	PartName PartName `xml:"part-name"`
}

// PartName represents the name of a part.
type PartName struct {
	XMLName xml.Name `xml:"part-name"`
	Text    string   `xml:",chardata"`
}

// Part represents the musical content of a single part.
type Part struct {
	XMLName  xml.Name  `xml:"part"`
	ID       string    `xml:"id,attr"`
	Measures []Measure `xml:"measure"`
}

// Measure represents a single measure in a part.
type Measure struct {
	XMLName    xml.Name    `xml:"measure"`
	Number     int         `xml:"number,attr"`
	Attributes *Attributes `xml:"attributes,omitempty"`
	Direction  *Direction  `xml:"direction,omitempty"`
	Notes      []NoteXML   `xml:"note"`
	Barline    *Barline    `xml:"barline,omitempty"`
}

// Attributes contains musical attributes like divisions, key, time, and clef.
type Attributes struct {
	XMLName   xml.Name `xml:"attributes"`
	Divisions int      `xml:"divisions,omitempty"`
	Key       *Key     `xml:"key,omitempty"`
	Time      *Time    `xml:"time,omitempty"`
	Clef      *Clef    `xml:"clef,omitempty"`
}

// Key represents the key signature.
type Key struct {
	XMLName xml.Name `xml:"key"`
	Fifths  int      `xml:"fifths"`
}

// Time represents the time signature.
type Time struct {
	XMLName  xml.Name `xml:"time"`
	Beats    string   `xml:"beats"`
	BeatType string   `xml:"beat-type"`
}

// Clef represents the clef.
type Clef struct {
	XMLName xml.Name `xml:"clef"`
	Sign    string   `xml:"sign"`
	Line    int      `xml:"line"`
}

// NoteXML represents a musical note within a measure.
type NoteXML struct {
	XMLName  xml.Name `xml:"note"`
	Pitch    Pitch    `xml:"pitch"`
	Duration int      `xml:"duration"`
	Type     string   `xml:"type"`
}

// Pitch represents the pitch of a note.
type Pitch struct {
	XMLName xml.Name `xml:"pitch"`
	Step    string   `xml:"step"`
	Alter   *int     `xml:"alter,omitempty"`
	Octave  int      `xml:"octave"`
}

// Barline represents a barline element.
type Barline struct {
	XMLName  xml.Name `xml:"barline"`
	Location string   `xml:"location,attr"`
	BarStyle BarStyle `xml:"bar-style"`
}

// BarStyle represents the style of the barline.
type BarStyle struct {
	XMLName xml.Name `xml:"bar-style"`
	Text    string   `xml:",chardata"`
}

// Direction represents a musical direction, e.g., tempo.
type Direction struct {
	XMLName       xml.Name      `xml:"direction"`
	Placement     string        `xml:"placement,attr"`
	DirectionType DirectionType `xml:"direction-type"`
	Sound         *Sound        `xml:"sound,omitempty"`
}

// DirectionType contains different types of directions.
type DirectionType struct {
	XMLName   xml.Name   `xml:"direction-type"`
	Metronome *Metronome `xml:"metronome,omitempty"`
}

// Metronome represents a metronome mark for tempo.
type Metronome struct {
	XMLName   xml.Name `xml:"metronome"`
	BeatUnit  string   `xml:"beat-unit"`
	PerMinute int      `xml:"per-minute"`
}

// Sound contains sound-related attributes, e.g., tempo.
type Sound struct {
	XMLName xml.Name `xml:"sound"`
	Tempo   float64  `xml:"tempo,attr"`
}

// Note represents a musical note for conversion to MusicXML.
type Note struct {
	Step       int
	Octave     int
	Alteration int
}

// ToMusicXML converts a slice of note sequences into a MusicXML string.
func ToMusicXML(sequences [][]Note) (string, error) {
	if len(sequences) == 0 {
		return "", errors.New("cannot create MusicXML from empty sequences")
	}

	// Check that all sequences have the same length
	expectedLength := len(sequences[0])
	for i, seq := range sequences {
		if len(seq) != expectedLength {
			return "", fmt.Errorf("sequence %d has length %d, expected %d", i+1, len(seq), expectedLength)
		}
	}

	stepMap := []string{"C", "D", "E", "F", "G", "A", "B"}

	var measures []Measure
	for measureNum, sequence := range sequences {
		var notesXML []NoteXML
		for _, n := range sequence {
			var alter *int
			if n.Alteration != 0 {
				a := n.Alteration
				alter = &a
			}

			notesXML = append(notesXML, NoteXML{
				Pitch: Pitch{
					Step:   stepMap[n.Step],
					Alter:  alter,
					Octave: n.Octave,
				},
				Duration: 4,
				Type:     "whole",
			})
		}

		measure := Measure{
			Number: measureNum + 1,
			Notes:  notesXML,
			Barline: &Barline{
				Location: "right",
				BarStyle: BarStyle{Text: "light-heavy"},
			},
		}

		if measureNum == 0 {
			beats := fmt.Sprintf("%d", len(sequence))
			measure.Attributes = &Attributes{
				Divisions: 4,
				Key:       &Key{Fifths: 0},
				Time: &Time{
					Beats:    beats,
					BeatType: "1",
				},
				Clef: &Clef{
					Sign: "G",
					Line: 2,
				},
			}
			measure.Direction = &Direction{
				Placement: "above",
				DirectionType: DirectionType{
					Metronome: &Metronome{
						BeatUnit:  "quarter",
						PerMinute: 300,
					},
				},
				Sound: &Sound{
					Tempo: 300.0,
				},
			}
		}

		measures = append(measures, measure)
	}

	score := ScorePartwise{
		PartList: PartList{
			ScorePart: ScorePart{
				ID:       "P1",
				PartName: PartName{Text: "Cantus Firmus"},
			},
		},
		Part: Part{
			ID:       "P1",
			Measures: measures,
		},
	}

	output, err := xml.MarshalIndent(score, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling MusicXML: %w", err)
	}

	return xml.Header + string(output), nil
}

// ConvertRealizationsToXMLNotes converts a slice of music.Realization to MusicXML Note format
func ConvertRealizationsToXMLNotes(realizations []music.Realization) [][]Note {
	var xmlSequences [][]Note
	for _, realization := range realizations {
		var xmlSeq []Note
		for _, note := range realization {
			xmlSeq = append(xmlSeq, Note{
				Step:       note.Step,
				Octave:     note.Octave,
				Alteration: note.Alteration,
			})
		}
		xmlSequences = append(xmlSequences, xmlSeq)
	}
	return xmlSequences
}

// GenerateAndSaveMusicXML generates MusicXML from note sequences and saves to file
func GenerateAndSaveMusicXML(sequences [][]Note, filename string) error {
	xmlString, err := ToMusicXML(sequences)
	if err != nil {
		return fmt.Errorf("error generating MusicXML: %w", err)
	}

	err = os.WriteFile(filename, []byte(xmlString), 0644)
	if err != nil {
		return fmt.Errorf("error writing MusicXML file: %w", err)
	}
	return nil
}
