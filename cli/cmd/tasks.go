package cmd

import (
	"errors"
	"fmt"
	"github.com/coffeemakr/ruck"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

var (
	taskCommand = &cobra.Command{
		Use: "task",
	}
	taskListCommand = &cobra.Command{
		Use:     "list",
		Run:     runTaskList,
		Aliases: []string{"ls"},
	}
	taskGetCommand = &cobra.Command{
		Use:  "get",
		Run:  runTaskGet,
		Args: cobra.ExactArgs(1),
	}
	taskAddCommand = &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
		Run:   runAddTask,
		Args:  cobra.ExactArgs(1),
	}
	taskAddOptionDaily, taskAddOptionWeekly, taskAddOptionMonthly, taskAddOptionYearly bool
	taskAddOptionInterval                                                              uint32

	taskDoneCommand = &cobra.Command{
		Use:   "complete",
		Short: "Complete a task",
		Aliases: []string{"done"},
		Run:   runCompleteTask,
		Args:  cobra.ExactArgs(1),
	}

)

func runTaskGet(cmd *cobra.Command, args []string) {
	taskId := args[0]
	var task *ruck.Task
	task, err := client.GetTaskDetails(taskId)
	if err != nil {
		log.Fatalln(err)
	}
	t, err := template.New("taskTemplate").Parse("ID    {{.ID}}\n" +
		"Name  {{.Name}}\n" +
		"Group {{.Group}}\n")
	if err != nil {
		log.Fatalln(err)
	}
	err = t.Execute(os.Stdout, task)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	taskAddCommand.PersistentFlags().BoolVarP(&taskAddOptionDaily, "daily", "d", false, "Repeat task daily")
	taskAddCommand.PersistentFlags().BoolVarP(&taskAddOptionWeekly, "weekly", "w", false, "Repeat task weekly")
	taskAddCommand.PersistentFlags().BoolVarP(&taskAddOptionMonthly, "monthly", "m", false, "Repeat task monthly")
	taskAddCommand.PersistentFlags().BoolVarP(&taskAddOptionYearly, "yearly", "y", false, "Repeat task yearly")
	taskAddCommand.PersistentFlags().Uint32Var(&taskAddOptionInterval, "interval", 1, "Interval number e.g. X weeks when --weeks flag is used")
	taskCommand.AddCommand(taskAddCommand, taskListCommand, taskGetCommand, taskDoneCommand)
}

func getDaysUntilTime(due time.Time) int {
	return int(due.Sub(time.Now()).Hours() / 24)
}

func formatDue(dueDate time.Time) string {
	daysLeft := getDaysUntilTime(dueDate)
	switch {
	case daysLeft < -1:
		return fmt.Sprintf("overdue (%d days)", daysLeft)
	case daysLeft == -1:
		return "overdue (1 day)"
	case daysLeft == 0:
		return "due today"
	case daysLeft == 1:
		return "1 day left"
	default:
		return fmt.Sprintf("%d days left", daysLeft)
	}
}

func runTaskList(cmd *cobra.Command, args []string) {
	tasks, err := client.GetTaskList()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sort.Sort(ruck.ByDueDate(tasks))
	for _, task := range tasks {
		ID := task.ID
		fmt.Printf("%s %-40s %-20s %s\n", ID, task.Name, task.AssigneeName, formatDue(task.DueDate))
	}
}

func getIntervalUnit() (unit ruck.IntervalUnit, err error) {
	var bits int
	const (
		DailyBit   = 1 << 0
		WeeklyBit  = 1 << 1
		MonthlyBit = 1 << 2
		YearlyBit  = 1 << 3
	)
	if taskAddOptionDaily {
		bits |= DailyBit
	}
	if taskAddOptionWeekly {
		bits |= WeeklyBit
	}

	if taskAddOptionMonthly {
		bits |= MonthlyBit
	}
	if taskAddOptionYearly {
		bits |= YearlyBit
	}
	switch bits {
	case YearlyBit:
		unit = ruck.Years
	case MonthlyBit:
		unit = ruck.Months
	case WeeklyBit:
		unit = ruck.Weeks
	case DailyBit:
		unit = ruck.Days
	case 0:
		err = errors.New("require at least one of weekly, daily, monthly or yearly flags")
	default:
		err = fmt.Errorf("got conflicting interval flags %d", bits)
	}
	return
}

func runAddTask(cmd *cobra.Command, args []string) {
	var (
		defaultGroupId = getDefaultGroup()
		task           ruck.Task
		err            error
	)

	task.Name = strings.TrimSpace(args[0])
	task.Interval.Amount = taskAddOptionInterval
	task.Interval.Unit, err = getIntervalUnit()
	task.GroupID = defaultGroupId
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if defaultGroupId == "" {
		log.Fatalln("no default group set")
	}

	err = client.CreateTask(&task)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Task created: %s\n", task)
}

func runCompleteTask(cmd *cobra.Command, args []string) {
	taskId := args[0]
	execution, err := client.CompleteTask(taskId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Task completed: %s\n", execution)
}
