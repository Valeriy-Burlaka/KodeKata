## "Cognitive" abilities:

I was not impressed by the result and I think that arriving to a working implementation would take more effort than if I were implementing the code myself.
What went wrong:
- The model didn't spot potential concurrency issues and didn't offer to use a lock mechanism to work with data files. It also didn't try to organize the code accessing the data files at all (hardcoded file names everywhere, repetitive "load" functions, etc.).
- The model decided to put all JavaScript into a single file, which led to issues later, when we added the 2nd page and the need for "detecting on what page we're on" arised.
- The model ignored the suggestion to use Web Components to organize HTML/JS entirely. This probably has too little weight in its parameters?
- It never proposed to use `gin` or other Go library to handle routes more effectively (instead of manually parsing `spaceId` from the URL).
- It started to produce garbage really quickly. Toward the end of the initial 8000-token window, it already tried to broadcast a single space update but parse it as a _list_ of spaces.

## Operational experience:

The 400 (input limit) errors are a bummer. I don't understand what is the point of the advertised "1M context size" for Gemini models, if they can't handle an input >8000 tokens.
