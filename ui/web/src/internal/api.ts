import Source from "./support/source.ts"

type API = {
    readonly Users: user.Gateway
}

export async function Construct(): Promise<API> {
    if (Source.server === "") {
        return await MockConstruct()
    } else {
        return await APIConstruct()
    }
}

export async function MockConstruct(): Promise<API> {
    const users = await import("./domain/user/gateway/mock.ts")

    return {
        Users: new users.UserGateway()
    }
}

export async function APIConstruct(): Promise<API> {
    const users = await import("./domain/user/gateway/api.ts")

    return {
        Users: new users.UserGateway(Source.From("./users", Source.server))
    }
}