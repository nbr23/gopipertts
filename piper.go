package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

func streamWavData(c *gin.Context, audioData io.Reader) {
	buffer := make([]byte, 4096)
	streaming := true

	done := make(chan bool)

	go func() {
		<-c.Request.Context().Done()
		done <- true
	}()

	for streaming {
		select {
		case <-done:
			streaming = false
			break
		default:
			n, err := audioData.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("error reading audio data: %v", err)
				}
				streaming = false
				break
			}

			if n > 0 {
				if _, err := c.Writer.Write(buffer[:n]); err != nil {
					log.Printf("error writing to client: %v", err)
					streaming = false
					break
				}
				c.Writer.Flush()
			}
		}
	}

}

func buildPiperCmd(voice string, speaker int) *exec.Cmd {
	cmdArgs := []string{
		PIPER_BINARY,
		"--model", fmt.Sprintf("%s/%s.onnx", VOICES_PATH, voice),
		"--config", fmt.Sprintf("%s/%s.onnx.json", VOICES_PATH, voice),
		"--json-input",
		"--output-raw",
	}
	if speaker > 0 {
		cmdArgs = append(cmdArgs, "--speaker-id", strconv.Itoa(speaker))
	}
	return exec.Command(PIPER_BINARY, cmdArgs...)
}

func writeInputToPiper(stdin io.WriteCloser, text string) error {
	ttsObj := map[string]interface{}{
		"text": text,
	}

	jsonStr, err := json.Marshal(ttsObj)
	if err != nil {
		return err
	}

	_, err = io.WriteString(stdin, string(jsonStr))
	if err != nil {
		return fmt.Errorf("failed writing to piper stdin: %v", err)
	}

	return nil
}

func buildFfmpegCmd(sampleRate int) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-f", "s16le",
		"-ar", strconv.Itoa(sampleRate),
		"-ac", "1",
		"-i", "pipe:0",
		"-f", "mp3",
		"-codec:a", "libmp3lame",
		"pipe:1",
	)
}

func streamTTSAsMp3(c *gin.Context, voice string, speaker int, text string, sampleRate int) error {
	if _, ok := DOWNLOADED_VOICES[voice]; !ok {
		return fmt.Errorf("voice not found: %s", voice)
	}
	if logInput {
		fmt.Println(strconv.Quote(text))
	}
	piperCmd := buildPiperCmd(voice, speaker)
	ffmpegCmd := buildFfmpegCmd(sampleRate)

	piperStdout, err := piperCmd.StdoutPipe()
	if err != nil {
		return err
	}
	ffmpegCmd.Stdin = piperStdout

	ffmpegStdout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return err
	}

	piperCmd.Stderr = os.Stderr
	ffmpegCmd.Stderr = os.Stderr

	piperStdin, err := piperCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
	}
	if err := piperCmd.Start(); err != nil {
		ffmpegCmd.Process.Kill()
		return fmt.Errorf("failed to start piper: %v", err)
	}

	if err := writeInputToPiper(piperStdin, text); err != nil {
		piperStdin.Close()
		piperCmd.Process.Kill()
		ffmpegCmd.Process.Kill()
		return err
	}
	piperStdin.Close()

	c.Header("Content-Type", "audio/mpeg")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	streamWavData(c, ffmpegStdout)

	piperCmd.Process.Kill()
	piperCmd.Wait()
	ffmpegCmd.Process.Kill()
	ffmpegCmd.Wait()
	return nil
}

func streamTTS(c *gin.Context, voice string, speaker int, text string, sampleRate int, channels int, bitsPerSample int) error {
	if _, ok := DOWNLOADED_VOICES[voice]; !ok {
		return fmt.Errorf("voice not found: %s", voice)
	}
	if logInput {
		fmt.Println(strconv.Quote(text))
	}
	cmd := buildPiperCmd(voice, speaker)
	log.Println("running piper command:", cmd)

	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start piper: %v", err)
	}

	if err := writeInputToPiper(stdin, text); err != nil {
		stdin.Close()
		log.Printf("error writing to piper: %v", err)
		return err
	}
	stdin.Close()

	streamWavData(c, stdout)

	cmd.Process.Kill()
	cmd.Wait()
	return nil
}
