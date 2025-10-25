package gotar_hero

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

type KV struct {
	key   string
	value []any
}

type Section struct {
	values []KV
}

// A parsed .chart file
type UnstructuredChart struct {
	sections map[string]Section
}

const (
	// Last parsed line was a header
	StateHeader = iota
	// Last parsed line was a {
	StateOpenBracket
	// Last parsed line was a }
	StateCloseBracket
	// Last parsed line was a key-value
	StateKV
)

func ParseRaw(input io.Reader) (*UnstructuredChart, error) {
	scanner := bufio.NewScanner(input)
	header_regex := regexp.MustCompile(`\[\w+\]`)
	var fields []KV
	var name string

	var chart UnstructuredChart
	chart.sections = map[string]Section{}

	// this state expects a header, which should be the first thing in the file
	state := StateCloseBracket

	line_number := 1

	for {
		not_end := scanner.Scan()
		if !not_end {
			return &chart, scanner.Err()
		}
		line := string(scanner.Bytes())
		log.Debug("got line", "line", line)
		if header_regex.MatchString(line) {
			// header
			if state != StateCloseBracket {
				return nil, fmt.Errorf("unexpected header at line %v", line_number)
			}
			// strip square brackets
			name = strings.Clone(string(line[1 : len(line)-1]))
			log.Debug("section name", "name", name)
			state = StateHeader
		} else if line == "{" {
			// start of section
			if state != StateHeader {
				return nil, fmt.Errorf("unexpected { at line %v", line_number)
			}
			state = StateOpenBracket
		} else if line == "}" {
			// end of section
			if state != StateKV {
				if state == StateOpenBracket {
					return nil, fmt.Errorf("unexpected } at line %v (empty section)", line_number)
				}
				return nil, fmt.Errorf("unexpected } at line %v", line_number)
			}

			// add section to chart, first checking against duplicates
			_, exists := chart.sections[name]
			if exists {
				return nil, fmt.Errorf("duplicate section %v, second section ends at %v", name, line_number)
			}
			chart.sections[name] = Section{fields}
			fields = []KV{}

			state = StateCloseBracket
		} else {
			// key-value
			if state != StateOpenBracket && state != StateKV {
				return nil, fmt.Errorf("unexpected kv at line %v", line_number)
			}

			key, value_str, found := strings.Cut(line, "=")
			if !found {
				return nil, fmt.Errorf("missing '=' in kv at line %v", line_number)
			}

			key = strings.TrimSpace(key)

			item_regex := regexp.MustCompile(`[\d\.]+|("[^"]+")|\w+`)
			items := item_regex.FindAllString(value_str, 32)

			var values []any
			for i := range items {
				value := items[i]

				if value == "true" {
					// boolean
					values = append(values, true)
				} else if value == "false" {
					values = append(values, false)
				} else if value[0] == '"' {
					// quoted string
					if value[len(value)-1] != '"' {
						return nil, fmt.Errorf("unterminated quoted string at line %v", line_number)
					}
					literal := strings.Clone(string(value[1 : len(value)-1]))
					values = append(values, literal)
				} else if parsed_float, err := strconv.ParseFloat(value, 64); err == nil {
					// float
					values = append(values, parsed_float)
				} else if parsed_int, err := strconv.ParseInt(value, 10, 64); err == nil {
					// int
					values = append(values, float64(parsed_int))
				} else {
					// log.Debug("adding bare string", "key", key, "value", value)
					values = append(values, value)
				}
			}

			fields = append(fields, KV{key, values})

			state = StateKV
		}
		line_number += 1
	}
}

type TSChange struct {
	tick        int
	numerator   int
	denominator int
}

type TempoChange struct {
	tick  int
	tempo float64
}

type Note struct {
	tick int
	typ  int
	len  int
}

type InstrumentTrack struct {
	name  string
	notes []Note
}

type Chart struct {
	Title                string
	Artist               string
	Album                string
	Genre                string
	Year                 string
	Charter              string
	Resolution           int
	Difficulty           int
	Length               float64
	Offset               float64
	PreviewStart         float64
	PreviewEnd           float64
	TimeSignatureChanges []TSChange
	TempoChanges         []TempoChange
	Tracks               []InstrumentTrack
}

func Parse(uchart *UnstructuredChart) (*Chart, error) {
	var chart Chart
	chart.TimeSignatureChanges = []TSChange{}
	chart.TempoChanges = []TempoChange{}

	metadata, exists := uchart.sections["Song"]
	if !exists {
		return nil, fmt.Errorf("chart is missing [Song] section")
	}

	for i := range metadata.values {
		kv := metadata.values[i]
		log.Debug(kv.key)
		log.Debug(kv.value[0])
		switch kv.key {
		case "Name":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Name is not a string")
			}
			log.Debug("setting name", "name", t)
			chart.Title = t
		case "Artist":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Artist is not a string")
			}
			chart.Artist = t
		case "Album":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Album is not a string")
			}
			chart.Album = t
		case "Genre":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Genre is not a string")
			}
			chart.Genre = t
		case "Year":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Year is not a string")
			}
			chart.Year = t
		case "Charter":
			t, ok := kv.value[0].(string)
			if !ok {
				return nil, fmt.Errorf("chart Charter is not a string")
			}
			chart.Charter = t
		case "Resolution":
			t, ok := kv.value[0].(float64)
			i := int(t)
			if !ok || float64(i) != t {
				return nil, fmt.Errorf("chart Resolution is not an int")
			}
			chart.Resolution = i
		case "Difficulty":
			t, ok := kv.value[0].(float64)
			i := int(t)
			if !ok || float64(i) != t {
				return nil, fmt.Errorf("chart Resolution is not an int")
			}
			chart.Resolution = i
		case "Length":
			t, ok := kv.value[0].(float64)
			if !ok {
				return nil, fmt.Errorf("chart Length is not a decimal")
			}
			chart.Length = t
		case "Offset":
			t, ok := kv.value[0].(float64)
			if !ok {
				log.Error(reflect.TypeOf(kv.value[0]))
				return nil, fmt.Errorf("chart Offset is not a decimal")
			}
			chart.Offset = t
		case "PreviewStart":
			t, ok := kv.value[0].(float64)
			if !ok {
				return nil, fmt.Errorf("chart PreviewStart is not a decimal")
			}
			chart.Offset = t
		case "PreviewEnd":
			t, ok := kv.value[0].(float64)
			if !ok {
				return nil, fmt.Errorf("chart PreviewEnd is not a decimal")
			}
			chart.Offset = t
		}
	}

	sync, exists := uchart.sections["SyncTrack"]
	if !exists {
		return nil, fmt.Errorf("chart is missing [Sync] section")
	}
	for i := range sync.values {
		kv := sync.values[i]
		tick, err := strconv.ParseInt(kv.key, 10, 64)
		if err != nil {
			return nil, err
		}
		switch kv.value[0] {
		case "B":
			// Tempo Change
			t, ok := kv.value[1].(float64)
			i := int(t)
			if !ok || float64(i) != t {
				return nil, fmt.Errorf("chart Tempo Change is not an int")
			}
			chart.TempoChanges = append(chart.TempoChanges, TempoChange{int(tick), float64(i) / 1000.0})

		case "TS":
			// Time Signature Change
			t, ok := kv.value[1].(float64)
			numerator := int(t)
			if !ok || float64(numerator) != t {
				return nil, fmt.Errorf("chart Time Signature Change numerator is not an int")
			}

			denominator := 4
			if len(kv.value) == 3 {
				t, ok = kv.value[2].(float64)
				denominator = int(t)
				if !ok || float64(denominator) != t {
					return nil, fmt.Errorf("chart Time Signature Change denominator is not an int")
				}
				denominator = int(math.Exp2(float64(denominator)))
			}

			chart.TimeSignatureChanges = append(chart.TimeSignatureChanges, TSChange{int(tick), numerator, denominator})
		}
	}

	for section_name := range uchart.sections {
		if section_name == "Song" || section_name == "SyncTrack" || section_name == "Events" {
			continue
		}
		track := InstrumentTrack{section_name, []Note{}}
		log.Info("parsing track", "track", section_name)
		section := uchart.sections[section_name]
		for i := range section.values {
			kv := section.values[i]
			tick, err := strconv.ParseInt(kv.key, 10, 64)
			if err != nil {
				return nil, err
			}

			switch kv.value[0].(string) {
			case "N":
				// note
				typ := int(kv.value[1].(float64))
				length := int(kv.value[2].(float64))

				track.notes = append(track.notes, Note{int(tick), typ, length})
			case "S":
				// special phrase

			}
		}
		chart.Tracks = append(chart.Tracks, track)
	}

	return &chart, nil
}

//file, err := os.Open("notes.chart")
//	if err != nil {
//		panic(err.Error())
//	}
// skip BOM
//	file.Seek(3, 0)
//	uchart, err := gotar_hero.ParseRaw(file)
//	if err != nil {
//		panic(err.Error())
//	}
// log.SetLevel(log.DebugLevel)
//	chart, err := gotar_hero.Parse(uchart)
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Println(chart)
