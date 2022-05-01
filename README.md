# StakeBackendGoTest

Welcome to the Stake Backend Developer Go Test. Thanks for taking the time to complete this, we're keen to see how you do!

## 1. Introduction
* This repo contains a Web application with a single endpoint: http://localhost:8080/api/equityPositions?token=
* The goal of this task is to complete the functionality required in order for the application to return an equity positions response in the format required by Stake's frontend application. This includes a transformation from a third party's data model into Stake's required response format, including some calculations. See Part 3 for everything that's required.
* Note: the gin web framework is referenced in the repo, however feel free to subsitute it for your preferred web framework library. Keen to hear your reasoning behind picking another lib too!

## 2. Getting Started
- Clone this repository:
	- `git clone https://github.com/stake-test/StakeBackendGoTest_XX.git` (update 'XX' to your initials)
- This test has been written using go version 1.18.1 in the GoLand IDE.
- Create a new branch with your fullname in Pascal case (e.g., `AndrewVassili`) - you'll need this to raise a PR against `master` (see step 4 below)

## 3. What you need to do
What we would like you to do is write the remaining code in order to complete this endpoint. This includes the following:
1. You'll notice that most of the files are in the `internal` directory - create a series of folders within the `internal` directory to improve the overall project structure. Keen to hear your reasoning behind setting up your project structure as you do.
2. Transform the `mockdata/mockpositions.json`/`mockdata/mockprices.json` to `response/positions.go`. Note that the `mockpositions.go` and `mockprices.go` are designed to mimic a response obtained from a third party broker source, whereas the `positions.go` is the response required by Stake's Frontend applications (ie, app and web).
3. As part of step 1's transformation process, you'll have to calculate the following values: `dayProfitOrLoss`, `dayProfitOrLossPercent`, `totalProfitOrLoss`, `totalProfitOrLossPercent`
4. Add sufficient logging at appropriate verbosity levels
5. Handle HTTP error responses
6. Add unit and integration tests as you see relevant
7. [Stretch Goal] Identify the security flaw in the API endpoint's design and fix it

In order for the API request to work, you'll need to pass through the following user token as a query parameter: `fJCoxhq8uR9GiUIgaIGfMgw7zCqxwDhQ`.

## 4. Submitting your work
Once you have performed the above to the best of your ability:
1. Run the application (noting the web service runs on port 8080), hit the endpoint, take a screenshot and include it in the repo at the top level directory
2. Please raise a PR against the `master` branch
