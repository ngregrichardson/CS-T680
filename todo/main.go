package main

import (
	"drexel.edu/todo/commands"
	"github.com/spf13/cobra"
)

// processCmdLineFlags parses the command line flags for our CLI
//
// TODO: This function uses the flag package to parse the command line
//		 flags.  The flag package is not very flexible and can lead to
//		 some confusing code.

//			 REQUIRED:     Study the code below, and make sure you understand
//						   how it works.  Go online and readup on how the
//						   flag package works.  Then, write a nice comment
//				  		   block to document this function that highights that
//						   you understand how it works.
//
//			 EXTRA CREDIT: The best CLI and command line processor for
//						   go is called Cobra.  Refactor this function to
//						   use it.  See github.com/spf13/cobra for information
//						   on how to use it.
//
//	 YOUR ANSWER: The flag package is used to parse the command line flags. Some global variables are created to hold the values of each flag,
//					the addresses of which are passed into respective flag functions that define the types of each flag. For example, since '-l' is
//					a boolean flag, the address of the global variable 'listFlag' is passed into the 'BoolVar' function. Then, 'flag.Parse' is used to parse
//					the flags. 'flag.Visit' is used to loop over the flags and check which ones are set, and set the appOpt to that flag. This appOpt
//					is used in a switch statement to call the appropriate database function based on which flags were set. The values of each flag can then be
//					gotten by the global variables for each flag. Whenever 'flag.Usage' is called, it prints the usage of the command with the details defined
//					for each flag.
func processCmdLineFlags(cmd *cobra.Command, args []string) int {
	// not used anymore
	return 0
}

// main is the entry point for our todo CLI application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	commands.Execute()
}
