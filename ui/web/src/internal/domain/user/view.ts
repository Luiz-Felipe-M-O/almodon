import { AsyncTry } from "../../../module/errors/try.ts"
import jsxmm from "../../../module/jsxmm/element.ts"
import Signal from "../../../module/jsxmm/signals.ts"
import { Role, RoleToString } from "../../auth/auth.ts"
import { ClickForKeys, ListenClickAndKeys } from "../../component/common.ts"
import { ErrorPage } from "../../component/status.ts"
import { Dialog } from "../../support/dialog.ts"
import { EmailInput, Form, FormField, FormInfo, SIAPEInput, TextInput } from "../../component/form.ts"
import { TableView } from "../../support/table.ts"

export default class UserView {
    #users: user.Gateway

    #build: HTMLElement | null
    #table: TableView<user.Response> | null

    constructor(gateway: user.Gateway) {
        this.#users = gateway

        this.#build = null
        this.#table = null
    }

    Final(): boolean {
        return false
    }

    HTML(): HTMLElement {
        if (this.#build === null) {
            this.build()
        }

        this.fetch()
        return this.#build!
    }

    private build(): void {
        const buttons = {
            create: jsxmm.Element("button", { className: "button create", textContent: "Novo Usuário" }),
            update: jsxmm.Element("button", { className: "button normal", textContent: "Editar" }),
            delete: jsxmm.Element("button", { className: "button destruct", textContent: "Excluir" }),
        }

        this.#table = new TableView(
            "Usuários",
            jsxmm.Element("thead", {},
                jsxmm.Element("tr", {},
                    jsxmm.Element("th", { textContent: "UUID" }),
                    jsxmm.Element("th", { textContent: "Nome" }),
                    jsxmm.Element("th", { textContent: "SIAPE" }),
                    jsxmm.Element("th", { textContent: "e-mail" }),
                    jsxmm.Element("th", { textContent: "Perfil" }),
                ),
            ),
            format,
            buttons.update,
            buttons.delete,
            buttons.create,
        )

        const selected = this.#table!.Selected()

        new Signal.Effect(() => {
            switch (selected.Read().size) {
            case 1:
                buttons.update.disabled = false
                buttons.delete.disabled = false
                break

            default:
                buttons.update.disabled = true
                buttons.delete.disabled = true
                break
            }
        })

        ListenClickAndKeys(buttons.create, async () => {
            buttons.create.disabled = true

            if (await create_form(this.create.bind(this))) {
                await this.fetch()
            } else {
                buttons.create.disabled = false
            }
        }, "Enter", " ")

        ListenClickAndKeys(buttons.update, async () => {
            if (selected.Read().size !== 1) {
                return
            }

            buttons.update.disabled = true

            const user = Array.from(selected.Read())[0]
            if (await update_form(user, this.update.bind(this))) {
                await this.fetch()
            } else {
                buttons.update.disabled = false
            }
        }, "Enter", " ")

        ListenClickAndKeys(buttons.delete, async () => {
            if (selected.Read().size !== 1) {
                return
            }

            buttons.delete.disabled = true

            const user = Array.from(selected.Read())[0]
            if (await delete_form(user, this.delete.bind(this))) {
                await this.fetch()
            } else {
                buttons.delete.disabled = false
            }
        }, "Enter", " ")

        this.#build = jsxmm.Element("div", { className: "view", id: "almodon" })
    }

    private async fetch(): Promise<void> {
        const [users, error] = await AsyncTry(() => this.#users.List(0, Number.MAX_SAFE_INTEGER))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                return
            }

            this.#build!.replaceWith(ErrorPage(error))
            return
        }

        this.#table!.Store(users.records)
        this.#build!.replaceChildren(this.#table!.HTML())
    }

    private async create(form: HTMLFormElement, reporter: HTMLElement): Promise<boolean> {
        const [name, email, siape, password, role] = FormInfo(form, "name", "email", "siape", "password", "role")
        const [, error] = await AsyncTry(() => this.#users.Create({ name, email, siape, password, role }))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                reporter.textContent = `Erro: ${error.message}`
                return false
            }

            reporter.textContent = `Erro: ${error.message}`
            return false
        }

        return true
    }

    private async update(uuid: UUID, form: HTMLFormElement, reporter: HTMLElement): Promise<boolean> {
        const [name, email] = FormInfo(form, "name", "email")
        const [, error] = await AsyncTry(() => this.#users.Patch(uuid, { name, email }))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                reporter.textContent = `Erro: ${error.message}`
                return false
            }

            reporter.textContent = `Erro: ${error.message}`
            return false
        }

        return true
    }

    private async delete(uuid: UUID, reporter: HTMLElement): Promise<boolean> {
        const [, error] = await AsyncTry(() => this.#users.Delete(uuid))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                reporter.textContent = `Erro: ${error.message}`
                return false
            }

            reporter.textContent = `Erro: ${error.message}`
            return false
        }

        return true
    }
}

function format(user: user.Response): HTMLTableRowElement {
    return jsxmm.Element("tr", {},
        jsxmm.Element("td", { className: "uuid monospace", textContent: user.uuid }),
        jsxmm.Element("td", { className: "name left", textContent: user.name }),
        jsxmm.Element("td", { className: "siape monospace", textContent: user.siape }),
        jsxmm.Element("td", { className: "email", textContent: user.email }),
        jsxmm.Element("td", { className: "role", textContent: RoleToString(user.role) }),
    )
}

async function create_form(create: (form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("Nome", "name", TextInput("name")),
        FormField("SIAPE", "siape", SIAPEInput("siape")),
        FormField("e-mail", "email", EmailInput("email")),
        FormField("Senha", "password", TextInput("password")),
        FormField("Perfil", "role", jsxmm.Element("select", { name: "role", required: true },
            jsxmm.Element("option", { value: "user", textContent: "Usuário" }),
            jsxmm.Element("option", { value: "admin", textContent: "Administrador" }),
            jsxmm.Element("option", { value: "chief", textContent: "Chefe" }),
        ))
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button create", textContent: "Criar" }),
    ]

    const dialog = new Dialog("Criar Novo Usuário", form, ...buttons)
    const element = dialog.HTML()

    document.body.append(element)
    dialog.Show()

    const result = await new Promise<boolean>(function (resolve) {
        ListenClickAndKeys(buttons[0], function () { resolve(false) }, "Enter", " ")
        ListenClickAndKeys(buttons[1], async function () {
            if (await create(form, reporter)) {
                resolve(true)
            }
        }, "Enter", " ")

        form.addEventListener("input", function () { reporter.replaceChildren() })
        form.addEventListener("keydown", async function (evt: KeyboardEvent) {
            if (evt.key === "Enter" && !evt.shiftKey) {
                if (await create(form, reporter)) {
                    resolve(true)
                }
            }
        })

        element.addEventListener("close", function () { resolve(false) })
    })

    dialog.Dispose()
    return result
}

async function update_form(user: user.Response, update: (uuid: UUID, form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("Nome", "name", TextInput("name", user.name)),
        FormField("e-mail", "email", EmailInput("email", user.email)),
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button create", textContent: "Salvar" }),
    ]

    const dialog = new Dialog("Editar Usuário", form, ...buttons)
    const element = dialog.HTML()

    document.body.append(element)
    dialog.Show()

    const result = await new Promise<boolean>(function (resolve) {
        ListenClickAndKeys(buttons[0], function () { resolve(false) }, "Enter", " ")
        ListenClickAndKeys(buttons[1], async function () {
            if (await update(user.uuid, form, reporter)) {
                resolve(true)
            }
        }, "Enter", " ")

        form.addEventListener("input", function () { reporter.replaceChildren() })
        form.addEventListener("keydown", async function (evt: KeyboardEvent) {
            if (evt.key === "Enter" && !evt.shiftKey) {
                if (await update(user.uuid, form, reporter)) {
                    resolve(true)
                }
            }
        })

        element.addEventListener("close", function () { resolve(false) })
    })

    dialog.Dispose()
    return result
}

async function delete_form(user: user.Response, delete_: (uuid: UUID, reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        jsxmm.Element("p", { textContent: `Tem certeza que deseja excluir o usuário "${user.name}"? Esta ação não pode ser desfeita.` }),
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button destruct", textContent: "Excluir" }),
    ]

    const dialog = new Dialog("Editar Usuário", form, ...buttons)
    const element = dialog.HTML()

    document.body.append(element)
    dialog.Show()

    const result = await new Promise<boolean>(function (resolve) {
        ListenClickAndKeys(buttons[0], function () { resolve(false) }, "Enter", " ")
        ListenClickAndKeys(buttons[1], async function () {
            if (await delete_(user.uuid, reporter)) {
                resolve(true)
            }
        }, "Enter", " ")

        element.addEventListener("close", function () { resolve(false) })
    })

    dialog.Dispose()
    return result
}