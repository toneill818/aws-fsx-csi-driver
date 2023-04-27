/*
Copyright 2020 The Kubernetes Authors.

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

package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/aws-fsx-csi-driver/cmd/options"
	"sigs.k8s.io/aws-fsx-csi-driver/pkg/driver"

	"k8s.io/component-base/featuregate"
	logsapi "k8s.io/component-base/logs/api/v1"
)

// Options is the combined set of options for all operating modes.
type Options struct {
	*options.ServerOptions
	*options.ControllerOptions
	*options.NodeOptions
}

// used for testing
var osExit = os.Exit

var featureGate = featuregate.NewFeatureGate()

func GetOptions(fs *flag.FlagSet) *Options {
	var (
		version = fs.Bool("version", false, "Print the version and exit.")
		args    = os.Args[1:]

		serverOptions     = options.ServerOptions{}
		controllerOptions = options.ControllerOptions{}
		nodeOptions       = options.NodeOptions{}
	)

	serverOptions.AddFlags(fs)

	c := logsapi.NewLoggingConfiguration()

	err := logsapi.AddFeatureGates(featureGate)
	if err != nil {
		klog.ErrorS(err, "failed to add feature gates")
	}

	logsapi.AddFlags(c, fs)

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	err = logsapi.ValidateAndApply(c, featureGate)
	if err != nil {
		klog.ErrorS(err, "failed to validate and apply logging configuration")
	}

	if *version {
		versionInfo, err := driver.GetVersionJSON()
		if err != nil {
			klog.ErrorS(err, "failed to get version")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}
		fmt.Println(versionInfo)
		osExit(0)
	}

	return &Options{
		ServerOptions:     &serverOptions,
		ControllerOptions: &controllerOptions,
		NodeOptions:       &nodeOptions,
	}
}
