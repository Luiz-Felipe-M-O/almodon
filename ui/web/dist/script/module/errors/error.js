export class APIError {
    kind;
    title;
    message;
    cause;
    static New(kind, title, message, cause) {
        const error = new APIError();
        error.kind = kind;
        error.title = title;
        error.message = message;
        error.cause = cause;
        return error;
    }
    static FromObject(object) {
        const error = new APIError();
        if (typeof object.kind === "string") {
            error.kind = object.kind;
        }
        else {
            return null;
        }
        if (typeof object.title === "string") {
            error.title = object.title;
        }
        else {
            return null;
        }
        if (typeof object.message === "string") {
            error.message = object.message;
        }
        else {
            return null;
        }
        const cause = object.cause;
        if (cause !== null && cause !== undefined) {
            const type = typeof cause;
            switch (true) {
                case Array.isArray(cause):
                    error.cause = [];
                    for (let i = 0; i < cause.length; i++) {
                        const suberror = this.FromObject(cause[i]);
                        if (suberror !== null) {
                            error.cause.push(suberror);
                        }
                    }
                    break;
                case type === "object":
                    const suberror = this.FromObject(cause);
                    if (suberror !== null) {
                        error.cause = suberror;
                    }
                    break;
                case type === "string":
                    error.cause = cause;
                    break;
            }
        }
        return error;
    }
    constructor() {
        this.kind = undefined;
        this.title = undefined;
        this.message = undefined;
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXJyb3IuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvbW9kdWxlL2Vycm9ycy9lcnJvci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxNQUFNLE9BQU8sUUFBUTtJQUNqQixJQUFJLENBQVE7SUFDWixLQUFLLENBQVE7SUFDYixPQUFPLENBQVE7SUFDZixLQUFLLENBQW1CO0lBRXhCLE1BQU0sQ0FBQyxHQUFHLENBQUMsSUFBWSxFQUFFLEtBQWEsRUFBRSxPQUFlLEVBQUUsS0FBYTtRQUNsRSxNQUFNLEtBQUssR0FBRyxJQUFJLFFBQVEsRUFBRSxDQUFBO1FBRTVCLEtBQUssQ0FBQyxJQUFJLEdBQUcsSUFBSSxDQUFBO1FBQ2pCLEtBQUssQ0FBQyxLQUFLLEdBQUcsS0FBSyxDQUFBO1FBQ25CLEtBQUssQ0FBQyxPQUFPLEdBQUcsT0FBTyxDQUFBO1FBQ3ZCLEtBQUssQ0FBQyxLQUFLLEdBQUcsS0FBSyxDQUFBO1FBRW5CLE9BQU8sS0FBSyxDQUFBO0lBQ2hCLENBQUM7SUFFRCxNQUFNLENBQUMsVUFBVSxDQUFDLE1BQVc7UUFDekIsTUFBTSxLQUFLLEdBQUcsSUFBSSxRQUFRLEVBQUUsQ0FBQTtRQUU1QixJQUFJLE9BQU8sTUFBTSxDQUFDLElBQUksS0FBSyxRQUFRLEVBQUUsQ0FBQztZQUNsQyxLQUFLLENBQUMsSUFBSSxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUE7UUFDNUIsQ0FBQzthQUFNLENBQUM7WUFDSixPQUFPLElBQUksQ0FBQTtRQUNmLENBQUM7UUFFRCxJQUFJLE9BQU8sTUFBTSxDQUFDLEtBQUssS0FBSyxRQUFRLEVBQUUsQ0FBQztZQUNuQyxLQUFLLENBQUMsS0FBSyxHQUFHLE1BQU0sQ0FBQyxLQUFLLENBQUE7UUFDOUIsQ0FBQzthQUFNLENBQUM7WUFDSixPQUFPLElBQUksQ0FBQTtRQUNmLENBQUM7UUFFRCxJQUFJLE9BQU8sTUFBTSxDQUFDLE9BQU8sS0FBSyxRQUFRLEVBQUUsQ0FBQztZQUNyQyxLQUFLLENBQUMsT0FBTyxHQUFHLE1BQU0sQ0FBQyxPQUFPLENBQUE7UUFDbEMsQ0FBQzthQUFNLENBQUM7WUFDSixPQUFPLElBQUksQ0FBQTtRQUNmLENBQUM7UUFFRCxNQUFNLEtBQUssR0FBRyxNQUFNLENBQUMsS0FBSyxDQUFBO1FBQzFCLElBQUksS0FBSyxLQUFLLElBQUksSUFBSSxLQUFLLEtBQUssU0FBUyxFQUFFLENBQUM7WUFDeEMsTUFBTSxJQUFJLEdBQUcsT0FBTyxLQUFLLENBQUE7WUFFekIsUUFBUSxJQUFJLEVBQUUsQ0FBQztnQkFDZixLQUFLLEtBQUssQ0FBQyxPQUFPLENBQUMsS0FBSyxDQUFDO29CQUNyQixLQUFLLENBQUMsS0FBSyxHQUFHLEVBQUUsQ0FBQTtvQkFDaEIsS0FBSyxJQUFJLENBQUMsR0FBRyxDQUFDLEVBQUUsQ0FBQyxHQUFHLEtBQUssQ0FBQyxNQUFNLEVBQUUsQ0FBQyxFQUFFLEVBQUUsQ0FBQzt3QkFDcEMsTUFBTSxRQUFRLEdBQUcsSUFBSSxDQUFDLFVBQVUsQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQTt3QkFDMUMsSUFBSSxRQUFRLEtBQUssSUFBSSxFQUFFLENBQUM7NEJBQ3BCLEtBQUssQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxDQUFBO3dCQUM5QixDQUFDO29CQUNMLENBQUM7b0JBQ0QsTUFBSztnQkFFVCxLQUFLLElBQUksS0FBSyxRQUFRO29CQUNsQixNQUFNLFFBQVEsR0FBRyxJQUFJLENBQUMsVUFBVSxDQUFDLEtBQUssQ0FBQyxDQUFBO29CQUN2QyxJQUFJLFFBQVEsS0FBSyxJQUFJLEVBQUUsQ0FBQzt3QkFDcEIsS0FBSyxDQUFDLEtBQUssR0FBRyxRQUFRLENBQUE7b0JBQzFCLENBQUM7b0JBQ0QsTUFBSztnQkFFVCxLQUFLLElBQUksS0FBSyxRQUFRO29CQUNsQixLQUFLLENBQUMsS0FBSyxHQUFHLEtBQUssQ0FBQTtvQkFDbkIsTUFBSztZQUNULENBQUM7UUFDTCxDQUFDO1FBRUQsT0FBTyxLQUFLLENBQUE7SUFDaEIsQ0FBQztJQUVEO1FBQ0ksSUFBSSxDQUFDLElBQUksR0FBRyxTQUFnQixDQUFBO1FBQzVCLElBQUksQ0FBQyxLQUFLLEdBQUcsU0FBZ0IsQ0FBQTtRQUM3QixJQUFJLENBQUMsT0FBTyxHQUFHLFNBQWdCLENBQUE7SUFDbkMsQ0FBQztDQUNKIn0=