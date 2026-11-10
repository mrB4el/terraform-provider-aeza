// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mrb4el/terraform-provider-aeza/internal/provider"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in provider: %v", r)
		}
	}()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.New,
	})
}
