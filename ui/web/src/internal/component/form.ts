import jsxmm from "../../module/jsxmm/element.ts"

export function Form(...children: HTMLElement[]) {
    const form = jsxmm.Element("form", { className: "form" }, ...children)
    form.addEventListener("submit", prevent_default)

    return form
}

export function FormField(label: string, name: string, input: HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement) {
    input.classList.add("value", "editable")
    input.name = name

    return jsxmm.Element("div", { className: "field" },
        jsxmm.Element("label", { className: "label", textContent: label }),
        input,
    )
}

export function FormInfo<T extends string[]>(form: HTMLFormElement, ...fields: T): T {
    const data = new FormData(form)
    const values = new Array<string>(fields.length)

    for (let i = 0; i < fields.length; i++) {
        let value = data.get(fields[i]) as string | null
        if (value === null) {
            value = ""
        }

        values[i] = value
    }

    return values as T
}

export function TextInput(name: string, value: string = ""): HTMLInputElement {
    return jsxmm.Element("input", {
        type: "text", name: name, value: value,
        autocomplete: "off",
        required: true,
    })
}

export function SIAPEInput(name: string, value: string = ""): HTMLInputElement {
    return jsxmm.Element("input", {
        type: "text", name: name, value: value,
        pattern: "^[0-9]{7}$",
        autocomplete: "username",
        required: true,
        className: "monospace",
    })
}

export function EmailInput(name: string, value: string = ""): HTMLInputElement {
    return jsxmm.Element("input", {
        type: "email", name: name, value: value,
        pattern: "^[0-9A-Za-z_%+\\-]+(\\.[0-9A-Za-z_%+\\-]+)*@[0-9A-Za-z\\-]+(\\.[0-9A-Za-zA-Z\\-]+)*\\.[A-Za-z]{2,}$",
        autocomplete: "email",
        required: true,
    })
}

export function PasswordInput(name: string): HTMLInputElement {
    return jsxmm.Element("input", {
        type: "password", name: name,
        autocomplete: "current-password",
        required: true,
    })
}

function prevent_default(evt: Event) {
    evt.preventDefault()
}
