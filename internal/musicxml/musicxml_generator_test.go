package musicxml

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestToMusicXML(t *testing.T) {
	tests := []struct {
		name        string
		sequences   [][]Note
		wantErr     bool
		errContains string
		wantXML     string // Partial XML string to check for key elements
	}{
		{
			name:        "empty sequences",
			sequences:   [][]Note{},
			wantErr:     true,
			errContains: "cannot create MusicXML from empty sequences",
		},
		{
			name: "inconsistent sequence lengths",
			sequences: [][]Note{
				{{Step: 0, Octave: 4, Alteration: 0}},
				{{Step: 0, Octave: 4, Alteration: 0}, {Step: 1, Octave: 4, Alteration: 0}},
			},
			wantErr:     true,
			errContains: "sequence 2 has length 2, expected 1",
		},
		{
			name: "single measure, single note, C4",
			sequences: [][]Note{
				{
					{Step: 0, Octave: 4, Alteration: 0}, // C4
				},
			},
			wantErr: false,
			wantXML: `<score-partwise>` +
				`<part-list>` +
				`<score-part id="P1">` +
				`<part-name>Cantus Firmus</part-name>` +
				`</score-part>` +
				`</part-list>` +
				`<part id="P1">` +
				`<measure number="1">` +
				`<attributes>` +
				`<divisions>4</divisions>` +
				`<key><fifths>0</fifths></key>` +
				`<time><beats>1</beats><beat-type>1</beat-type></time>` +
				`<clef><sign>G</sign><line>2</line></clef>` +
				`</attributes>` +
				`<direction placement="above">` +
				`<direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>240</per-minute></metronome></direction-type>` +
				`<sound tempo="240"></sound>` +
				`</direction>` +
				`<note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<barline location="right"><bar-style>light-heavy</bar-style></barline>` +
				`</measure>` +
				`</part>` +
				`</score-partwise>`,
		},
		{
			name: "single measure, multiple notes, C4 D4 E4",
			sequences: [][]Note{
				{
					{Step: 0, Octave: 4, Alteration: 0}, // C4
					{Step: 1, Octave: 4, Alteration: 0}, // D4
					{Step: 2, Octave: 4, Alteration: 0}, // E4
				},
			},
			wantErr: false,
			wantXML: `<score-partwise>` +
				`<part-list>` +
				`<score-part id="P1">` +
				`<part-name>Cantus Firmus</part-name>` +
				`</score-part>` +
				`</part-list>` +
				`<part id="P1">` +
				`<measure number="1">` +
				`<attributes>` +
				`<divisions>4</divisions>` +
				`<key><fifths>0</fifths></key>` +
				`<time><beats>3</beats><beat-type>1</beat-type></time>` +
				`<clef><sign>G</sign><line>2</line></clef>` +
				`</attributes>` +
				`<direction placement="above">` +
				`<direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>240</per-minute></metronome></direction-type>` +
				`<sound tempo="240"></sound>` +
				`</direction>` +
				`<note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<note><pitch><step>D</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<note><pitch><step>E</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<barline location="right"><bar-style>light-heavy</bar-style></barline>` +
				`</measure>` +
				`</part>` +
				`</score-partwise>`,
		},
		{
			name: "multiple measures, C4 in first, D4 in second",
			sequences: [][]Note{
				{{Step: 0, Octave: 4, Alteration: 0}}, // Measure 1: C4
				{{Step: 1, Octave: 4, Alteration: 0}}, // Measure 2: D4
			},
			wantErr: false,
			wantXML: `<score-partwise>` +
				`<part-list>` +
				`<score-part id="P1">` +
				`<part-name>Cantus Firmus</part-name>` +
				`</score-part>` +
				`</part-list>` +
				`<part id="P1">` +
				`<measure number="1">` +
				`<attributes>` +
				`<divisions>4</divisions>` +
				`<key><fifths>0</fifths></key>` +
				`<time><beats>1</beats><beat-type>1</beat-type></time>` +
				`<clef><sign>G</sign><line>2</line></clef>` +
				`</attributes>` +
				`<direction placement="above">` +
				`<direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>240</per-minute></metronome></direction-type>` +
				`<sound tempo="240"></sound>` +
				`</direction>` +
				`<note><pitch><step>C</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<barline location="right"><bar-style>light-heavy</bar-style></barline>` +
				`</measure>` +
				`<measure number="2">` +
				`<note><pitch><step>D</step><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>` +
				`<barline location="right"><bar-style>light-heavy</bar-style></barline>` +
				`</measure>` +
				`</part>` +
				`</score-partwise>`,
		},
		{
			name: "note with alteration (C#4)",
			sequences: [][]Note{
				{
					{Step: 0, Octave: 4, Alteration: 1}, // C#4
				},
			},
			wantErr: false,
			wantXML: `<note><pitch><step>C</step><alter>1</alter><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>`,
		},
		{
			name: "note with negative alteration (Db4)",
			sequences: [][]Note{
				{
					{Step: 1, Octave: 4, Alteration: -1}, // Db4
				},
			},
			wantErr: false,
			wantXML: `<note><pitch><step>D</step><alter>-1</alter><octave>4</octave></pitch><duration>4</duration><type>whole</type></note>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotXML, err := ToMusicXML(tt.sequences)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ToMusicXML() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ToMusicXML() error = %v, want containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ToMusicXML() unexpected error: %v", err)
				return
			}

			gotXML = strings.TrimPrefix(gotXML, xml.Header)
			gotXML = strings.ReplaceAll(gotXML, " ", "")
			gotXML = strings.ReplaceAll(gotXML, "\n", "")

			wantXMLFormatted := strings.ReplaceAll(tt.wantXML, " ", "")
			wantXMLFormatted = strings.ReplaceAll(wantXMLFormatted, "\n", "")

			if !strings.Contains(gotXML, wantXMLFormatted) {
				t.Errorf("ToMusicXML() got XML does not contain expected part.\nGot:\n%s\nWant part:\n%s", gotXML, wantXMLFormatted)
			}
		})
	}
}
