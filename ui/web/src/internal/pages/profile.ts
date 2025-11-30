import { AsyncTry } from "../../module/errors/try.ts"
import jsxmm from "../../module/jsxmm/element.ts"
import { Dialog } from "../support/dialog.ts"
import { ErrorPage } from "../component/status.ts"
import { Form, FormField, Info } from "../support/form.ts"
import { FormatDate, LoadCSSFile, Skeleton } from "../support/pages.ts"
import Source from "../support/source.ts"
import Signal from "../../module/jsxmm/signals.ts"
import { RoleToString } from "../auth/auth.ts"

export default class ProfileView {
    #users: user.Gateway

    #build: HTMLElement | null
    #user: Signal.Value<user.Response | null>

    #logged_in_page: HTMLElement | null
    #logged_out_page: HTMLElement | null
    #login_form: Dialog | null

    #edit_form: Dialog | null

    constructor(users: user.Gateway) {
        this.#users = users

        this.#build = null
        this.#user = new Signal.Value<user.Response | null>(null)

        this.#logged_in_page = null
        this.#logged_out_page = null
        this.#login_form = null

        this.#edit_form = null
    }

    User(): Signal.RValue<user.Response | null> {
        return this.#user
    }

    Final(): boolean {
        return false
    }

    HTML(): HTMLElement {
        if (this.#build === null) {
            LoadCSSFile(Source.From("./style/profile.css"))
            this.#build = jsxmm.Element("div", { className: "profile", id: "almodon" })
        }

        this.verify()
        return this.#build
    }

    private async verify(): Promise<void> {
        const [user, error] = await AsyncTry(() => this.#users.Me())
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                return
            }

            if (error.kind !== "not found") {
                this.#build?.replaceChildren(ErrorPage(error))
                return
            }

            this.#user.Write(null)
            this.logged_out()
            return
        }

        this.#user?.Write(user)
        this.logged_in(user)
    }

    private logged_in(user: user.Response): void {
        if (this.#logged_in_page === null) {
            const profile = Skeleton(
                ["uuid", "UUID", "monospace"],
                ["siape", "SIAPE", "monospace"],
                ["name", "Nome"],
                ["email", "e-mail"],
                ["role", "Perfil"],
                ["created", "Criado em"],
                ["updated", "Atualizado em"],
            )

            const buttons = [
                jsxmm.Element("button", { className: "button normal logout", textContent: "Sair" }),
                jsxmm.Element("button", { className: "button create", textContent: "Editar" }),
            ]

            this.#logged_in_page = jsxmm.Element("div", { className: "logged-in" },
                jsxmm.Element("div", { className: "body" }, profile),
                jsxmm.Element("div", { className: "actions" }, ...buttons),
            )

            buttons[0].addEventListener("click", async () => {
                await this.#users.Logout()
                await this.verify()

                this.#login_form?.Close()
            })

            buttons[1].addEventListener("click", () => {
                this.editting()
            })
        }

        const profile = this.#logged_in_page.children[0].children[0]

        profile.querySelector<HTMLElement>(".value.uuid")!.textContent = user.uuid
        profile.querySelector<HTMLElement>(".value.siape")!.textContent = user.siape
        profile.querySelector<HTMLElement>(".value.name")!.textContent = user.name
        profile.querySelector<HTMLElement>(".value.email")!.textContent = user.email
        profile.querySelector<HTMLElement>(".value.role")!.textContent = RoleToString(user.role)
        profile.querySelector<HTMLElement>(".value.created")!.textContent = FormatDate(user.created)
        profile.querySelector<HTMLElement>(".value.updated")!.textContent = FormatDate(user.updated)

        this.#build?.replaceChildren(this.#logged_in_page)
    }

    private logged_out(): void {
        if (this.#login_form === null) {
            this.#login_form = login_form(this.login.bind(this))
        }

        const form = this.#login_form.HTML().querySelector("form")!
        form.querySelector<HTMLElement>(".reporter")?.replaceChildren()
        form.reset()

        this.#login_form.Show()

        if (this.#logged_out_page === null) {
            const login = jsxmm.Element("button", { className: "button create", textContent: "Entrar" })
            login.addEventListener("click", () => { this.#login_form!.Show() })

            this.#logged_out_page = jsxmm.Element("div", { className: "logged-out" },
                jsxmm.Element("div", { textContent: "Usuário não logado" }),
                jsxmm.Element("div", { className: "actions" }, login),
            )
        }

        this.#build?.replaceChildren(this.#logged_out_page)
    }

    private async login(form: HTMLFormElement, reporter: HTMLElement): Promise<boolean> {
        const [siape, password] = Info(form, "siape", "password")
        const [, error] = await AsyncTry(() => this.#users.Autheticate(siape, password))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                reporter.textContent = `Erro: ${error.message}`
                return false
            }

            reporter.textContent = `Erro: ${error.message}`
            return false
        }

        this.verify()
        return true
    }

    private editting(): void {
        const user = this.#user.Read()
        if (user === null) {
            this.verify()
            return
        }

        if (this.#edit_form === null) {
            this.#edit_form = edit_form(this.edit.bind(this))
        }

        const form = this.#edit_form.HTML().querySelector("form")!
        form.querySelector<HTMLElement>(".reporter")?.replaceChildren()
        form.reset()

        form.querySelector<HTMLInputElement>(".value.name")!.value = user.name
        form.querySelector<HTMLInputElement>(".value.email")!.value = user.email

        this.#edit_form.Show()
    }

    private async edit(form: HTMLFormElement, reporter: HTMLElement): Promise<boolean> {
        const user = this.#user.Read()
        if (user === null) {
            reporter.textContent = "Erro: usuário não está logado"
            return false
        }

        const [name, email] = Info(form, "name", "email")
        const [, error] = await AsyncTry(() => this.#users.Patch(user.uuid, { name, email }))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                reporter.textContent = `Erro: ${error.message}`
                return false
            }

            reporter.textContent = `Erro: ${error.message}`
            return false
        }

        this.verify()
        return true
    }
}

function login_form(login: (form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Dialog {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("SIAPE", jsxmm.Element("input", { type: "text", name: "siape", pattern: "\\d{7}", className: "editable monospace" })),
        FormField("Senha", jsxmm.Element("input", { type: "password", name: "password", className: "editable" })),
    )

    const buttons = [
        jsxmm.Element("button", { type: "button", className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { type: "button", className: "button create", textContent: "Entrar" }),
    ]

    const dialog = new Dialog("Login", form, ...buttons)

    form.addEventListener("input", () => { reporter.replaceChildren() })

    form.addEventListener("keydown", async (evt: KeyboardEvent) => {
        if (evt.key === "Enter" && !evt.shiftKey) {
            if (await login(form, reporter)) {
                dialog.Close()
            }
        }
    })

    buttons[0].addEventListener("click", dialog.Close.bind(dialog))

    buttons[1].addEventListener("click", async () => {
        buttons[1].disabled = true

        if (await login(form, reporter)) {
            dialog.Close()
        }

        buttons[1].disabled = false
    })

    document.body.append(dialog.HTML())
    return dialog
}

function edit_form(edit: (form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Dialog {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("Nome", jsxmm.Element("input", { type: "text", name: "name", className: "name editable", required: true })),
        FormField("e-mail", jsxmm.Element("input", { type: "email", name: "email", className: "email editable", pattern: ".+@.+\\..+", required: true })),
    )

    const buttons = [
        jsxmm.Element("button", { type: "button", className: "button normal", textContent: "Cancelar" },),
        jsxmm.Element("button", { type: "button", className: "button create", textContent: "Aplicar" }),
    ]

    const dialog = new Dialog("Login", form, ...buttons)

    form.addEventListener("input", () => { reporter.replaceChildren() })

    form.addEventListener("keydown", async (evt: KeyboardEvent) => {
        if (evt.key === "Enter" && !evt.shiftKey) {
            buttons[1].disabled = true

            if (await edit(form, reporter)) {
                dialog.Close()
            }

            buttons[1].disabled = false
        }
    })

    buttons[0].addEventListener("click", dialog.Close.bind(dialog))

    buttons[1].addEventListener("click", async () => {
        buttons[1].disabled = true

        if (await edit(form, reporter)) {
            dialog.Close()
        }

        buttons[1].disabled = false
    })

    document.body.append(dialog.HTML())
    return dialog
}