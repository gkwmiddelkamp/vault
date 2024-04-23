package vault

type Result struct {
	Done bool
}

func NotDone() Result {
	return Result{}
}

func Done() Result {
	return Result{Done: true}
}
