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
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/wrapper"
	wgin "github.com/shiningrush/droplet/wrapper/gin"

	"github.com/apisix/manager-api/internal/core/entity"
	"github.com/apisix/manager-api/internal/core/store"
	"github.com/apisix/manager-api/internal/handler"
)

type Handler struct {
	routeStore    store.Interface
	serviceStore  store.Interface
	upstreamStore store.Interface
	sslStore      store.Interface
	consumerStore store.Interface
}

func NewHandler() (handler.RouteRegister, error) {
	return &Handler{
		routeStore:    store.GetStore(store.HubKeyRoute),
		serviceStore:  store.GetStore(store.HubKeyService),
		upstreamStore: store.GetStore(store.HubKeyUpstream),
		sslStore:      store.GetStore(store.HubKeySsl),
		consumerStore: store.GetStore(store.HubKeyConsumer),
	}, nil
}

func (h *Handler) ApplyRoute(r *gin.Engine) {
	r.GET("/api/labels/:type", wgin.Wraps(h.List,
		wrapper.InputType(reflect.TypeOf(ListInput{}))))
}

type ListInput struct {
	Type   string `auto_read:"type,path" validate:"required"`
	Labels string `auto_read:"label,query" validate:"required"`
	store.Pagination
}

func genLabelMap(label string) map[string]string {
	mp := make(map[string]string)

	if label == "" {
		return mp
	}

	labels := strings.Split(label, ",")
	for _, l := range labels {
		kv := strings.Split(l, ":")
		if len(kv) == 2 {
			mp[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			mp[kv[0]] = ""
		}
	}

	return mp
}

func checkMatch(reqLabels, labels map[string]string) bool {
	if len(reqLabels) == 0 {
		return true
	}

	for k, v := range labels {
		l, exist := reqLabels[k]
		if exist && ((l == "") || v == l) {
			return true
		}
	}

	return false
}

func getMatch(reqLabels, labels map[string]string) map[string]string {
	if len(reqLabels) == 0 {
		return labels
	}

	var res = make(map[string]string)
	for k, v := range labels {
		l, exist := reqLabels[k]
		if exist && ((l == "") || v == l) {
			res[k] = v
		}
	}

	return res
}

func (h *Handler) List(c droplet.Context) (interface{}, error) {
	input := c.Input().(*ListInput)

	typ := input.Type
	reqLabels := genLabelMap(input.Labels)

	var items []interface{}

	switch typ {
	case "route":
		items = append(items, h.routeStore)
	case "service":
		items = append(items, h.serviceStore)
	case "consumer":
		items = append(items, h.consumerStore)
	case "ssl":
		items = append(items, h.sslStore)
	case "upstream":
		items = append(items, h.upstreamStore)
	case "all":
		items = append(items, h.routeStore, h.serviceStore, h.upstreamStore,
			h.sslStore, h.consumerStore)
	}

	predicate := func(obj interface{}) bool {
		var ls map[string]string

		switch obj.(type) {
		case *entity.Route:
			ls = obj.(*entity.Route).Labels
		case *entity.Consumer:
			ls = obj.(*entity.Consumer).Labels
		case *entity.SSL:
			ls = obj.(*entity.SSL).Labels
		case *entity.Service:
			ls = obj.(*entity.Service).Labels
		case *entity.Upstream:
			ls = obj.(*entity.Upstream).Labels
		default:
			return false
		}

		v := checkMatch(reqLabels, ls)
		return v
	}

	format := func(obj interface{}) interface{} {
		val := reflect.ValueOf(obj).Elem()
		l := val.FieldByName("Labels")
		if l.IsNil() {
			return nil
		}

		ls := l.Interface().(map[string]string)
		return getMatch(reqLabels, ls)
	}

	var retSum = new(store.ListOutput)
	for _, item := range items {
		ret, err := item.(store.Interface).List(
			store.ListInput{
				Predicate:  predicate,
				Format:     format,
				PageSize:   input.PageSize,
				PageNumber: input.PageNumber},
		)

		if err != nil {
			return nil, err
		}

		for _, r := range ret.Rows {
			for k, v := range r.(map[string]string) {
				new := make(map[string]string)
				new[k] = v
				retSum.Rows = append(retSum.Rows, new)
				retSum.TotalSize += 1
			}
		}
	}

	return retSum, nil
}
