package main

const INDEX_HTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TTS Voice Selection</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 20px auto;
            padding: 0 20px;
        }
        .form-group {
            margin-bottom: 20px;
            display: flex;
            align-items: center;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        select, textarea {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        textarea {
            min-height: 150px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #45a049;
        }
        .copy-button {
            background-color: #008CBA;
            margin-left: 10px;
            flex-shrink: 0;
        }
        .copy-button:hover {
            background-color: #007BB5;
        }
    </style>
</head>
<body>
    <h1>Text-to-Speech</h1>
    <form id="ttsForm">
        <label for="textInput">Text to speak:</label>
        <div class="form-group">
            <textarea id="textInput"></textarea>
        </div>
        
        <label for="languageSelect">Language:</label>
        <div class="form-group">
            <select id="languageSelect"></select>
        </div>
        
        <label for="voiceSelect">Voice:</label>
        <div class="form-group">
            <select id="voiceSelect"></select>
            <button type="button" class="copy-button" onclick="copyVoice()">Copy</button>
        </div>
        
        <label for="speakerSelect">Speaker:</label>
        <div class="form-group">
            <select id="speakerSelect"></select>
            <button type="button" class="copy-button" onclick="copySpeaker()">Copy</button>
        </div>
        
        <button type="submit">Speak!</button>
        
        <div class="form-group" style="margin-top: 20px;">
            <audio id="audioPlayer" controls style="width: 100%; display: none;"></audio>
        </div>
    </form>

    <script>
        let voices = [];
        
        async function fetchVoices() {
            try {
                const response = await fetch(` + "`${window.location.pathname}api/voices`" + `);
                voices = Object.values(await response.json());
                initializeForm();
            } catch (error) {
                console.error('Error fetching voices:', error);
            }
        }

        function initializeForm() {
            const languages = [...new Set(voices.flatMap(voice => voice.language.code))];
            const languageSelect = document.getElementById('languageSelect');
            languages.forEach(lang => {
                const option = new Option(lang, lang);
                languageSelect.add(option);
            });
            languageSelect.value = 'en_US';
            updateVoiceSelect();
        }

        function updateVoiceSelect() {
            const selectedLanguage = document.getElementById('languageSelect').value;
            const voiceSelect = document.getElementById('voiceSelect');
            
            voiceSelect.innerHTML = '';

            const filteredVoices = voices.filter(voice => 
                voice.language.code === selectedLanguage
            );
            filteredVoices.forEach(voice => {
                const option = new Option(voice.key, voice.key);
                voiceSelect.add(option);
            });
            updateSpeakerSelect();
        }

        function updateSpeakerSelect() {
            const selectedVoice = document.getElementById('voiceSelect').value;
            const speakerSelect = document.getElementById('speakerSelect');
            speakerSelect.innerHTML = '';
            const voice = voices.find(v => v.key === selectedVoice);
            if (voice && voice.speaker_id_map) {
				console.log(voice.speaker_id_map);
                Object.keys(voice.speaker_id_map).forEach(speaker => {
				console.log(speaker);
                    const option = new Option(speaker, speaker);
                    speakerSelect.add(option);
                });
                speakerSelect.disabled = false;
            } else {
                speakerSelect.add(new Option('No speakers available', ''));
                speakerSelect.disabled = true;
            }
        }

        function copyVoice() {
            const selectedVoice = document.getElementById('voiceSelect').value;
            navigator.clipboard.writeText(selectedVoice).catch(err => {
                console.error('Error copying voice:', err);
            });
        }

        function copySpeaker() {
            const selectedSpeaker = document.getElementById('speakerSelect').value;
            navigator.clipboard.writeText(selectedSpeaker).catch(err => {
                console.error('Error copying speaker:', err);
            });
        }

        document.getElementById('languageSelect').addEventListener('change', updateVoiceSelect);
        document.getElementById('voiceSelect').addEventListener('change', updateSpeakerSelect);

        document.getElementById('ttsForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const formData = {
                text: document.getElementById('textInput').value,
                language: document.getElementById('languageSelect').value,
                voice: document.getElementById('voiceSelect').value,
                speaker: document.getElementById('speakerSelect').value || null
            };

            try {
                const response = await fetch(` + "`${window.location.pathname}api/tts`" + `, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) {
                    throw new Error('Synthesis failed');
                }

                const audioBlob = await response.blob();
                const audioUrl = URL.createObjectURL(audioBlob);
                const audioElement = document.getElementById('audioPlayer');
                audioElement.style.display = 'block';
                audioElement.src = audioUrl;
                await audioElement.play();
                audioElement.onended = () => {
                    URL.revokeObjectURL(audioUrl);
                };
            } catch (error) {
                console.error('Error during synthesis request:', error);
            }
        });
        fetchVoices();
    </script>
</body>
</html>
`
