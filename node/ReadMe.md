# Project info

## Errors policy

### Don't use errors for buisness logic

At least from the bottom to service layer. It's forbidden to pass useful info in errors, use errors.As, errors.Is for your error processing. Error is just a text.

### All errors must contains callstack

So when pass errors from layer to layer you don't need to add caller function
to error, you should add only meaningful context (about params, for example).