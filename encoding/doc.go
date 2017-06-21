// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package encoding is responsible for encoding and decoding of structs.
// `squash` option assumes that field names are unique in structs, otherwise results are unpredictable
// `cast` option casts value to specified type on encoding to string-interface map
package encoding
