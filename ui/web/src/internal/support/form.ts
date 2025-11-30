import jsxmm from "../../module/jsxmm/element.ts"

export function Form(...children: HTMLElement[]) {
    const form = jsxmm.Element("form", { className: "form" }, ...children)
    form.addEventListener("submit", prevent_default)

    return form
}

export function FormField(label: string, input: HTMLInputElement) {
    input.classList.add("value")

    return jsxmm.Element("div", { className: "field" },
        jsxmm.Element("label", { className: "label", textContent: label }),
        input,
    )
}

export function Info<T extends string[]>(form: HTMLFormElement, ...fields: T): T {
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

function prevent_default(evt: Event) {
    evt.preventDefault()
}
