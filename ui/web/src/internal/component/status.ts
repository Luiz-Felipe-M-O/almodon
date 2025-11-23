import jsxmm from "../../module/jsxmm/element.ts"
import Source from "../support/source.ts"

const Messages: Record<number, string> = {
    204: "Página Vazia",
    404: "Página não Encontrada",
}

export function StatusPage(status: number): HTMLElement {
    return (
        jsxmm.Element("div", { id: "almodon", className: "status" },
            jsxmm.Element("img", { className: "wrong-icon", src: Source.From("/dist/assets/wrong.png") }),
            jsxmm.Element("div", { className: "status-text", textContent: Messages[status] }),
            jsxmm.Element("div", { className: "status-code", textContent: `Erro ${status}` }),
            jsxmm.Element("a", { className: "go-back", href: "./", textContent: `Voltar para a Página Inicial` }),
        )
    )
}

export function ErrConnectionPage(status: number): HTMLElement {
    return (
        jsxmm.Element("div", { id: "almodon", className: "status" },
            jsxmm.Element("img", { className: "wrong-icon", src: Source.From("/dist/assets/wrong.png") }),
            jsxmm.Element("div", { className: "status-text", textContent: Messages[status] }),
            jsxmm.Element("div", { className: "status-code", textContent: `Erro ${status}` }),
            jsxmm.Element("a", { className: "go-back", href: "./", textContent: `Voltar para a Página Inicial` }),
        )
    )
}
