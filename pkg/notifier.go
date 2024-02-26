package pkg

// Notifier interface
type Notifier interface {
	Notify(p Product, i ReleaseInfo)
}

func NewNilNotifier() *NilNotifier {
	return &NilNotifier{}
}

// NilNotifier implementation of Notifier that does nothing
type NilNotifier struct{}

// Notify is a Nil operation
func (n NilNotifier) Notify(_ Product, _ ReleaseInfo) {}
