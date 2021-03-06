// Package interfaces is auto generated
package interfaces

import tarantool "github.com/tarantool/go-tarantool"

// IConnection is an autogenerated interface to tarantool.Connection
type IConnection interface {
	// Ping sends empty request to Tarantool to check connection.
	Ping() (resp *tarantool.Response, err error)
	// Select performs select to box space.
	//
	// It is equal to conn.SelectAsync(...).Get()
	Select(space, index interface{}, offset, limit, iterator uint32, key interface{}) (resp *tarantool.Response, err error)
	// Insert performs insertion to box space.
	// Tarantool will reject Insert when tuple with same primary key exists.
	//
	// It is equal to conn.InsertAsync(space, tuple).Get().
	Insert(space interface{}, tuple interface{}) (resp *tarantool.Response, err error)
	// Replace performs "insert or replace" action to box space.
	// If tuple with same primary key exists, it will be replaced.
	//
	// It is equal to conn.ReplaceAsync(space, tuple).Get().
	Replace(space interface{}, tuple interface{}) (resp *tarantool.Response, err error)
	// Delete performs deletion of a tuple by key.
	// Result will contain array with deleted tuple.
	//
	// It is equal to conn.DeleteAsync(space, tuple).Get().
	Delete(space, index interface{}, key interface{}) (resp *tarantool.Response, err error)
	// Update performs update of a tuple by key.
	// Result will contain array with updated tuple.
	//
	// It is equal to conn.UpdateAsync(space, tuple).Get().
	Update(space, index interface{}, key, ops interface{}) (resp *tarantool.Response, err error)
	// Upsert performs "update or insert" action of a tuple by key.
	// Result will not contain any tuple.
	//
	// It is equal to conn.UpsertAsync(space, tuple, ops).Get().
	Upsert(space interface{}, tuple, ops interface{}) (resp *tarantool.Response, err error)
	// Call calls registered tarantool function.
	// It uses request code for tarantool 1.6, so result is converted to array of arrays
	//
	// It is equal to conn.CallAsync(functionName, args).Get().
	Call(functionName string, args interface{}) (resp *tarantool.Response, err error)
	// Call17 calls registered tarantool function.
	// It uses request code for tarantool 1.7, so result is not converted
	// (though, keep in mind, result is always array)
	//
	// It is equal to conn.Call17Async(functionName, args).Get().
	Call17(functionName string, args interface{}) (resp *tarantool.Response, err error)
	// Eval passes lua expression for evaluation.
	//
	// It is equal to conn.EvalAsync(space, tuple).Get().
	Eval(expr string, args interface{}) (resp *tarantool.Response, err error)
	// GetTyped performs select (with limit = 1 and offset = 0)
	// to box space and fills typed result.
	//
	// It is equal to conn.SelectAsync(space, index, 0, 1, IterEq, key).GetTyped(&result)
	GetTyped(space, index interface{}, key interface{}, result interface{}) (err error)
	// SelectTyped performs select to box space and fills typed result.
	//
	// It is equal to conn.SelectAsync(space, index, offset, limit, iterator, key).GetTyped(&result)
	SelectTyped(space, index interface{}, offset, limit, iterator uint32, key interface{}, result interface{}) (err error)
	// InsertTyped performs insertion to box space.
	// Tarantool will reject Insert when tuple with same primary key exists.
	//
	// It is equal to conn.InsertAsync(space, tuple).GetTyped(&result).
	InsertTyped(space interface{}, tuple interface{}, result interface{}) (err error)
	// ReplaceTyped performs "insert or replace" action to box space.
	// If tuple with same primary key exists, it will be replaced.
	//
	// It is equal to conn.ReplaceAsync(space, tuple).GetTyped(&result).
	ReplaceTyped(space interface{}, tuple interface{}, result interface{}) (err error)
	// DeleteTyped performs deletion of a tuple by key and fills result with deleted tuple.
	//
	// It is equal to conn.DeleteAsync(space, tuple).GetTyped(&result).
	DeleteTyped(space, index interface{}, key interface{}, result interface{}) (err error)
	// UpdateTyped performs update of a tuple by key and fills result with updated tuple.
	//
	// It is equal to conn.UpdateAsync(space, tuple, ops).GetTyped(&result).
	UpdateTyped(space, index interface{}, key, ops interface{}, result interface{}) (err error)
	// CallTyped calls registered function.
	// It uses request code for tarantool 1.6, so result is converted to array of arrays
	//
	// It is equal to conn.CallAsync(functionName, args).GetTyped(&result).
	CallTyped(functionName string, args interface{}, result interface{}) (err error)
	// Call17Typed calls registered function.
	// It uses request code for tarantool 1.7, so result is not converted
	// (though, keep in mind, result is always array)
	//
	// It is equal to conn.Call17Async(functionName, args).GetTyped(&result).
	Call17Typed(functionName string, args interface{}, result interface{}) (err error)
	// EvalTyped passes lua expression for evaluation.
	//
	// It is equal to conn.EvalAsync(space, tuple).GetTyped(&result).
	EvalTyped(expr string, args interface{}, result interface{}) (err error)
	// SelectAsync sends select request to tarantool and returns Future.
	SelectAsync(space, index interface{}, offset, limit, iterator uint32, key interface{}) *tarantool.Future
	// InsertAsync sends insert action to tarantool and returns Future.
	// Tarantool will reject Insert when tuple with same primary key exists.
	InsertAsync(space interface{}, tuple interface{}) *tarantool.Future
	// ReplaceAsync sends "insert or replace" action to tarantool and returns Future.
	// If tuple with same primary key exists, it will be replaced.
	ReplaceAsync(space interface{}, tuple interface{}) *tarantool.Future
	// DeleteAsync sends deletion action to tarantool and returns Future.
	// Future's result will contain array with deleted tuple.
	DeleteAsync(space, index interface{}, key interface{}) *tarantool.Future
	// Update sends deletion of a tuple by key and returns Future.
	// Future's result will contain array with updated tuple.
	UpdateAsync(space, index interface{}, key, ops interface{}) *tarantool.Future
	// UpsertAsync sends "update or insert" action to tarantool and returns Future.
	// Future's sesult will not contain any tuple.
	UpsertAsync(space interface{}, tuple interface{}, ops interface{}) *tarantool.Future
	// CallAsync sends a call to registered tarantool function and returns Future.
	// It uses request code for tarantool 1.6, so future's result is always array of arrays
	CallAsync(functionName string, args interface{}) *tarantool.Future
	// Call17Async sends a call to registered tarantool function and returns Future.
	// It uses request code for tarantool 1.7, so future's result will not be converted
	// (though, keep in mind, result is always array)
	Call17Async(functionName string, args interface{}) *tarantool.Future
	// EvalAsync sends a lua expression for evaluation and returns Future.
	EvalAsync(expr string, args interface{}) *tarantool.Future
}
