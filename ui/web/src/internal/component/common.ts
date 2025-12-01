export function ClickForKeys(element: HTMLElement, ...keys: string[]): void {
    element.addEventListener("keydown", function (evt: KeyboardEvent) {
        if (keys.includes(evt.key)) {
            (evt.target as HTMLElement).click()
            evt.preventDefault()
        }
    })
}

export function ListenClickAndKeys(element: HTMLElement, listener: (evt: Event) => void, ...keys: string[]): void {
    element.addEventListener("click", listener)
    ClickForKeys(element, ...keys)
}
