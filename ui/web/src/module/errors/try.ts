import { APIError } from "./error.ts"

type Function = (...args: any) => any

export function Try<T extends Function>(func: T, ...params: Parameters<T>): [ReturnType<T>, null] | [undefined, APIError | Error] {
    try {
        return [func(...params), null]
    } catch (error) {
        return resonable(error)
    }
}

type AsyncFunction = (...args: any) => Promise<any>

export async function AsyncTry<T extends AsyncFunction>(func: T, ...params: Parameters<T>): Promise<[ReturnType<T> extends Promise<infer P> ? P : never, null] | [undefined, APIError | Error]> {
    try {
        return [await func(...params), null]
    } catch (error) {
        return resonable(error)
    }
}

function resonable(error: unknown): [undefined, APIError | Error] {
    if (error instanceof APIError) {
        return [undefined, error]
    }

    if (error instanceof Error) {
        return [undefined, error]
    }

    return [undefined, new Error(String(error))]
}