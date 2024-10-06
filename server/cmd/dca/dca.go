// This code belongs to https://github.com/bwmarrin
//
// I included this code to modify it for my own use.
// It is copied from https://github.com/bwmarrin/dca/tree/master/cmd/dca
//
// Usage: ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca

package dca

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"layeh.com/gopus"
)

var (

	// AudioChannels sets the ops encoder channel value.
	// Must be set to 1 for mono, 2 for stereo
	AudioChannels int

	// AudioFrameRate sets the opus encoder Frame Rate value.
	// Must be one of 8000, 12000, 16000, 24000, or 48000.
	// Discord only uses 48000 currently.
	AudioFrameRate int

	// AudioBitrate sets the opus encoder bitrate (quality) value.
	// Must be within 500 to 512000 bits per second are meaningful.
	// Discord only uses 8000 to 128000 and default is 64000.
	AudioBitrate int

	// AudioApplication sets the opus encoder Application value.
	// Must be one of voip, audio, or lowdelay.
	// DCA defaults to audio which is ideal for music.
	// Not sure what Discord uses here, probably voip.
	AudioApplication string

	// AudioFrameSize sets the opus encoder frame size value.
	// The Frame Size is the length or amount of milliseconds each Opus frame
	// will be.
	// Must be one of 960 (20ms), 1920 (40ms), or 2880 (60ms)
	AudioFrameSize int

	// MaxBytes is a calculated value of the largest possible size that an
	// opus frame could be.
	MaxBytes int

	// OpusEncoder holds an instance of an gopus Encoder
	OpusEncoder *gopus.Encoder

	// EncodeChan is used for sending data to the encoder goroutine
	EncodeChan chan []int16
	// OutputChan is used for sending data to the writer goroutine
	OutputChan chan []byte

	// WaitGroup is used to wait untill all goroutines have finished.
	WaitGroup sync.WaitGroup
)

// reader reads from the input
func reader() {

	var err error

	defer func() {
		close(EncodeChan)
		WaitGroup.Done()
	}()

	// Create a 16KB input buffer
	stdin := bufio.NewReaderSize(os.Stdin, 16384)

	// Loop over the stdin input and pass the data to the encoder.
	for {

		buf := make([]int16, AudioFrameSize*AudioChannels)

		err = binary.Read(stdin, binary.LittleEndian, &buf)
		if err == io.EOF {
			// Okay! There's nothing left, time to quit.
			return
		}

		if err == io.ErrUnexpectedEOF {
			// Well there's just a tiny bit left, lets encode it, then quit.
			EncodeChan <- buf
			return
		}

		if err != nil {
			// Oh no, something went wrong!
			log.Println("error reading from stdin,", err)
			return
		}

		// write pcm data to the EncodeChan
		EncodeChan <- buf
	}

}

// encoder listens on the EncodeChan and encodes provided PCM16 data
// to opus, then sends the encoded data to the OutputChan
func encoder() {

	defer func() {
		close(OutputChan)
		WaitGroup.Done()
	}()

	for {
		pcm, ok := <-EncodeChan
		if !ok {
			// if chan closed, exit
			return
		}

		// try encoding pcm frame with Opus
		opus, err := OpusEncoder.Encode(pcm, AudioFrameSize, MaxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		// write opus data to OutputChan
		OutputChan <- opus
	}
}

// writer listens on the OutputChan and writes the output to stdout pipe
func writer() {

	defer WaitGroup.Done()

	var opuslen int16
	var err error

	// 16KB output buffer
	stdout := bufio.NewWriterSize(os.Stdout, 16384)
	defer func() {
		err := stdout.Flush()
		if err != nil {
			log.Println("error flushing stdout, ", err)
		}
	}()

	for {
		opus, ok := <-OutputChan
		if !ok {
			// if chan closed, exit
			return
		}

		// write header
		opuslen = int16(len(opus))
		err = binary.Write(stdout, binary.LittleEndian, &opuslen)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}

		// write opus data to stdout
		err = binary.Write(stdout, binary.LittleEndian, &opus)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}
	}
}

// Parse command line arguments and setup a couple of variables.
func init() {

	// Opus Encoding Options
	flag.IntVar(&AudioChannels, "ac", 2, "audio channels")
	flag.IntVar(&AudioFrameRate, "ar", 48000, "audio sampling rate")
	flag.IntVar(&AudioFrameSize, "as", 960, "audio frame size can be 960 (20ms), 1920 (40ms), or 2880 (60ms)")
	flag.IntVar(&AudioBitrate, "ab", 64, "audio encoding bitrate in kb/s can be 8 - 128")
	flag.StringVar(&AudioApplication, "aa", "audio", "audio application can be voip, audio, or lowdelay")

	flag.Parse()

	MaxBytes = (AudioFrameSize * AudioChannels) * 2 // max size of opus data
}

func Main() {

	//////////////////////////////////////////////////////////////////////////
	// Basic validation
	//////////////////////////////////////////////////////////////////////////

	// Make sure the stdin pipe is open
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
	} else {
		fmt.Println("error: stdin is not a pipe.")
		flag.Usage()
		return
	}

	//////////////////////////////////////////////////////////////////////////
	// Create channels, buffers, and encoder for use
	//////////////////////////////////////////////////////////////////////////

	// Create an Open Encoder to use
	OpusEncoder, err = gopus.NewEncoder(AudioFrameRate, AudioChannels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	// Set Opus Encoder Bitrate
	if AudioBitrate < 1 || AudioBitrate > 512 {
		AudioBitrate = 64 // Set to Discord default
	}
	OpusEncoder.SetBitrate(AudioBitrate * 1000)

	// Set Opus Encoder Application
	switch AudioApplication {
	case "voip":
		OpusEncoder.SetApplication(gopus.Voip)
	case "audio":
		OpusEncoder.SetApplication(gopus.Audio)
	case "lowdelay":
		OpusEncoder.SetApplication(gopus.RestrictedLowDelay)
	default:
		OpusEncoder.SetApplication(gopus.Audio)
	}

	// Create channels used by the reader/encoder/writer go routines
	EncodeChan = make(chan []int16, 10)
	OutputChan = make(chan []byte, 10)

	//////////////////////////////////////////////////////////////////////////
	// Start reader, encoder, and writer workers.  These add the DCA0 format
	// audio content to the file.
	//////////////////////////////////////////////////////////////////////////

	WaitGroup.Add(1)
	go reader()

	WaitGroup.Add(1)
	go encoder()

	WaitGroup.Add(1)
	go writer()

	// wait for above goroutines to finish, then exit.
	WaitGroup.Wait()
}
