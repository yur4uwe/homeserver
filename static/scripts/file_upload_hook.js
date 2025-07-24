function addFilePreviewWithProgress(file) {
    const filePreview = document.createElement('div');
    filePreview.classList.add('file-preview');

    const fileName = document.createElement('p');
    fileName.textContent = file.name;
    fileName.classList.add('file-name');

    const fileSize = document.createElement('p');
    fileSize.classList.add('file-size');
    const sizeUnits = ['bytes', 'KB', 'MB', 'GB'];

    let i = 0;
    while (file.size >= Math.pow(1024, i + 1) && i < sizeUnits.length - 1) {
        i++;
    }

    fileSize.textContent = `Size: ${(file.size / Math.pow(1024, i)).toFixed(2)} ${sizeUnits[i]}`;

    const progressBar = document.createElement('progress');
    progressBar.value = 0;
    progressBar.max = 100;
    progressBar.classList.add('file-progress-bar');

    const statusMessage = document.createElement('p');
    statusMessage.textContent = 'Waiting to upload...';
    statusMessage.classList.add('status-message');

    const removeButton = document.createElement('img');
    removeButton.src = '/static/assets/trash.svg';
    removeButton.alt = 'Remove';
    removeButton.style.width = '20px';
    removeButton.style.height = '20px';
    removeButton.classList.add('remove-button');
    removeButton.onclick = function () {
        fetch('/api/v1/upload?fileName=' + encodeURIComponent(file.name), {
            method: 'DELETE',
        })
            .then(response => {
                if (response.ok) {
                    filePreview.remove();
                } else {
                    console.error('Failed to remove file:', response.statusText);
                }
            })
            .catch(error => {
                console.error('Error removing file:', error);
            });
    };

    filePreview.appendChild(fileName);
    filePreview.appendChild(fileSize);
    filePreview.appendChild(progressBar);
    filePreview.appendChild(statusMessage);
    filePreview.appendChild(removeButton);

    const formBody = document.getElementById('form-body');

    const isEmpty = formBody.getElementsByClassName('empty').length > 0;
    if (isEmpty) {
        formBody.innerHTML = ''; // Clear the empty message if it exists
    }

    formBody.appendChild(filePreview);

    return { progressBar, statusMessage };
}

document.getElementById('file').addEventListener('change', function () {
    const fileInput = this;
    const form = document.getElementById('upload-form');

    // Iterate over each selected file
    Array.from(fileInput.files).forEach(file => {
        // Add file preview with progress bar
        const { progressBar, statusMessage } = addFilePreviewWithProgress(file);

        // Prepare FormData for the current file
        const formData = new FormData();
        formData.append('file', file);

        // Perform the upload using XMLHttpRequest
        const xhr = new XMLHttpRequest();
        xhr.open('POST', form.action, true);

        // Update progress bar for the current file
        xhr.upload.addEventListener('progress', function (event) {
            if (event.lengthComputable) {
                const percentComplete = (event.loaded / event.total) * 100;
                progressBar.value = percentComplete;
            }
        });

        // Handle upload completion for the current file
        xhr.onload = function () {
            if (xhr.status === 200) {
                statusMessage.textContent = 'Upload complete!';
                statusMessage.classList.add('success');
                progressBar.value = 100;
            } else {
                statusMessage.textContent = 'Upload failed. Please try again.';
                statusMessage.classList.add('error');
            }
        };

        // Handle upload errors for the current file
        xhr.onerror = function () {
            statusMessage.textContent = 'An error occurred during the upload.';
            statusMessage.classList.add('error');
        };

        // Send the form data for the current file
        xhr.send(formData);
    });
});