export class APIError {
    kind: string
    title: string
    message: string
    cause: Cause | undefined

    static New(kind: string, title: string, message: string, cause?: Cause): APIError {
        const error = new APIError()

        error.kind = kind
        error.title = title
        error.message = message
        error.cause = cause

        return error
    }

    static FromObject(object: any): APIError | null {
        const error = new APIError()

        if (typeof object.kind === "string") {
            error.kind = object.kind
        } else {
            return null
        }

        if (typeof object.title === "string") {
            error.title = object.title
        } else {
            return null
        }

        if (typeof object.message === "string") {
            error.message = object.message
        } else {
            return null
        }

        const cause = object.cause
        if (cause !== null && cause !== undefined) {
            const type = typeof cause

            switch (true) {
            case Array.isArray(cause):
                error.cause = []
                for (let i = 0; i < cause.length; i++) {
                    const suberror = this.FromObject(cause[i])
                    if (suberror !== null) {
                        error.cause.push(suberror)
                    }
                }
                break

            case type === "object":
                const suberror = this.FromObject(cause)
                if (suberror !== null) {
                    error.cause = suberror
                }
                break

            case type === "string":
                error.cause = cause
                break
            }
        }

        return error
    }

    private constructor() {
        this.kind = undefined as any
        this.title = undefined as any
        this.message = undefined as any
    }
}

type Cause = APIError | string | Cause[]
