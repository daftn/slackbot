package slackbot

import "github.com/robfig/cron"

type cronScheduler interface {
	Schedule(cron.Schedule, cron.Job)
	Start()
}

type (
	// ScheduledTask is used to run the Task on a scheduled cron using the string Schedule
	ScheduledTask struct {
		Schedule string
		Task     taskFunc
	}

	scheduler struct {
		cronScheduler
	}

	// wrapping the taskFunc to allow passing the Bot to the Task
	taskFuncWrapper struct {
		taskFunc taskFunc
		bot      *Bot
	}

	taskFunc func(*Bot)
)

func (t taskFuncWrapper) Run() {
	t.taskFunc(t.bot)
}

func (sc *scheduler) scheduleTasks(bot *Bot, tasks []ScheduledTask) error {
	for _, t := range tasks {
		s, err := cron.ParseStandard(t.Schedule)
		if err != nil {
			return err
		}

		tw := taskFuncWrapper{
			bot:      bot,
			taskFunc: t.Task,
		}
		sc.Schedule(s, tw)
	}
	sc.Start()

	return nil
}
