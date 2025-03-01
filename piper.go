package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		"/usr/share/piper/piper",
		"--model", fmt.Sprintf("%s/%s.onnx", VOICES_PATH, voice),
		"--config", fmt.Sprintf("%s/%s.onnx.json", VOICES_PATH, voice),
		"--json-input",
		"--output-raw",
	}
	if speaker > 0 {
		cmdArgs = append(cmdArgs, "--speaker-id", strconv.Itoa(speaker))
	}
	return exec.Command("/usr/share/piper/piper", cmdArgs...)
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

func streamTTS(c *gin.Context, voice string, speaker int, text string, sampleRate int, channels int, bitsPerSample int) error {
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
