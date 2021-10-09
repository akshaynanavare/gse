/*
Copyright © 2021 Madhav Jivrajani madhav.jiv@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/MadhavJivrajani/gse/pkg/utils"
)

// RunFromConfig runs the target binary and serves the scheduler traces.
func RunFromConfig(ctx context.Context, config *utils.Config) error {
	schedTrace := NewSchedTrace()
	for line := range streamExecutionOutput(ctx, config) {
		schedTrace.UpdateSchedTraceFromRawTrace(line)
	}

	return nil
}

func streamExecutionOutput(ctx context.Context, config *utils.Config) <-chan string {
	outChan := make(chan string, 1)
	go func() {
		cmd := exec.CommandContext(ctx, "sh", "-c", constructCommandFromConfig(config))
		stdErrPipe, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(stdErrPipe)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			outChan <- scanner.Text()
		}
	}()

	return outChan
}

func constructCommandFromConfig(config *utils.Config) string {
	return fmt.Sprintf(
		"GODEBUG=schedtrace=%d %s",
		config.SchedTrace.Interval,
		config.PathToBinary,
	)
}
