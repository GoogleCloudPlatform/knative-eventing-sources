/*
Copyright 2020 Google LLC

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

package utils

import (
	"strings"
)

const (
	// This is the CloudEvent extension which holds the timeout duration of the forward client.
	ProbeEventTimeoutExtension = "timeout"

	// This is the pair of CloudEvent extensions which hold the path of probe requests.
	// The path of receiver probe requests is injected into the receiverpath extension
	// by the receiver client. The corresponding forward probe requests should include
	// the targetpath extension with the same path. For example:
	//
	// If a source is configured to sink events to:
	//   http://probe-helper-receiver.events-system-probe.svc.cluster.local/some-path-goes-here
	//
	// then, each event generated by the source and passed to the probe helper
	// receiver client will include the header 'Ce-Receiverpath: /some-path-goes-here'.
	// And in order for the probe to succeed, the corresponding forward probe should
	// be made to include the header 'Ce-Targetpath: /some-path-goes-here'.
	ProbeEventTargetPathExtension   = "targetpath"
	ProbeEventReceiverPathExtension = "receiverpath"
)

var (
	ProbeEventTargetPathHeader   = "Ce-" + strings.Title(ProbeEventTargetPathExtension)
	ProbeEventReceiverPathHeader = "Ce-" + strings.Title(ProbeEventReceiverPathExtension)
)
