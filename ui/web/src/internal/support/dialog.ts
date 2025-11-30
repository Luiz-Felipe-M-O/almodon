import jsxmm from "../../module/jsxmm/element.ts"
import Source from "./source.ts"

export class Dialog {
    #build: HTMLDialogElement | null

    #title: string
    #body: HTMLElement
    #buttons: HTMLElement[]

    constructor(title: string, body: HTMLElement, ...buttons: HTMLElement[]) {
        this.#build = null
        this.#title = title
        this.#body = body
        this.#buttons = buttons
    }

    Show(): void {
        if (this.#build === null) {
            this.HTML()
        }

        this.#build!.showModal()
    }

    HTML(): HTMLDialogElement {
        if (this.#build !== null) {
            return this.#build
        }

        this.#body.classList.add("body")

        const close = jsxmm.Element("button", { className: "close" },
            jsxmm.Element("img", { className: "icon", src: Source.From("./assets/icons/close.svg"), alt: "Close" }),
        )

        const dialog = jsxmm.Element("dialog", { className: "dialog" },
            jsxmm.Element("div", { className: "title" },
                jsxmm.Element("div", { textContent: this.#title }),
                close,
            ),
            this.#body,
        )

        if (this.#buttons.length > 0) {
            const buttons = jsxmm.Element("div", { className: "actions" })
            for (const button of this.#buttons) {
                buttons.append(button)
            }

            dialog.append(buttons)
        }

        close.addEventListener("click", this.Close.bind(this))

        this.#build = dialog
        return this.#build
    }

    Close(): void {
        if (this.#build === null) {
            return
        }

        this.#build.close()
    }

    Dispose(): void {
        if (this.#build === null) {
            return
        }

        this.#build.close()
        this.#build.remove()
        this.#build = null
    }
}
