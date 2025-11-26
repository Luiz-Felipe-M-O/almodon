import { Context, Orquestrator, resouce, StaticPage } from "./internal/context/context.ts"
import UserView from "./internal/domain/user/view.ts"
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

    const api = await setup_api()

    setup_navigation(sidebar, main,
        ["home", new StaticPage(Source.From("./dist/pages/home.html"))],
        ["users", new UserView(api.users)],
        ["table", new StaticPage(Source.From("./dist/pages/table.html"))],
        ["about", new StaticPage(Source.From("./dist/pages/about.html"))],
    )
}

async function setup_api(): Promise<{ users: user.Gateway }> {
    if (Source.server === "") {
        const users = await import("./internal/domain/user/gateway/mock.ts")
        return {
            users: new users.UserGateway(),
        }
    }

    const users = await import("./internal/domain/user/gateway/api.ts")
    return {
        users: new users.UserGateway(Source.From("./users", Source.server)),
    }
}

function setup_navigation(sidebar: HTMLElement, room: HTMLElement, ...namespaces: [string, Context][]): void {
    if (namespaces.length === 0) {
        throw new Error("No namespaces are not allowed")
    }

    const progress = document.getElementById("progress")

    const orq = new Orquestrator(room)
    for (const [namespace, context] of namespaces) {
        if (progress !== null && context instanceof StaticPage) {
            context.onpreload = () => { progress.classList.remove("complete") }
            context.onload = () => { progress.classList.add("complete") }
        }

        orq.Link(namespace, context)
    }

    let current = orq.Current()
    if (current === undefined || namespaces.findIndex(([namespace]) => namespace === current) === -1) {
        current = namespaces[0][0]
    }

    const options = sidebar.querySelectorAll<HTMLElement>(".option")
    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        option.addEventListener("click", handler_sidebar_click.bind(null, option, options, orq, progress))
        click_for_keys(option, "Enter", " ")
    }

    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        const namespace = option.dataset["namespace"]
        if (namespace === current) {
            handler_sidebar_click(option, options, orq, progress)
            break
        }
    }
}

function handler_sidebar_click(option: HTMLElement, options: NodeListOf<HTMLElement>, orq: Orquestrator, progress: HTMLElement | null): void {
    for (let j = 0; j < options.length; j++) {
        options[j].classList.remove("selected")
    }

    const namespace = option.dataset["namespace"]
    if (namespace === undefined) {
        return
    }

    if (progress !== null) {
        progress.classList.add("complete")
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