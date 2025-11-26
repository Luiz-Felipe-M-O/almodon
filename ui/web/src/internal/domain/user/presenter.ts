import { APIError } from "../../../module/errors/error.ts"
import jsxmm from "../../../module/jsxmm/element.ts"

namespace UserPresenter {
    const Header: HTMLElement = (
        jsxmm.Element("thead", {},
            jsxmm.Element("th", { textContent: "UUID" }),
            jsxmm.Element("th", { textContent: "Nome" }),
            jsxmm.Element("th", { textContent: "SIAPE" }),
            jsxmm.Element("th", { textContent: "e-mail" }),
            jsxmm.Element("th", { textContent: "Perfil" }),
        )
    )

    const EmptyTable: HTMLElement = (
        jsxmm.Element("table", { className: "user" },
            Header.cloneNode(true),
            jsxmm.Element("tfoot", { textContent: "nenhum registro encontrado" })
        )
    )

    export function Table(users: user.ListResponse): HTMLElement {
        if (users.length === 0) {
            return EmptyTable.cloneNode(true) as HTMLElement
        }

        return jsxmm.Element("div", { className: "table user" },
            Header.cloneNode(true),
            jsxmm.Element("tbody", {},
                ...users.records.map(user => row(user)),
            ),
        )
    }

    function row(user: user.Response): HTMLElement {
        return (
            jsxmm.Element("tr", { className: "record user" },
                jsxmm.Element("td", { className: "uuid", textContent: user.uuid }),
                jsxmm.Element("td", { className: "name", textContent: user.name }),
                jsxmm.Element("td", { className: "siape", textContent: `${user.siape}`.padStart(7, '0') }),
                jsxmm.Element("td", { className: "email", textContent: user.email }),
                jsxmm.Element("td", { className: "role", textContent: user.role }),
            )
        )
    }
}

export default UserPresenter