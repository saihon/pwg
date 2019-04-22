package pwg

// Options for password,username generation
type Options struct {
	Length     int    // Password length
	Generate   int    // Number of passwords to generate
	All        bool   // Use all kinds of chars (same as -Llns)
	Evenly     bool   // Use character types as evenly as possible
	LowerCase  bool   // Use lower case letters
	UpperCase  bool   // Use upper case letters
	Number     bool   // Use numbers in passwords
	Symbol     bool   // Use symbols in passwords
	Custom     string // Specify any character string to be used for password
	Random     bool   // Random username length
	Capitalize bool   // Capitalize the beginning of a user name
}
