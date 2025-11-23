import { context, Orquestrator } from "./internal/context/context.ts"
import Source from "./internal/support/source.ts"

async function main(): Promise<void> {
    const sidebar = document.getElementById("sidebar")
    if (sidebar === null) {
        throw new Error("There must be a #sidebar element")
    }

    const main = document.getElementById("main")
    if (main === null) {
        throw new Error("There must be a #main element")
    }

    setup_navigation(sidebar, main, "home", "users", "users-2", "about")
}

function setup_navigation(sidebar: HTMLElement, room: HTMLElement, ...namespaces: string[]): void {
    if (namespaces.length === 0) {
        throw new Error("No namespaces are not allowed")
    }
    
    const progress = document.getElementById("progress")

    const orq = new Orquestrator(room)
    for (const namespace of namespaces) {
        const ctx = new context(Source.From(`./dist/pages/${namespace}.html`))
        
        if (progress !== null) {   
            ctx.onpreload = () => { progress.classList.remove("complete") }
            ctx.onload = () => { progress.classList.add("complete") }
        }

        orq.Link(namespace, ctx)
    }

    let current = orq.SwapperCurrent()
    if (current === undefined || !namespaces.includes(current)) {
        current = namespaces[0]
    }

    const options = sidebar.querySelectorAll<HTMLElement>(".option")
    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        option.addEventListener("click", handler_sidebar_click.bind(null, option, options, orq))
        click_for_keys(option, "Enter", " ")
    }

    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        const namespace = option.dataset["namespace"]
        if (namespace === current) {
            handler_sidebar_click(option, options, orq)
            break
        }
    }
}

function handler_sidebar_click(option: HTMLElement, options: NodeListOf<HTMLElement>, orq: Orquestrator): void {
    for (let j = 0; j < options.length; j++) {
        options[j].classList.remove("selected")
    }

    const namespace = option.dataset["namespace"]
    if (namespace === undefined) {
        return
    }

    if (!orq.SwapTo(namespace)) {
        return
    }

    option.classList.add("selected")
    option.scrollIntoView()
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