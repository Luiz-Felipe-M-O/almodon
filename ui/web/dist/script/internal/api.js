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
    throw new Error("not implemented");
    const users = await import("./domain/user/gateway/mock.js");
    return {
        Materials: null,
        Users: new users.UserGateway()
    };
}
export async function APIConstruct() {
    const material = await import("./domain/material/gateway/api.js");
    const users = await import("./domain/user/gateway/api.js");
    return {
        Materials: new material.MaterialGateway(Source.From("./materials", Source.server)),
        Users: new users.UserGateway(Source.From("./users", Source.server))
    };
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYXBpLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vc3JjL2ludGVybmFsL2FwaS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxPQUFPLE1BQU0sTUFBTSxxQkFBcUIsQ0FBQTtBQU94QyxNQUFNLENBQUMsS0FBSyxVQUFVLFNBQVM7SUFDM0IsSUFBSSxNQUFNLENBQUMsTUFBTSxLQUFLLEVBQUUsRUFBRSxDQUFDO1FBQ3ZCLE9BQU8sTUFBTSxhQUFhLEVBQUUsQ0FBQTtJQUNoQyxDQUFDO1NBQU0sQ0FBQztRQUNKLE9BQU8sTUFBTSxZQUFZLEVBQUUsQ0FBQTtJQUMvQixDQUFDO0FBQ0wsQ0FBQztBQUVELE1BQU0sQ0FBQyxLQUFLLFVBQVUsYUFBYTtJQUMvQixNQUFNLElBQUksS0FBSyxDQUFDLGlCQUFpQixDQUFDLENBQUE7SUFHbEMsTUFBTSxLQUFLLEdBQUcsTUFBTSxNQUFNLENBQUMsK0JBQStCLENBQUMsQ0FBQTtJQUUzRCxPQUFPO1FBQ0gsU0FBUyxFQUFxQyxJQUFXO1FBQ3pELEtBQUssRUFBRSxJQUFJLEtBQUssQ0FBQyxXQUFXLEVBQUU7S0FDakMsQ0FBQTtBQUNMLENBQUM7QUFFRCxNQUFNLENBQUMsS0FBSyxVQUFVLFlBQVk7SUFDOUIsTUFBTSxRQUFRLEdBQUcsTUFBTSxNQUFNLENBQUMsa0NBQWtDLENBQUMsQ0FBQTtJQUNqRSxNQUFNLEtBQUssR0FBRyxNQUFNLE1BQU0sQ0FBQyw4QkFBOEIsQ0FBQyxDQUFBO0lBRTFELE9BQU87UUFDSCxTQUFTLEVBQUUsSUFBSSxRQUFRLENBQUMsZUFBZSxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsYUFBYSxFQUFFLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQztRQUNsRixLQUFLLEVBQUUsSUFBSSxLQUFLLENBQUMsV0FBVyxDQUFDLE1BQU0sQ0FBQyxJQUFJLENBQUMsU0FBUyxFQUFFLE1BQU0sQ0FBQyxNQUFNLENBQUMsQ0FBQztLQUN0RSxDQUFBO0FBQ0wsQ0FBQyJ9