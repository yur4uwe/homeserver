document.getElementById('search-form').addEventListener('input', function () {
    const inputValue = encodeURIComponent(document.getElementById('search-input').value); // URL-encode the input value
    const apiUrl = `.?query=${inputValue}`; // Construct the URL with the query parameter

    this.setAttribute('hx-get', apiUrl);
});

document.addEventListener('click', function (event) {
    // Find all open action menus
    const openMenus = document.querySelectorAll('.actions-menu-toggle:checked');

    openMenus.forEach(menu => {
        const menuContent = menu.nextElementSibling; // The .actions-menu-content element
        const menuLabel = menu.previousElementSibling; // The label that toggles the menu

        // Check if the click happened outside the menu content
        if (!menuContent.contains(event.target) && !menuLabel.contains(event.target) && event.target !== menu) {
            menu.checked = false; // Uncheck the checkbox to close the menu
        }
    });
});

document.addEventListener('keydown', function (event) {
    if (event.key === 'Escape') {
        const openMenus = document.querySelectorAll('.actions-menu-toggle:checked');
        openMenus.forEach(menu => {
            menu.checked = false; // Uncheck all open menus
        });
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const searchResults = document.querySelectorAll('.search-result');

    searchResults.forEach(result => {
        result.addEventListener('click', function () {
            result.classList.toggle('selected');
            result.draggable = document.querySelectorAll('.search-result.selected').length > 0; // Enable dragging only if something is selected
        });
    });

    const renameAction = document.querySelectorAll('.rename-action');

    renameAction.forEach(item => {
        item.addEventListener('click', function () {
            const row = item.closest('tr');
            const fileNameCell = row.querySelector('.file-name');
            if (!fileNameCell) {
                console.error('No .file-name cell found in row');
                return;
            }
            const fileSig = fileNameCell.textContent.trim().split('.');
            const fileName = fileSig.shift();
            const fileExtension = fileSig.pop();
            let newName = prompt('Enter new name for ' + fileName, fileName);
            if (newName) {
                newName = encodeURI(newName + '.' + fileExtension);
                console.log('Renaming ' + fileName + ' to ' + newName);
                const filePathData = row.getAttribute('data-path');
                let filePath = filePathData ? filePathData.trim() : '';
                if (filePath) {
                    filePath = encodeURI(filePath);
                }

                const origFileName = encodeURI(fileName + '.' + fileExtension);

                htmx.ajax('PATCH', '/api/v1/file-actions/' + filePath + origFileName + "?desired_name=" + newName, {
                    target: "tr[data-id='" + row.getAttribute('data-id') + "'] td.file-name"
                });
            }
        });
    });

    const shareActions = document.querySelectorAll('.share-action');

    for (const item of shareActions) {
        const shareCallback = () => {
            const fileId = item.getAttribute('data-id');
            const fileType = item.getAttribute('data-is-folder');

            fetch(`/api/v1/file-actions/share?is_folder=${fileType}&id=${fileId}`)
                .then(response => {
                    (async () => {
                        const hxTrigger = JSON.parse(response.headers.get('HX-Trigger'));
                        if (!hxTrigger) {
                            return;
                        }

                        const evt = new CustomEvent('raise-alert', { detail: hxTrigger['raise-alert'] });

                        document.dispatchEvent(evt);
                    })();

                    return response.text();
                })
                .then(signedUrl => {

                    console.log('File shared successfully:', signedUrl);
                    navigator.clipboard.writeText(signedUrl);
                });
        }

        item.addEventListener('click', shareCallback);
    }
});
