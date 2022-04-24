package lambda

// LambdaExecutionEnvironment is used to determine if code is running
// inside of AWS or in some other fashion (e.g., local test, Travis, test, etc)
const LambdaExecutionEnvironment = `AWS_EXECUTION_ENV`
