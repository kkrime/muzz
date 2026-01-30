# Muzz Backend Engineering Assignment
This repository is the solution - [**this**](https://github.com/kkrime/muzz/blob/main/Muzz-Go_Backend_Exercise.pdf) is the problem.

## Getting Started
Everything is containerized using docker-compose, so from the root of the repoistory you can run the db and backend by;
`docker-compose up`<br>
This should build the db and the backend docker images and run them<br>

The db image includes the schema, so when you first run the db imgage it will set everything up for the db<p>

I did not set an expiary for the jtw token.

### Design
`/discover` fills all the criteria including returning results on order of attractiveness and filters by gender and age; however, the age it filters at are hard coded to between 21 - 35 inclusive.<br>
In a real life system there would be a `user_prefrences` table where such values would be pulled from, but for the sake of demontrating how filtering by age can be done, I've used hard coded values.
<br><br>
I intentionally over engineered `/login` slightly just to use more of Gos features to demonstrate my understanding of them, I've added a comment about this in the code.

### Testing
#### Unit Tests
You can find unit tests in `internal/service` from inside that directory you can run the UTs; `go test .`.<br>
These tests are just here to deomonstrate I can write UTs and I understand go conventions around UTs.
#### End To End Tests
I've included a script to help you test out my code, it's in the `demo` folder in the root.<br>
This is a simple `python` script and does the following;
- creates 10 users
- logs in the 10 users
- all 10 users swipe on each (randomly right or left), any matches will be displayed

This script is more for demonstatative purposes than actual e2e tests, although it does serves that purpose.

## To Run The Demo
1. start the main service using `docker-compose up`
2. from inside the `demo` folder, run: `docker build . -t demo && docker run --network host -t demo`
