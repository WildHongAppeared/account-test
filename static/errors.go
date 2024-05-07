package static

const (
	EmptyPort           = "PORT cannot be empty"
	ErrUnableToReadBody = "Failed to read request body"

	// Business Logic Specific Error - Account
	ErrAccountAlreadyExist     = "Account already exist"
	ErrAccountDoesNotExist     = "Account does not exist"
	ErrBalanceNotValidNumber   = "initial_balance is not a valid number"
	ErrBalanceCannotBeNegative = "initial_balance cannot be a negative number"
	ErrBalanceTooLarge         = "initial_balance value is too large"
	ErrCreatingAccount         = "Error creating new account"
	ErrIDLengthCannotBeZero    = "ID must be at least one character long"
	ErrIDLengthTooLong         = "ID length must be not be longer than 32 characters"
	ErrUnableToRetrieveAccount = "Error retrieving account balance"

	//Business Logic Specific Error - Transaction
	ErrSourceAccountDoesNotExist       = "Source account does not exist"
	ErrDestinationAccountDoesNotExist  = "Destination account does not exist"
	ErrSourceDestinationSame           = "Source account and destination account cannot be the same"
	ErrAmountNotValidNumber            = "amount is not a valid number"
	ErrAmountCannotBeNegative          = "amount cannot be a negative number"
	ErrAmountTooLarge                  = "amount value is too large"
	ErrGetSourceAccount                = "Error retrieving source account"
	ErrGetDestinationAccount           = "Error retrieving destination account"
	ErrTransferAmountLargerThanAccount = "amount cannot be larger than source account's balance"
	ErrUnableToCompleteTransaction     = "Error - unable to complete transaction"
)

var ()
