const OriginToRoot: Record<string, { client: string, server: string }> = {
    "http://localhost:4545": {
        client: "http://localhost:4545/",
        server: "http://localhost:4545/api/v1/",
    },
    "https://alan-b-lima.github.io": {
        client: "https://alan-b-lima.github.io/ui/web",
        server: "",
    },
}

function urls() {
    const origin = location.origin
    if (!Object.hasOwn(OriginToRoot, origin)) {
        throw new Error("Unknown location " + origin)
    }

    const root = OriginToRoot[origin]
    return root
}

namespace Source {
    export function From(path: string, origin: string = client): string {
        return new URL(path, origin).href
    }

    const source = urls()

    export const client = source.client
    export const server = source.server
}

export default Source
