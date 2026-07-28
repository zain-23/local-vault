package reqctx

import "context"

type Info struct {
	ActorID		string
	DeviceID	string
	IP			string
}

type Key struct {}

func With(ctx context.Context, i Info) context.Context {
	return context.WithValue(ctx, Key{}, i)
}


func From(ctx context.Context) Info {
	if i, ok := ctx.Value(Key{}).(Info); ok {
		return i
	}

	return Info{}
}
