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

package label

import (
	"github.com/shiningrush/droplet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
)

type testCase struct {
	caseDesc  string
	giveInput *ListInput
	giveData  []interface{}
	giveErr   error
	wantErr   error
	wantInput store.ListInput
	wantRet   interface{}
	called    bool
}

func init() {
	testing.Init()
}

func isEqualMap(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		v2, exist := m2[k]
		if !exist || v != v2 {
			return false
		}
	}

	return true
}

func genMockStore(t *testing.T, tc *testCase) *store.MockInterface {
	mStore := &store.MockInterface{}
	mStore.On("List", mock.Anything).Run(func(args mock.Arguments) {
		tc.called = true
		input := args.Get(0).(store.ListInput)
		assert.Equal(t, input.PageSize, tc.wantInput.PageSize)
		assert.Equal(t, input.PageNumber, tc.wantInput.PageNumber)
	}).Return(func(input store.ListInput) *store.ListOutput {
		var returnData []interface{}
		for _, c := range tc.giveData {
			if input.Predicate(c) {
				returnData = append(returnData, input.Format(c))
			}
		}
		return &store.ListOutput{
			Rows:      returnData,
			TotalSize: len(returnData),
		}
	}, tc.giveErr)

	return mStore
}

func TestRoute(t *testing.T) {
	handler := &Handler{}
	assert.NotNil(t, handler)

	testCases := []testCase{
		{
			giveInput: &ListInput{
				Type: "route",
				Pagination: store.Pagination{
					PageSize:   10,
					PageNumber: 10,
				},
			},
			wantInput: store.ListInput{
				PageSize:   10,
				PageNumber: 10,
			},
			giveData: []interface{}{
				&entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "1",
					},
					URI: "/test/r1",
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
				},
				&entity.Route{
					BaseInfo: entity.BaseInfo{
						ID: "2",
					},
					URI: "/test/r2",
					Labels: map[string]string{
						"label1": "value2",
					},
				},
			},
			wantRet: &store.ListOutput{
				Rows: []interface{}{
					map[string]string{
						"label1": "value1",
					},
					map[string]string{
						"label1": "value2",
					},
					map[string]string{
						"label2": "value2",
					},
				},
				TotalSize: 3,
			},
		},
	}

	for _, tc := range testCases {
		handler.routeStore = genMockStore(t, &tc)
		ctx := droplet.NewContext()
		ctx.SetInput(tc.giveInput)
		handler.List(ctx)
		assert.True(t, tc.called)
	}
}
