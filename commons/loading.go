/*
 * Copyright 2021 Wang Min Xiang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commons

import (
	"fmt"
	"time"
)

func NewLoading(title string, duration time.Duration) *Loading {
	return &Loading{
		title:    title,
		duration: duration,
		stop:     make(chan struct{}, 1),
	}
}

type Loading struct {
	title    string
	duration time.Duration
	stop     chan struct{}
}

func (x *Loading) Show() {
	go func(x *Loading) {
		c := []string{"/", "-", "\\", "|"}
		i := 0
		for {
			stopped := false
			select {
			case <-x.stop:
				stopped = true
				break
			case <-time.After(x.duration):
				p := c[i%len(c)]
				fmt.Printf("\r%s %v", x.title, p)
				i++
			}
			if stopped {
				fmt.Println()
				break
			}
		}
	}(x)
}

func (x *Loading) Close() {
	close(x.stop)
}
