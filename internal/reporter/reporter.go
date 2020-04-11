package reporter

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/google/goterm/term"
)

// Message Holds a message with arguments
type Message struct {
	title string
	args  map[string]string
}

// Reporter support reporting functions
type Reporter struct {
	message *Message
	level   int
}

// Task unit of work
type Task struct {
	reporter *Reporter
	level    int
}

func printMessageAndArgs(indent int, title string, args map[string]string) {
	var tab = ""
	for i := 0; i < indent; i++ {
		tab += "  "
	}
	fmt.Print(tab)
	fmt.Println(term.Cyan(title))
	if args != nil {
		for name, value := range args {
			fmt.Printf(term.Bluef("  %s- %s: %s\n", tab, name, value))
		}
	}
}

func printError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		fmt.Println(term.Red("AWS Error is:"))
		fmt.Println(term.Red(" - code:" + aerr.Code()))
		fmt.Println(term.Red(" - error:" + aerr.Error()))
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(term.Red("Error is:"))
		fmt.Println(term.Red(err.Error()))
	}

}

// NewMessage create a message
func NewMessage(title string) *Message {
	var result = Message{
		title: title,
		args:  make(map[string]string),
	}
	return &result
}

// WithArg add arg to task description
func (message *Message) WithArg(name string, value string) *Message {
	message.args[name] = value
	return message
}

func (message *Message) print(indent int, withArgs bool, extra string) {
	var args map[string]string = nil
	if withArgs {
		args = message.args
	}
	printMessageAndArgs(indent, message.title+extra, args)
}

// NewReporterM constructor
func NewReporterM(message *Message) *Reporter {
	var result = Reporter{
		message: message,
		level:   0,
	}
	return &result
}

// NewReporter simple ctor
func NewReporter(title string) *Reporter {
	return NewReporterM(NewMessage(title))
}

// Start start the task being reported
func (reporter *Reporter) Start() *Task {
	// we log the reporter title
	reporter.message.print(reporter.level, true, "")

	var result = Task{
		level:    1 + reporter.level,
		reporter: reporter,
	}
	return &result
}

// Ok display the fact that the reporter ends successfully
func (reporter *Reporter) Ok() {
	// we log the reporter title
	reporter.message.print(reporter.level, false, ": SUCCESS")
}

// Fail indicates that the reporter failed
func (reporter *Reporter) Fail(err error) {
	// we log the reporter title
	reporter.message.print(reporter.level, false, ": FAILED")
	printError(err)
}

// SubM create sub reporter
func (task *Task) SubM(message *Message) *Task {
	var result = Reporter{
		message: message,
		level:   task.level + 1,
	}
	return result.Start()
}

// Sub create sub reporter
func (task *Task) Sub(title string) *Task {
	return task.SubM(NewMessage(title))
}

// LogM a message in the task
func (task *Task) LogM(message *Message) *Task {
	// we log the message title
	message.print(task.level, true, "")
	return task
}

// Log a message in the task
func (task *Task) Log(title string) *Task {
	// we log the message title
	task.LogM(NewMessage(title))
	return task
}

// Okr indicates success and print some values
func (task *Task) Okr(result map[string]string) {
	printMessageAndArgs(task.level, "SUCCESS:", result)
}

// Ok indicates success
func (task *Task) Ok() {
	printMessageAndArgs(task.level, "SUCCESS", nil)
}

// Fail indicates failure
func (task *Task) Fail(err error) {
	printMessageAndArgs(task.level, "FAILURE", nil)
	printError(err)
}

// Experiment experiment
func Experiment() {
	r := NewReporterM(NewMessage("my first reporter").WithArg("arg1", "42"))
	t0 := r.Start()
	t0.Log("Let's go")
	t1 := t0.Sub("Sub task...")
	t1.Log("step 1")
	t1.Okr(map[string]string{"a": "42"})
	t1.Log("step 2")
	t1.Fail(errors.New("There is an error"))
	r.Ok()
}
