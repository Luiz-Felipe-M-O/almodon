import jsxmm from "../../module/jsxmm/element.ts"

export async function LoadCSSFile(path: string): Promise<void> {
    const style = document.createElement("link")

    style.rel = "stylesheet"
    style.href = path

    document.head.append(style)

    await new Promise(resolve => {
        style.onload = resolve
    })
}

type Field = [name: string, label: string, className?: string]

export function Skeleton(...fields: Field[]): HTMLElement {
    const skeleton = jsxmm.Element("div", { className: "form skeleton" })

    for (const [name, label, className] of fields) {
        const value = jsxmm.Element("div", { className: `value ${name}`, textContent: "..." })
        if (className !== undefined) {
            value.className += " " + className
        }

        const element = jsxmm.Element("div", { className: `field ${name}` },
            jsxmm.Element("div", { className: "label", textContent: label }),
            value,
        )

        skeleton.append(element)
    }

    return skeleton
}

const date_formatter = new Intl.DateTimeFormat("pt-BR", {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
})

export function FormatDate(date: Date): string {
    return date_formatter.format(date)
}