package app

func (receiver *Application) source() {
	receiver.telegram.ForwardTo(receiver.sourceChan)
}
