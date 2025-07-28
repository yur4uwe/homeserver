document.getElementById('search-input').addEventListener('input', function () {
    const inputValue = encodeURIComponent(this.value); // URL-encode the input value
    const apiUrl = `/search/?query=${inputValue}`; // Construct the URL with the query parameter

    this.setAttribute('hx-get', apiUrl);
    this.value = inputValue; // Update the input value to the encoded one

    document.querySelectorAll('.search-result.search-folder').forEach(element => {
        const encodedFolderPath = encodeURIComponent(element.getAttribute('data-folder-path'));
        element.setAttribute('hx-get', `${apiUrl}&directory=${encodedFolderPath}`);
    });
});

document.addEventListener('DOMContentLoaded', function () {
    /**
     * @type {Array<HTMLTableRowElement>} searchResults
     */
    const searchResults = document.querySelectorAll('.search-result');

    searchResults.forEach(result => {
        result.addEventListener('click', function () {
            result.classList.toggle('selected');
            result.draggable = !result.draggable;
        });

        result.addEventListener('dragstart', function (event) {
            const selectedResults = document.querySelectorAll('.search-result.selected');

            result.classList.add('dragging');
            selectedResults.forEach(selected => {
                selected.classList.add('dragging');
            });
        });

        result.addEventListener('dragend', function () {
            result.classList.remove('dragging');
        });
    });

    const dropZone = document.querySelector('#search-results');

    dropZone.addEventListener('dragover', function (event) {
        event.preventDefault();
        dropZone.classList.add('drag-over');
    });

    dropZone.addEventListener('dragleave', function () {
        dropZone.classList.remove('drag-over');
    });

    dropZone.addEventListener('drop', function (event) {
        event.preventDefault();
        dropZone.classList.remove('drag-over');

        const folderPath = event.dataTransfer.getData('text/plain');
        console.log('Dropped folder path:', folderPath);
    });
});
