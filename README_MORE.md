# Project structure
```
.
├── CHANGELOG.md
├── Dockerfile
├── Makefile
├── README.md
├── README_MORE.md
├── api
│   ├── api_handler.go
│   ├── error_code.go
│   ├── form
│   ├── middleware
│   │   ├── auth.go
│   │   └── trace.go
│   ├── response
│   │   ├── error_code.go
│   │   ├── position.go
│   │   └── response.go
│   └── router.go
├── build
├── cmd
│   └── service
│       ├── main.go
│       ├── main_test.go
│       └── service
├── configs
│   └── config.go
├── controller
│   └── engine.go
├── docs
├── etc
│   └── settings.yaml
├── example
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   ├── adapter.go
│   │   ├── adapter_manager.go
│   │   ├── mock_adapter.go
│   │   ├── mock_adapter_test.go
│   │   └── mockdata
│   │       ├── mockpositions.json
│   │       └── mockprices.json
│   ├── entity
│   │   ├── README.md
│   │   ├── mock
│   │   │   ├── position.go
│   │   │   └── price.go
│   │   └── stake
│   │       ├── position.go
│   │       └── price.go
│   ├── model
│   │   ├── data_manager.go
│   │   └── data_manager_test.go
│   └── pkg
│       └── converters
│           ├── README.md
│           └── mock
│               ├── converter.go
│               └── converter_test.go
├── pkg
│   ├── const
│   │   └── def.go
│   └── log
│       └── log.go
├── scripts
└── test
```

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
- API, responsible for setting routers, format input request, response messages, etc.
- Controller, responsible for starting http server, creating AdapterManager, collecting data from AdapterManager, and managing their life.
- Internal, do something that is transparent to clients
    - DataManager, responsible for data management (calculation, merging, etc)
    - AdapterManager, responsible for specific adapters management
        - MockerAdapter, an adapter for demonstration use
       - ...
### Float calculation/display
- The built-in `float` or `double` are not suitable for currency calculation as the known precision issue in computer, so `shopstring.decimal` is used.
- In crypto, 1e-8 is referred as a `Satoshi`, the smallest unit, so we use 8 decimal places for price, volume, and value display though it's integer in Mock for volume.

### Error handling
- If a function consumes errors it gets from its callee, it should suppress them locally with logging; if it forwards errors to its caller, it should not do logging.

### Authentication
- The system uses Gin's middleware to authenticate users. The token is put in the url as required. For security, putting it in the request header probably is better.
- The token value is hard coded. We could use other solutions, e.g., jwt-go, to generate dynamic tokens.

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
- We get the current volume in the positionWith from `PortfolioUnits` in mock Position.
- 2 decimal place are used for percent showning, so if a value in terms of percentage is, e.g., 10.12%, the shown value is "10.12".
## Field mapping
- `averagePrice` in `equityPositions` are assumed to be the open price of the positions
- Use field definitions in json file, so `backOfficeAvailableUnits` and `backOfficePortfolioUnits` are used, `availableUnits` and `portfolioUnits` are changed accordingly
## Position identification
We assume there is only one account, so no account/user information for positions we have. It's easy to scale to support multiple accounts as it mentions in data_manager.go

# More to do
- There are TODOs which mention they are not finised yet but are on the radar
- Fill Dockerfile
- Create a CI file (say, CircleCI) so that the repo could be managed by CI for deployment, testing, etc.
- Fill Makefile to help users/CI build/test/deploy the system
- The system returns the response according to the format required. There is response definition in `response.go` with `err_code` and `err_message` to have a more friendly message for clients.