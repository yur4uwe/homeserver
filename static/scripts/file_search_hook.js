document.getElementById('search-input').addEventListener('input', function () {
    const inputValue = encodeURIComponent(this.value); // URL-encode the input value
    const apiUrl = `/search/?query=${inputValue}`; // Construct the URL with the query parameter

    this.setAttribute('hx-get', apiUrl);
    this.value = inputValue; // Update the input value to the encoded one

    document.querySelectorAll('.search-result.search-folder').forEach(element => {
        const encodedFolderPath = element.getAttribute('data-folder-path');
        element.setAttribute('hx-get', `${apiUrl}&directory=${encodedFolderPath}`);
    });
});

document.addEventListener('DOMContentLoaded', function () {
    const searchResults = document.querySelectorAll('.search-result');

    searchResults.forEach(result => {
        // Toggle selection on click
        result.addEventListener('click', function () {
            result.classList.toggle('selected');
            result.draggable = document.querySelectorAll('.search-result.selected').length > 0; // Enable dragging only if something is selected
        });

        // Handle dragstart for selected files
        result.addEventListener('dragstart', function (event) {
            const selectedResults = document.querySelectorAll('.search-result.selected');
            const selectedPaths = Array.from(selectedResults).map(selected => {
                const path = selected.getAttribute('data-path');
                const name = selected.getAttribute('data-name');

                let fullPath = path;
                if (!fullPath) {
                    fullPath = name;
                } else {
                    fullPath = `${fullPath}/${name}`;
                }

                return fullPath;
            });

            // Store all selected paths in the drag event
            event.dataTransfer.setData('application/json', JSON.stringify(selectedPaths));

            // Add dragging class to all selected items
            selectedResults.forEach(selected => {
                selected.classList.add('dragging');
            });
        });

        // Remove dragging class on dragend
        result.addEventListener('dragend', function () {
            const selectedResults = document.querySelectorAll('.search-result.dragging');
            selectedResults.forEach(selected => {
                selected.classList.remove('dragging');
            });
        });
    });

    const dropZone = document.querySelector('#search-results');

    // Handle dragover event
    dropZone.addEventListener('dragover', function (event) {
        event.preventDefault(); // Allow dropping
        dropZone.classList.add('drag-over');
    });

    // Handle dragleave event
    dropZone.addEventListener('dragleave', function () {
        dropZone.classList.remove('drag-over');
    });

    // Handle drop event
    dropZone.addEventListener('drop', function (event) {
        event.preventDefault();
        dropZone.classList.remove('drag-over');

        // Retrieve all dropped paths
        const droppedPaths = JSON.parse(event.dataTransfer.getData('application/json'));
        console.log('Dropped paths:', droppedPaths);

        // Perform your action here (e.g., move files, update UI)
    });
});
