/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func stringRepresentation(value interface{}) (string, error) {
	var result string
	switch value.(type) {
	case string:
		result = value.(string) // use string value as-is
	default:
		json, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		result = string(json) // return JSON text representation of value object
	}
	return result, nil
}

func divertStdoutToString(fn func() error) (string, error) {
	previous := os.Stdout
	defer func() {
		os.Stdout = previous
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()
	if err != nil {
		return "", err
	}

	errc := make(chan error)
	var buf bytes.Buffer

	go func() {
		_, err := io.Copy(&buf, r)
		errc <- err
	}()
	w.Close()
	err = <-errc
	output := strings.TrimSpace(buf.String())
	return output, err
}
