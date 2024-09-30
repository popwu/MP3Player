# This project is a music player written in golang, based on web RESTful API and web frontend.

Project name: mp3-player

### Web:
- The frontend is written using Tailwind CSS
- Uses Alpine.js for partial updates

### Backend:
- Uses the Gin framework

### Data Storage:
- Uses file storage for music library directory information
- Scans the music library directory on each startup
- Updates music library directory information whenever there's a change

### Features:
- Add music library directories, web retrieves path list via API for selection
- Index all music files in the music library directories
- Web retrieves all music files and displays the playlist
- Player functions: play music, pause, previous track, next track, shuffle play, loop play, volume control