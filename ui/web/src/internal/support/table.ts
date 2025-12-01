import jsxmm from "../../module/jsxmm/element.ts"
import Signal from "../../module/jsxmm/signals.ts"
import { ClickForKeys } from "../component/common.ts"

export class TableView<T> {
    #build: HTMLElement

    #title: string
    #action: HTMLElement
    #body: HTMLElement
    #note: HTMLElement

    #content: HTMLTableSectionElement
    #format: (value: T) => HTMLTableRowElement

    #set: Set<T>
    #selected: Signal.Value<Set<T>>

    constructor(title: string, header: HTMLTableSectionElement, format: (value: T) => HTMLTableRowElement, ...buttons: HTMLElement[]) {
        this.#title = title
        this.#action = jsxmm.Element("div", { className: "actions" }, ...buttons)
        this.#note = jsxmm.Element("div", { className: "note", textContent: "nenhum registro encontrado" })

        this.#content = jsxmm.Element("tbody", {})
        this.#format = format

        this.#body = jsxmm.Element("div", { className: "table body" },
            jsxmm.Element("table", {}, header, this.#content),
            this.#note,
        )

        this.#build = jsxmm.Element("div", {},
            jsxmm.Element("h2", { className: "title", textContent: this.#title }),
            this.#action,
            this.#body,
        )

        this.#set = new Set<T>()
        this.#selected = new Signal.Value(this.#set, () => false)
    }

    Selected(): Signal.RValue<Set<T>> {
        return this.#selected
    }

    HTML(): HTMLElement {
        return this.#build
    }

    Store(values: T[]): void {
        this.#content.replaceChildren()
        
        this.#selected.Read().clear()
        this.#selected.Write(this.#set)

        for (const value of values) {
            const row = this.#format(value)

            row.addEventListener("click", handle_row_selection.bind(null, this.#selected as any, value, row))
            ClickForKeys(row, "Enter", " ")

            this.#content.append(row)
        }
    }
}

function handle_row_selection<T>(selected: Signal.Value<Set<T>>, value: T, row: HTMLElement): void {
    const set = selected.Read()
    if (set.has(value)) {
        set.delete(value)
        row.classList.remove("selected")
    } else {
        set.add(value)
        row.classList.add("selected")
    }

    selected.Write(set)
}