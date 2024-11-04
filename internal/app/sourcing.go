package app

func (receiver *Application) source() {
	receiver.telegramSource.ForwardTo(receiver.sourceChan)
}
