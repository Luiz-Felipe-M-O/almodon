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
            jsxmm.Element("img", { className: "wrong-icon", src: Source.From("./dist/assets/wrong.png") }),
            jsxmm.Element("div", { className: "status-text", textContent: message }),
            jsxmm.Element("div", { className: "status-code", textContent: `Erro ${status}` }),
            jsxmm.Element("a", { className: "go-back", href: "./", textContent: `Voltar para a Página Inicial` }),
        )
    )
}
