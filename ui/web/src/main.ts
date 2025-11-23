import { context, Orquestrator } from "./internal/context/context.ts"

async function main(): Promise<void> {
    const sidebar = document.getElementById("sidebar")
    if (sidebar === null) {
        throw new Error("There must be a #sidebar element")
    }

    const content = document.getElementById("content")
    if (content === null) {
        throw new Error("There must be a #content element")
    }

    const namespacing = setup_menu(content, "home", "users", "users-2", "about")
    setup_sidebar(sidebar, namespacing)
}

function setup_menu(room: HTMLElement, ...namespaces: string[]): (namespace: string) => boolean {
    if (namespaces.length === 0) {
        throw new Error("No namespaces are not allowed")
    }

    const orq = new Orquestrator(room)

    for (const namespace of namespaces) {
        orq.Link(namespace, new context(`./pages/${namespace}.html`))
    }

    let hash = location.hash.slice(1)
    if (!namespaces.includes(hash)) {
        hash = namespaces[0]
    }

    if (orq.SwapTo(hash)) {
        const target = document.querySelector(`[data-namespace="${hash}"]`)
        if (target !== null) {
            target.classList.add("selected")
            location.hash = hash
        }
    }

    return function (namespace: string): boolean {
        if (orq.SwapTo(namespace)) {
            location.hash = namespace
            return true
        }

        return false
    }
}

function setup_sidebar(sidebar: HTMLElement, namespacing: (namespace: string) => boolean): void {
    const options = sidebar.querySelectorAll<HTMLElement>(".option")
    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        option.addEventListener("click", function (evt: MouseEvent) {
            for (let j = 0; j < options.length; j++) {
                options[j].classList.remove("selected")
            }

            const namespace = option.dataset["namespace"]
            if (namespace !== undefined && namespacing(namespace)) {
                option.scrollIntoView()
                option.classList.add("selected")
            }
        })

        click_for_keys(option, "Enter", " ")
    }
}

function click_for_keys(element: HTMLElement, ...keys: string[]): void {
    element.addEventListener("keydown", function (evt: KeyboardEvent) {
        if (keys.includes(evt.key)) {
            (evt.target as HTMLElement).click()
            evt.preventDefault()
        }
    })
}

window.addEventListener("DOMContentLoaded", main)