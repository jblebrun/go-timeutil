# go-timeutil
Utilities to make it easier to safely use timers and test timer-based code

* duration.go: Duration rounding and unmarshaling
* helpers.go: Safe timer stop and reset based on best practices from documentation
* module.go: Wraps timer in code to allow for choreographed tests
* span.go: Represents a timespan via a start time and duration pair
