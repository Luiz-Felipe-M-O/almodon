export const Unlogged = "unlogged" as const
export const Chief = "chief" as const
export const Promoted = "promoted-admin" as const
export const Admin = "admin" as const
export const User = "user" as const

export type Role =
    | typeof Unlogged
    | typeof Chief
    | typeof Promoted
    | typeof Admin
    | typeof User

export type Hierarchy = (r0: Role, r1: Role) => boolean

export class Permission {
    #classes: Role[]
    #hierarchy: Hierarchy

    constructor(classes: Role[], hierarchy: Hierarchy = DefaultHierarchy) {
        this.#classes = classes
        this.#hierarchy = hierarchy
    }

    Authorize(role: Role): boolean {
        for (const clazz of this.#classes) {
            if (this.#hierarchy(clazz, role)) {
                return true
            }
        }

        return false
    }

    String(): string {
        return JSON.stringify(this.#classes)
    }
}

const RoleOrder = [User, Admin, Promoted, Chief]

export function DefaultHierarchy(r0: Role, r1: Role): boolean {
    if (r0 === Unlogged) {
        return true
    }

    if (r1 === Unlogged) {
        return false
    }

    return RoleOrder.indexOf(r0) >= RoleOrder.indexOf(r1)
}

export function RoleFromString(string: string): Role | undefined {
    return RoleOrder.find(r => r === string)
}

export function IsRoleString(value: string): value is Role {
    return RoleFromString(value) !== undefined
}

const RoleString: Record<Role, string> = {
    [Unlogged]: "não logado",
    [User]: "usuário",
    [Admin]: "administrador",
    [Promoted]: "promovido",
    [Chief]: "chefe",
}

export function RoleToString(role: Role): string {
    return RoleString[role]
}
