package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"drexel.edu/todo/db"
)

// Global variables to hold the command line flags to drive the todo CLI
// application
var (
	dbFileNameFlag string
	listFlag       bool
	itemStatusFlag bool
	queryFlag      int
	addFlag        string
	updateFlag     string
	deleteFlag     int
)

type AppOptType int

// To make the code a little more clean we will use the following
// constants as basically enumerations for different options.  This
// allows us to use a switch statement in main to process the command line
// flags effectively
const (
	LIST_DB_ITEM AppOptType = iota
	QUERY_DB_ITEM
	ADD_DB_ITEM
	UPDATE_DB_ITEM
	DELETE_DB_ITEM
	CHANGE_ITEM_STATUS
	NOT_IMPLEMENTED
	INVALID_APP_OPT
)

func contains(opts []AppOptType, toFind AppOptType) bool {
	for _, opt := range opts {
		if opt == toFind {
			return true
		}
	}

	return false
}

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
//	 YOUR ANSWER: This function first sets the flag variables we care about. These are specified
//				  by type, read into a specific variable (passed in as a reference to the var),
//				  contain a default value, and the description of the flag.
//				  Then these are parsed which basically initializes the flag object with these flags.
//				  We then have an AppOptType which is our stand in for an enum used to specify and
//				  return what type of flag the user has passed in. The visit function goes through
//				  each user provided flag in alphabetical order and runs the function passed into it.
//				  Our code essentially only can handle one user passed in flag at a time in its current state.
//				  Trying to pass in multiple flags at once will overwrite our returned enum stand-in with the
//				  last recognized flag value and therefore our main code will only respond to a single flag.
//				  Regardless, the switch block reads the flag name and matches it to the case we care about.
//				  Once matched, the corresponding enum value is assigned to AppOptType. If the
//				  passed in flag is invalid or no flag is passed in, there is a block which returns an error
//				  and prints a warning. If a valid flag is recognized, it returns the enum and a nil error.
func processCmdLineFlags() ([]AppOptType, error) {
	flag.StringVar(&dbFileNameFlag, "db", "./data/todo.json", "Name of the database file")

	flag.BoolVar(&listFlag, "l", false, "List all the items in the database")
	flag.IntVar(&queryFlag, "q", 0, "Query an item in the database")
	flag.StringVar(&addFlag, "a", "", "Add an item to the database")
	flag.StringVar(&updateFlag, "u", "", "Update an item in the database")
	flag.IntVar(&deleteFlag, "d", 0, "Delete an item from the database")
	flag.BoolVar(&itemStatusFlag, "s", false, "Change item 'done' status to true or false")

	flag.Parse()

	var appOpts []AppOptType

	//show help if no flags are set
	if len(os.Args) == 1 {
		flag.Usage()
		return appOpts, errors.New("no flags were set")
	}

	// Loop over the flags and check which ones are set, set appOpt
	// accordingly
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "l":
			appOpts = append(appOpts, LIST_DB_ITEM)
		case "q":
			appOpts = append(appOpts, QUERY_DB_ITEM)
		case "a":
			appOpts = append(appOpts, ADD_DB_ITEM)
		case "u":
			appOpts = append(appOpts, UPDATE_DB_ITEM)
		case "d":
			appOpts = append(appOpts, DELETE_DB_ITEM)

		//TODO: EXTRA CREDIT - Implment the -s flag that changes the
		//done status of an item in the database.  For example -s=true
		//will set the done status for a particular item to true, and
		//-s=false will set the done states for a particular item to
		//false.
		//
		//HINT FOR EXTRA CREDIT
		//Note the -s option also requires an id for the item to that
		//you want to change.  I recommend you use the -q option to
		//specify the item id.  Therefore, the -s option is only valid
		//if the -q option is also set
		case "s":
			//For extra credit you will need to change some things here
			//and also in main under the CHANGE_ITEM_STATUS case
			appOpts = append(appOpts, CHANGE_ITEM_STATUS)
		default:
			appOpts = append(appOpts, INVALID_APP_OPT)
		}
	})

	if len(appOpts) == 0 || contains(appOpts, INVALID_APP_OPT) {
		fmt.Println("Invalid option set or the desired option is not currently implemented")
		flag.Usage()
		return appOpts, errors.New("no flags or unimplemented were set")
	}

	return appOpts, nil
}

// main is the entry point for our todo CLI application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {

	//Process the command line flags
	opts, err := processCmdLineFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Create a new db object
	todo, err := db.New(dbFileNameFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Switch over the command line flags and call the appropriate
	//function in the db package
	for _, opt := range opts {
		switch AppOptType(opt) {
		case LIST_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			todoList, err := todo.GetAllItems()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			for _, item := range todoList {
				todo.PrintItem(item)
			}
			fmt.Println("THERE ARE", len(todoList), "ITEMS IN THE DB")
			fmt.Println("Ok")

		case QUERY_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			item, err := todo.GetItem(queryFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			todo.PrintItem(item)
			fmt.Println("Ok")
		case ADD_DB_ITEM:
			fmt.Println("Running ADD_DB_ITEM...")
			item, err := todo.JsonToItem(addFlag)
			if err != nil {
				fmt.Println("Add option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.AddItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case UPDATE_DB_ITEM:
			fmt.Println("Running UPDATE_DB_ITEM...")
			item, err := todo.JsonToItem(updateFlag)
			if err != nil {
				fmt.Println("Update option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.UpdateItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case DELETE_DB_ITEM:
			fmt.Println("Running DELETE_DB_ITEM...")
			err := todo.DeleteItem(deleteFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case CHANGE_ITEM_STATUS:
			//For the CHANGE_ITEM_STATUS extra credit you will also
			//need to add some code here
			fmt.Println("Running CHANGE_ITEM_STATUS...")
			if len(opts) > 1 {
				if contains(opts, QUERY_DB_ITEM) {
					item, err := todo.GetItem(queryFlag)
					if err != nil {
						fmt.Println("Failed to get query item for status update")
					} else {
						item.IsDone = itemStatusFlag
						err = todo.UpdateItem(item)
						if err != nil {
							fmt.Println("Failed to update item status")
						}
						fmt.Printf("Updated item id %d done status to %t\n", queryFlag, itemStatusFlag)
					}

				} else {
					fmt.Println("Attempting to change item status without querying for item 1")
				}
			} else {
				fmt.Println("Attempting to change item status without querying for item")
			}
		default:
			fmt.Println("INVALID_APP_OPT")
		}
	}
}
