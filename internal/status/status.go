// Copyright 2022 Marek Dalewski

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package status

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type statusContextKeyType string

const statusContextKey statusContextKeyType = "status"

type Status struct {
	ErrorOccurred bool
}

type statusStore struct {
	sync.Mutex
	status Status
}

func Track(ctx context.Context, mods ...func(Status) Status) context.Context {
	ss := unpack(ctx)
	if ss == nil {
		ss = &statusStore{}
		ctx = context.WithValue(ctx, statusContextKey, ss)
	}

	ss.Lock()
	defer ss.Unlock()
	for _, m := range mods {
		ss.status = m(ss.status)
	}
	return ctx
}

func Unpack(ctx context.Context) Status {
	ss := unpack(ctx)
	if ss == nil {
		return Status{}
	}
	ss.Lock()
	defer ss.Unlock()
	return ss.status
}

func unpack(ctx context.Context) *statusStore {
	x := ctx.Value(statusContextKey)
	if x == nil {
		return nil
	}
	s, ok := x.(*statusStore)
	if !ok {
		return nil
	}
	return s
}

func ShowErr(ctx context.Context, msg interface{}) {
	if msg != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", msg)
		Track(ctx, func(s Status) Status { s.ErrorOccurred = true; return s })
	}
}

func CheckErr(ctx context.Context, msg interface{}) {
	if msg != nil {
		ShowErr(ctx, msg)
		Exit(ctx, 1)
	}
}

func Exit(ctx context.Context, ec int) {
	if ec != 0 {
		os.Exit(ec)
	}
	if Unpack(ctx).ErrorOccurred {
		os.Exit(1)
	}
	os.Exit(0)
}
