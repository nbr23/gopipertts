# gopipertts

A small HTTP API wrapper for piper's texttospeech

## Usage

The provided docker image embeds piper for easier use. Just run:
```bash
docker run \
    -v ./voices:/voices \ # This is where the voices will be saved to. Mounted for persistence
    -p 8080:8080 \
    nbr23/gopipertts
```

You can then navigate to `http://localhost:8080/` to access the webui and start speaking text.

## API Reference

### List voices

`GET /api/voices` will return a json list of voices available for download and usage

### Process text into speech

`/api/tts` will convert the text passed into a wav file.
This endpoint accepts POST and GET requests.

POST requests expect a json body like the following:
```json
{
    "text": "Hello World",
    "speed": 1.0,
    "voice": "en_US-amy-low",
    "speaker": ""               // only available for select voices
}
```

GET requests expect the parameters `text` and optionally `speed`, `voice` and `speaker` to be passed as url query parameters.

Some usage examples:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"text": "happy text to speaching!"}' 'http://localhost:8080/api/tts' | mpv -
```

```bash
curl -X POST -H "Content-Type: application/json" -d '{"text": "happy text to speaching!", "voice":"en_US-amy-low", "speed": 1.2}' 'http://localhost:8080/api/tts' | mpv -
```

```bash
curl -X POST -H "Content-Type: application/json" -d '{"text": "happy text to speaching!"}' 'http://localhost:8080/api/tts?speed=1.2&voice=en_US-amy-low' | mpv -
```

```bash
curl 'http://localhost:8080/api/tts?speed=1.1&voice=en_US-amy-low&text=Happy%20text%20to%20speeching' | mpv -
```

Leverages [piper](https://github.com/rhasspy/piper) for TTS and voices from [rhasspy/piper-voices](https://huggingface.co/rhasspy/piper-voices/tree/main)