import { APIError, Cause } from "../../module/errors/error.ts"
import jsxmm from "../../module/jsxmm/element.ts"
import Source from "../support/source.ts"

const Messages: Record<number, string> = {
    204: "Página Vazia",
    403: "Você não tem Autorização para ver essa Página",
    404: "Página não Encontrada",
    500: "Erro no Servidor",
}

export function StatusPage(status: number): HTMLElement {
    let message = Messages[status]
    if (!Object.hasOwn(Messages, status)) {
        message = "Erro Desconhecido"
    }

    return (
        jsxmm.Element("div", { id: "almodon", className: "status" },
            jsxmm.Element("img", { className: "wrong-icon", src: Source.From("./assets/wrong.png") }),
            jsxmm.Element("div", { className: "status-text", textContent: message }),
            jsxmm.Element("div", { className: "status-code", textContent: `Erro ${status}` }),
            jsxmm.Element("a", { className: "go-back", href: "./", textContent: `Voltar para a Página Inicial` }),
        )
    )
}

export function ErrorPage(error: APIError): HTMLElement {
    const element = jsxmm.Element("div", { className: "error duck-tape" },
        jsxmm.Element("div", { className: "duck", textContent: "I'm a duck-taped solution" }),
        jsxmm.Element("div", { className: "kind", textContent: error.kind }),
        jsxmm.Element("div", { className: "title", textContent: error.title }),
        jsxmm.Element("div", { className: "message", textContent: error.message }),
    )

    if (error.cause !== undefined) {
        const cause = error_cause(error.cause)
        cause.classList.add("cause")
        element.append(cause)
    }

    return element
}

function error_cause(cause: Cause): HTMLElement {
    if (Array.isArray(cause)) {
        const element = jsxmm.Element("div")
        for (let i = 0; i < cause.length; i++) {
            const subcause = error_cause(cause[i])
            element.append(subcause)
        }

        return element
    }

    if (typeof cause === "string") {
        return jsxmm.Element("div", { textContent: cause })
    }

    const element = jsxmm.Element("div", { className: "error duck-tape" },
        jsxmm.Element("div", { className: "kind", textContent: cause.kind }),
        jsxmm.Element("div", { className: "title", textContent: cause.title }),
        jsxmm.Element("div", { className: "message", textContent: cause.message }),
    )

    if (cause.cause !== undefined) {
        const subcause = error_cause(cause.cause)
        element.append(subcause)
    }

    return element
}