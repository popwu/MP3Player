function settingsApp() {
    return {
        libraries: [],
        newLibraryPath: '',

        init() {
            this.fetchLibraries();
        },

        fetchLibraries() {
            fetch('/api/libraries')
                .then(response => response.json())
                .then(data => {
                    this.libraries = data;
                });
        },

        addLibrary() {
            if (this.newLibraryPath.trim() === '') return;

            fetch('/api/libraries', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `libraryPath=${encodeURIComponent(this.newLibraryPath)}`
            })
            .then(response => response.json())
            .then(() => {
                this.fetchLibraries();
                this.newLibraryPath = '';
            });
        }
    }
}