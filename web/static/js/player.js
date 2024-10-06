function playerApp() {
    return {
        songs: [],
        currentSongIndex: -1,
        isPlaying: false,
        isShuffled: false,
        isLooped: false,
        playHistory: [],
        historyIndex: -1,

        init() {
            this.fetchSongs();
        },

        fetchSongs() {
            fetch('/api/songs')
                .then(response => response.json())
                .then(data => {
                    this.songs = data;
                });
        },

        playSong(song, addToHistory = true) {
            const audioPlayer = this.$refs.audioPlayer;
            audioPlayer.src = `/api/stream/${encodeURIComponent(song.path)}`;
            audioPlayer.play();
            this.isPlaying = true;
            this.currentSongIndex = this.songs.findIndex(s => s.path === song.path);
            
            if (addToHistory) {
                this.playHistory = this.playHistory.slice(0, this.historyIndex + 1);
                this.playHistory.push(this.currentSongIndex);
                this.historyIndex = this.playHistory.length - 1;
            }
        },

        prevSong() {
            if (this.isShuffled) {
                if (this.historyIndex > 0) {
                    this.historyIndex--;
                    this.playSong(this.songs[this.playHistory[this.historyIndex]], false);
                } else {
                    this.playRandomSong();
                }
            } else {
                this.currentSongIndex = (this.currentSongIndex - 1 + this.songs.length) % this.songs.length;
                this.playSong(this.songs[this.currentSongIndex]);
            }
        },

        nextSong() {
            if (this.isShuffled) {
                this.playRandomSong();
            } else {
                this.currentSongIndex = (this.currentSongIndex + 1) % this.songs.length;
                this.playSong(this.songs[this.currentSongIndex]);
            }
        },

        playRandomSong() {
            const availableIndices = this.songs.map((_, index) => index).filter(index => index !== this.currentSongIndex);
            const randomIndex = availableIndices[Math.floor(Math.random() * availableIndices.length)];
            this.playSong(this.songs[randomIndex]);
        },

        toggleShuffle() {
            this.isShuffled = !this.isShuffled;
        },

        toggleLoop() {
            this.isLooped = !this.isLooped;
        },

        togglePlay() {
            if (this.currentSongIndex === -1 && this.songs.length > 0) {
                this.playSong(this.songs[0]);
            } else if (this.currentSongIndex !== -1) {
                this.isPlaying = !this.isPlaying;
                if (this.isPlaying) {
                    this.$refs.audioPlayer.play();
                } else {
                    this.$refs.audioPlayer.pause();
                }
            }
        },

        onSongEnded() {
            if (this.isLooped) {
                this.$refs.audioPlayer.play();
            } else if (this.isShuffled) {
                this.playRandomSong();
            } else {
                this.nextSong();
            }
        },

        playButtonClicked() {
            if (this.currentSongIndex === -1 && this.songs.length > 0) {
                this.playSong(this.songs[0]);
            } else {
                this.togglePlay();
            }
        }
    }
}