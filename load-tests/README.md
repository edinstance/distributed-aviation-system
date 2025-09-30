# Load Testing Documentation

## Overview

This load testing suite is designed to test the distributed aviation system's performance using k6, a modern load
testing tool. It tests both individual microservices (Aircraft and Flight services) and the complete system through the
GraphQL router.

## Architecture

The load testing framework consists of:

- **k6**: Load testing runtime
- **GraphQL Code Generation**: Type-safe GraphQL operations
- **TypeScript**: Type safety and modern JavaScript features
- **ESBuild**: Fast bundling for k6 execution

## Installation

To run the tests, you will need to install both k6 and the node dependencies. To install k6, follow the instructions
[here](https://grafana.com/docs/k6/latest/set-up/install-k6/). Once that is complete, install the node dependencies by
running ```npm install```.

## Configuration

### Environment Variables

The tests can be configured using environment variables, they are all defined in [this](./src/config.ts) config file. They can be set in the following ways:

1. By exporting the values before running the tests.
```
export FLIGHT_URL=http://localhost:8081/graphql
export AIRCRAFT_URL=http://localhost:8082/graphql
export ROUTER_URL=http://localhost:4000/graphql
```

2. Or you can add a -e flag to the run configurations. 

```
-e FLIGHT_URL=http://localhost:8081/graphql
```

## GraphQL Integration

### Code Generation

GraphQL operations are type-safe using `@graphql-codegen/cli`:

```bash
npx graphql-codegen 
```

### Request Helper

The `graphql()` helper function (`helpers/graphql_request.ts`) provides:

- Type-safe GraphQL requests
- Automatic response validation
- Error handling and logging
- k6 check assertions for:
    - HTTP 200 status
    - Valid JSON response
    - No GraphQL errors

## Running Tests

### Prerequisites

1. Ensure services are running:
    - Aircraft Service (port 8080)
    - Flight Service (port 8081)
    - Apollo Router (port 4000)

2. Install dependencies:
   ```bash
   npm install
   ```

### Individual Service Tests

```bash
# Test Aircraft service only
npm run test:aircraft

# Test Flight service only
npm run test:flights

# Test through Apollo Router
npm run test:router
```

### All Tests

```bash
# Run all tests sequentially
npm run test:all

# Run aircraft and flight tests in parallel
npm run test:all:parallel
```

### Manual k6 Execution

```bash
# Build first
npm run build

# Run specific test with custom options
k6 run dist/aircraft.js --vus 10 --duration 30s

# Run with environment variables
AIRCRAFT_URL=http://prod.example.com/graphql k6 run dist/aircraft.js
```

## Monitoring and Validation

### Built-in Checks

All tests include k6 checks for:

- HTTP response codes (200)
- JSON response format
- GraphQL error absence
- Business logic validation (e.g., "aircraft created")

### Output

k6 provides detailed metrics including:

- Request rates and response times
- Error rates and types
- Custom check success rates
- Resource utilization

## Development

### Adding New Tests

1. Create GraphQL operations in `queries/` or `mutations/`
2. Run ` npx graphql-codegen` to generate types
3. Create test scenario in `tests/`
4. Add npm script to `package.json`

### Modifying Load Profiles

Update the `options` export in test files:

```typescript
export let options: Options = {
    stages: [
        {duration: "2m", target: 10},  // Ramp up
        {duration: "5m", target: 10},  // Stay at 10 users
        {duration: "2m", target: 0},   // Ramp down
    ],
};
```

### Code Quality

- **Linting**: `npm run lint`
- **Formatting**: `npm run prettier`
- **Type Checking**: Automatic via TypeScript compilation

## Troubleshooting

### Common Issues

1. **GraphQL Schema Mismatch**: Regenerate types with `npx graphql-codegen`
2. **Service Connection**: Check service URLs and availability
3. **Build Failures**: Verify TypeScript compilation
4. **Test Failures**: Check GraphQL errors in k6 output

### Debug Mode

Add logging to test scenarios:

```typescript
console.log("Response:", JSON.stringify(response, null, 2));
```

Run with verbose k6 output:

```bash
k6 run dist/aircraft.js --verbose
```