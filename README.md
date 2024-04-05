# breakstuff

Toy project that is similar to pixel-wars.com

The website returns a grid of cells that are each in an "OK" or "broken" state. Depending on the state that they're in, a different GIF is displayed. Clicking on a given cell alters its state in the DB.

## components

`clausius` contains the lambdas acting as the backend (i.e, treating the reads and writes in DB)

`funes` is the fancy name for the dynamoDB storing the grid state

The `show` folder contains the frontend
