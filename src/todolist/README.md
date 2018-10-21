# todolist

todolist is a simple server that responds to Twilio SMS commands

* "add <task>": adds a task to the list
* "list": returns the current lists of tasks, as a numbered ordered list
* "remove <#>": removes a task, given the task number

## Development notes

Lots of hackery here - storing things as a slice in memory (rather than, say, [bbolt](https://github.com/etcd-io/bbolt) or sqlite), remove doesn't check for out of bounds, etc.