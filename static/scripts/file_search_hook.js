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