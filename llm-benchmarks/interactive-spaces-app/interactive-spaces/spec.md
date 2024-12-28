Implement a web app that allows users to create spaces that stream updates to many users.
After the space owner "starts" the space, it begins to "evolve", adding one new "active cell" each second.
When a new client visit the space page, it gets the most recent state of this space and also subscribes to space updates. The idea is that all clients subscribed to a given space see all evolution happening to this space simultaneously and in a real time.

## General Considerations

- The app doesn't need to be deployed yet: both the client and the server can run on a localhost initially.
- The backend controls the eveolution state and logic of a space.
- When evolution of the space is finished (i.e., all cells in the space have finished transferring to an "active" state), the space transfer to "evolution completed" state and we stop.

## Client

General:

- Use Vanilla JS/HTML/CSS and Web Components for implementation
- Use toasts for displaying errors to the client.

Index /spaces page:

- Shows the list of existing spaces with a link to that space. The information about the space includes its name, creation time, running time, number of clients connected, and number of "active cells" in this space
- "New space" button that sends a request to /spaces/new
- Various UI elements that improve UX
- Connects to /spaces/events server-side events (SSE) endpoint to display updates to a space list in real time (new spaces added, space info updated, etc.)

Space /spaces/{id} page:

- If the space was just created (a first visit to a space by its author), displays the space password (only one time). We use a modal window for this, which warns the user that the password will be displayed only once and suggests to copy it using the "Copy" button.
- Allows navigation back to the index
- A 100x100 grid, scrollable. Uses 5mm cell size
- When space is in "started" state, the cells in the grid are constantly "evolving", adding a new "active cell" approximately each second or so.
- "inactive" and "active" cells have a different appearance (transparent and colored respectively)
- When teh space is in "stopped" state, no evolution happens
- For simplicity, the evolution happens to cells sequentially, starting from the top left corner and going right. When reaching the end of the line, it moves to the next line in the grid and goes to the right again.
- Connects to /spaces/{id}/events SSE endpoint to receive updates to the grid (activated cells) in real time
- Has "Unlock" button that shows a modal window accepting the space password and unlocks the space for modifications if the password is correct
- Modifications include editing the space name and stop/pausing the space evolution (activating new cells)

## Backend

General:

- Use Golang for implementation. The usage of 3rd-party libraries is allowed where it makes sense.
- Store space data in a local JSON file. Ensure correct file handling with concurrent writes.
- Use structured logging to the console and file for logs.

Routes:

/

- Returns the html template with spaces

/spaces/events

- Subscribes the client to updates to a space list

POST /spaces/new

- Creates a new space and a password for it. The password should be stored securely (hash + salt), and in a separate file. We show the password in a clear text only once, when the space is created, and offer the client to copy and save it.
- Redirects to the newly created space at /spaces/{id} with password included, so that the newly created space is unlocked initially.
- The new space is in unlocked and "not started" state, so the user can edit its name and "start" the space by posting to `POST /spaces/{id}` route

GET /spaces/{id}

- Gets the current state of a space.

POST /spaces/{id}

- Updates the space at id {id}
- Propagates the update made to the space to all clients subscribed to this /spaces/{id}/events and /spaces/events using SSE

GET /spaces/{id}/events

- Subscribes a client to space updates using SSE.

GET /spaces/events

- Subscribes a client to the general updates made to any of the existing spaces. This is needed to update the content on the index page dynamically, like a new space added or a space state has changed (e.g., space name, space state — started/stopped, number of active cells in the space, evolution round, etc.).

## Process

General:

- Minimize friendliness and useless chatting. Stick to the business to save input and output tokens.
- Ask clarifying questions if you find gaps in the design or an implementation. This is valid at the start of the process and on every stage.

### Initial phase

- Start with outlining the general design and plan for implementation.
- Outline the proposed project/repository structure. i.e., the files you want to create, the tests you'd want to create, etc.
- Outline the testing strategy that will ensure the correctness of the implementation.
- Do not start actual implementation before clarifying any outstanding questions and receiving my confirmation to proceed further.

### Next phase

- After receiving the confirmation for starting the implementation, propose the step or the file you'd like to implement next.
- Implement only one source file at a time. Two if the 2nd file is an accompanying test file.
- Proceed to implementing the next file only after receiving my confirmation that the previous one worked as intended. We first fix all issues and only then move to the next step

