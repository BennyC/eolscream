package pkg

// Notifier interface
type Notifier interface {
	Notify()
}

func NewNilNotifier() *NilNotifier {
	return &NilNotifier{}
}

// NilNotifier implementation of Notifier that does nothing
type NilNotifier struct{}

// Notify is a Nil operation
func (n NilNotifier) Notify() {}
