/* Copyright 2022 Zinc Labs Inc. and Contributors
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

package upgrade

import (
	"os"
	"path"

	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/config"
)

// UpgradeFromV024 upgrades from version <= 0.2.4
// upgrade steps:
// range ZINC_DATA_PATH/
// -- mv    index index_old
// -- mkdir index
// -- mv    index_old index/000000
func UpgradeFromV024() error {
	rootPath := config.Global.DataPath
	fs, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if !f.IsDir() {
			continue
		}
		if f.Name() == "_metadata.db" {
			continue
		}
		log.Info().Msgf("Upgrading index: %s", f.Name())
		if err := UpgradeFromV024Index(f.Name()); err != nil {
			return err
		}
	}
	return nil
}

func UpgradeFromV024Index(indexName string) error {
	rootPath := config.Global.DataPath
	if err := os.Rename(path.Join(rootPath, indexName), path.Join(rootPath, indexName+"_old")); err != nil {
		return err
	}
	if err := os.Mkdir(path.Join(rootPath, indexName), 0755); err != nil {
		return err
	}
	if err := os.Rename(path.Join(rootPath, indexName+"_old"), path.Join(rootPath, indexName, "000000")); err != nil {
		return err
	}
	return nil
}
