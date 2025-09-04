const renameCallback = (item) => {
    const row = item.closest('tr');
    const fileName = row.getAttribute('data-name');
    const fileID = row.getAttribute('data-id');

    let newName = prompt('Enter new name for ' + fileName, fileName);
    if (!newName) {
        return;
    }

    htmx.ajax(
        'PATCH',
        '/api/v1/file-actions/' + fileID +
        "?desired_name=" + encodeURI(newName) +
        "&is_folder=" + row.classList.contains('search-folder') +
        "&orig_name=" + encodeURI(fileName),
        {
            target: "tr[data-id='" + row.getAttribute('data-id') + "'] td.file-name"
        }
    );
};

const shareCallback = (item) => {
    const fileId = item.parentElement.getAttribute('data-id');
    const fileType = item.parentElement.getAttribute('data-is-folder');

    fetch(`/api/v1/file-actions/${fileId}?is_folder=${fileType}`)
        .then(response => {
            (async () => {
                const hxTrigger = JSON.parse(response.headers.get('HX-Trigger'));
                if (!hxTrigger) {
                    return;
                }

                showAlert(hxTrigger['raise-alert']);
            })();

            return response.text();
        })
        .then(signedUrl => {
            copyWhenFocused(signedUrl);
        });
};

// document.getElementById('search-form').addEventListener('input', function () {
//     const inputValue = encodeURIComponent(document.getElementById('search-input').value); // URL-encode the input value
//     const apiUrl = `.?query=${inputValue}`; // Construct the URL with the query parameter

//     this.setAttribute('hx-get', apiUrl);
// });

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

function copyWhenFocused(text) {
    if (document.hasFocus()) {
        navigator.clipboard.writeText(text);
    } else {
        const onFocus = () => {
            navigator.clipboard.writeText(text);
            window.removeEventListener('focus', onFocus);
        };
        window.addEventListener('focus', onFocus);
    }
}
