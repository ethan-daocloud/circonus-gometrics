// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).

import (
	"fmt"
)

// GaugeWithTags sets a gauge metric with tags to a value
func (m *CirconusMetrics) GaugeWithTags(metric string, tags Tags, val interface{}) {
	m.SetGaugeWithTags(metric, tags, val)
}

// Gauge sets a gauge to a value
func (m *CirconusMetrics) Gauge(metric string, val interface{}) {
	m.SetGauge(metric, val)
}

// SetGaugeWithTags sets a gauge metric with tags to a value
func (m *CirconusMetrics) SetGaugeWithTags(metric string, tags Tags, val interface{}) {
	m.SetGauge(MetricNameWithStreamTags(metric, tags), val)
}

// SetGauge sets a gauge to a value
func (m *CirconusMetrics) SetGauge(metric string, val interface{}) {
	m.gm.Lock()
	defer m.gm.Unlock()
	m.gauges[metric] = val
}

// AddGaugeWithTags adds value to existing gauge metric with tags
func (m *CirconusMetrics) AddGaugeWithTags(metric string, tags Tags, val interface{}) {
	m.AddGauge(MetricNameWithStreamTags(metric, tags), val)
}

// AddGauge adds value to existing gauge
func (m *CirconusMetrics) AddGauge(metric string, val interface{}) {
	m.gm.Lock()
	defer m.gm.Unlock()

	v, ok := m.gauges[metric]
	if !ok {
		m.gauges[metric] = val
		return
	}

	switch val.(type) {
	default:
		// ignore it, unsupported type
	case int:
		m.gauges[metric] = v.(int) + val.(int)
	case int8:
		m.gauges[metric] = v.(int8) + val.(int8)
	case int16:
		m.gauges[metric] = v.(int16) + val.(int16)
	case int32:
		m.gauges[metric] = v.(int32) + val.(int32)
	case int64:
		m.gauges[metric] = v.(int64) + val.(int64)
	case uint:
		m.gauges[metric] = v.(uint) + val.(uint)
	case uint8:
		m.gauges[metric] = v.(uint8) + val.(uint8)
	case uint16:
		m.gauges[metric] = v.(uint16) + val.(uint16)
	case uint32:
		m.gauges[metric] = v.(uint32) + val.(uint32)
	case uint64:
		m.gauges[metric] = v.(uint64) + val.(uint64)
	case float32:
		m.gauges[metric] = v.(float32) + val.(float32)
	case float64:
		m.gauges[metric] = v.(float64) + val.(float64)
	}
}

// RemoveGaugeWithTags removes a gauge metric with tags
func (m *CirconusMetrics) RemoveGaugeWithTags(metric string, tags Tags) {
	m.RemoveGauge(MetricNameWithStreamTags(metric, tags))
}

// RemoveGauge removes a gauge
func (m *CirconusMetrics) RemoveGauge(metric string) {
	m.gm.Lock()
	defer m.gm.Unlock()
	delete(m.gauges, metric)
}

// GetGaugeTest returns the current value for a gauge. (note: it is a function specifically for "testing", disable automatic submission during testing.)
func (m *CirconusMetrics) GetGaugeTest(metric string) (interface{}, error) {
	m.gm.Lock()
	defer m.gm.Unlock()

	if val, ok := m.gauges[metric]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Gauge metric '%s' not found", metric)
}

// SetGaugeFuncWithTags sets a gauge metric with tags to a function [called at flush interval]
func (m *CirconusMetrics) SetGaugeFuncWithTags(metric string, tags Tags, fn func() int64) {
	m.SetGaugeFunc(MetricNameWithStreamTags(metric, tags), fn)
}

// SetGaugeFunc sets a gauge to a function [called at flush interval]
func (m *CirconusMetrics) SetGaugeFunc(metric string, fn func() int64) {
	m.gfm.Lock()
	defer m.gfm.Unlock()
	m.gaugeFuncs[metric] = fn
}

// RemoveGaugeFuncWithTags removes a gauge metric with tags function
func (m *CirconusMetrics) RemoveGaugeFuncWithTags(metric string, tags Tags) {
	m.RemoveGaugeFunc(MetricNameWithStreamTags(metric, tags))
}

// RemoveGaugeFunc removes a gauge function
func (m *CirconusMetrics) RemoveGaugeFunc(metric string) {
	m.gfm.Lock()
	defer m.gfm.Unlock()
	delete(m.gaugeFuncs, metric)
}

// getGaugeType returns accurate resmon type for underlying type of gauge value
func (m *CirconusMetrics) getGaugeType(v interface{}) string {
	mt := "n"
	switch v.(type) {
	case int:
		mt = "i"
	case int8:
		mt = "i"
	case int16:
		mt = "i"
	case int32:
		mt = "i"
	case uint:
		mt = "I"
	case uint8:
		mt = "I"
	case uint16:
		mt = "I"
	case uint32:
		mt = "I"
	case int64:
		mt = "l"
	case uint64:
		mt = "L"
	}

	return mt
}
