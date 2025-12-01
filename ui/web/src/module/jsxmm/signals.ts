type ReadFunction<T> = () => T
type WriteFunction<T> = (v: T) => void
type ModifyFunction<T> = (func: (prev: T) => T) => void
type DisposeFunction = () => void

namespace Signal {
	export class Value<T> {
		Read: ReadFunction<T>
		Write: WriteFunction<T>
		Modify: ModifyFunction<T>

		constructor(value: T, equal_func: EqualFunction<T> = triple_equal) {
			const s = new valued(value, equal_func)

			this.Read = valued.prototype.Read.bind(s)
			this.Write = valued.prototype.Write.bind(s)
			this.Modify = valued.prototype.Modify.bind(s)
		}
	}

	export interface RValue<T> {
		Read: ReadFunction<T>
	}

	export interface WValue<T> {
		Write: WriteFunction<T>
		Modify: ModifyFunction<T>
	}

	export class Compute<T> {
		Read: ReadFunction<T>
		Dispose: DisposeFunction

		constructor(value_func: () => T, equal_func: EqualFunction<T> = triple_equal) {
			const d = new computed(value_func, equal_func)

			this.Read = computed.prototype.Read.bind(d)
			this.Dispose = computed.prototype.Dispose.bind(d)
		}
	}

	export class Effect {
		Dispose: DisposeFunction

		constructor(callback: () => void) {
			const e = new effect(callback)

			this.Dispose = effect.prototype.Dispose.bind(e)
		}
	}
}

export default Signal

interface Source {
	Attach(dep: Sink): void
	Detach(dep: Sink): void
}

interface Sink {
	Execute(): void
	Acknowledge(sub: Source): void
}

type EqualFunction<T> = (o0: T, o1: T) => boolean

let active_observers: Sink | undefined = undefined
const NotSet = Symbol()

class valued<T> {
	#observers: Set<Sink>
	#equal_func: EqualFunction<T>
	#value: T

	constructor(value: T, equal_func: EqualFunction<T>) {
		this.#equal_func = equal_func
		this.#observers = new Set<Sink>()
		this.#value = value
	}

	Read(): T {
		if (active_observers !== undefined) {
			this.#observers.add(active_observers)
			active_observers.Acknowledge(this)
		}

		return this.#value
	}

	Write(value: T): void {
		if (this.#equal_func(value, this.#value)) {
			return
		}

		this.#value = value
		this.Notify()
	}

	Modify(func: (prev: T) => T): void {
		const new_value = func(this.#value)
		if (this.#equal_func(new_value, this.#value)) {
			return
		}

		this.#value = new_value
		this.Notify()
	}

	Attach(dep: Sink): void {
		this.#observers.add(dep)
	}

	Detach(dep: Sink): void {
		this.#observers.delete(dep)
	}

	Notify(): void {
		for (const observer of [...this.#observers]) {
			observer.Execute()
		}
	}
}

class computed<T> {
	#subjects: Set<Source>
	#observers: Set<Sink>
	#equal_func: EqualFunction<T>
	#value_func: () => T
	#value: T

	constructor(value_func: () => T, equal_func: EqualFunction<T> = triple_equal<T>) {
		this.#equal_func = equal_func
		this.#subjects = new Set<Source>()
		this.#observers = new Set<Sink>()
		this.#value_func = value_func

		this.#value = NotSet as any
		this.Execute()
	}

	Read(): T {
		if (active_observers !== undefined) {
			this.#observers.add(active_observers)
			active_observers.Acknowledge(this)
		}

		return this.#value
	}

	Acknowledge(sub: Source): void {
		this.#subjects.add(sub)
	}

	Execute(): void {
		this.Dispose()
		active_observers = this

		const new_value = this.#value_func()
		if (!this.#equal_func(this.#value, new_value)) {
			this.#value = new_value
			this.Notify()
		}

		active_observers = undefined
	}

	Dispose(): void {
		for (const subject of this.#subjects) {
			subject.Detach(this)
		}
		this.#subjects.clear()
	}

	Attach(dep: Sink): void {
		this.#observers.add(dep)
	}

	Detach(dep: Sink): void {
		this.#observers.delete(dep)
	}

	Notify(): void {
		for (const observer of Array.from(this.#observers)) {
			observer.Execute()
		}
	}
}

class effect {
	#effect: () => void
	#subjects: Set<Source>

	constructor(effect: () => void) {
		this.#subjects = new Set<Source>()
		this.#effect = effect

		this.Execute()
	}

	Acknowledge(sub: Source): void {
		this.#subjects.add(sub)
	}

	Execute(): void {
		this.Dispose()
		active_observers = this
		this.#effect()
		active_observers = undefined
	}

	Dispose(): void {
		for (const subject of this.#subjects) {
			subject.Detach(this)
		}
		this.#subjects.clear()
	}
}

function triple_equal<T>(o0: T, o1: T): boolean {
	return o0 === o1
}
