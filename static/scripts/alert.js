document.addEventListener("raise-alert", (event) => {
    const { title, value: message, status, icon } = event.detail;
    showAlert({ title, message, status, icon });
});

function showAlert({ title, message, status, icon }) {
    const alertContainer = document.createElement("div");
    alertContainer.className = `alert alert-${status}`;
    alertContainer.setAttribute("hx-get", "/nothing")
    alertContainer.setAttribute("hx-trigger", "load delay:5s")
    alertContainer.setAttribute("hx-on::before-request", "this.classList.add('removing')")
    alertContainer.setAttribute("hx-on::after-request", "this.remove()")

    const header = document.createElement("div")
    header.classList.add("alert-header");

    const iconWrapper = document.createElement("span");
    iconWrapper.style.width = "40px";
    iconWrapper.style.height = "30px";
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

    htmx.process(modal);
}
