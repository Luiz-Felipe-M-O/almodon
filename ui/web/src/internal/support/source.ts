const OriginToRoot = {
    ":3000": {
        client: "/dist/",
        server: "",
    },
    ":4545": {
        client: "/",
        server: "/api/v1/",
    },
    ":80": {
        client: "/almodon/ui/web/dist/",
        server: "",
    },
} as const

namespace Source {
    export function From(path: string, origin: string = client): string {
        return new URL(path, origin).href
    }

    const source = urls()

    export const client = source.client
    export const server = source.server
}

export default Source

function urls() {
    const origin = location.origin
    let root

    Find: {
        if (origin.endsWith(":4545")) {
            root = OriginToRoot[":4545"]
            break Find
        }

        if (origin.endsWith(":3000")) {
            root = OriginToRoot[":3000"]
            break Find
        }

        if (origin.endsWith(":80") || !origin.includes(":")) {
            root = OriginToRoot[":80"]
            break Find
        }

        throw new Error("Unknown location " + origin)
    }

    let server = ""
    if (root.server !== "") {
        server = Source.From(root.server, origin)
    }

    return {
        client: Source.From(root.client, origin),
        server: server,
    }
}
