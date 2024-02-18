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

package global

import (
	"github.com/creasty/defaults"

	"go.uber.org/atomic"
)

// ShutdownConfig is used as configuration for graceful shutdown
type ShutdownConfig struct {
	/*
	 * Total timeout. Even though we don't release all resources,
	 * the applicationConfig will shutdown if the costing time is over this configuration. The unit is ms.
	 * default value is 60 * 1000 ms = 1 minutes
	 * In general, it should be bigger than 3 * StepTimeout.
	 */
	Timeout string `default:"60s" yaml:"timeout" json:"timeout,omitempty" property:"timeout"`
	/*
	 * the timeout on each step. You should evaluate the response time of request
	 * and the time that client noticed that server shutdown.
	 * For example, if your client will received the notification within 10s when you start to close server,
	 * and the 99.9% requests will return response in 2s, so the StepTimeout will be bigger than(10+2) * 1000ms,
	 * maybe (10 + 2*3) * 1000ms is a good choice.
	 */
	StepTimeout string `default:"3s" yaml:"step-timeout" json:"step.timeout,omitempty" property:"step.timeout"`

	/*
	 * ConsumerUpdateWaitTime means when provider is shutting down, after the unregister, time to wait for client to
	 * update invokers. During this time, incoming invocation can be treated normally.
	 */
	ConsumerUpdateWaitTime string `default:"3s" yaml:"consumer-update-wait-time" json:"consumerUpdate.waitTIme,omitempty" property:"consumerUpdate.waitTIme"`
	// when we try to shutdown the applicationConfig, we will reject the new requests. In most cases, you don't need to configure this.
	RejectRequestHandler string `yaml:"reject-handler" json:"reject-handler,omitempty" property:"reject_handler"`
	// internal listen kill signal，the default is true.
	InternalSignal *bool `default:"true" yaml:"internal-signal" json:"internal.signal,omitempty" property:"internal.signal"`
	// offline request window length
	OfflineRequestWindowTimeout string `yaml:"offline-request-window-timeout" json:"offlineRequestWindowTimeout,omitempty" property:"offlineRequestWindowTimeout"`
	// true -> new request will be rejected.
	RejectRequest atomic.Bool
	// active invocation
	ConsumerActiveCount atomic.Int32
	ProviderActiveCount atomic.Int32

	// provider last received request timestamp
	ProviderLastReceivedRequestTime atomic.Time
}

func DefaultShutdownConfig() *ShutdownConfig {
	cfg := &ShutdownConfig{}
	defaults.MustSet(cfg)

	return cfg
}

// Clone a new ShutdownConfig
func (c *ShutdownConfig) Clone() *ShutdownConfig {
	var newInternalSignal *bool
	if c.InternalSignal != nil {
		newInternalSignal = new(bool)
		*newInternalSignal = *c.InternalSignal
	}

	newShutdownConfig := &ShutdownConfig{
		Timeout:                     c.Timeout,
		StepTimeout:                 c.StepTimeout,
		ConsumerUpdateWaitTime:      c.ConsumerUpdateWaitTime,
		RejectRequestHandler:        c.RejectRequestHandler,
		InternalSignal:              newInternalSignal,
		OfflineRequestWindowTimeout: c.OfflineRequestWindowTimeout,
	}

	newShutdownConfig.RejectRequest.Store(c.RejectRequest.Load())
	newShutdownConfig.ConsumerActiveCount.Store(c.ConsumerActiveCount.Load())
	newShutdownConfig.ProviderActiveCount.Store(c.ProviderActiveCount.Load())
	newShutdownConfig.ProviderLastReceivedRequestTime.Store(c.ProviderLastReceivedRequestTime.Load())

	return newShutdownConfig
}
