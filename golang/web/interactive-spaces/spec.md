Implement a web app that allows users to create spaces that stream updates to many users.
After the space owner "starts" the space, it begins to "evolve", adding one new "active cell" each second

## Client:

General:

- Use Vanilla JS/HTML/CSS for implementation
- Use Web Components if it makes sense

Index /spaces page:

- Shows the list of existing spaces with a link to that space. The information about the space includes its name, creation time, running time, number of clients connected, and number of "active cells" in this space
- "New space" button that sends a request to /spaces/new
- Various UI elements that improve UX
- Connects to /spaces/events server-side events (SSE) endpoint to display updates to a space list in real time (new spaces added, space info updated, etc.)

Space /spaces/{id} page:

- If the space was just created (a first visit to a space by its author), displays the space password (only one time)
- Allows navigation back to the index
- A 100x100 grid, scrollable. Uses 5mm cell size
- When space is in "started" state, the cells in the grid are constantly "evolving", adding a new "active cell" approximately each second or so.
- "inactive" and "active" cells have a different appearance (transparent and colored respectively)
- When teh space is in "stopped" state, no evolution happens
- For simplicity, the evolution happens to cells sequentially, starting from the top left corner and going right. When reaching the end of the line, it moves to the next line in the grid and goes right again.
- Connects to /spaces/{id}/events SSE endpoint to receive updates to the grid (activated cells) in real time
- Has "Unlock" button that shows a modal window accepting the space password and unlocks the space for modifications if the password is correct
- Modifications include editing the space name and stop/pausing the space evolution (activating new cells)


## Backend:

General:

- This doesn't need to be deployed, — the server can run on localhost
- Use Golang for implementation
- Store space data in a local JSON file

Routes:

/

- Returns the html template with spaces

/spaces/events

- Subscribes the client to updates to a space list

/spaces/new

- Creates a new space and a password for it
- Redirects to the newly created space at /spaces/{id} with password included, so that the newly created space opens unlocked initially
- The new space is in unlocked and "not started" state, so the user can edit its name and "start" the space by posting to `POST /spaces/{id}` route

POST /spaces/{id}

- Updates the space at id {id}
- Propagates the update made to the space to all clients subscribed to this /spaces/{id}/events and /spaces/events using SSE

/spaces/{id}/events

- Propagates the update made to the space to all clients subscribed to this space events

/spaces/events

- Propagates the updates made to all spaces that need to be visible on the index page (e.g., a new space added or a space state has changed, like space has been started or a number of "evolved" cells has changed)

## Process

General:

- Minimize friendliness and useless chatting. Stick to the business to save input and output tokens.
- Ask clarifying questions if you find gaps in the design or an implementation. This is valid at the start of the process and on every stage.

### Initial phase

- Start with outlining the general design and plan for implementation.
- Outline the proposed project/repository structure. i.e., the files you want to create, the tests you'd want to create, etc.
- Outline the testing strategy that will ensure the correctness of the implementation.
- Do not start actual implementation before receiving my agreement and confirmation

### Next phase

- After receiving the confirmation for starting the implementation, propose the step or the file you'd like to implement next.
- Implement only one source file at a time. Two if the 2nd file is an accompanying test file.
- Proceed to implementing the next file only after receiving my confirmation that the previous one worked as intended. We first fix all issues and only then move to the next step

