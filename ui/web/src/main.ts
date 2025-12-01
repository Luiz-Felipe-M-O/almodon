import { Construct } from "./internal/api.ts"
import { Admin, Chief, Promoted, Role, Unlogged, User } from "./internal/auth/auth.ts"
import { ClickForKeys, ListenClickAndKeys } from "./internal/component/common.ts"
import { Context, Orquestrator, resouce, StaticPage } from "./internal/context/context.ts"
import MaterialView from "./internal/domain/material/view.ts"
import UserView from "./internal/domain/user/view.ts"
import ProfileView from "./internal/pages/profile.ts"
import Source from "./internal/support/source.ts"
import Signal from "./module/jsxmm/signals.ts"

async function main(): Promise<void> {
    const api = await Construct()

    const sidebar = document.getElementById("sidebar")
    if (sidebar === null) {
        throw new Error("There must be a #sidebar element")
    }

    const main = document.getElementById("main")
    if (main === null) {
        throw new Error("There must be a #main element")
    }

    const profile = new ProfileView(api.Users)
    const user = profile.User()

    const profiles = document.querySelectorAll<HTMLElement>(".profile .name")
    for (let i = 0; i < profiles.length; i++) {
        new Signal.Effect(effect_user.bind(null, profiles[i], user))
    }

    await profile.AsyncHTML()

    setup_navigation(sidebar, main, user,
        ["home", new StaticPage(Source.From("./pages/home.html"))],
        ["materials", new MaterialView(api.Materials)],
        ["users", new UserView(api.Users)],
        ["about", new StaticPage(Source.From("./pages/about.html"))],
        ["profile", profile],
    )
}

function setup_navigation(sidebar: HTMLElement, room: HTMLElement, user: Signal.RValue<user.Response | null>, ...namespaces: [string, Context][]): void {
    if (namespaces.length === 0) {
        throw new Error("No namespaces are not allowed")
    }

    const progress = document.getElementById("progress")

    let orq: Orquestrator
    if (Source.server === "") {
        orq = new Orquestrator(room)
    } else {
        orq = new Orquestrator(room, new resouce())
    }

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

        ListenClickAndKeys(option, handler_sidebar_click.bind(null, option, options, orq, progress), "Enter", " ")
    }

    for (let i = 0; i < options.length; i++) {
        const option = options[i]

        const namespace = option.dataset["namespace"]
        if (namespace === current) {
            handler_sidebar_click(option, options, orq, progress)
        }
    }

    new Signal.Effect(effect_namespace.bind(null, options, orq, user))
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

function effect_user(profile: HTMLElement, signal: Signal.RValue<user.Response | null>): void {
    const user = signal.Read()
    if (user === null) {
        profile.textContent = "Entrar"
        return
    }

    const names = user.name.split(" ")
    switch (names.length) {
    case 0:
        profile.textContent = "..."
        break

    case 1:
        profile.textContent = names[0]
        break

    case 2:
        profile.textContent = names[0] + " " + names[1]
        break

    default:
        if (names[names.length - 2].length <= 3) {
            profile.textContent = names[0] + " " + names[names.length - 2] + " " + names[names.length - 1]
        } else {
            profile.textContent = names[0] + " " + names[names.length - 1]
        }

        break
    }
}

const NamespaceForRoles: Record<Role, string[]> = {
    [Chief]: ["home", "materials", "users", "about", "profile"],
    [Promoted]: ["home", "materials", "about", "profile"],
    [Admin]: ["home", "materials", "about", "profile"],
    [User]: ["home", "materials", "about", "profile"],
    [Unlogged]: ["home", "about", "profile"],
}

function effect_namespace(options: NodeListOf<HTMLElement>, orq: Orquestrator, signal: Signal.RValue<user.Response | null>): void {
    const user = signal.Read()
    let namespaces

    if (user === null) {
        namespaces = NamespaceForRoles[Unlogged]
    } else {
        namespaces = NamespaceForRoles[user.role as Role]
    }

    for (let i = 0; i < options.length; i++) {
        const option = options[i]
        option.style.display = "none"

        const namespace = option.dataset["namespace"]
        if (namespace === undefined) {
            continue
        }

        if (namespaces.includes(namespace)) {
            option.style.display = ""
        }
    }

    const current = orq.Current()
    if (current === undefined || !namespaces.includes(current)) {
        options[0].click()
    }
}

window.addEventListener("DOMContentLoaded", main, { once: true })