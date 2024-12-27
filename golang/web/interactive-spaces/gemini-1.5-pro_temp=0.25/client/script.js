const spaceList = document.getElementById('space-list');
const newSpaceButton = document.getElementById('new-space');

function fetchSpaces() {
    fetch('/spaces') // We'll implement this server route next
        .then(response => response.json())
        .then(spaces => {
            spaceList.innerHTML = ''; // Clear existing list
            spaces.forEach(space => {
                const listItem = document.createElement('li');
                // For now, just display the name and ID. We'll add more details later.
                listItem.textContent = `${space.name} (${space.id})`;
                spaceList.appendChild(listItem);
            });
        })
        .catch(error => {
            alert('Error fetching spaces: ' + error);
            console.error(error);
        });
}

newSpaceButton.addEventListener('click', () => {
    fetch('/spaces/new', { method: 'POST' })
        .then(response => {
            if (response.ok) {
                // Redirect to the new space page (we'll implement this later)
                window.location.href = response.url;
            } else {
                alert('Error creating space');
                console.error('Error creating space:', response.status, response.statusText);
            }
        })
        .catch(error => {
            alert('Error creating space: ' + error);
            console.error(error);
        });
});

fetchSpaces(); // Initial fetch

const eventSource = new EventSource('/spaces/events');

eventSource.onmessage = event => {
    const spaces = JSON.parse(event.data);
    // Update the space list in the UI (same logic as in fetchSpaces)
    spaceList.innerHTML = '';
    spaces.forEach(space => {
        const listItem = document.createElement('li');
        listItem.textContent = `${space.name} (${space.id})`;
        spaceList.appendChild(listItem);
    });
};

eventSource.onerror = error => {
    console.error('SSE Error:', error);
    alert('Error connecting to space updates. Please refresh the page.');
};
