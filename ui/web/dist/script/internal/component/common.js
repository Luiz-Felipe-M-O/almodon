export function ClickForKeys(element, ...keys) {
    element.addEventListener("keydown", function (evt) {
        if (keys.includes(evt.key)) {
            evt.target.click();
            evt.preventDefault();
        }
    });
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tbW9uLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2ludGVybmFsL2NvbXBvbmVudC9jb21tb24udHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsTUFBTSxVQUFVLFlBQVksQ0FBQyxPQUFvQixFQUFFLEdBQUcsSUFBYztJQUNoRSxPQUFPLENBQUMsZ0JBQWdCLENBQUMsU0FBUyxFQUFFLFVBQVUsR0FBa0I7UUFDNUQsSUFBSSxJQUFJLENBQUMsUUFBUSxDQUFDLEdBQUcsQ0FBQyxHQUFHLENBQUMsRUFBRSxDQUFDO1lBQ3hCLEdBQUcsQ0FBQyxNQUFzQixDQUFDLEtBQUssRUFBRSxDQUFBO1lBQ25DLEdBQUcsQ0FBQyxjQUFjLEVBQUUsQ0FBQTtRQUN4QixDQUFDO0lBQ0wsQ0FBQyxDQUFDLENBQUE7QUFDTixDQUFDIn0=