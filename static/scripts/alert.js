document.addEventListener("raise-alert", (event) => {
    console.log(event.detail);
    showAlert(event.detail);
});

function showAlert({ title, message, status, icon }) {
    const alertContainer = document.createElement("div");
    alertContainer.className = `alert alert-${status}`;

    const header = document.createElement("div")
    header.classList.add("alert-header");

    const iconWrapper = document.createElement("span");
    iconWrapper.style.width = "40px";
    iconWrapper.style.height = icon === 'x' ? "30px" : "40px";
    header.appendChild(iconWrapper);

    fetch("/static/assets/" + icon + ".svg")
        .then(resp => resp.ok ? resp.text() : "")
        .then(svgText => {
            if (svgText) {
                iconWrapper.innerHTML = svgText;
            }
        })
        .catch((error) => {
            console.error("Error loading icon", error);
        });


    const titleElement = document.createElement("h4");
    titleElement.innerText = title;
    header.appendChild(titleElement);

    alertContainer.appendChild(header);

    const messageContainer = document.createElement("p");
    messageContainer.innerText = message;
    alertContainer.appendChild(messageContainer);

    const modal = document.getElementById("modal-content")

    modal.appendChild(alertContainer);

    setTimeout(() => alertContainer.classList.add("removing"), 5000);
    setTimeout(() => alertContainer.remove(), 6000);
}
