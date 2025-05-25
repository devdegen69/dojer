package downloader

import "io"

type Queue struct {
	Items        []DownloadItem
	Limit        int
	Events       chan struct{}
	RunningTasks int
	Callback     func()
}

func NewQueue() Queue {
	limit := 5
	q := Queue{
		Items:        []DownloadItem{},
		Limit:        limit,
		RunningTasks: 0,
		Events:       make(chan struct{}, limit)}
	return q
}

func (q *Queue) AddItem(item DownloadItem) {
	q.Items = append(q.Items, item)
	q.Events <- struct{}{}
}

func (q *Queue) DoneTask() {
	q.RunningTasks -= 1
	if len(q.Items) > 0 {
		q.Events <- struct{}{}
	}
}

func (q *Queue) StartTask() {
	q.RunningTasks += 1
}

func (q *Queue) PopItem() {
	q.Items = q.Items[:len(q.Items)-1]
}

func (q *Queue) ShiftItem() DownloadItem {
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

func (q *Queue) RmItem(index int) {
	q.Items = append(q.Items[:index], q.Items[index+1:]...)
}

var DownloadQueue = NewQueue()

func ListenDownloadQueue() {
	go func() {
		for {
			select {
			case <-DownloadQueue.Events:
				if DownloadQueue.RunningTasks < DownloadQueue.Limit && len(DownloadQueue.Items) > 0 {
					DownloadQueue.StartTask()
					go func(item DownloadItem) {
						defer DownloadQueue.DoneTask()
						Download(item)
					}(DownloadQueue.ShiftItem())
				}
			default:
				break
			}
		}
	}()

	go func() {
		for {
			pr, _ := GetPipe()
			io.Copy(io.Discard, pr)
		}
	}()

	<-make(chan struct{})
}
