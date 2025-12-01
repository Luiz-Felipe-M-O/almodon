import { AsyncTry } from "../../../module/errors/try.ts"
import jsxmm from "../../../module/jsxmm/element.ts"
import Signal from "../../../module/jsxmm/signals.ts"
import { ClickForKeys, ListenClickAndKeys } from "../../component/common.ts"
import { ErrorPage } from "../../component/status.ts"
import { Dialog } from "../../support/dialog.ts"
import { Form, FormField, FormInfo, TextInput } from "../../component/form.ts"
import { TableView } from "../../support/table.ts"

export default class MaterialView {
    #materials: material.Gateway

    #build: HTMLElement | null
    #table: TableView<material.Response> | null

    constructor(gateway: material.Gateway) {
        this.#materials = gateway

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
            create: jsxmm.Element("button", { className: "button create", textContent: "Novo Material" }),
            update: jsxmm.Element("button", { className: "button normal", textContent: "Editar" }),
            delete: jsxmm.Element("button", { className: "button destruct", textContent: "Excluir" }),
        }

        this.#table = new TableView(
            "Materiais",
            jsxmm.Element("thead", {},
                jsxmm.Element("tr", {},
                    jsxmm.Element("th", { textContent: "UUID" }),
                    jsxmm.Element("th", { textContent: "Nome" }),
                    jsxmm.Element("th", { textContent: "SIADS" }),
                    jsxmm.Element("th", { textContent: "CATMAT" }),
                    jsxmm.Element("th", { textContent: "e-Campus" }),
                    jsxmm.Element("th", { textContent: "Qtd. Mín." }),
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
            case 0:
                buttons.update.disabled = true
                buttons.delete.disabled = true
                break

            case 1:
                buttons.update.disabled = false
                buttons.delete.disabled = false
                break

            default:
                buttons.update.disabled = true
                buttons.delete.disabled = false
                break
            }
        })

        ListenClickAndKeys(buttons.create, async () => {
            buttons.create.disabled = true

            if (await create_form(this.create.bind(this))) {
                await this.fetch()
            }

            buttons.create.disabled = false
        }, "Enter", " ")

        ListenClickAndKeys(buttons.update, async () => {
            if (selected.Read().size !== 1) {
                return
            }

            buttons.update.disabled = true

            const material = Array.from(selected.Read())[0]
            if (await update_form(material, this.update.bind(this))) {
                await this.fetch()
            } else {
                buttons.update.disabled = false
            }
        }, "Enter", " ")

        ListenClickAndKeys(buttons.delete, async () => {
            if (selected.Read().size === 0) {
                return
            }

            buttons.delete.disabled = true

            const materials = Array.from(selected.Read())
            if (await delete_form(materials, this.delete.bind(this))) {
                await this.fetch()
            } else {
                buttons.delete.disabled = false
            }
        }, "Enter", " ")

        this.#build = jsxmm.Element("div", { className: "view", id: "almodon" })
    }

    private async fetch(): Promise<void> {
        const [materials, error] = await AsyncTry(() => this.#materials.List(0, Number.MAX_SAFE_INTEGER))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                return
            }

            this.#build!.replaceWith(ErrorPage(error))
            return
        }

        this.#table!.Store(materials.records)
        this.#build!.replaceChildren(this.#table!.HTML())
    }

    private async create(form: HTMLFormElement, reporter: HTMLElement): Promise<boolean> {
        const [name, siads, catmat, ecampus, description, unit, min_quantity] = FormInfo(form, "name", "siads", "catmat", "ecampus", "description", "unit", "min_quantity")
        const [, error] = await AsyncTry(() => this.#materials.Create({
            name,
            siads,
            catmat,
            ecampus,
            description,
            unit,
            min_quantity: parseFloat(min_quantity)
        }))
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
        const [name, siads, catmat, ecampus, description, unit, min_quantity] = FormInfo(form, "name", "siads", "catmat", "ecampus", "description", "unit", "min_quantity")
        const data: material.PartialEntity = {
            name,
            siads,
            catmat,
            ecampus,
            description,
            unit,
            min_quantity: parseFloat(min_quantity),
        }

        const [, error] = await AsyncTry(() => this.#materials.Patch(uuid, data))
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

    private async delete(uuids: UUID[], reporter: HTMLElement): Promise<boolean> {
        const tasks = uuids.map(uuid => this.#materials.Delete(uuid))
        const [, error] = await AsyncTry(() => Promise.all(tasks))
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

function format(material: material.Response): HTMLTableRowElement {
    return jsxmm.Element("tr", {},
        jsxmm.Element("td", { className: "uuid monospace", textContent: material.uuid }),
        jsxmm.Element("td", { className: "name left", textContent: material.name }),
        jsxmm.Element("td", { className: "siads monospace", textContent: material.siads }),
        jsxmm.Element("td", { className: "catmat monospace", textContent: material.catmat }),
        jsxmm.Element("td", { className: "ecampus monospace", textContent: material.ecampus }),
        jsxmm.Element("td", { className: "min-quantity number", textContent: `${material.min_quantity} ${material.unit}` }),
    )
}

async function create_form(create: (form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("Nome", "name", TextInput("name")),
        FormField("SIADS", "siads", TextInput("siads", "")),
        FormField("CATMAT", "catmat", TextInput("catmat", "")),
        FormField("e-Campus", "ecampus", TextInput("ecampus", "")),
        FormField("Descrição", "description", jsxmm.Element("textarea", { name: "description", maxLength: 2048, rows: 3 })),
        FormField("Unidade", "unit", TextInput("unit", "")),
        FormField("Quantidade Mínima", "min_quantity", jsxmm.Element("input", { type: "number", min: "0", step: "0.01" })),
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button create", textContent: "Criar" }),
    ]

    const dialog = new Dialog("Criar Novo Material", form, ...buttons)
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

async function update_form(material: material.Response, update: (uuid: UUID, form: HTMLFormElement, reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        FormField("Nome", "name", TextInput("name", material.name)),
        FormField("SIADS", "siads", TextInput("siads", material.siads)),
        FormField("CATMAT", "catmat", TextInput("catmat", material.catmat)),
        FormField("e-Campus", "ecampus", TextInput("ecampus", material.ecampus)),
        FormField("Descrição", "description", jsxmm.Element("textarea", {
            name: "description",
            maxLength: 2048,
            rows: 3,
            textContent: material.description
        })),
        FormField("Unidade", "unit", TextInput("unit", material.unit)),
        FormField("Quantidade Mínima", "min_quantity", jsxmm.Element("input", { type: "number", min: "0", step: "0.01", value: `${material.min_quantity}` })),
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button create", textContent: "Salvar" }),
    ]

    const dialog = new Dialog("Editar Material", form, ...buttons)
    const element = dialog.HTML()

    document.body.append(element)
    dialog.Show()

    const result = await new Promise<boolean>(function (resolve) {
        ListenClickAndKeys(buttons[0], function () { resolve(false) }, "Enter", " ")
        ListenClickAndKeys(buttons[1], async function () {
            if (await update(material.uuid, form, reporter)) {
                resolve(true)
            }
        }, "Enter", " ")

        form.addEventListener("input", function () { reporter.replaceChildren() })
        form.addEventListener("keydown", async function (evt: KeyboardEvent) {
            if (evt.key === "Enter" && !evt.shiftKey) {
                if (await update(material.uuid, form, reporter)) {
                    resolve(true)
                }
            }
        })

        element.addEventListener("close", function () { resolve(false) })
    })

    dialog.Dispose()
    return result
}

async function delete_form(materials: material.Response[], delete_: (uuid: UUID[], reporter: HTMLElement) => Promise<boolean>): Promise<boolean> {
    const reporter = jsxmm.Element("div", { className: "reporter" })

    const form = Form(
        reporter,
        jsxmm.Element("p", { textContent: `Tem certeza que deseja excluir esses materiais? Esta ação não pode ser desfeita.` }),
    )

    const buttons = [
        jsxmm.Element("button", { className: "button normal", textContent: "Cancelar" }),
        jsxmm.Element("button", { className: "button destruct", textContent: "Excluir" }),
    ]

    const dialog = new Dialog("Excluir Material", form, ...buttons)
    const element = dialog.HTML()

    document.body.append(element)
    dialog.Show()

    const result = await new Promise<boolean>(function (resolve) {
        ListenClickAndKeys(buttons[0], function () { resolve(false) }, "Enter", " ")
        ListenClickAndKeys(buttons[1], async function () {
            if (await delete_(materials.map(m => m.uuid), reporter)) {
                resolve(true)
            }
        }, "Enter", " ")

        element.addEventListener("close", function () { resolve(false) })
    })

    dialog.Dispose()
    return result
}
