import Source from "./support/source.ts"

type API = {
    readonly Materials: material.Gateway
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
    throw new Error("not implemented")

    // const material = await import("./domain/material/gateway/mock.ts")
    const users = await import("./domain/user/gateway/mock.ts")

    return {
        Materials: /*new material.MaterialGateway()*/ null as any,
        Users: new users.UserGateway()
    }
}

export async function APIConstruct(): Promise<API> {
    const material = await import("./domain/material/gateway/api.ts")
    const users = await import("./domain/user/gateway/api.ts")

    return {
        Materials: new material.MaterialGateway(Source.From("./materials", Source.server)),
        Users: new users.UserGateway(Source.From("./users", Source.server))
    }
}