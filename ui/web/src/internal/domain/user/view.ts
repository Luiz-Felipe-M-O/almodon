import { AsyncTry } from "../../../module/errors/try.ts"
import jsxmm from "../../../module/jsxmm/element.ts"
import { ErrorPage } from "../../component/status.ts"

export default class UserView {
    #api: user.Gateway

    #build: HTMLElement | null

    constructor(gateway: user.Gateway) {
        this.#api = gateway

        this.#build = null
    }

    Final(): boolean {
        return false
    }

    HTML(): HTMLElement {
        if (this.#build === null) {
            this.#build = jsxmm.Element("div", { className: "user", id: "almodon" })
        }

        this.load()
        return this.#build
    }

    private async load(): Promise<void> {
        const [response, error] = await AsyncTry(() => this.#api.List(0, Number.MAX_SAFE_INTEGER))
        if (error !== null) {
            if (error instanceof Error) {
                console.error(error)
                return
            }

            this.#build?.replaceChildren(ErrorPage(error))
            return
        }

        this.#build?.replaceChildren(Table(response))
    }
}

const Header: HTMLElement = jsxmm.Element("thead", {},
    // jsxmm.Element("th"),
    jsxmm.Element("th", { textContent: "UUID" }),
    jsxmm.Element("th", { textContent: "Nome" }),
    jsxmm.Element("th", { textContent: "SIAPE" }),
    jsxmm.Element("th", { textContent: "e-mail" }),
    jsxmm.Element("th", { textContent: "Perfil" }),
)


const EmptyTable: HTMLElement = (
    jsxmm.Element("table", { className: "user" },
        Header.cloneNode(true),
        jsxmm.Element("tfoot", { textContent: "nenhum registro encontrado" })
    )
)

function Table(users: user.ListResponse): HTMLElement {
    if (users.length === 0) {
        return EmptyTable.cloneNode(true) as HTMLElement
    }

    return jsxmm.Element("div", { className: "table" },
        jsxmm.Element("table", {},
            Header.cloneNode(true),
            jsxmm.Element("tbody", {},
                ...users.records.map(user => row(user)),
            ),
        )
    )
}

function row(user: user.Response): HTMLElement {
    return (
        jsxmm.Element("tr", { className: "record user" },
            // jsxmm.Element("td", {}, jsxmm.Element("input", { type: "checkbox" })),
            jsxmm.Element("td", { className: "uuid monospace", textContent: user.uuid }),
            jsxmm.Element("td", { className: "name left", textContent: user.name }),
            jsxmm.Element("td", { className: "siape monospace", textContent: user.siape }),
            jsxmm.Element("td", { className: "email", textContent: user.email }),
            jsxmm.Element("td", { className: "role", textContent: user.role }),
        )
    )
}