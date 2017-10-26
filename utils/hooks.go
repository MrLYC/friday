package utils

type HookAction interface{}

const (
	HookActionNone     = iota
	HookActionContinue = iota
	HookActionBreak    = iota
	HookActionTrue     = true
	HookActionFalse    = false
)

type IHookContext interface {
	Copy() IHookContext
}

type IHook interface {
	Fire(IHookContext) HookAction
	Wrap(IHook)
}

type Hook struct {
	Next IHook
}

func (h *Hook) Fire(context IHookContext) HookAction {
	if h.Next != nil {
		return h.Next.Fire(context)
	}
	return HookActionNone
}

func (h *Hook) Wrap(hook IHook) {
	h.Next = hook
}
