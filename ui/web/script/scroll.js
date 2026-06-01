function main() {
    const element = document.querySelector("main")
    if (element === null) {
        return
    }

    const scroll = sessionStorage.getItem("pos")
    if (scroll !== null) {
        element.scrollTop = parseInt(scroll, 10)
    }

    window.addEventListener("beforeunload", function () {
        sessionStorage.setItem("pos", element.scrollTop)
    })
}

document.addEventListener("DOMContentLoaded", main, { once: true })