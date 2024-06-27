package stores

type Sender interface {
	Send([]string) error
}
