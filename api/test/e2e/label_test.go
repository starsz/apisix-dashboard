/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package e2e

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabel(t *testing.T) {
	// build test body
	testCert, err := ioutil.ReadFile("../certs/test2.crt")
	assert.Nil(t, err)
	testKey, err := ioutil.ReadFile("../certs/test2.key")
	assert.Nil(t, err)
	body, err := json.Marshal(map[string]interface{}{
		"id":   "1",
		"cert": string(testCert),
		"key":  string(testKey),
		"label": map[string]string{
			"build": "16",
			"env":   "production",
			"extra": "ssl",
		},
	})
	assert.Nil(t, err)

	tests := []HttpTestCase{
		{
			caseDesc: "config route",
			Object:   ManagerApiExpect(t),
			Path:     "/apisix/admin/routes/r1",
			Method:   http.MethodPut,
			Body: `{
					"uri": "/hello",
					"labels": {
						"build":"16",
						"env":"production",
						"version":"v2"
					},
					"upstream": {
						"type": "roundrobin",
						"nodes": [{
							"host": "172.16.238.20",
							"port": 1980,
							"weight": 1
						}]
					}
				}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create consumer",
			Object:   ManagerApiExpect(t),
			Path:     "/apisix/admin/consumers",
			Method:   http.MethodPut,
			Body: `{
				"username": "jack",
				"plugins": {
					"key-auth": {
						"key": "auth-one"
					}
				},
				"labels": {
					"build":"16",
					"env":"production",
					"version":"v2"
				},
				"desc": "test description"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create upstream",
			Object:   ManagerApiExpect(t),
			Method:   http.MethodPut,
			Path:     "/apisix/admin/upstreams/1",
			Body: `{
				"nodes": [{
					"host": "172.16.238.20",
					"port": 1980,
					"weight": 1
				}],
				"labels": {
					"build":"16",
					"env":"production",
					"version":"v2"
				},
				"type": "roundrobin"
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "create ssl",
			Object:       ManagerApiExpect(t),
			Method:       http.MethodPost,
			Path:         "/apisix/admin/ssl",
			Body:         string(body),
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc: "create service",
			Object:   ManagerApiExpect(t),
			Method:   http.MethodPost,
			Path:     "/apisix/admin/service/s1",
			Body: `{
				"id": "1",
				"plugins": {
					"limit-count": {
						"count": 2,
						"time_window": 60,
						"rejected_code": 503,
						"key": "remote_addr"
					}
				},
				"upstream": {
					"type": "roundrobin",
					"nodes": [{
						"host": "39.97.63.215",
						"port": 80,
						"weight": 1
					}]
				},
				"labels": {
					"build":"16",
					"env":"production",
					"version":"v2"
				},
			}`,
			Headers:      map[string]string{"Authorization": token},
			ExpectStatus: http.StatusOK,
		},
		{
			caseDesc:     "get route label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/route",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
		{
			caseDesc:     "get ssl label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/ssl",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
		{
			caseDesc:     "get consumer label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/consumer",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
		{
			caseDesc:     "get service label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/service",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
		{
			caseDesc:     "get upstream label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/upstream",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
		{
			caseDesc:     "get all label",
			Object:       APISIXExpect(t),
			Method:       http.MethodGet,
			Path:         "/api/labels/all",
			ExpectStatus: http.StatusOK,
			ExpectBody:   "Missing API key found in request",
			Sleep:        sleepTime * 2,
		},
	}

	for _, tc := range tests {
		testCaseCheck(tc)
	}
}
