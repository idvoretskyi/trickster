/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"context"
	"testing"

	"github.com/tricksterproxy/trickster/pkg/backends/rule/options"
)

func TestHops(t *testing.T) {

	ctx := context.Background()
	_, j := Hops(ctx)
	if j != options.DefaultMaxRuleExecutions {
		t.Errorf("expected %d got %d", options.DefaultMaxRuleExecutions, j)
	}

	ctx = WithHops(ctx, 0, 1)
	i, j := Hops(ctx)

	if i != 0 {
		t.Errorf("expected %d got %d", 0, i)
	}

	if j != 1 {
		t.Errorf("expected %d got %d", 1, j)
	}

}
