class FormField extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        const fieldset = document.createElement('fieldset');
        const legend = document.createElement('legend');
        const input = document.createElement('input');

        const placeholder = this.getAttribute('placeholder') || '';
        const name = this.getAttribute('name') || '';
        const type = this.getAttribute('type') || 'text';

        if (this.hasAttribute('required')) {
            input.setAttribute('required', '');
        }
        input.setAttribute('placeholder', placeholder);
        input.setAttribute('name', name);
        input.setAttribute('type', type);

        fieldset.appendChild(input);

        legend.textContent = this.getAttribute('legend') || 'Field';

        fieldset.appendChild(legend);

        this.appendChild(fieldset);
    }
}

customElements.define('form-field', FormField);