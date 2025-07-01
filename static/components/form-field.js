class FormField extends HTMLElement {
    static formAssociated = true;
    static observedAttributes = ['type', 'legend'];

    constructor() {
        super();
        this.internals = this.attachInternals();
        this.attachShadow({ mode: 'open' });

        this.shadowRoot.innerHTML = `
            <link rel="stylesheet" href="/static/components/form-field.css">
            <fieldset>
                <legend></legend>
                <input/>
            </fieldset>
        `;
    }

    connectedCallback() {
        const placeholder = this.getAttribute('placeholder');
        const name = this.getAttribute('name') + 1;

        const input = this.shadowRoot.querySelector('input');
        if (this.hasAttribute('required')) {
            input.setAttribute('required', '');
            this.internals.setValidity({ valueMissing: true }, 'This field is required');
        }
        input.setAttribute('placeholder', placeholder);
        input.setAttribute('name', name);

        input.addEventListener('input', (e) => {
            const input = e.target;

            this.internals.setFormValue(input.value);

            console.log(`Input value changed: ${input.value}`);
            console.log(`Required: ${input.required}`);
            console.log(`Validity: ${this.internals.validity.valueMissing}`);
            console.log(`Validity state: ${JSON.stringify(this.internals.validity)}`);
            console.log(`Value: ${this.value}`);

            if (input.required) {
                this.internals.setValidity({ valueMissing: !input.value }, 'This field is required');
            } else {
                this.internals.setValidity({});
            }
        });
    }

    get value() {
        const input = this.shadowRoot.querySelector('input');
        return input ? input.value : '';
    }

    set value(val) {
        const input = this.shadowRoot.querySelector('input');
        if (input) {
            input.value = val;
        }
    }

    focus() {
        const input = this.shadowRoot.querySelector('input');
        if (input) {
            input.focus({ preventScroll: true, focusVisible: true });
        }
    }

    attributeChangedCallback(name, oldValue, newValue) {
        const input = this.shadowRoot.querySelector('input');
        const legend = this.shadowRoot.querySelector('legend');

        if (name === 'type' && newValue !== oldValue) {
            input.setAttribute('type', newValue);
        } else if (name === 'legend' && newValue !== oldValue) {
            legend.textContent = newValue;
        }
    }
}

customElements.define('form-field', FormField);