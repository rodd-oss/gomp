/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"github.com/negrel/assert"
)

type AnyAssetLibrary interface {
	LoadAll()
	Unload(path string)
	UnloadAll()
}

func CreateAssetLibrary[T any](loader func(path string) T, unloader func(path string, asset *T)) AssetLibrary[T] {
	return AssetLibrary[T]{
		data:        make(map[string]*T),
		loader:      loader,
		unloader:    unloader,
		loaderQueue: make([]string, 0, 1024),
	}
}

type AssetLibrary[T any] struct {
	data        map[string]*T
	loader      func(path string) T
	unloader    func(path string, asset *T)
	loaderQueue []string
}

func (r *AssetLibrary[T]) Get(path string) *T {
	value, ok := r.data[path]
	if !ok {
		r.loaderQueue = append(r.loaderQueue, path)
		value = new(T)
		r.data[path] = value
	}

	return value
}

func (r *AssetLibrary[T]) Load(path string) {
	_, ok := r.data[path]
	assert.False(ok, "asset already loaded")

	resource := r.loader(path)
	r.data[path] = &resource
}

func (r *AssetLibrary[T]) LoadAll() {
	if len(r.loaderQueue) == 0 {
		return
	}

	for _, path := range r.loaderQueue {
		resource := r.loader(path)
		*r.data[path] = resource
	}

	r.loaderQueue = r.loaderQueue[:0]
}

func (r *AssetLibrary[T]) Unload(path string) {
	value, ok := r.data[path]
	assert.True(ok, "asset not loaded")
	r.unloader(path, value)
	delete(r.data, path)
}

func (r *AssetLibrary[T]) UnloadAll() {
	for path, value := range r.data {
		r.unloader(path, value)
	}

	r.data = make(map[string]*T)
}
