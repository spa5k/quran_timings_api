# Quran Reciters Data

This repository contains data about various Quran reciters and their recitations. The data is organized in JSON format and includes detailed information about each reciter, their styles, and the audio files for each chapter of the Quran.

## Reciters Data Format

The `reciters.json` file contains an array of reciters, each with detailed information. Below is the structure of the data:

### Reciter Object

- **id**: Unique identifier for the reciter.
- **reciter_id**: Another unique identifier for the reciter.
- **name**: Name of the reciter.
- **translated_name**: Object containing the translated name and language.
  - **name**: Translated name of the reciter.
  - **language_name**: Language of the translated name.
- **style**: Object containing the style of recitation.
  - **name**: Name of the style.
  - **language_name**: Language of the style name.
  - **description**: Description of the style.
- **qirat**: Object containing the qirat information.
  - **name**: Name of the qirat.
  - **language_name**: Language of the qirat name.
- **slug**: Unique slug for the reciter

### Example Reciter Object

```json
{
	"id": 1,
	"reciter_id": 1,
	"name": "AbdulBaset AbdulSamad",
	"translated_name": {
		"name": "AbdulBaset AbdulSamad",
		"language_name": "english"
	},
	"style": {
		"name": "Mujawwad",
		"language_name": "english",
		"description": "Mujawwad is a melodic style of Holy Quran recitation."
	},
	"qirat": {
		"name": "Hafs",
		"language_name": "english"
	},
	"slug": "abdul_baset_mujawwad"
}
```

## Fetching and Saving Chapter Data

The program fetches chapter data for each reciter from an API and saves it in the format `data/<style>/<slug>/chapter_number.json`. The chapter data includes audio files and verse timings.

### Chapter Data Format

- **audio_files**: Array of audio files for the chapter.
  - **id**: Unique identifier for the audio file.
  - **chapter_id**: Identifier for the chapter.
  - **file_size**: Size of the audio file.
  - **format**: Format of the audio file (e.g., mp3).
  - **audio_url**: URL to download the audio file.
  - **duration**: Duration of the audio file.
  - **verse_timings**: Array of verse timings.
    - **verse_key**: Key for the verse (e.g., "1:1").
    - **timestamp_from**: Start timestamp of the verse.
    - **timestamp_to**: End timestamp of the verse.
    - **duration**: Duration of the verse.
    - **segments**: Array of segments within the verse.

### Example Chapter Data

```json
{
	"audio_files": [
		{
			"id": 8149,
			"chapter_id": 1,
			"file_size": 2848896.0,
			"format": "mp3",
			"audio_url": "https://download.quranicaudio.com/qdc/siddiq_al-minshawi/mujawwad/001.mp3",
			"duration": 130952,
			"verse_timings": [
				{
					"verse_key": "1:1",
					"timestamp_from": 0,
					"timestamp_to": 13320,
					"duration": -13320,
					"segments": [
						[1, 0.0, 4840.0],
						[2, 4840.0, 5370.0],
						[3, 5370.0, 6880.0],
						[4, 6880.0, 12250.0]
					]
				},
				{
					"verse_key": "1:2",
					"timestamp_from": 13320,
					"timestamp_to": 26380,
					"duration": -13060,
					"segments": [
						[1, 13320.0, 17870.0],
						[2, 17870.0, 18240.0],
						[3, 18240.0, 21390.0],
						[4, 21390.0, 25210.0]
					]
				}
			]
		}
	]
}
```

## Usage

### Fetching Reciters Data

To fetch the reciters data and save it to `reciters.json`, run the following function:

```go
FetchQuranTimingReciters()
```

### Fetching and Saving Chapter Data

To fetch and save chapter data for each reciter, run the following function:

```go
AyahTimingsPerReciter()
```

This will create directories based on the reciter's style and slug, and save each chapter's data in the corresponding directory.

## Contributing

If you would like to contribute to this project, please fork the repository and submit a pull request. We welcome all contributions!

## References

- [Quran.com](https://quran.com/)
- [Quran.com API](https://api-docs.quran.com/)
- [Every Ayah](https://everyayah.com/)
- [Verse-by-verse timings](http://versebyversequran.com)