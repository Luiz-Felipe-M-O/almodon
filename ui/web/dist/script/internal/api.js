import Source from "./support/source.js";
export async function Construct() {
    if (Source.server === "") {
        return await MockConstruct();
    }
    else {
        return await APIConstruct();
    }
}
export async function MockConstruct() {
    const users = await import("./domain/user/gateway/mock.js");
    return {
        Users: new users.UserGateway()
    };
}
export async function APIConstruct() {
    const users = await import("./domain/user/gateway/api.js");
    return {
        Users: new users.UserGateway(Source.From("./users", Source.server))
    };
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2ludGVybmFsL2FwaS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxPQUFPLE1BQU0sTUFBTSxxQkFBcUIsQ0FBQTtBQU14QyxNQUFNLENBQUMsS0FBSyxVQUFVLFNBQVM7SUFDM0IsSUFBSSxNQUFNLENBQUMsTUFBTSxLQUFLLEVBQUUsRUFBRSxDQUFDO1FBQ3ZCLE9BQU8sTUFBTSxhQUFhLEVBQUUsQ0FBQTtJQUNoQyxDQUFDO1NBQU0sQ0FBQztRQUNKLE9BQU8sTUFBTSxZQUFZLEVBQUUsQ0FBQTtJQUMvQixDQUFDO0FBQ0wsQ0FBQztBQUVELE1BQU0sQ0FBQyxLQUFLLFVBQVUsYUFBYTtJQUMvQixNQUFNLEtBQUssR0FBRyxNQUFNLE1BQU0sQ0FBQywrQkFBK0IsQ0FBQyxDQUFBO0lBRTNELE9BQU87UUFDSCxLQUFLLEVBQUUsSUFBSSxLQUFLLENBQUMsV0FBVyxFQUFFO0tBQ2pDLENBQUE7QUFDTCxDQUFDO0FBRUQsTUFBTSxDQUFDLEtBQUssVUFBVSxZQUFZO0lBQzlCLE1BQU0sS0FBSyxHQUFHLE1BQU0sTUFBTSxDQUFDLDhCQUE4QixDQUFDLENBQUE7SUFFMUQsT0FBTztRQUNILEtBQUssRUFBRSxJQUFJLEtBQUssQ0FBQyxXQUFXLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxTQUFTLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxDQUFDO0tBQ3RFLENBQUE7QUFDTCxDQUFDIn0=