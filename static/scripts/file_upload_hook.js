const emptyFormBody = `
<div class="empty">
    <svg width="110" height="95" viewBox="0 0 22 19" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M5 12.0001L6.5 9.1001C6.66307 8.77624 6.91112 8.5028 7.21761 8.30904C7.5241 8.11529 7.8775 8.0085 8.24 8.0001H19M19 8.0001C19.3055 7.99956 19.6071 8.06903 19.8816 8.20318C20.1561 8.33733 20.3963 8.53258 20.5836 8.77396C20.7709 9.01534 20.9004 9.29643 20.9622 9.59565C21.024 9.89487 21.0164 10.2043 20.94 10.5001L19.4 16.5001C19.2886 16.9317 19.0362 17.3136 18.6829 17.5854C18.3296 17.8571 17.8957 18.0031 17.45 18.0001H3C2.46957 18.0001 1.96086 17.7894 1.58579 17.4143C1.21071 17.0392 1 16.5305 1 16.0001V3.0001C1 2.46966 1.21071 1.96096 1.58579 1.58588C1.96086 1.21081 2.46957 1.0001 3 1.0001H6.9C7.23449 0.996817 7.56445 1.07748 7.8597 1.23472C8.15495 1.39195 8.40604 1.62072 8.59 1.9001L9.4 3.1001C9.58211 3.37663 9.83003 3.60362 10.1215 3.7607C10.413 3.91778 10.7389 4.00004 11.07 4.0001H17C17.5304 4.0001 18.0391 4.21081 18.4142 4.58588C18.7893 4.96096 19 5.46966 19 6.0001V8.0001Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
    <p class="no-files">No file chosen</p>
</div>`

/**
 * @param {HTMLParagraphElement} fileName 
 * @param {string | null} path 
 * @returns 
 */
function renameFile(fileName, path = null) {
    let newFileName = prompt('Enter new file name:', fileName.textContent);
    if (!newFileName || newFileName.trim() === '') {
        return;
    }

    let oldFileName = fileName.textContent;

    if (path) {
        newFileName = path.split('/').slice(0, -1).join('/') + '/' + newFileName
        console.log(newFileName);
        oldFileName = path.split('/').slice(0, -1).join('/') + '/' + oldFileName;
        console.log(oldFileName);
    }

    fetch(`/api/v1/upload?new_name=${encodeURIComponent(newFileName)}&old_name=${encodeURIComponent(oldFileName)}`, {
        method: 'PATCH',
    }).then(response => {
        if (!response.ok) {
            alert('Failed to rename file. Please try again.');
            console.error('Error renaming file:', response.statusText);
            return;
        }
        return response.text();
    }).then(newFilePath => {
        newFilePath = newFilePath.split('/').pop();

        fileName.textContent = newFilePath;
        fileName.title = newFilePath;
    });
}

function removeFile(xhr, filePreview, fileName, relativePath = null) {
    const cancel = confirm(`Are you sure you want to remove the file "${fileName.textContent}"?`);

    if (!cancel) {
        return;
    }

    let filename = ""

    if (relativePath) {
        filename = relativePath.split('/').slice(0, -1).join('/') + '/' + fileName.textContent;
    } else {
        filename = fileName.textContent;
    }

    if (xhr) {
        xhr.abort();
    }

    fetch('/api/v1/upload?filename=' + encodeURIComponent(filename), {
        method: 'DELETE',
    }).then(response => {
        if (!response.ok) {
            console.error('Failed to remove file:', response.statusText);
        }
        filePreview.remove();

        const formBody = document.getElementById('form-body');

        if (formBody.children.length === 0) {
            formBody.innerHTML = emptyFormBody;
        }
    }).catch(error => {
        console.error('Error removing file:', error);
    });
}


/**
 * @param {HTMLDivElement} formBody 
 * @param {string[]} folders 
 */
function addFolderPreview(formBody, folders) {
    folders.forEach((folder, index) => {
        path = folders.slice(0, index + 1).join('/');
        if (document.getElementById(path)) {
            return;
        }

        const folderPreview = document.createElement('div');
        folderPreview.classList.add('folder-preview');
        folderPreview.classList.add('expanded');
        folderPreview.id = path;

        const folderName = document.createElement('p');
        folderName.textContent = folder;
        folderName.classList.add('folder-name');

        const chevronIcon = document.createElement('img');
        chevronIcon.src = '/static/assets/chevron-right.svg';
        chevronIcon.alt = 'Expand';
        chevronIcon.classList.add('chevron-icon');


        const folderHeader = document.createElement('div');
        folderHeader.classList.add('folder-header');

        folderHeader.onclick = () => {
            if (folderPreview.classList.contains('expanded')) {
                folderPreview.classList.remove('expanded');
            } else {
                folderPreview.classList.add('expanded');
            }
        };

        folderHeader.appendChild(folderName);
        folderHeader.appendChild(chevronIcon);

        const fileContainer = document.createElement('div');
        fileContainer.classList.add('file-container');
        fileContainer.id = path + '-files';

        folderPreview.appendChild(folderHeader);
        folderPreview.appendChild(fileContainer);

        const parentFolders = path.split('/').slice(0, -1);

        if (parentFolders.length > 0) {
            const parentFolder = parentFolders.join('/');
            const parentContainer = document.getElementById(parentFolder + '-files');
            if (parentContainer) {
                folderPreview.style.marginLeft = parentFolders.length * 20 + 'px';
                folderPreview.style.width = 'calc(100% - ' + ((parentFolders.length * 20) + 2) + 'px)';
                folderPreview.style.borderLeft = '1px solid #cecece';
                parentContainer.appendChild(folderPreview);
            }
        } else {
            formBody.appendChild(folderPreview);
        }
    });
}

function addFilePreviewWithProgress(file) {
    const formBody = document.getElementById('form-body');

    if (file.webkitRelativePath && file.webkitRelativePath !== '') {
        const folderPath = file.webkitRelativePath.split('/').slice(0, -1);
        addFolderPreview(formBody, folderPath);
    }

    const filePreview = document.createElement('div');
    filePreview.classList.add('file-preview');

    const fileName = document.createElement('p');
    fileName.textContent = file.name;
    fileName.classList.add('file-name');


    const editIcon = document.createElement('img');
    editIcon.src = '/static/assets/pencil.svg';
    editIcon.alt = 'Edit';
    editIcon.style.width = '20px';
    editIcon.style.height = '20px';
    editIcon.classList.add('edit-icon');
    editIcon.onclick = () => renameFile(fileName, file.webkitRelativePath || null);



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

    let associatedXHR = null;
    const bindXHR = (xhr) => {
        associatedXHR = xhr;
    }

    const removeButton = document.createElement('img');
    removeButton.src = '/static/assets/trash.svg';
    removeButton.alt = 'Remove';
    removeButton.style.width = '20px';
    removeButton.style.height = '20px';
    removeButton.classList.add('remove-button');
    removeButton.onclick = () => removeFile(associatedXHR, filePreview, fileName, file.webkitRelativePath || null);

    filePreview.appendChild(fileName);
    filePreview.appendChild(fileSize);
    filePreview.appendChild(progressBar);
    filePreview.appendChild(statusMessage);
    filePreview.appendChild(editIcon);
    filePreview.appendChild(removeButton);

    if (file.webkitRelativePath && file.webkitRelativePath !== '') {
        const fileContainer = file.webkitRelativePath.split('/').slice(0, -1).join('/') + '-files';

        document.getElementById(fileContainer).appendChild(filePreview);
    } else {
        formBody.appendChild(filePreview);
    }

    return { fileName, progressBar, statusMessage, bindXHR };
}

function requestQueue(tasks) {
    let currentIndex = 0;

    return new Promise((resolve) => {
        const handleResult = () => {
            currentIndex++;
            if (currentIndex < tasks.length) {
                startNext();
            } else {
                resolve();
            }
        }

        const startNext = () => {
            if (currentIndex < tasks.length) {
                tasks[currentIndex]().then(handleResult).catch(handleResult);
            } else if (currentIndex === tasks.length) {
                resolve();
            }
        }

        startNext();
    });
}

function sendXHRRequest(xhr, formData) {
    return new Promise(() => {
        xhr.send(formData);
    });
}

function buildXHRRequest(url, fileName, progressBar, statusMessage) {
    const xhr = new XMLHttpRequest();
    xhr.open('POST', url);

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
            console.log('File uploaded successfully:', xhr.response);

            fileName.textContent = xhr.response;

            resolve();
        } else {
            statusMessage.textContent = 'Upload failed. Please try again.';
            statusMessage.classList.add('error');
        }
    };

    xhr.onerror = function () {
        statusMessage.textContent = 'An error occurred during the upload.';
        statusMessage.classList.add('error');
    };

    return xhr;
}

function fileUploadHook() {
    const fileInput = this;
    const form = document.getElementById('upload-form');
    const formBody = document.getElementById('form-body');

    const isEmpty = formBody.getElementsByClassName('empty').length > 0;
    if (isEmpty) {
        formBody.innerHTML = ''; // Clear the empty message if it exists
    }

    // Create a queue of functions that return promises
    const xhrRequestQueue = Array.from(fileInput.files).map(file => {
        const { fileName, progressBar, statusMessage, bindXHR } = addFilePreviewWithProgress(file);

        const xhr = buildXHRRequest(form.action, fileName, progressBar, statusMessage);

        bindXHR(xhr);

        const formData = new FormData();
        formData.append('file', file);

        if (fileInput.id === 'folder') {
            formData.append('folder', file.webkitRelativePath); // Include folder path if needed
        }

        return () => sendXHRRequest(xhr, formData);
    });

    requestQueue(xhrRequestQueue)
}

document.getElementById('file').addEventListener('change', (event) => {
    fileUploadHook.call(event.target);
});

document.getElementById('folder').addEventListener('change', (event) => {
    fileUploadHook.call(event.target);
});