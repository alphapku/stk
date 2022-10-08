# Project structure
The task was done from the point of view of implementing a system which could be scaled or maintaned earsier, instead of only for the purpose of finishing the task. So some empty directories were created though probably there are no files in them.
  - docs
  - etc
  - example
  - scripts
  - test

# Design
## Architecture
### Key Components
Here are the key components of the system
- AddRouters, responsible for setting routers, all specific will be done by DataManager at present as our system is a data service based system.
- Engine, responsible for start http server, create AdapterManager, collect data from AdapterManager and their life management
- DataManager, responsible for data management (calculation, merging, etc)
- AdapterManager, responsible for specific adapters management
  - MockerAdapter, an adapter for demonstration use
  - ...

# Test
 - Try to cover corner cases
 - No test cases for mocker JSON parsing in `mock_adapter_test.go`, but for real adapters, data parsing has to be verified

# Documentation
 To make it eary to understand the system, not all documentation work are done in the root level README.md. There are some README.md files in sub-directories for details.

# Arguments
Here list items that could be done in other ways instead of ways used in the system while trying to give some explanations.
## Position Theoretical Value Calcuation
- When calculating PNL, there are a couple of ways to calculate the value of the positions
    - Use the mid-price of the bid and the ask
    - Use the worse price: If you are long, the bid is used; if you are short, the ask is used
    - Use the last traded price. In our system, this is used
## Field mapping
- `averagePrice` in `equityPositions` are assumed to be the open price of the positions
- Use field definitions in json file, so `backOfficeAvailableUnits` and `backOfficePortfolioUnits` are used, `availableUnits` and `portfolioUnits` are changed accordingly


# More to do
- There are TODOs which mention they are not finised yet but are on the radar
- Fill Dockerfile
- Create a CI file (say, CircleCI) so that the repo could be managed by CI for deployment, testing, etc.
- Fill Makefile to help users/CI build/test/deploy the system