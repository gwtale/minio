/*
 * Minimalist Object Storage, (C) 2015 Minio, Inc.
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

package donut

import (
	"github.com/minio/minio/pkg/iodine"
	"github.com/minio/minio/pkg/storage/donut/disk"
)

// donut struct internal data
type donut struct {
	name    string
	buckets map[string]bucket
	nodes   map[string]Node
}

// config files used inside Donut
const (
	// donut object metadata and config
	donutObjectMetadataConfig = "donutObjectMetadata.json"
	donutConfig               = "donutMetadata.json"

	// bucket, object metadata
	bucketMetadataConfig = "bucketMetadata.json"
	objectMetadataConfig = "objectMetadata.json"

	// versions
	objectMetadataVersion      = "1.0"
	donutObjectMetadataVersion = "1.0"
)

// attachDonutNode - wrapper function to instantiate a new node for associated donut
// based on the provided configuration
func (d donut) attachDonutNode(hostname string, disks []string) error {
	node, err := NewNode(hostname)
	if err != nil {
		return iodine.New(err, nil)
	}
	donutName := d.name
	for i, d := range disks {
		// Order is necessary for maps, keep order number separately
		newDisk, err := disk.New(d)
		if err != nil {
			return iodine.New(err, nil)
		}
		if err := newDisk.MakeDir(donutName); err != nil {
			return iodine.New(err, nil)
		}
		if err := node.AttachDisk(newDisk, i); err != nil {
			return iodine.New(err, nil)
		}
	}
	if err := d.AttachNode(node); err != nil {
		return iodine.New(err, nil)
	}
	return nil
}

// NewDonut - instantiate a new donut
func NewDonut(donutName string, nodeDiskMap map[string][]string) (Donut, error) {
	if donutName == "" || len(nodeDiskMap) == 0 {
		return nil, iodine.New(InvalidArgument{}, nil)
	}
	nodes := make(map[string]Node)
	buckets := make(map[string]bucket)
	d := donut{
		name:    donutName,
		nodes:   nodes,
		buckets: buckets,
	}
	for k, v := range nodeDiskMap {
		if len(v) == 0 {
			return nil, iodine.New(InvalidDisksArgument{}, nil)
		}
		err := d.attachDonutNode(k, v)
		if err != nil {
			return nil, iodine.New(err, nil)
		}
	}
	return d, nil
}
