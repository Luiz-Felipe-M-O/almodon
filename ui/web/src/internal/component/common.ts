export function ClickForKeys(element: HTMLElement, ...keys: string[]): void {
    element.addEventListener("keydown", function (evt: KeyboardEvent) {
        if (keys.includes(evt.key)) {
            (evt.target as HTMLElement).click()
            evt.preventDefault()
        }
    })
}