# Quran Timings API

This is a Go program that fetches and saves timings data for each recitation in the Quran.

## Usage

To run the program, use the following command:

```bash
go run cmd/main.go
```

This will fetch and save timings data for all recitations in the Quran.

## Data

The data is stored in the `data` folder. Each recitation has its own folder, and each chapter has its own JSON file.

The JSON file contains an array of objects, where each object represents a verse. Each object has the following properties:

- `url`: The URL of the audio file for the verse.
- `verse`: The verse number.
- `chapter`: The chapter number.
- `segments`: An array of arrays, where each sub-array represents a segment of the audio file. Each sub-array contains integers representing the start and end times of the segment in milliseconds.
