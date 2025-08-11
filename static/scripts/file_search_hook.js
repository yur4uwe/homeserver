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

        // Check if the click happened outside the menu content
        if (!menuContent.contains(event.target) && event.target !== menu) {
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
            const newName = encodeURI(prompt('Enter new name for ' + fileName, fileName) + '.' + fileExtension);
            if (newName) {
                console.log('Renaming ' + fileName + ' to ' + newName);
                const filePathCell = row.querySelector('td.file-path');
                let filePath = filePathCell ? filePathCell.textContent.trim() : '';
                if (filePath) {
                    filePath += "/"
                    filePath = encodeURI(filePath);
                }

                const origFileName = encodeURI(fileName + '.' + fileExtension);

                htmx.ajax('PATCH', '/api/v1/file-actions/' + filePath + origFileName + "?desired_name=" + newName);
            }
        });
    });
});
