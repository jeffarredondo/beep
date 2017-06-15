// beep - Sound notifier with music note engine
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

const (
	sampleRate     = 44100
	sampleRate64   = float64(sampleRate)
	bitsPerSample  = 16
	sample16bit    = bitsPerSample == 16
	sampleAmp16bit = 32767.0
	sampleAmp8bit  = 127.0
	noteC5         = 523.25
)

var (
	flagHelp      = flag.Bool("h", false, "help")
	flagCount     = flag.Int("c", 1, "beep count")
	flagFreq      = flag.Float64("f", noteC5, "frequency in Hertz (1-22050)")
	flagVolume    = flag.Int("v", 100, "volume (1-100)")
	flagDuration  = flag.Int("t", 250, "beep time duration in millisecond (1-60000)")
	flagDevice    = flag.String("d", "default", "audio device, Linux example: hw:0,0")
	flagLine      = flag.Bool("l", false, "beep per line from stdin")
	flagMusic     = flag.Bool("m", false, "play music notes from stdin (see beep notation)")
	flagPrintDemo = flag.Bool("p", false, "print demo music sheet (Mozart K33b)")
	flagBell      = flag.Bool("b", false, "send bell to PC speaker")
	flagQuiet     = flag.Bool("q", false, "quiet stdout while playing music")
	flagNotes     = flag.Bool("n", false, "print notes while playing music")
	flagOutput    = flag.String("o", "", "output music waveform to file. Use '-' for stdout")
	flagWeb       = flag.Bool("w", false, "start beep web server")
	flagWebIP     = flag.String("a", "127.0.0.1:4444", "web server address")
	flagVoiceDl   = flag.Bool("vd", false, "download voice files, by default downloads all voices")
	flagMidiPlay  = flag.String("mp", "", "play MIDI file")
	flagMidiNote  = flag.String("mn", "", "parses MIDI file and print notes")
	flagPlayNotes = flag.String("play", "", "play notes from command argument")
)

const intro = `beep - Sound notifier with music note engine

Batbold Dashzeveg 2014-12-31`

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover: ", r)
		}
	}()
	flag.Parse()

	help := *flagHelp
	freqHertz := *flagFreq
	count := *flagCount
	volume := *flagVolume
	duration := *flagDuration
	device := *flagDevice
	lineBeep := *flagLine
	playMusic := *flagMusic
	printDemo := *flagPrintDemo
	writeBell := *flagBell
	webServer := *flagWeb
	downloadVoices := *flagVoiceDl
	midiPlay := *flagMidiPlay
	midiNote := *flagMidiNote
	musicNotes := *flagPlayNotes

	if help {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\n%s\n%s\n%s\n%s\n",
			intro,
			beepNotation,
			builtinMusic[0].Notation, //demoMusic,
			demoHelp,
		)
		return
	}
	if printDemo {
		fmt.Print(demoMusic)
		return
	}
	if volume < 1 || volume > 100 {
		volume = 100
	}
	if duration < 1 || duration > 1000*60 {
		duration = 250
	}
	if freqHertz < 1 || freqHertz > sampleRate64/2 {
		fmt.Fprintf(os.Stderr, "Invalid frequency. Must be 1-22050")
		os.Exit(1)
	}
	freq := hertzToFreq(freqHertz)

	music = &Music{
		played:     make(chan bool),
		stopped:    make(chan bool),
		linePlayed: make(chan bool),
		output:     *flagOutput,
	}

	openSoundDevice(device)
	initSoundDevice()
	defer closeSoundDevice()

	if lineBeep {
		beepPerLine(volume, freq)
		return
	}

	if playMusic {
		playMusicSheet(volume)
		return
	}

	if writeBell {
		sendBell()
		return
	}

	if webServer {
		startWebServer(*flagWebIP)
		return
	}

	if downloadVoices {
		var names []string
		for i, arg := range os.Args {
			if i == 0 || strings.HasPrefix(arg, "-") {
				continue
			}
			names = append(names, arg)
		}
		downloadVoiceFiles(os.Stdout, names)
		return
	}

	if len(midiPlay) > 0 {
		parseMidiBeep(midiPlay)
		return
	}

	if len(midiNote) > 0 {
		parseMidiNote(midiNote)
		return
	}

	if len(musicNotes) > 0 {
		playMusicNotesFromCL(musicNotes, volume)
		return
	}

	beep(volume, duration, count, freq)
}

func beepDefault() {
	freq := hertzToFreq(noteC5)
	beep(100, 250, 1, freq)
}

func beep(volume, duration, count int, freq float64) {
	bar := sampleAmp16bit * (float64(volume) / 100.0)
	samples := int(sampleRate64 * (float64(duration) / 1000.0))
	rest := 0
	if count > 1 {
		rest = (sampleRate / 20) * 4 // 200ms
	}
	buf := make([]int16, samples+rest)
	var last int16
	var fade = 1024
	if samples < fade {
		fade = 1
	}
	for i := range buf {
		if i < samples-fade {
			buf[i] = int16(bar * math.Sin(float64(i)*freq))
			last = buf[i]
		} else {
			if last > 0 {
				last -= 31
			} else {
				last += 31
			}
			buf[i] = last
		}
	}
	initSoundDevice()
	for i := 0; i < count; i++ {
		go playback(buf, buf)
		<-music.linePlayed
	}
	flushSoundBuffer()
}

// Beeps per line read from Stdin
func beepPerLine(volume int, freq float64) {
	buf := make([]int16, sampleRate/5)
	bar := sampleAmp16bit * (float64(volume) / 100.0)
	gap := sampleRate / 6
	var last int16
	for i := range buf {
		if i < gap {
			buf[i] = int16(bar * math.Sin(float64(i)*freq))
			last = buf[i]
		} else {
			if last > 0 {
				last -= 31
			} else {
				last += 31
			}
			buf[i] = last
		}
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		fmt.Print(string(line))
		if !isPrefix {
			fmt.Println()
			go playback(buf, buf)
			<-music.linePlayed
		}
	}
	flushSoundBuffer()
}

// Reads music sheet files given as arguments. If no files given, reads Stdin.
func playMusicSheet(volume int) {
	var files []io.Reader
	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Error: misplaced switch: '%s'\n", arg)
			os.Exit(1)
		}
		if arg == "demo" {
			demo := bytes.NewBuffer([]byte(demoMusic))
			files = append(files, demo)
			continue
		}
		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		files = append(files, os.Stdin)
	}
	for i, file := range files {
		reader := bufio.NewReader(file)
		if i > 0 {
			fmt.Println()
			time.Sleep(time.Second * 1)
		}
		initSoundDevice()
		go playMusicNotes(reader, volume)
		<-music.played
		flushSoundBuffer()
	}
	for _, file := range files {
		if file != os.Stdin {
			if closer, ok := file.(io.ReadCloser); ok {
				closer.Close()
			}
		}
	}
}

// Play a MIDI file
func parseMidiBeep(filename string) {
	midi, err := ParseMidi(filename, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	} else {
		midi.play()
	}
}

// Parses a MIDI file and print notes
func parseMidiNote(filename string) {
	_, err := ParseMidi(filename, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func playMusicNotesFromCL(musicNotes string, volume int) {
	reader := bufio.NewReader(strings.NewReader(musicNotes))
	initSoundDevice()
	go playMusicNotes(reader, volume)
	<-music.played
	flushSoundBuffer()
}
