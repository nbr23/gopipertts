package main

// WAV header structure (44 bytes for standard PCM WAV)
func generateWAVHeader(sampleRate, channels, bitsPerSample int) []byte {
	header := make([]byte, 44)

	// RIFF chunk
	copy(header[0:4], []byte("RIFF"))
	copy(header[4:7], []byte{0xFF, 0xFF, 0xFF, 0xFF}) // Size -1 for streaming
	copy(header[8:12], []byte("WAVE"))
	copy(header[12:16], []byte("fmt "))
	// fmt chunk size (16 for PCM)
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	// Audio format (1 for PCM)
	header[20] = 1
	header[21] = 0
	// Number of channels
	header[22] = byte(channels)
	header[23] = 0
	// Sample rate
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)
	// Byte rate
	byteRate := sampleRate * channels * bitsPerSample / 8
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)
	// Block align
	blockAlign := channels * bitsPerSample / 8
	header[32] = byte(blockAlign)
	header[33] = byte(blockAlign >> 8)
	// Bits per sample
	header[34] = byte(bitsPerSample)
	header[35] = byte(bitsPerSample >> 8)

	// data chunk
	copy(header[36:40], []byte("data"))
	// data size (unknown for streaming)
	copy(header[40:43], []byte{0xFF, 0xFF, 0xFF, 0xFF})
	return header
}
